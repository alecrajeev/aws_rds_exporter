package collector

import (
	"fmt"
	"math"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/version"
)

const (
	namespace = "aws_rds"
)

// Metrics descriptions
var (

	// labels are the static labels that come with every metric
	labels = []string{"region", "instance"}

	storage = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "storage"),
		"Amount of storage in bytes for the RDS instance",
		labels,
		nil,
	)

	iops = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "iops"),
		"Amount of IOPS (I/O operations per second) value) for the RDS instance",
		labels,
		nil,
	)
)

// RDSService represents a service on an RDS instance
type DBInstance struct {
	Identifier               string  // Instance Identifier
	AllocatedStorage         float64 // allocated storage
	Iops                     float64 // iops
}

// RDSClient is a wrapper for AWS rds client that implements helpers to get RDS metrics
type RDSClient struct {
	client        rdsiface.RDSAPI
	apiMaxResults int64
}

// RDSGatherer is the interface that implements the methods required to gather RDS data
type RDSGatherer interface {
	GetRDSInstances() ([]*DBInstance, error)
}

// NewRDSClient will return an initialized RDSClient
func NewRDSClient(awsRegion string) (*RDSClient, error) {
	// Create AWS session
	s := session.New(&aws.Config{Region: aws.String(awsRegion)})
	if s == nil {
		return nil, fmt.Errorf("error creating aws session")
	}

	return &RDSClient{
		client:        rds.New(s),
		apiMaxResults: 100,
	}, nil
}

// GetRDSInstances will get the instances from the RDS API
func (e *RDSClient) GetRDSInstances() ([]*DBInstance, error) {
	rs := []*DBInstance{}
	params := &rds.DescribeDBInstancesInput{
	}

	resp, err := e.client.DescribeDBInstances(params)
	if err != nil {
		return nil, err
	}

	for _, rdsInstance := range resp.DBInstances {

		var c = 0.0
		if (rdsInstance.Iops) != nil {
			c = float64(*(rdsInstance.Iops))
		}

		var b = (float64(*(rdsInstance.AllocatedStorage)))*math.Pow(10,9)
		db := &DBInstance{
			Identifier: aws.StringValue(rdsInstance.DBInstanceIdentifier),
			AllocatedStorage: b,
			Iops: c,
		}

		rs = append(rs, db)
	}

	return rs, nil
}

func NewExporter(awsRegion string) (*exporter, error) {

	RdsClient, err := NewRDSClient(awsRegion)

	if err != nil {
		fmt.Printf("Error with rds client")
		return nil, fmt.Errorf("Error with rds client")
	}

	return &exporter{
		client: RdsClient,
		region: awsRegion,
	},nil
}

type exporter struct {
	client    RDSGatherer
	region string
}

// Describe describes the metrics exported by the RDS exporter. It
// implements prometheus.Collector.
func (e *exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- storage
	ch <- iops
}

// Collect fetches the stats from the configured RDS and delivers them
// as Prometheus metrics. It implements prometheus.Collector
func (e *exporter) Collect(ch chan<- prometheus.Metric) {


	rs, err := e.client.GetRDSInstances()

	if err != nil {
		return
	}

	for _, r := range rs {
		ch <- prometheus.MustNewConstMetric(
			storage, prometheus.GaugeValue, r.AllocatedStorage, e.region, r.Identifier,
		)
		ch <- prometheus.MustNewConstMetric(
			iops, prometheus.GaugeValue, r.Iops, e.region, r.Identifier,
		)
	}
}

func init() {
	prometheus.MustRegister(version.NewCollector("aws_rds_exporter"))
}
