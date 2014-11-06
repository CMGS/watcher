package main

import (
	"flag"

	"./common"
	"./defines"
	"./metrics"
)

var Metrics *metrics.MetricsRecorder

func main() {
	var config defines.MetricsConfig
	var hostname string
	var endpoint string

	flag.IntVar(&config.ReportInterval, "report", 10, "report interval")
	flag.StringVar(&config.Host, "host", "10.1.201.42:8086", "influxdb host")
	flag.StringVar(&config.Username, "username", "root", "user name")
	flag.StringVar(&config.Password, "password", "root", "user password")
	flag.StringVar(&config.Database, "database", "test", "database name")
	flag.StringVar(&endpoint, "endpoint", "tcp://192.168.59.103:2375", "docker endpoint")

	Metrics = metrics.NewMetricsRecorder(hostname, config)
	common.Docker = defines.NewDocker(endpoint)

}
