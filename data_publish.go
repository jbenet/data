package data

import (
	"fmt"
	"github.com/gonuts/flag"
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
	Run:  publishCmd,
	Flag: *flag.NewFlagSet("data-pack-make", flag.ExitOnError),
}

func init() {
	cmd_data_publish.Flag.Bool("clean", true,
		"rebuild manifest (data pack make --clean)")
}

func publishCmd(c *commander.Command, args []string) error {
	pOut("==> Guided Data Package Publishing.\n")

	u := configUser()
	if !isNamedUser(u) {
		return fmt.Errorf(NotLoggedInErr)
	}

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

const NotLoggedInErr = `You are not logged in. First, either:

- Run 'data user add' to create a new user account.
- Run 'data user auth' to log in to an existing user account.


Why does publishing require a registered user account (and email)? The index
service needs to distinguish users to perform many of its tasks. For example:

- Verify who can or cannot publish datasets, or modify already published ones.
  (i.e. the creator + collaborators should be able to, others should not).
- Profiles credit people for the datasets they have published.
- Malicious users can be removed, and their email addresses blacklisted to
  prevent further abuse.
`
