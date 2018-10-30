package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"gopkg.in/yaml.v2"
)

type ContributionCollection []Contribution

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
	handle(err)

	defer result.Body.Close()

	if err := yaml.NewDecoder(result.Body).Decode(&fileData); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}

	return &fileData
}

func collectContributions() ContributionCollection {

	cf := getConfigFile()
	usrs := cf.Handles
	orgs := cf.Orgs
	repos := cf.Repos

	nc := newClient()

	var contributions ContributionCollection

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

func (c ContributionCollection) filterMonthlyContributions() ContributionCollection {
	var filtered ContributionCollection

	startDate, err := time.Parse("01-02-2006", yesterdayFrom())
	handle(err)
	endDate, err := time.Parse("01-02-2006", yesterdayUntil())
	handle(err)

	for _, cont := range c {
		if cont.Date.After(startDate) && cont.Date.Before(endDate) {
			filtered = append(filtered, cont)
		}
	}
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Date.Before(filtered[j].Date)
	})
	return filtered
}

func (c ContributionCollection) filterContributions(startDateString string, endDateString string) ContributionCollection {
	var filtered ContributionCollection

	startDate, err := time.Parse("01-02-2006", startDateString)
	handle(err)
	endDate, err := time.Parse("01-02-2006", endDateString)
	handle(err)

	for _, cont := range c {
		if cont.Date.After(startDate) && cont.Date.Before(endDate) {
			filtered = append(filtered, cont)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Date.Before(filtered[j].Date)
	})

	return filtered
}
