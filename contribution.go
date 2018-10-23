package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"gopkg.in/yaml.v2"
)

type Contribution struct {
	Date    time.Time
	Project string
	Type    string
	User    string
	URL     string
}

type ConfigFile struct {
	Handles []string
	Orgs    []string
	Repos   map[string]string
}

func Testing() string {
	return "Testing Contribution Package"
}

func getConfigFile() *ConfigFile {
	bucket := "osscontributions-eksteam"
	item := "configs.yaml"
	var fileData ConfigFile

	sess := session.New()
	svc := s3.New(sess, aws.NewConfig().WithRegion("us-west-2"))

	result, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(item),
	})
	if err != nil {
		fmt.Printf("ERROR GETTING ITEM: %v\n", err)
	}
	defer result.Body.Close()

	if err := yaml.NewDecoder(result.Body).Decode(&fileData); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}

	return &fileData
}

func collectContributions() []Contribution {

	cf := getConfigFile()
	usrs := cf.Handles
	orgs := cf.Orgs
	repos := cf.Repos

	nc := newClient()

	var contributions []Contribution

	for _, o := range orgs {
		repos := getRepos(nc, o)

		for _, r := range repos {

			for _, usr := range usrs {
				iss := getIssues(nc, r, usr)
				for _, i := range iss {
					if i.IsPullRequest() {
						newCont := Contribution{
							Date:    i.GetCreatedAt(),
							Project: r.GetFullName(),
							Type:    "Pull Request",
							User:    i.User.GetLogin(),
							URL:     i.GetHTMLURL(),
						}
						contributions = append(contributions, newCont)
					} else {
						newCont := Contribution{
							Date:    i.GetCreatedAt(),
							Project: r.GetFullName(),
							Type:    "Issue",
							User:    i.User.GetLogin(),
							URL:     i.GetHTMLURL(),
						}
						contributions = append(contributions, newCont)
					}

				}
			}
		}
	}

	for owner, repo := range repos {

		r := getRepo(nc, owner, repo)

		for _, usr := range usrs {
			iss := getIssues(nc, r, usr)
			for _, i := range iss {

				if i.IsPullRequest() {
					newCont := Contribution{
						Date:    i.GetCreatedAt(),
						Project: r.GetFullName(),
						Type:    "Pull Request",
						User:    i.User.GetLogin(),
						URL:     i.GetHTMLURL(),
					}
					contributions = append(contributions, newCont)
				} else {
					newCont := Contribution{
						Date:    i.GetCreatedAt(),
						Project: r.GetFullName(),
						Type:    "Issue",
						User:    i.User.GetLogin(),
						URL:     i.GetHTMLURL(),
					}
					contributions = append(contributions, newCont)
				}
			}
		}

	}

	return contributions
}
