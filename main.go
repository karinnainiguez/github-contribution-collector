package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	Message string `json: "message"`
	Ok      bool   `json: "ok"`
}

func Handler() (Response, error) {
	nc := collectContributions()
	nc.sendMonthlyEmail()
	return Response{
		Message: fmt.Sprint("Monthly Email Sent Successfully"),
		Ok:      true,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
