package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/ecs"
)

func main() {
	// 引数処理
	if !validateArguments(os.Args) {
		log.Fatal("Invalid Argument. Required task-definition name, commit id and container name")
	}
	taskDefinitionName, commitId, containerName := retrieveArg(os.Args)

	// AWSインスタンス初期化
	awsSession, err := initializeAWSSession()
	if err != nil {
		log.Fatal(err)
	}
	svc := ecs.New(awsSession)

	latestTaskDefinition, err := getSpecifiedTaskDefinition(svc, taskDefinitionName)
	if err != nil {
		log.Fatal(err)
	}

	taskDefinition, err := describeTaskDefinition(svc, latestTaskDefinition)
	if err != nil {
		log.Fatal(err)
	}
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
func validateArguments(args []string) bool {
	return len(args) == 4
}

// 引数の取得
func retrieveArg(args []string) (string, string, string) {
	return args[1], args[2], args[3]
}

// AWSセッション初期化
func initializeAWSSession() (*session.Session, error) {
	return session.NewSession(
		&aws.Config{
			Region: aws.String(endpoints.ApNortheast1RegionID),
		})
}

// 指定の名前ののタスク定義を取得
func getSpecifiedTaskDefinition(svc *ecs.ECS, specifiedTaskDefinitionName string) (*ecs.DescribeTaskDefinitionInput, error) {
	var taskDefinitionArns []string
	nextToken := ""
	for {
		result, err := svc.ListTaskDefinitions(
			&ecs.ListTaskDefinitionsInput{
				NextToken: aws.String(nextToken),
			})
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		for _, arn := range result.TaskDefinitionArns {
			taskDefinitionArns = append(taskDefinitionArns, *arn)
		}
		if result.NextToken == nil {
			break
		}
		nextToken = *result.NextToken
	}

	for _, taskDefinitionArn := range taskDefinitionArns {
		if strings.Contains(taskDefinitionArn, specifiedTaskDefinitionName) {
			return &ecs.DescribeTaskDefinitionInput{
				TaskDefinition: aws.String(taskDefinitionArn),
			}, nil
		}
	}
	return nil, fmt.Errorf("TaskDefinition Not Found")
}

func describeTaskDefinition(svc *ecs.ECS, input *ecs.DescribeTaskDefinitionInput) (*ecs.TaskDefinition, error) {
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
		return nil, err
	}

	return result.TaskDefinition, nil
}
