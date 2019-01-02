package main

import (
	"os"

	"github.com/spf13/cobra"
)

/*

// Response struct used for Lambda Response.
type Response struct {
	Message string `json: "message"`
	Ok      bool   `json: "ok"`
}

// Handler Function used as lamda.Start function in main.
func Handler() (Response, error) {
	nc, err := collectContributionsConcurrently()
	filtered := nc.filterMonthlyContributions()
	newMessage(filtered)
	return Response{
		Message: fmt.Sprint("Monthly Email Sent Successfully"),
		Ok:      true,
	}, err
}

// Main Function used for Lambda
func main() {
	lambda.Start(Handler)
}
*/

// Main Function used for CLI
func main() {
	cmd := &cobra.Command{
		Use:   "eks-oss-contributions",
		Short: "\nReport team contributions into open source software using as little as one command.",
	}

	cmd.AddCommand(reportContributions())

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}

}
