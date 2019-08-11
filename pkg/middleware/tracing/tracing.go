package tracing

import (
	"fmt"
	"os"

	"contrib.go.opencensus.io/exporter/jaeger"
)

// NewExporter : return Jaeger Exporter (OpenCensus)
func NewExporter() (*jaeger.Exporter, error) {
	agentEndpointURI := os.Getenv("OT_AGENT_URI")
	collectorEndpointURI := os.Getenv("OT_COL_URI")

	fmt.Println(agentEndpointURI)
	fmt.Println(collectorEndpointURI)

	je, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint:     agentEndpointURI,
		CollectorEndpoint: collectorEndpointURI,
		ServiceName:       "user-service",
	})
	if err != nil {
		return nil, err
	}

	return je, nil
}
