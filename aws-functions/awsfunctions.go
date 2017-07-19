package awsfunctions

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

//Function to terminate an instance
func terminateInstance(instanceID string, s session.Session, c aws.Config) bool {
	svc := ec2.New(s, c)
	resp, err := svc.TerminateInstances(&ec2.TerminateInstancesInput{InstanceIds: []*string{aws.String(instanceID)}})
	if err != nil {
		fmt.Println(resp)
		panic(err)
	}
	return true
}
