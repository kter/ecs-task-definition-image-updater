package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/service/ecs"
)

func main() {
	// 引数処理
	if !validateArg(os.Args) {
		log.Fatal("Invalid Argument. Required taskDefinitionName and commit id")
	}
	taskDefinitionName, commitId, containerName := retrieveArg(os.Args)

	// AWSインスタンス初期化
	awsSession, err := initAwsSession()
	if err != nil {
		log.Fatal(err)
	}
	svc := ecs.New(awsSession)

	var taskDefinitions []string
	nextToken := ""
	for {
		result, err := svc.ListTaskDefinitions(
			&ecs.ListTaskDefinitionsInput{
				NextToken: aws.String(nextToken),
			})
		if err != nil {
			log.Fatal(err.Error())
		}
		for _, arn := range result.TaskDefinitionArns {
			taskDefinitions = append(taskDefinitions, *arn)
		}
		if result.NextToken == nil {
			break
		}
		nextToken = *result.NextToken
	}
	input := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(taskDefinitions[len(taskDefinitions)-1]),
	}

	result, err := svc.DescribeTaskDefinition(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecs.ErrCodeServerException:
				fmt.Println(ecs.ErrCodeServerException, aerr.Error())
			case ecs.ErrCodeClientException:
				fmt.Println(ecs.ErrCodeClientException, aerr.Error())
			case ecs.ErrCodeInvalidParameterException:
				fmt.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	taskDefinition := result.TaskDefinition
	for _, container := range taskDefinition.ContainerDefinitions {
		if *container.Name == containerName {
			container.Image = aws.String(taskDefinitionName + ":" + commitId)
		}
	}
	registerTaskInput := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: taskDefinition.ContainerDefinitions,
		Family:               taskDefinition.Family,
	}
	registerResult, err := svc.RegisterTaskDefinition(registerTaskInput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(*registerResult.TaskDefinition.TaskDefinitionArn)
}

// 引数のバリデート
func validateArg(args []string) bool {
	if len(args) != 4 {
		return false
	}
	return true
}

func retrieveArg(args []string) (string, string, string) {
	return args[1], args[2], args[3]
}

func initAwsSession() (*session.Session, error) {
	return session.NewSession(
		&aws.Config{
			Region: aws.String(endpoints.ApNortheast1RegionID),
		})
}
