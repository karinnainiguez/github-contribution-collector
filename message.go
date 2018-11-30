package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

func newMessage(c ContributionCollection) error {
	htmlBody, textBody := buildMonthlyEmail(c)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	handle(err)

	svc := ses.New(sess)

	charset := "UTF-8"

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(os.Getenv("SESVerifiedEmail"))},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(charset),
					Data:    aws.String(htmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(charset),
					Data:    aws.String(textBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charset),
				Data:    aws.String(buildMonthlySubject()),
			},
		},
		Source: aws.String(os.Getenv("SESVerifiedEmail")),
	}

	_, err = svc.SendEmail(input)
	return err
}

func newMessageTo(c ContributionCollection, emailAddress string, startDate string, endDate string) error {
	htmlBody, textBody := buildEmail(c)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	handle(err)

	svc := ses.New(sess)

	charset := "UTF-8"

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(emailAddress)},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(charset),
					Data:    aws.String(htmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(charset),
					Data:    aws.String(textBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charset),
				Data:    aws.String(buildSubject(startDate, endDate)),
			},
		},
		Source: aws.String(os.Getenv("SESVerifiedEmail")),
	}

	_, err = svc.SendEmail(input)
	return err
}

func buildMonthlyEmail(c ContributionCollection) (string, string) {
	var body strings.Builder
	body.WriteString("<br/><h1>GitHub Open Source Contributions Report</h1><br/>This month the team had ")
	body.WriteString(strconv.Itoa(len(c)))
	body.WriteString(" contributions into open source projects.<br/>Below is a table of all contributions for the month<br/><br/>")

	body.WriteString(createTable(c))

	var textBody strings.Builder
	textBody.WriteString("Open Source Contributions Report: Team GitHub Contributions")
	textBody.WriteString("\n\nThis month the team had ")
	textBody.WriteString(strconv.Itoa(len(c)))

	return body.String(), textBody.String()
}

func buildEmail(c ContributionCollection) (string, string) {
	var body strings.Builder
	body.WriteString("<br/><h1>Open Source Contributions Report - Team GitHub Contributions</h1><br/>During this time period the team had ")
	body.WriteString(strconv.Itoa(len(c)))
	body.WriteString(" contributions into open source projects.<br/>Below is a table of all contributions<br/><br/>")

	body.WriteString(createTable(c))

	var textBody strings.Builder
	textBody.WriteString("Open Source Contributions Report: Team GitHub Contributions")
	textBody.WriteString("\n\nThis month the team had ")
	textBody.WriteString(strconv.Itoa(len(c)))

	return body.String(), textBody.String()
}

func createTable(c ContributionCollection) string {

	var contTable strings.Builder
	contTable.WriteString("<style>table,td { border: 1px solid black; padding: 2px} </style>")

	contTable.WriteString("<table><tr><th>Date</th><th>Project</th><th>Type</th><th>User</th><th>Link</th><th>Description</th></tr>")
	for _, cont := range c {
		contTable.WriteString("<tr>")

		contTable.WriteString("<td>")
		contTable.WriteString(cont.Date.Format("01/02/2006"))
		contTable.WriteString("</td>")
		contTable.WriteString("<td>")
		contTable.WriteString(cont.Project)
		contTable.WriteString("</td>")
		contTable.WriteString("<td>")
		contTable.WriteString(cont.Type)
		contTable.WriteString("</td>")
		contTable.WriteString("<td>")
		contTable.WriteString(cont.User)
		contTable.WriteString("</td>")
		contTable.WriteString("<td>")
		contTable.WriteString(cont.URL)
		contTable.WriteString("</td>")
		contTable.WriteString("<td>")
		contTable.WriteString(cont.Description)
		contTable.WriteString("</td>")

		contTable.WriteString("</tr>")

	}
	contTable.WriteString("</table>")
	return contTable.String()
}

func buildMonthlySubject() string {
	startDate, err := time.Parse("01-02-2006", yesterdayFrom())
	handle(err)

	month := startDate.Month()
	year := startDate.Year()
	subject := fmt.Sprint("Team GitHub Contributions - ", month, year)
	return subject
}

func buildSubject(startDateString string, endDateString string) string {
	startDate, err := time.Parse("01-02-2006", startDateString)
	handle(err)
	endDate, err := time.Parse("01-02-2006", endDateString)
	handle(err)

	sd := startDate.Format("January 2, 2006")
	ed := endDate.Format("January 2, 2006")
	subject := fmt.Sprintf("Team GitHub Contributions - From %v - %v", sd, ed)
	return subject
}
