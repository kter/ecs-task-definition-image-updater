package main

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
)

func TestValidArguments(t *testing.T) {
	args := []string{"a", "b", "c", "d"}
	if !validateArguments(args) {
		t.Error("args did not validate when they should have")
	}
}

func TestInvalidArguments(t *testing.T) {
	args := []string{"a", "b", "c"}
	if validateArguments(args) {
		t.Error("args validated when they should not have")
	}
}

func TestInitializeAWSSession(t *testing.T) {
	initializeAWSSessionResult, err := initializeAWSSession()

	if err != nil {
		t.Errorf("initializeAWSSession returned unexpected error %s", err.Error())
	}

	// Check that the result is not nil
	if initializeAWSSessionResult == nil {
		t.Errorf("initializeAWSSession did not return a session")
	}

	// Check that the region is correct
	if initializeAWSSessionResult != nil && aws.StringValue(initializeAWSSessionResult.Config.Region) != endpoints.ApNortheast1RegionID {
		t.Errorf("expected region to be %s, but was %s", endpoints.ApNortheast1RegionID, aws.StringValue(initializeAWSSessionResult.Config.Region))
	}
}

func TestGetTaskDefinitionInput(t *testing.T) {
	arns := []string{
		"arn:aws:ecs:ap-northeast-1:848738341109:task-definition/test2:1",
		"arn:aws:ecs:ap-northeast-1:848738341109:task-definition/test2:2",
		"arn:aws:ecs:ap-northeast-1:848738341109:task-definition/test:1",
		"arn:aws:ecs:ap-northeast-1:848738341109:task-definition/test:2",
	}
	searchTaskDefinition := "test2"

	result, err := getTaskDefinitionInput(arns, searchTaskDefinition)
	if err != nil {
		t.Errorf("getTaskDefinitionInput returned unexpected error %s", err.Error())
	}

	// Check that the searchTaskDefinition is returned
	if result != nil && *result.TaskDefinition != "arn:aws:ecs:ap-northeast-1:848738341109:task-definition/test2:2" {
		fmt.Println(aws.String(*result.TaskDefinition))
		t.Errorf("Expected TaskDefinition to be %s but got %s", searchTaskDefinition, *result.TaskDefinition)
	}
}
