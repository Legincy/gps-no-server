package influx

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"gps-no-server/internal/common/config"
)

type InfluxConnection struct {
	client   influxdb2.Client
	writeAPI api.WriteAPI
	config   *config.InfluxConfig
}

func NewInfluxConnection(cfg *config.InfluxConfig) (*InfluxConnection, error) {
	client := influxdb2.NewClient(cfg.URL, cfg.Token)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	health, err := client.Health(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to InfluxDB: %w", err)
	}

	if health.Status != "pass" {
		return nil, fmt.Errorf("InfluxDB health check failed: %s", health.Status)
	}

	writeAPI := client.WriteAPI(cfg.Org, cfg.Bucket)

	return &InfluxConnection{
		client:   client,
		writeAPI: writeAPI,
		config:   cfg,
	}, nil
}

func (i *InfluxConnection) Connect() error {
	return nil
}

func (i *InfluxConnection) Disconnect() error {
	i.writeAPI.Flush()
	i.client.Close()
	return nil
}

func (i *InfluxConnection) WritePoint(measurement string, tags map[string]string, fields map[string]interface{}) error {
	point := influxdb2.NewPoint(measurement, tags, fields, time.Now())
	i.writeAPI.WritePoint(point)
	return nil
}

func (i *InfluxConnection) IsConnected() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	health, err := i.client.Health(ctx)
	return err == nil && health.Status == "pass"
}
