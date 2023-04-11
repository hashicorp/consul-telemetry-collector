# Adding a new component

There are three steps to adding a new component. Create a proposal, add 
the component to the collector, add the configuration for the compoent.

## Proposal

Each proposal needs to have 3 parts. What is the use case for the new 
proposal. How are users expected to take advantage of it. How will users 
configure the component.

## Adding the component to the collector

All components are included in `pkg/otelcol/components.go`. The component 
should be created in the expected section.

## Adding the configuration
The configuration for the component should be included in 
`confresolver/inmem/provider.go`.