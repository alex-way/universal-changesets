package version

import (
	"fmt"

	"github.com/alex-way/changesets/pkg/changeset"
	"github.com/alex-way/changesets/pkg/version"
	"github.com/urfave/cli/v2"
)

func Run(cCtx *cli.Context) error {
	changes, err := changeset.GetChanges()
	if err != nil {
		return cli.Exit(err, 1)
	}

	if len(changes) == 0 {
		println("No changesets found. Please run 'changeset add' to add changes.")
		return nil
	}

	_changeset := changeset.Changeset{
		CurrentVersion: version.Version{Major: 0, Minor: 0, Patch: 0},
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

	println("Changeset consumed successfully.")

	return nil
}
