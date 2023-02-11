package main

import (
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
