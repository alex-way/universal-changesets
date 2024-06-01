package add

import (
	"github.com/alex-way/changesets/pkg/changeset"
	"github.com/alex-way/changesets/pkg/version"
	"github.com/charmbracelet/huh"
	"github.com/urfave/cli/v2"
)

type AddCtx struct {
	BumpType version.BumpType
	Message  string
}

func getMessageOrPrompt(cCtx *cli.Context) (string, error) {
	message := cCtx.String("message")
	if message == "" {
		err := huh.NewInput().
			Title("Message").
			Value(&message).
			Run()
		if err != nil {
			return "", err
		}
	}
	return message, nil
}

func getBumpTypeOrPrompt(cCtx *cli.Context) (version.BumpType, error) {
	type_ := cCtx.String("type")
	var bump_type version.BumpType

	if type_ != "" {
		parsed_type, err := version.ParseBumpType(type_)
		if err != nil {
			return 0, cli.Exit(err, 1)
		}
		bump_type = parsed_type
	} else {

		err := huh.NewSelect[version.BumpType]().
			Title("Type of change").
			Options(
				huh.NewOption("Major", version.Major),
				huh.NewOption("Minor", version.Minor),
				huh.NewOption("Patch", version.Patch),
				huh.NewOption("Other", version.None),
			).
			Value(&bump_type).
			Run()

		if err != nil {
			return 0, cli.Exit(err, 1)
		}
	}
	return bump_type, nil
}

func Run(cCtx *cli.Context) error {
	bump_type, err := getBumpTypeOrPrompt(cCtx)
	if err != nil {
		return cli.Exit(err, 1)
	}
	message, err := getMessageOrPrompt(cCtx)
	if err != nil {
		return cli.Exit(err, 1)
	}

	changeset_filepath, err := changeset.CreateChangeFile(bump_type, message)
	if err != nil {
		return cli.Exit(err, 1)
	}

	println("Created changeset " + changeset_filepath)
	println("You can now edit the file and commit it to version control.")
	return nil
}
