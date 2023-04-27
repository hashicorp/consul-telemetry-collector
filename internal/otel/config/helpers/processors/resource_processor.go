package processors

import (
	"errors"

	"go.opentelemetry.io/collector/component"
)

const resourceProcessorName = "resource"

// FilterProcessorID is the component id of the filter processor
var ResourceProcessorID component.ID = component.NewID(resourceProcessorName)

type ResourceProcessorConfig struct {
	Attributes []Actions `mapstructure:"attributes"`
}

type Action string

const (
	insert Action = "insert" // insert
	update        = "update" // update
	upsert        = "upsert" // upsert
	delete        = "delete" // delete

	clusterKey = "cluster"
)

type Actions struct {
	// Key specifies the attribute to act upon.
	// This is a required field.
	Key string `mapstructure:"key"`

	// Value specifies the value to populate for the key.
	// The type of the value is inferred from the configuration.
	Value interface{} `mapstructure:"value"`

	// Action specifies the type of action to perform.
	// The set of values are {INSERT, UPDATE, UPSERT, DELETE, HASH}.
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
	Action Action `mapstructure:"action"`
}

var ErrUnknownAction = errors.New("unknown action")

// FilterProcessorCfg generates the config for a filter processor
func ResourcesProcessorCfg(resourceID string) ResourceProcessorConfig {
	return ResourceProcessorConfig{Attributes: []Actions{
		{
			Key:    clusterKey,
			Value:  resourceID,
			Action: upsert,
		},
	}}
}
