package aws

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	rds "github.com/aws/aws-sdk-go/service/rds"
	"github.com/golang/mock/gomock"

	"github.com/alecrajeev/aws_rds_exporter/mock/aws/sdk"
	"github.com/alecrajeev/aws_rds_exporter/types"
)

// MockDescribeDBInstances mocks describing the RDS Instances
func MockDescribeDBInstances(t *testing.T, mockMatcher *sdk.MockRDSAPI, wantError bool, testInstances ...types.DBInstance) {
	var err error
	if wantError {
		err = errors.New("DescribeDBInstances wrong!")
	}
	rIds := []*(rds.DBInstance){}

	for _, instance := range testInstances {

		b := int64(instance.AllocatedStorage)
		c := int64(instance.Iops)

		rdsInstance := &rds.DBInstance{
			AllocatedStorage:     &b,
			DBInstanceIdentifier: aws.String(instance.Identifier),
			Iops:                 &c,
		}

		rIds = append(rIds, rdsInstance)
	}

	// builds mock output based on the input
	result := &rds.DescribeDBInstancesOutput{
		DBInstances: rIds,
	}
	mockMatcher.EXPECT().DescribeDBInstances(gomock.Any()).Do(func(input interface{}) {
	}).AnyTimes().Return(result, err)

}
