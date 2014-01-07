package data

import (
	"github.com/jbenet/commander"
)

var cmd_data_publish = &commander.Command{
	UsageLine: "publish",
	Short:     "Guided dataset publishing.",
	Long: `data publish - Guided dataset publishing.

    This command guides the user through the necessary steps to
    create a data package (Datafile and Manifest), uploads it,
    and publishes it to the dataset index.

    See 'data pack'.
  `,
	Run: publishCmd,
}

func publishCmd(c *commander.Command, args []string) error {

	pOut("==> Guided Data Package Publishing.\n")

	pOut("\n==> Step 1/3: Creating the package.\n")
	err := packMakeCmd(c, []string{})
	if err != nil {
		return err
	}

	pOut("\n==> Step 2/3: Uploading the package contents.\n")
	err = packUploadCmd(c, []string{})
	if err != nil {
		return err
	}

	pOut("\n==> Step 3/3: Publishing the package to the index.\n")
	return packPublishCmd(c, []string{})
}
