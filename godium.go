package main

import (
	"fmt"

	"io/ioutil"

	"os"

	"errors"

	"github.com/ericaro/frontmatter"
	medium "github.com/medium/medium-sdk-go"
	"github.com/mitchellh/go-homedir"
	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli"
)

const TokenFileName string = "~/.godium"

type postData struct {
	Title   string   `fm:"title"`
	Tags    []string `fm:"tags"`
	Content string   `fm:"content"`
}

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		addTokenCommand(),
		infoCommand(),
		publishCommand(),
	}
	app.HideVersion = true
	app.Copyright = "MIT"
	app.Usage = "Interact with the Medium publishing platform through the command line."

	app.Run(os.Args)
}

func addTokenCommand() cli.Command {
	return cli.Command{
		Name:    "set-token",
		Aliases: []string{"st"},
		Usage:   "Add an integration token to be used by this application",
		Action: func(c *cli.Context) error {
			filename, err := homedir.Expand(TokenFileName)
			if err != nil {
				return err
			}

			token := c.Args().First()
			err = ioutil.WriteFile(filename, []byte(token), 0644)
			if err != nil {
				return err
			}

			fmt.Printf("Token was written into %s successfully.", filename)
			return nil
		},
	}
}

func infoCommand() cli.Command {
	return cli.Command{
		Name:  "info",
		Usage: "Get the information for the owner of the access token",
		Action: func(c *cli.Context) error {
			accessToken, err := getAccessToken()
			if err != nil {
				return err
			}

			m := medium.NewClientWithAccessToken(accessToken)
			user, err := m.GetUser()
			if err != nil {
				return err
			}

			fmt.Println("User ID:", user.ID)
			fmt.Println("Username:", user.Username)
			fmt.Println("Name:", user.Name)
			fmt.Println("Profile URL:", user.URL)

			return nil
		},
	}
}

func publishCommand() cli.Command {
	return cli.Command{
		Name:    "publish",
		Aliases: []string{"p"},
		Usage:   "Publish a markdown file to Medium with the status set to " + string(medium.PublishStatusDraft) + " and open the post editor page in the browser",
		Action:  publishAction,
	}
}

func getAccessToken() (string, error) {
	filename, err := homedir.Expand(TokenFileName)
	if err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadFile(filename)

	// If the file is not found, it means we have no tokens. Try and make it more user-friendly.
	if err != nil && os.IsNotExist(err) {
		err := open.Run("https://medium.com/me/settings")
		if err != nil {
			return "", errors.New("We tried to open your browser for you automatically, but for some reason it failed. Please manually browse to https://medium.com/me/settings, generate an integration token, and use the `godium set-token <token>` command to add it.")
		}

		return "", errors.New("Could not find the token. We have opened your browser for you. Please generate an integration token in your browser window, and use `godium set-token <token>` to add it.")
	}

	// It's not a file not found error, so just return it
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func createPostOptionsFromFile(filename string) (*medium.CreatePostOptions, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	postData := &postData{}
	if err := frontmatter.Unmarshal(contents, postData); err != nil {
		return nil, err
	}

	createPostOptions := &medium.CreatePostOptions{
		Title:         postData.Title,
		ContentFormat: medium.ContentFormatMarkdown,
		Content:       postData.Content,
		Tags:          postData.Tags,
		PublishStatus: medium.PublishStatusDraft,
	}

	return createPostOptions, nil
}

func publishAction(c *cli.Context) error {
	createPostOptions, err := createPostOptionsFromFile(c.Args().First())
	if err != nil {
		return err
	}

	accessToken, err := getAccessToken()
	if err != nil {
		return err
	}

	m := medium.NewClientWithAccessToken(accessToken)
	user, err := m.GetUser()
	if err != nil {
		return err
	}

	createPostOptions.UserID = user.ID
	post, err := m.CreatePost(*createPostOptions)
	if err != nil {
		return err
	}

	fmt.Println("Post created. URL is ", post.URL)
	open.Run(post.URL)

	return nil
}
