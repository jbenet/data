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
      auth <username>    Re-authenticate as user.
      pass [<username>]  Change user password.
      info [<username>]  Show (or edit) public user information.
      url [<username>]   Output user profile url.

    User accounts are needed in order to publish dataset packages to the
    dataset index. Packages are listed under their owner's username:
    '<owner>/<dataset>'.
  `,
	Subcommands: []*commander.Command{
		cmd_data_user_add,
		cmd_data_user_pass,
		cmd_data_user_info,
		cmd_data_user_url,
	},
}

var cmd_data_user_add = &commander.Command{
	UsageLine: "add [<username>]",
	Short:     "Register new user with index.",
	Long: `data user add - Register new user with index.

    Guided process to register a new user account with dataset index.

    See data user.
  `,
	Run: userAddCmd,
}

var cmd_data_user_pass = &commander.Command{
	UsageLine: "pass [<username>]",
	Short:     "Change user password.",
	Long: `data user pass - Change user password.

    Guided process to change user account password with dataset index.

    See data user.
  `,
	Run: userPassCmd,
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

func userAddCmd(c *commander.Command, args []string) error {
	ui, err := userCmdUserIndex(args)
	if err != nil {
		return err
	}

	pass, err := inputNewPassword()
	if err != nil {
		return err
	}

	email, err := inputNewEmail()
	if err != nil {
		return err
	}

	err = ui.Add(pass, email)
	if err != nil {
		return err
	}

	pOut("%s registered.\n", ui.User)
	return nil
}

func userPassCmd(c *commander.Command, args []string) error {
	ui, err := userCmdUserIndex(args)
	if err != nil {
		return err
	}

	pOut("Current Password: ")
	curp, err := readInputSilent()
	if err != nil {
		return err
	}

	pOut("New ")
	newp, err := inputNewPassword()
	if err != nil {
		return err
	}

	err = ui.Pass(curp, newp)
	if err != nil {
		return err
	}

	pOut("Password changed. You will receive an email notification.\n")
	return nil
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

const PasswordMinLength = 6

func inputNewPassword() (string, error) {
	var pass string
	for len(pass) < PasswordMinLength {
		pOut("Password (%d char min): ", PasswordMinLength)
		var err error
		pass, err = readInputSilent()
		if err != nil {
			return "", err
		}
	}
	return pass, nil
}

func inputNewEmail() (string, error) {
	var email string

	for !EmailRegexp.MatchString(email) {
		pOut("Email (for security): ")
		var err error

		email, err = readInput()
		if err != nil {
			return "", err
		}
	}
	return email, nil
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

func (i UserIndex) Passhash(pass string) (string, error) {
	// additional hashing of the password before sending.
	// this resulting `passhash` is really the user's password.
	// this is so that passwords are never seen by the server as plaintext
	return stringHash(pass + i.User)
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
	return i.post("user/info", p)
}

func (i *UserIndex) post(url string, body interface{}) error {
	r, err := Marshal(body)
	if err != nil {
		return err
	}

	resp, err := httpPost(i.Url(url), "application/yaml", r)
	if err != nil {
		return err
	}

	resp.Body.Close()
	return nil
}

func (i *UserIndex) Pass(cp string, np string) error {
	cph, err := i.Passhash(cp)
	if err != nil {
		return err
	}

	nph, err := i.Passhash(np)
	if err != nil {
		return err
	}

	return i.post("user/pass", &NewPassMsg{cph, nph})
}

func (i *UserIndex) Auth(pass string) error {
	ph, err := i.Passhash(pass)
	if err != nil {
		return err
	}

	return i.post("user/auth", ph)
}

func (i *UserIndex) Add(pass string, email string) error {
	ph, err := i.Passhash(pass)
	if err != nil {
		return err
	}

	err = i.post("user/add", &NewUserMsg{ph, email})
	if err != nil {
		if strings.Contains(err.Error(), "user exists") {
			m := "Error: username '%s' already in use. Try another."
			return fmt.Errorf(m, i.User)
		}
	}
	return err
}

// DataIndex extension to generate a UserIndex
func (d *DataIndex) UserIndex(user string) *UserIndex {
	return &UserIndex{
		User:    user,
		BaseUrl: d.Url,
	}
}

type NewUserMsg struct {
	Pass  string
	Email string
}

type NewPassMsg struct {
	Current string
	New     string
}
