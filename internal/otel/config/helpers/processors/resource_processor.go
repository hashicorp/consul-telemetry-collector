package processors

import (
	"go.opentelemetry.io/collector/component"

	"github.com/hashicorp/consul-telemetry-collector/internal/hcp"
	"github.com/hashicorp/go-hclog"
)

const resourceProcessorName = "resource"

// ResourceProcessorID is the component id of the resource processor.
var ResourceProcessorID component.ID = component.NewID(resourceProcessorName)

// ResourceProcessorConfig configures the Resource Processor.
type ResourceProcessorConfig struct {
	Attributes []Actions `mapstructure:"attributes"`
}

type action string

const upsertAction action = "upsert"

// Actions specifys the key, value and action that should be acted upon.
type Actions struct {
	// Key specifies the attribute to act upon.
	// This is a required field.
	Key string `mapstructure:"key"`

	// Value specifies the value to populate for the key.
	// The type of the value is inferred from the configuration.
	Value interface{} `mapstructure:"value"`

	// Action specifies the type of action to perform.
	// The set of values are {INSERT, UPDATE, UPSERT, DELETE}.
	// Both lower case and upper case are supported.
	// INSERT -  Inserts the key/value to attributes when the key does not exist.
	//           No action is applied to attributes where the key already exists.
	//           Either Value, FromAttribute or FromContext must be set.
	// UPDATE -  Updates an existing key with a value. No action is applied
	//           to attributes where the key does not exist.
	//           Either Value, FromAttribute or FromContext must be set.
	// UPSERT -  Performs insert or update action depending on the attributes
	//           containing the key. The key/value is inserted to attributes
	//           that did not originally have the key. The key/value is updated
	//           for attributes where the key already existed.
	//           Either Value, FromAttribute or FromContext must be set.
	// DELETE  - Deletes the attribute. If the key doesn't exist,
	//           no action is performed.
	Action action `mapstructure:"action"`
}

// ResourcesProcessorCfg generates the config for a resource processor.
// The cluster's ResourceID is upserted as a label in all metrics.
func ResourcesProcessorCfg(client hcp.TelemetryClient) ResourceProcessorConfig {
	logger := hclog.Default().Named("config/helpers")

	metricAttributes, err := client.MetricAttributes()
	if err != nil {
		logger.Error("failed to retrieve metric labels from HCP", "error", err)
		return ResourceProcessorConfig{}
	}

	actions := make([]Actions, 0, len(metricAttributes))
	for key, value := range metricAttributes {
		actions = append(actions, Actions{
			Key:    key,
			Value:  value,
			Action: upsertAction,
		})
	}

	return ResourceProcessorConfig{
		Attributes: actions,
	}
}
