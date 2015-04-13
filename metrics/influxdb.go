package metrics

import (
	"net/http"

	"../defines"
	"../logs"
	"github.com/influxdb/influxdb/client"
)

type InfluxDBClient struct {
	hostname string
	client   *client.Client
	series   []*client.Series
}

var influxdb_columns []string = []string{"host", "apptype", "appid", "metric", "value"}

func NewInfluxDBClient(hostname string, config defines.MetricsConfig) *InfluxDBClient {
	c := &client.ClientConfig{
		Host:       config.Host,
		Username:   config.Username,
		Password:   config.Password,
		Database:   config.Database,
		HttpClient: http.DefaultClient,
		IsSecure:   false,
		IsUDP:      false,
	}
	i, err := client.New(c)
	if err != nil {
		logs.Assert(err, "InfluxDB")
	}
	return &InfluxDBClient{hostname, i, []*client.Series{}}
}

func (self *InfluxDBClient) GenSeries(cid string, app *MetricData) {
	points := [][]interface{}{
		{self.hostname, app.apptype, cid, "cpu_usage", app.cpu_usage},
		{self.hostname, app.apptype, cid, "cpu_system", app.cpu_system},
		{self.hostname, app.apptype, cid, "cpu_user", app.cpu_user},
		{self.hostname, app.apptype, cid, "mem_usage", app.mem_usage},
		{self.hostname, app.apptype, cid, "mem_rss", app.mem_rss},
		{self.hostname, app.apptype, cid, "mem_max_usage", app.mem_max_usage},
		{self.hostname, app.apptype, cid, "cpu_usage_rate", app.cpu_user_rate},
		{self.hostname, app.apptype, cid, "cpu_system_rate", app.cpu_system_rate},
		{self.hostname, app.apptype, cid, "cpu_user_rate", app.cpu_user_rate},
	}
	if app.isapp {
		for key, data := range app.network {
			p2 := []interface{}{self.hostname, app.apptype, cid, key, data}
			points = append(points, p2)
		}
		for key, data := range app.network_rate {
			p2 := []interface{}{self.hostname, app.apptype, cid, key, data}
			points = append(points, p2)
		}
	}
	series := &client.Series{
		Name:    app.appname,
		Columns: influxdb_columns,
		Points:  points,
	}
	self.series = append(self.series, series)
}

func (self *InfluxDBClient) Send() {
	if err := self.client.WriteSeries(self.series); err != nil {
		logs.Info("Write to InfluxDB Failed", err)
	}
	self.series = []*client.Series{}
}
