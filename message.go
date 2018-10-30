package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/now"

	"gopkg.in/gomail.v2"
)

func newMessage(c ContributionCollection) error {
	body := buildEmail(c)
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SESVerifiedEmail"))
	m.SetHeader("To", os.Getenv("SESVerifiedEmail"))
	m.SetHeader("Subject", buildSubject())
	m.SetBody("text/html", body)
	d := gomail.NewDialer(os.Getenv("SESServerName"), 465, os.Getenv("SESUserName"), os.Getenv("SESPassword"))

	err := d.DialAndSend(m)
	return err
}

func newMessageTo(c ContributionCollection, emailAddress string) error {
	body := buildEmail(c)
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SESVerifiedEmail"))
	m.SetHeader("To", emailAddress)
	m.SetHeader("Subject", buildSubject())
	m.SetBody("text/html", body)
	d := gomail.NewDialer(os.Getenv("SESServerName"), 465, os.Getenv("SESUserName"), os.Getenv("SESPassword"))

	err := d.DialAndSend(m)
	return err
}

func buildEmail(c ContributionCollection) string {
	var body strings.Builder
	body.WriteString("<br/>This month the team had ")
	body.WriteString(strconv.Itoa(len(c)))
	body.WriteString(" contributions into open source projects.<br/>Below is a table of all contributions for the month<br/><br/>")

	body.WriteString(createTable(c))

	return body.String()
}

func createTable(c ContributionCollection) string {

	var contTable strings.Builder
	contTable.WriteString("<style>table,td { border: 1px solid black; padding: 2px} </style>")

	contTable.WriteString("<table><tr><th>Date</th><th>Project</th><th>Type</th><th>User</th><th>Link</th></tr>")
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

		contTable.WriteString("</tr>")

	}
	contTable.WriteString("</table>")
	return contTable.String()
}

func buildSubject() string {
	var startDate time.Time
	if os.Getenv("since") != "" {
		startDate, _ = time.Parse("01/02/06", os.Getenv("since"))
	} else {
		yesterday := time.Now().AddDate(0, 0, -1)

		startDate = now.New(yesterday).BeginningOfMonth()
	}
	month := startDate.Month()
	year := startDate.Year()
	subject := fmt.Sprint("Amazon EKS OSS Contributions - ", month, year)
	return subject
}
