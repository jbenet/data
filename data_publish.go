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
	Flag: *flag.NewFlagSet("data-pack-publish", flag.ExitOnError),
}

func init() {
	cmd_data_publish.Flag.Bool("clean", true,
		"rebuild manifest (data pack make --clean)")
	cmd_data_publish.Flag.Bool("force", false,
		"force publish (data pack publish --force)")
}

func publishCmd(c *commander.Command, args []string) error {
	u := configUser()
	if !isNamedUser(u) {
		return fmt.Errorf(NotLoggedInErr)
	}

	pOut("==> Guided Data Package Publishing.\n")
	pOut(PublishMsgWelcome)

	pOut("\n==> Step 1/3: Creating the package.\n")
	pOut(PublishMsgDatafile)
	err := packMakeCmd(c, []string{})
	if err != nil {
		return err
	}

	pOut("\n==> Step 2/3: Uploading the package contents.\n")
	pOut(PublishMsgUpload)
	err = packUploadCmd(c, []string{})
	if err != nil {
		return err
	}

	pOut("\n==> Step 3/3: Publishing the package to the index.\n")
	pOut(PublishMsgPublish)
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

const PublishMsgWelcome = `
Welcome to Data Package Publishing. You should read these short
messages carefully, as they contain important information about
how data works, and how your data package will be published.

First, a 'data package' is a collection of files, containing:
- various files with your data, in any format.
- 'Datafile', a file with descriptive information about the package.
- 'Manifest', a file listing the other files in the package and their checksums.

This tool will automatically:
1. Create the package
  - Generate a 'Datafile', with information you will provide.
  - Generate a 'Manifest', with all the files in the current directory.
2. Upload the package contents
3. Publish the package to the index

(Note: to specify which files are part of the package, and other advanced
 features, use the 'data pack' command directly. See 'data pack help'.)

`

const PublishMsgDatafile = `
First, let's write the package's Datafile, which contains important
information about the package. The 'owner id' is the username of the
package's owner (usually your username). The 'dataset id' is the identifier
which defines this dataset. Good 'dataset ids' are like names: short, unique,
and memorable. For example: "mnist" or "cifar". Choose it carefully.

`

const PublishMsgUpload = `
Now, data will upload the contents of the package (this directory) to the index
sotrage service. This may take a while, if the files are large (over 100MB).

`

const PublishMsgPublish = `
Finally, data will publish the package to the index, where others can find
and download your package. The index is available through data, and on the web.

`
