package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	namespace = "rds"
)

// Metrics descriptions
var (

	// labels are the static labels that come with every metric
	labels = []string{"region", "instance"}

	storage = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "storage"),
		"Amount of storage for the RDS instance",
		labels,
		nil,
	)
)

// RDSService represents a service on an RDS instance
type DBInstance struct {
	ID         string // Instance Identifier
	AS         int // allocated storage
}

type promHTTPLogger struct {
	logger log.Logger
}

func (l promHTTPLogger) Println(v ...interface{}) {
	level.Error(l.logger).Log("msg", fmt.Sprint(v...))
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

	for _, pa := range resp.DBInstances {

		var b = int(*(pa.AllocatedStorage))
		db := &DBInstance{
			ID: aws.StringValue(pa.DBInstanceIdentifier),
			AS: b,
		}

		rs = append(rs, db)
	}

	// log.Info("Got %d clusters", len(rs))
	return rs, nil
}

type exporter struct {
	client    RDSGatherer
	region string
}

// Describe describes the metrics exported by the RDS exporter. It
// implements prometheus.Collector.
func (e *exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- storage
}

// Collect fetches the stats from the configured RDS and delivers them
// as Prometheus metrics. It implements prometheus.Collector
func (e *exporter) Collect(ch chan<- prometheus.Metric) {


	rs, err := e.client.GetRDSInstances()

	if err != nil {
		// log.Error("Error collecting rds metrics")
		return
	}

	for _, r := range rs {
		ch <- prometheus.MustNewConstMetric(
			storage, prometheus.GaugeValue, float64(r.AS), e.region, r.ID,
		)
	}
}

func init() {
	prometheus.MustRegister(version.NewCollector("aws_rds_exporter"))
}

func run() int {

	var (
		listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9785").String()
		metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	)


	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	fmt.Printf("Starting aws-rds-exporter...")
	fmt.Printf("\n")

	RdsClient, err := NewRDSClient("us-east-1")

	if (err != nil) {
		fmt.Printf("Error with rds client")
		return 1
	}

	// exporter := &exporter{api: rds.New(sess)}
	exporter := &exporter{
		client: RdsClient,
		region: "us-east-1",
	}
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath,
		promhttp.InstrumentMetricHandler(
			prometheus.DefaultRegisterer,
			promhttp.HandlerFor(
				prometheus.DefaultGatherer,
				promhttp.HandlerOpts{
					ErrorLog: &promHTTPLogger{
						logger: logger,
					},
				},
			),
		),
	)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>AWS RDS Exporter</title></head>
             <body>
             <h1>AWS RDS Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             <h2>Options</h2>
             <h2>Build</h2>
             </body>
             </html>`))
	})
	http.HandleFunc("/-/healthy", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})
	http.HandleFunc("/-/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	fmt.Printf("\n")

	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
		os.Exit(1)
	}

	return 0
}


func main() {

	exCode := run()
	os.Exit(exCode)
}




