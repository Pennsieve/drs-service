package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	//TODO update import
	"github.com/pennsieve/drs-service/service/handler"
)

func main() {
	lambda.Start(handler.DrsServiceHandler)
}
