package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	"github.com/urfave/cli/v2"
)

func main() {
	var filename, owner, repo string
	var delay int64
	app := &cli.App{
		Name:        "repo-admins",
		Usage:       "query repository admins",
		Version:     "1.0.0",
		Description: "gh repo-admins --owner [owner] --repo [repo]\ngh repo-admins --help",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Aliases:     []string{"o"},
				Destination: &owner,
				Name:        "owner",
				Required:    true,
				Usage:       "organization or user that owns the repository",
			},
			&cli.StringFlag{
				Aliases:     []string{"r"},
				Destination: &repo,
				Name:        "repo",
				Required:    true,
				Usage:       "repository name",
			},
			&cli.StringFlag{
				Aliases:     []string{"f"},
				Destination: &filename,
				Name:        "file",
				Required:    true,
				Usage:       "name of output file",
			},
			&cli.Int64Flag{
				Aliases:     []string{"d"},
				Destination: &delay,
				Name:        "delay",
				Required:    false,
				Usage:       "delay between GitHub API requests in milliseconds. If you are hitting API rate limits, increase this value.",
				Value:       500,
			},
		},
		Action: func(*cli.Context) error {
			delay := time.Duration(delay) * time.Millisecond

			restClient, err := gh.RESTClient(nil)
			if err != nil {
				return fmt.Errorf("error creating REST client: %v", err)
			}

			log.Printf("Retrieving teams for %s/%s", owner, repo)
			teams, err := retrieveTeams(restClient, owner, repo, delay)
			if err != nil {
				return err
			}

			var users []*user
			for _, team := range teams {
				if team.Permission == "admin" {
					log.Printf("Retrieving members for team: %s", team.Name)
					members, err := retrieveTeamMembers(restClient, owner, team.Slug, delay)
					if err != nil {
						return err
					}
					for _, member := range members {
						log.Printf("Retrieving user: %s", member.Login)
						user, err := retrieveUser(restClient, member.Login, delay)
						if err != nil {
							return err
						}
						users = append(users, user)
					}
				}
			}

			if len(users) == 0 {
				return fmt.Errorf("no users found for %s/%s", owner, repo)
			}

			log.Printf("Writing users to %s", filename)
			return writeUsersToCSV(filename, users)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

type user struct {
	Login string
	Name  string
	Email string
}

type team struct {
	Name       string
	Slug       string
	Permission string
}

func retrieveTeams(client api.RESTClient, owner, repo string, delay time.Duration) ([]*team, error) {
	var teams []*team
	url := fmt.Sprintf("repos/%s/%s/teams?per_page=100", owner, repo)
	err := client.Get(url, &teams)
	if err != nil {
		return nil, fmt.Errorf("error retrieving teams: %v", err)
	}
	time.Sleep(delay)
	return teams, nil
}

func retrieveTeamMembers(client api.RESTClient, org, slug string, delay time.Duration) ([]*user, error) {
	var allMembers []*user
	var page = 1
	for {
		var members []*user
		url := fmt.Sprintf("orgs/%s/teams/%s/members?page=%d&per_page=100", org, slug, page)
		err := client.Get(url, &members)
		if err != nil {
			return nil, fmt.Errorf("error retrieving team members: %v", err)
		}
		time.Sleep(delay)
		if len(members) == 0 {
			break
		}
		allMembers = append(allMembers, members...)
		page = page + 1
	}
	return allMembers, nil
}

func retrieveUser(client api.RESTClient, login string, delay time.Duration) (*user, error) {
	var user *user
	url := fmt.Sprintf("users/%s", login)
	err := client.Get(url, &user)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user: %v", err)
	}
	time.Sleep(delay)
	return user, nil
}

func writeUsersToCSV(filename string, users []*user) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	err = w.Write([]string{"Username", "Name", "Email"})
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	for _, user := range users {
		err = w.Write([]string{user.Login, user.Name, user.Email})
		if err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}
	}
	return nil
}
