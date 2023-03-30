package agent

const (
	// HCPClientID is the environment variable for the hcp client id
	HCPClientID = "HCP_CLIENT_ID"

	// HCPClientID is the cli opt for the hcp client id
	HCPClientIDOpt = "hcp-client-id"

	// HCPClientSecret is the environment variable for the hcp client secret
	HCPClientSecret = "HCP_CLIENT_SECRET"

	// HCPClientSecretOpt is the cli opt for the hcp client secret
	HCPClientSecretOpt = "hcp-client-secret"

	// HCPResourceID is the environment variable for the hcp resource id
	HCPResourceID = "HCP_RESOURCE_ID"
	// HCPResourceIDOpt is the cli opt for the hcp resource id
	HCPResourceIDOpt = "hcp-resource-id"

	// COOtelHTTPEndpoint is the environment variable for  OpenTelemetry HTTP Endpoints where metrics are forwarded
	COOtelHTTPEndpoint = "CO_OTEL_HTTP_ENDPOINT"

	// COOtelHTTPEndpointOpt is the cli opt for the OpenTelemetry HTTP Endpoints where metrics are forwarded
	COOtelHTTPEndpointOpt = "http-collector-endpoint"

	// COOConfigPath is the environment variable for path to the config
	COOConfigPath = "COO_CONFIG_PATH"

	// COOConfigPathOpt is the cli opt for path to the config
	COOConfigPathOpt = "config-file-path"
)
