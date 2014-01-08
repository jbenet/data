package data

import (
	"fmt"
	"github.com/jbenet/commander"
	"io"
	"os"
	"strings"
)

var cmd_data_user = &commander.Command{
	UsageLine: "user <command> <username>",
	Short:     "Manage users and credentials.",
	Long: `data user - Manage users and credentials.

    Usage:

      data user <command> <username>

    Commands:

      add [<username>]   Register new user with index.
      auth <username>    Re-authenticates as user.
      pass [<username>]  Changes user password.
      info [<username>]  Show (or edit) public user information.
      url [<username>]   Output user profile url.

    User accounts are needed in order to publish dataset packages to the
    dataset index. Packages are listed under their owner's username:
    '<owner>/<dataset>'.
  `,
	Subcommands: []*commander.Command{
		cmd_data_user_url,
		cmd_data_user_info,
	},
}

var cmd_data_user_info = &commander.Command{
	UsageLine: "info [<username>]",
	Short:     "Show (or edit) public user information.",
	Long: `data user info - Show (or edit) public user information.

    Output or edit the profile information of a user. Note that profiles
    are publicly viewable. User profiles include:

      Full Name
      Email Address
      Github Username
      Twitter Username
      Homepage Url
      Packages List

    See data user.
  `,
	Run: userInfoCmd,
}

var cmd_data_user_url = &commander.Command{
	UsageLine: "url [<username>]",
	Short:     "Output user profile url.",
	Long: `data user url - Output user profile url.

    Output the dataset index url for the profile of user named by <username>.

    See data user.
  `,
	Run: userUrlCmd,
}

func userCmdUserIndex(args []string) (*UserIndex, error) {
	var user string
	var err error

	if len(args) > 0 && len(args[0]) > 0 {
		user = args[0]
	}

	for !UserRegexp.MatchString(user) {
		pOut("Username: ")
		user, err = readInput()
		if err != nil {
			return nil, err
		}
	}

	di, err := NewMainDataIndex()
	if err != nil {
		return nil, err
	}

	ui := di.UserIndex(user)
	return ui, nil
}

func userInfoCmd(c *commander.Command, args []string) error {
	ui, err := userCmdUserIndex(args)
	if err != nil {
		return err
	}

	p, err := ui.GetInfo()
	if err != nil {
		return err
	}

	// entered username. lookup and print out info.
	if len(args) > 0 {
		rdr, err := Marshal(p)
		if err != nil {
			return err
		}

		_, err = io.Copy(os.Stdout, rdr)
		return err
	}

	// no username. edit own profile.
	err = fillOutUserProfile(p)
	if err != nil {
		return err
	}

	err = ui.PostInfo(p)
	if err != nil {
		return err
	}

	pOut("Profile saved.\n")
	return nil
}

func userUrlCmd(c *commander.Command, args []string) error {
	ui, err := userCmdUserIndex(args)
	if err != nil {
		return err
	}

	pOut("%s\n", ui.Url(""))
	return nil
}

// serializable into YAML
type UserProfile struct {
	Name     string
	Email    string
	Github   string   ",omitempty"
	Twitter  string   ",omitempty"
	Homepage string   ",omitempty"
	Packages []string ",omitempty"
}

type UserIndex struct {
	User    string
	BaseUrl string
	Refs    *DatasetRefs
}

func (i UserIndex) Url(url string) string {
	return i.BaseUrl + "/" + i.User + "/" + url
}

func (i *UserIndex) GetInfo() (*UserProfile, error) {
	resp, err := httpGet(i.Url("user/info"))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	profile := &UserProfile{}
	err = Unmarshal(resp.Body, profile)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (i *UserIndex) PostInfo(p *UserProfile) error {
	r, err := Marshal(p)
	if err != nil {
		return err
	}

	resp, err := httpPost(i.Url("user/info"), "application/yaml", r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (i *UserIndex) postPass(url string, pass string) error {

	// additional hashing of the password before sending.
	// this resulting `passhash` is really the user's password.
	// this is so that passwords are never seen by the server as plaintext
	passhash, err := stringHash(pass + i.User)
	if err != nil {
		return err
	}

	resp, err := httpPost(i.Url(url), "text", strings.NewReader(passhash))
	if err != nil {
		return err
	}

	resp.Body.Close()
	return nil
}

func (i *UserIndex) Pass(pass string) error {
	if len(pass) < 8 {
		return fmt.Errorf("data user: password too short. 8 character min.")
	}

	return i.postPass("user/pass", pass)
}

func (i *UserIndex) Auth(pass string) error {
	return i.postPass("user/auth", pass)
}

func (i *UserIndex) Add(pass string) error {
	return i.postPass("user/add", pass)
}

// DataIndex extension to generate a UserIndex
func (d *DataIndex) UserIndex(user string) *UserIndex {
	return &UserIndex{
		User:    user,
		BaseUrl: d.Url,
	}
}
