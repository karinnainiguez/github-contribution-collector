package main

import (
	"github.com/spf13/cobra"
)

type CommandConfig struct {
	From  string
	Until string
	Email string
}

func reportContributions() *cobra.Command {
	cc := &CommandConfig{}
	cmd := &cobra.Command{
		Use:       "report",
		ValidArgs: []string{"from", "until", "email"},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := doReportContributions(cc)
			handle(err)
			return nil
		},
	}

	fs := cmd.Flags()

	fs.StringVarP(&cc.From, "from", "f", defaultFrom(), `Date from which to begin reporting (Default is beginning of current month)`)
	fs.StringVarP(&cc.Until, "until", "u", defaultUntil(), `Date until which to run reporting (Default is today)`)
	fs.StringVarP(&cc.Email, "email", "e", "", "Email address to send report.")
	return cmd
}

func doReportContributions(c *CommandConfig) error {
	contributions := collectContributions()
	filtered := contributions.filterContributions(c.From, c.Until)

	filtered.renderTable()
	// send email if specified
	if c.Email != "" {
		// TODO: Verify email is sent, gracefully exit otherwise.
		err := newMessageTo(filtered, c.Email, c.From, c.Until)
		handle(err)
	}

	return nil
}
