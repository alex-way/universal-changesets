package wasm

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"github.com/tetratelabs/wazero/sys"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/alex-way/changesets/pkg/config"
	"github.com/alex-way/changesets/pkg/plugin"
)

var flight singleflight.Group

type runtimeAndCode struct {
	rt   wazero.Runtime
	code wazero.CompiledModule
}

type Runner struct {
	Plugin config.Plugin
}

// RestrictedFS is a custom file system implementation that restricts access to a specific file.
type RestrictedFS struct {
	target string
}

// NewRestrictedFS creates a new instance of RestrictedFS for the specified directory and target file.
func NewRestrictedFS(target string) fs.FS {
	return &RestrictedFS{target: target}
}

// Open opens the specified file if it matches the target file; otherwise, returns an error.
func (fsys *RestrictedFS) Open(name string) (fs.File, error) {
	if name != "." && name != fsys.target {
		return nil, fmt.Errorf("access denied to file: %s", name)
	}

	// if the thing trying to be opened is a directory, return a readable file
	if stat, err := os.Stat(name); err == nil && stat.IsDir() {
		println("Opening a directory")
		return os.Open(name)
	}
	println("Opening a file")
	file, err := os.OpenFile(name, os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return file, nil
}

func (fsys *RestrictedFS) Write(name string, data []byte) (int, error) {
	return 0, fmt.Errorf("write not allowed")
}

// Attempts to fetch the wasm file from either a URL or a local file depending on the prefix of the URL
// Returns the bytes of the wasm file, the sha256 of the wasm file, and any error
func (r *Runner) fetch(ctx context.Context, uri string) ([]byte, string, error) {
	var body io.ReadCloser

	switch {
	case strings.HasPrefix(uri, "file://"):
		file, err := os.Open(strings.TrimPrefix(uri, "file://"))
		if err != nil {
			return nil, "", fmt.Errorf("os.Open: %s %w", uri, err)
		}
		body = file

	case strings.HasPrefix(uri, "https://"):
		req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
		if err != nil {
			return nil, "", fmt.Errorf("http.Get: %s %w", uri, err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, "", fmt.Errorf("http.Get: %s %w", r.Plugin.URL, err)
		}
		body = resp.Body

	default:
		return nil, "", fmt.Errorf("unknown scheme: %s", r.Plugin.URL)
	}

	defer body.Close()

	wmod, err := io.ReadAll(body)
	if err != nil {
		return nil, "", fmt.Errorf("readall: %w", err)
	}

	sum := sha256.Sum256(wmod)
	actual_sha := fmt.Sprintf("%x", sum)

	return wmod, actual_sha, nil
}

// Verify the provided sha256 is valid.
func (r *Runner) getChecksum(ctx context.Context) (string, error) {
	if r.Plugin.SHA256 != "" {
		return r.Plugin.SHA256, nil
	}
	// TODO: Add a log line here about something
	_, sum, err := r.fetch(ctx, r.Plugin.URL)
	if err != nil {
		return "", err
	}
	slog.Warn("fetching WASM binary to calculate sha256. Set this value in your config file to prevent unneeded work", "sha256", sum)
	return sum, nil
}

func (r *Runner) loadAndCompile(ctx context.Context) (*runtimeAndCode, error) {
	expected_sha, err := r.getChecksum(ctx)
	if err != nil {
		return nil, err
	}

	currentUser, err := user.Current()
	if err != nil {
		slog.Error("Error:", err)
		return nil, err
	}

	home_dir := currentUser.HomeDir

	cacheDir := filepath.Join(home_dir, ".cache", "changesets")
	value, err, _ := flight.Do(expected_sha, func() (interface{}, error) {
		return r.loadAndCompileWASM(ctx, cacheDir, expected_sha)
	})
	if err != nil {
		return nil, err
	}
	data, ok := value.(*runtimeAndCode)
	if !ok {
		return nil, fmt.Errorf("returned value was not a compiled module")
	}
	return data, nil
}

func (r *Runner) loadAndCompileWASM(ctx context.Context, cache string, expected_sha string) (*runtimeAndCode, error) {
	pluginDir := filepath.Join(cache, expected_sha)
	pluginPath := filepath.Join(pluginDir, "plugin.wasm")
	_, staterr := os.Stat(pluginPath)

	uri := r.Plugin.URL
	if staterr == nil {
		// Load the plugin from the cache instead of fetching it
		uri = "file://" + pluginPath
	}

	wmod, actual_sha, err := r.fetch(ctx, uri)
	if err != nil {
		return nil, err
	}

	if expected_sha != actual_sha {
		return nil, fmt.Errorf("invalid checksum: expected %s, got %s", expected_sha, actual_sha)
	}

	if staterr != nil {
		slog.Debug("plugin not cached, caching now")
		err := os.MkdirAll(pluginDir, 0755)
		if err != nil && !os.IsExist(err) {
			return nil, fmt.Errorf("mkdirall: %w", err)
		}
		if err := os.WriteFile(pluginPath, wmod, 0444); err != nil {
			return nil, fmt.Errorf("cache wasm: %w", err)
		}
	}

	wazeroCache, err := wazero.NewCompilationCacheWithDir(filepath.Join(cache, "wazero"))
	if err != nil {
		return nil, fmt.Errorf("wazero.NewCompilationCacheWithDir: %w", err)
	}

	config := wazero.NewRuntimeConfig().WithCompilationCache(wazeroCache)
	rt := wazero.NewRuntimeWithConfig(ctx, config)

	if _, err := wasi_snapshot_preview1.Instantiate(ctx, rt); err != nil {
		return nil, fmt.Errorf("wasi_snapshot_preview1 instantiate: %w", err)
	}

	// Compile the Wasm binary once so that we can skip the entire compilation
	// time during instantiation.
	code, err := rt.CompileModule(ctx, wmod)
	if err != nil {
		return nil, fmt.Errorf("compile module: %w", err)
	}

	return &runtimeAndCode{rt: rt, code: code}, nil
}

func (r *Runner) Invoke(ctx context.Context, method string, args any, reply any, opts ...grpc.CallOption) error {
	req, ok := args.(protoreflect.ProtoMessage)
	if !ok {
		return status.Error(codes.InvalidArgument, "args isn't a protoreflect.ProtoMessage")
	}

	genReq, ok := req.(*plugin.RequestMessage)
	if ok {
		genReq.ProtoMessage()
		req = genReq
	}

	stdinBlob, err := proto.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to encode codegen request: %w", err)
	}

	runtimeAndCode, err := r.loadAndCompile(ctx)
	if err != nil {
		return fmt.Errorf("loadBytes: %w", err)
	}

	var stderr, stdout bytes.Buffer

	// fileSystem := NewRestrictedFS(r.Plugin.VersionedFile)

	//List all files in the rooted directory
	conf := wazero.NewModuleConfig().
		WithName(r.Plugin.Name).
		WithArgs("plugin.wasm", method).
		WithStdin(bytes.NewReader(stdinBlob)).
		WithStdout(&stdout).
		WithStderr(&stderr).WithFSConfig(wazero.NewFSConfig().WithDirMount(".", "."))

	result, err := runtimeAndCode.rt.InstantiateModule(ctx, runtimeAndCode.code, conf)
	if result != nil {
		defer result.Close(ctx)
	}
	if cerr := checkError(err, stderr); cerr != nil {
		return cerr
	}

	stdoutBlob := stdout.Bytes()

	resp, ok := reply.(protoreflect.ProtoMessage)
	if !ok {
		return fmt.Errorf("reply isn't a GenerateResponse")
	}

	if err := proto.Unmarshal(stdoutBlob, resp); err != nil {
		return err
	}

	return nil
}

func (r *Runner) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func checkError(err error, stderr bytes.Buffer) error {
	if err == nil {
		return err
	}

	if exitErr, ok := err.(*sys.ExitError); ok {
		if exitErr.ExitCode() == 0 {
			return nil
		}
	}

	stderrBlob := stderr.String()
	if len(stderrBlob) > 0 {
		return errors.New(stderrBlob)
	}
	return fmt.Errorf("call: %w", err)
}
