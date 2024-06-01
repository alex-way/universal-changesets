package get_version

import (
	"fmt"

	cmdver "github.com/alex-way/changesets/cmd/version"
	"github.com/urfave/cli/v2"
)

func Run(cCtx *cli.Context) error {
	current_version, err := cmdver.ReadVersionFile()
	if err != nil {
		return cli.Exit(err, 1)
	}

	fmt.Println(current_version.String())
	return nil
}
