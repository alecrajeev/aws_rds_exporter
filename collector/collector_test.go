package collector

import (
	"math"
	"testing"

	awsMock "github.com/alecrajeev/aws_rds_exporter/mock/aws"
	"github.com/alecrajeev/aws_rds_exporter/mock/aws/sdk"
	"github.com/alecrajeev/aws_rds_exporter/types"
	"github.com/golang/mock/gomock"
)

type TestDBInstance struct {
	instance          types.DBInstance
	wantErrorList     bool
	wantErrorDescribe bool
	expectError       bool
}

func TestGetRDSInstances(t *testing.T) {

	a1 := types.DBInstance{Identifier: "rds-dbinstance-1", AllocatedStorage: 55.0, Iops: 0.0}

	testInstance := &TestDBInstance{
		instance:          a1,
		wantErrorList:     false,
		wantErrorDescribe: false,
		expectError:       false,
	}

	// Mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRDS := sdk.NewMockRDSAPI(ctrl)
	awsMock.MockDescribeDBInstances(t, mockRDS, testInstance.wantErrorList, testInstance.instance)

	e := &RDSClient{
		client: mockRDS,
	}

	rdsInstances, err := e.GetRDSInstances()
	if !testInstance.expectError {
		if err != nil {
			t.Errorf("\n- %v\n-  Shouldn't return an error, but it did: %v", testInstance, err)
		}

		if len(rdsInstances) != 1 {
			t.Errorf("\n- %v\n-  Length in returned number of RDS Instances differs than expected, want: %d; got: %d", testInstance, 1, len(rdsInstances))
		}

		for _, got := range rdsInstances {
			wantInstanceIdentifier := testInstance.instance.Identifier
			wantAllocatedStorage := float64(testInstance.instance.AllocatedStorage) * math.Pow(10, 9)
			wantIops := float64(testInstance.instance.Iops)
			if wantInstanceIdentifier != got.Identifier {
				t.Errorf("\n- %v\n- Wanted an InstanceIdentifer of %v, got %v", testInstance, wantInstanceIdentifier, got.Identifier)
			}
			if wantAllocatedStorage != got.AllocatedStorage {
				t.Errorf("\n- %v\n- Wanted an AllocatedStorage of %v, got %v", testInstance, wantAllocatedStorage, got.AllocatedStorage)
			}
			if wantIops != float64(got.Iops) {
				t.Errorf("\n- %v\n- Wanted an Iops of %v, got %v", testInstance, wantIops, got.Iops)
			}
		}

	} else {
		if err == nil {
			t.Errorf("\n- %v\n-  Should return an error, it didn't", testInstance)
		}
	}
}
