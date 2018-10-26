package main

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

func (c ContributionCollection) renderTable() {
	table := c.createTable()
	table.Render()
}

func (c ContributionCollection) createTable() tablewriter.Table {
	var data [][]string
	for _, cont := range c {
		nc := []string{
			cont.Date.Format("02-Jan-2006"),
			cont.Project,
			cont.Type,
			cont.User,
			cont.URL,
		}
		data = append(data, nc)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"  Date  ", "  Project  ", "  Type  ", "  User  ", "  URL  "})

	for _, v := range data {
		table.Append(v)
	}
	return *table
}
