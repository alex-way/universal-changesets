package version

import (
	"context"
	"fmt"

	"github.com/alex-way/changesets/cmd/get_version"
	wasm "github.com/alex-way/changesets/pkg"
	"github.com/alex-way/changesets/pkg/changeset"
	"github.com/alex-way/changesets/pkg/config"
	"github.com/alex-way/changesets/pkg/plugin"
	"github.com/alex-way/changesets/pkg/version"
	"github.com/urfave/cli/v2"
)

func setVersion(version version.Version) error {
	_config, err := config.GetConfig()
	if err != nil {
		return cli.Exit(err, 1)
	}

	handler := &wasm.Runner{
		Plugin: _config.Plugin,
	}
	client := plugin.NewVersionGetterSetterServiceClient(handler)

	req := &plugin.RequestMessage{
		Request: &plugin.RequestMessage_SetVersion{
			SetVersion: &plugin.SetVersionRequest{
				FilePath: _config.Plugin.VersionedFile,
				Version:  version.String(),
			},
		},
	}

	ctx := context.Background()
	resp, err := client.Request(ctx, req)
	if err != nil {
		message := fmt.Sprintf("failed to set version: %v", err)
		return cli.Exit(message, 1)
	}

	if resp.Status.Code != 0 {
		message := fmt.Sprintf(resp.Status.Message)
		return cli.Exit(message, 1)
	}

	println(("You got it boss"))

	return nil
}

func Run(cCtx *cli.Context) error {
	changes, err := changeset.GetChanges()
	if err != nil {
		return cli.Exit(err, 1)
	}

	if len(changes) == 0 {
		println("No changesets found. Please run 'changeset add' to add changes.")
		return nil
	}

	current_version, err := get_version.GetVersion()
	if err != nil {
		return cli.Exit(err, 1)
	}

	_changeset := changeset.Changeset{
		CurrentVersion: current_version,
		Changes:        changes,
	}

	final_bump_type := _changeset.DetermineFinalBumpType()

	if final_bump_type == version.None {
		println(fmt.Sprintf("The version will remain at %s as all changes are not version impacting.", _changeset.CurrentVersion.String()))
		return nil
	}

	next_version := _changeset.DetermineNextVersion()
	println(fmt.Sprintf("The version will be bumped to: `%s` because a %s change was determined from the changes.", next_version.String(), final_bump_type.String()))

	if cCtx.Bool("dry-run") {
		return nil
	}

	_, err = _changeset.ConsumeChanges()

	if err != nil {
		return cli.Exit(err, 1)
	}

	if err := setVersion(next_version); err != nil {
		return cli.Exit(err, 1)
	}

	println("Changeset consumed successfully.")

	return nil
}
