package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/fsouza/go-dockerclient"

	"./common"
	"./defines"
	"./logs"
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

	containers, err := common.Docker.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		logs.Assert(err, "Load")
	}

	logs.Info("Load container")
	for _, container := range containers {
		sid := container.ID[:12]
		if strings.HasPrefix(container.Status, "Exit") {
			logs.Info(container.Names, sid, "die")
			continue
		}
		name := strings.Split(container.Names[0], "_")[0]
		Metrics.Add(name, sid, common.DEFAULT_TYPE)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	signal.Notify(sc, syscall.SIGTERM)
	signal.Notify(sc, syscall.SIGHUP)
	signal.Notify(sc, syscall.SIGKILL)
	signal.Notify(sc, syscall.SIGQUIT)
	logs.Info("Got <-", <-sc)
}
