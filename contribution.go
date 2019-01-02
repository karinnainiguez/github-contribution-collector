package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	yaml "gopkg.in/yaml.v2"
)

type ContributionCollection []Contribution

type Contribution struct {
	Date        time.Time
	Project     string
	Type        string
	User        string
	URL         string
	Description string
}

type ConfigFile struct {
	Handles []string
	Orgs    []string
	Repos   map[string]string
}

func getConfigFile() *ConfigFile {
	bucket := os.Getenv("S3BucketName")
	item := os.Getenv("S3ObjectName")
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

func getLocalConfigFile(pathString string) *ConfigFile {
	var cf ConfigFile
	absPath, _ := filepath.Abs(pathString)

	reader, err := os.Open(absPath)
	handle(err)

	buf, _ := ioutil.ReadAll(reader)

	err = yaml.Unmarshal(buf, &cf)
	handle(err)

	return &cf

}

func collectContributionsConcurrently() (contributions ContributionCollection, err error) {

	cf := getConfigFile()
	usrs := cf.Handles
	orgs := cf.Orgs
	repos := cf.Repos
	nc := newClient()

	// make tokens channel
	tokens := make(chan string, 200)
	// make respChannel
	respChan := make(chan issueResponse)

	// loop through and send goroutine for each user/repo and keep track of num
	routines := 0
	for _, o := range orgs {
		repos := getRepos(nc, o)
		for _, r := range repos {
			for _, usr := range usrs {
				routines++
				go getIssuesConcurrently(nc, r, usr, respChan, tokens)
			}
		}
	}

	for repo, owner := range repos {
		r := getRepo(nc, owner, repo)
		for _, usr := range usrs {
			routines++
			go getIssuesConcurrently(nc, r, usr, respChan, tokens)
		}

	}

	// create new contributionCollection
	possibleErrors := make(map[string]int)

	// combine all data
	combineResponses(respChan, &contributions, routines, possibleErrors)
	if len(possibleErrors) > 0 {
		var sb strings.Builder
		for strErr, occ := range possibleErrors {
			str := fmt.Sprintf("%v: %v\n", strErr, occ)
			sb.WriteString(str)
		}
		err = errors.New(sb.String())
	}

	// return combined data
	return contributions, err
}

func collectContributionsLocallyConcurrently(pathString string) (contributions ContributionCollection, err error) {

	cf := getLocalConfigFile(pathString)
	usrs := cf.Handles
	orgs := cf.Orgs
	repos := cf.Repos
	nc := newClient()

	// make tokens channel
	tokens := make(chan string, 200)
	// make respChannel
	respChan := make(chan issueResponse)

	// loop through and send goroutine for each user/repo and keep track of num
	routines := 0
	for _, o := range orgs {
		repos := getRepos(nc, o)
		for _, r := range repos {
			for _, usr := range usrs {
				routines++
				go getIssuesConcurrently(nc, r, usr, respChan, tokens)
			}
		}
	}

	for owner, repo := range repos {
		r := getRepo(nc, owner, repo)
		for _, usr := range usrs {
			routines++
			go getIssuesConcurrently(nc, r, usr, respChan, tokens)
		}

	}

	// create new contributionCollection and possibleErrors
	possibleErrors := make(map[string]int)

	// combine all data
	combineResponses(respChan, &contributions, routines, possibleErrors)
	if len(possibleErrors) > 0 {
		var sb strings.Builder
		for strErr, occ := range possibleErrors {
			str := fmt.Sprintf("%v: %v\n", strErr, occ)
			sb.WriteString(str)
		}
		err = errors.New(sb.String())
	}

	// return combined data
	return contributions, err
}

func combineResponses(
	respChan chan issueResponse,
	collection *ContributionCollection,
	desired int,
	errCollection map[string]int) {

	for i := 0; i < desired; i++ {
		resp := <-respChan
		for _, i := range resp.issues {
			if i.IsPullRequest() {
				newCont := Contribution{
					Date:        i.GetCreatedAt(),
					Project:     obtainRepoName(i.GetHTMLURL()),
					Type:        "Pull Request",
					User:        i.User.GetLogin(),
					URL:         i.GetHTMLURL(),
					Description: i.GetTitle(),
				}
				*collection = append(*collection, newCont)

			} else {
				newCont := Contribution{
					Date:        i.GetCreatedAt(),
					Project:     obtainRepoName(i.GetHTMLURL()),
					Type:        "Issue",
					User:        i.User.GetLogin(),
					URL:         i.GetHTMLURL(),
					Description: i.GetTitle(),
				}
				*collection = append(*collection, newCont)
			}
		}
		if resp.err != nil {
			errCollection[resp.err.Error()]++
		}
	}
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

func obtainRepoName(url string) string {
	arr := strings.Split(url, "/")
	if len(arr) < 5 {
		return ""
	}
	return arr[3] + "/" + arr[4]
}
