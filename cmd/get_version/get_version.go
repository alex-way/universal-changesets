package get_version

import (
	"context"
	"fmt"

	wasm "github.com/alex-way/changesets/pkg"
	"github.com/urfave/cli/v2"

	"github.com/alex-way/changesets/pkg/config"
	"github.com/alex-way/changesets/pkg/plugin"
	"github.com/alex-way/changesets/pkg/version"
)

func GetVersion() (version.Version, error) {
	_config, err := config.GetConfig()
	if err != nil {
		return version.Version{}, err
	}

	handler := &wasm.Runner{
		Plugin: _config.Plugin,
	}
	client := plugin.NewVersionGetterSetterServiceClient(handler)

	req := &plugin.RequestMessage{
		Request: &plugin.RequestMessage_GetVersion{
			GetVersion: &plugin.GetVersionRequest{
				FilePath: _config.Plugin.VersionedFile,
			},
		},
	}

	ctx := context.Background()
	resp, err := client.Request(ctx, req)
	if err != nil {
		message := fmt.Sprintf("failed to get version: %v", err)
		return version.Version{}, fmt.Errorf(message)
	}

	if resp.Status.Code != 0 {
		message := fmt.Sprintf(resp.Status.Message)
		return version.Version{}, fmt.Errorf(message)
	}

	unparsed_version := resp.Response.(*plugin.Response_GetVersion).GetVersion.Version
	return version.ParseVersion(unparsed_version)
}

func Run(cCtx *cli.Context) error {
	version, err := GetVersion()
	if err != nil {
		return cli.Exit(err, 1)
	}

	println(version.String())

	return nil
}
