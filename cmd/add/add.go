// package main

//	func main() {
//		println("Hello, World!")
//	}
package add

import (
	"github.com/alex-way/changesets/pkg/changeset"
	"github.com/alex-way/changesets/pkg/version"
	"github.com/charmbracelet/huh"
	"github.com/urfave/cli/v2"
)

type AddCtx struct {
	Type    version.IncrementType
	Message string
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

func getTypeOrPrompt(cCtx *cli.Context) (version.IncrementType, error) {
	type_ := cCtx.String("type")
	var changeset_type version.IncrementType

	if type_ != "" {
		parsed_type, err := version.ParseIncrementType(type_)
		if err != nil {
			return 0, cli.Exit(err, 1)
		}
		changeset_type = parsed_type
	} else {

		err := huh.NewSelect[version.IncrementType]().
			Title("Type of change").
			Options(
				huh.NewOption("Major", version.Major),
				huh.NewOption("Minor", version.Minor),
				huh.NewOption("Patch", version.Patch),
			).
			Value(&changeset_type).
			Run()

		if err != nil {
			return 0, cli.Exit(err, 1)
		}
	}
	return changeset_type, nil
}

func Run(cCtx *cli.Context) error {
	changeset_type, err := getTypeOrPrompt(cCtx)
	if err != nil {
		return cli.Exit(err, 1)
	}
	message, err := getMessageOrPrompt(cCtx)
	if err != nil {
		return cli.Exit(err, 1)
	}

	changeset_filepath := changeset.CreateChangeset(changeset_type, message)

	println("Created changeset " + changeset_filepath)
	println("You can now edit the file and commit it to version control.")
	return nil
}
