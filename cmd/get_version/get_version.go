package get_version

import (
	"context"
	"fmt"

	wasm "github.com/alex-way/changesets/pkg"
	"github.com/urfave/cli/v2"

	"github.com/alex-way/changesets/pkg/config"
	"github.com/alex-way/changesets/pkg/plugin"
)

func Run(cCtx *cli.Context) error {
	_config, err := config.GetConfig()
	if err != nil {
		return cli.Exit(err, 1)
	}

	handler := &wasm.Runner{
		Plugin: _config.Plugin,
	}
	client := plugin.NewVersionGetterSetterServiceClient(handler)

	input := &plugin.File{
		Path: _config.Plugin.VersionedFile,
	}

	req := &plugin.GetVersionRequest{
		Inputs: input,
	}

	ctx := context.Background()
	resp, err := client.GetVersion(ctx, req)
	if err != nil {
		message := fmt.Sprintf("failed to get version: %v", err)
		return cli.Exit(message, 1)
	}

	print(resp.Version)

	return nil
}
