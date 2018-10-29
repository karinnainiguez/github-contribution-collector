package main

import "github.com/spf13/cobra"

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
			cmd.Println("Printing from Report.  This is your report")
			cmd.Println("The Flag is: ")
			cmd.Println(cc.From)
			// nc := collectContributions()
			// nc.createTable()
			return nil
		},
	}

	fs := cmd.Flags()

	fs.StringVarP(&cc.From, "from", "f", defaultFrom(), `Date from which to begin reporting (Default is beginning of current month)`)
	fs.StringVarP(&cc.From, "until", "u", defaultUntil(), `Date until which to run reporting (Default is today)`)
	fs.StringVarP(&cc.From, "email", "e", "", "Email address to send report.")
	return cmd
}

func init() {

}
