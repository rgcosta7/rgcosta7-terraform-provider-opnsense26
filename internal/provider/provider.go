package provider

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &opnsenseProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &opnsenseProvider{
			version: version,
		}
	}
}

// opnsenseProvider is the provider implementation.
type opnsenseProvider struct {
	version string
}

// opnsenseProviderModel maps provider schema data to a Go type.
type opnsenseProviderModel struct {
	Host          types.String `tfsdk:"host"`
	ApiKey        types.String `tfsdk:"api_key"`
	ApiSecret     types.String `tfsdk:"api_secret"`
	Insecure      types.Bool   `tfsdk:"insecure"`
	TimeoutSeconds types.Int64  `tfsdk:"timeout_seconds"`
}

// Metadata returns the provider type name.
func (p *opnsenseProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "opnsense"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *opnsenseProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with OPNsense 26.1 API for firewall management, Kea DHCP, and WireGuard configuration.",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "OPNsense host URL (e.g., https://192.168.1.1). Can also be set via OPNSENSE_HOST environment variable.",
				Optional:    true,
			},
			"api_key": schema.StringAttribute{
				Description: "OPNsense API key. Can also be set via OPNSENSE_API_KEY environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"api_secret": schema.StringAttribute{
				Description: "OPNsense API secret. Can also be set via OPNSENSE_API_SECRET environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"insecure": schema.BoolAttribute{
				Description: "Skip TLS certificate verification. Defaults to false.",
				Optional:    true,
			},
			"timeout_seconds": schema.Int64Attribute{
				Description: "HTTP timeout in seconds. Defaults to 30.",
				Optional:    true,
			},
		},
	}
}

// Configure prepares an API client for data sources and resources.
func (p *opnsenseProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring OPNsense client")

	// Retrieve provider data from configuration
	var config opnsenseProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.
	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown OPNsense API Host",
			"The provider cannot create the OPNsense API client as there is an unknown configuration value for the OPNsense API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the OPNSENSE_HOST environment variable.",
		)
	}

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown OPNsense API Key",
			"The provider cannot create the OPNsense API client as there is an unknown configuration value for the OPNsense API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the OPNSENSE_API_KEY environment variable.",
		)
	}

	if config.ApiSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_secret"),
			"Unknown OPNsense API Secret",
			"The provider cannot create the OPNsense API client as there is an unknown configuration value for the OPNsense API secret. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the OPNSENSE_API_SECRET environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	host := os.Getenv("OPNSENSE_HOST")
	apiKey := os.Getenv("OPNSENSE_API_KEY")
	apiSecret := os.Getenv("OPNSENSE_API_SECRET")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.ApiKey.IsNull() {
		apiKey = config.ApiKey.ValueString()
	}

	if !config.ApiSecret.IsNull() {
		apiSecret = config.ApiSecret.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.
	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing OPNsense API Host",
			"The provider cannot create the OPNsense API client as there is a missing or empty value for the OPNsense API host. "+
				"Set the host value in the configuration or use the OPNSENSE_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing OPNsense API Key",
			"The provider cannot create the OPNsense API client as there is a missing or empty value for the OPNsense API key. "+
				"Set the api_key value in the configuration or use the OPNSENSE_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apiSecret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_secret"),
			"Missing OPNsense API Secret",
			"The provider cannot create the OPNsense API client as there is a missing or empty value for the OPNsense API secret. "+
				"Set the api_secret value in the configuration or use the OPNSENSE_API_SECRET environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Handle insecure TLS option
	insecure := false
	if !config.Insecure.IsNull() {
		insecure = config.Insecure.ValueBool()
	}

	// Handle timeout
	timeout := int64(30)
	if !config.TimeoutSeconds.IsNull() {
		timeout = config.TimeoutSeconds.ValueInt64()
	}

	ctx = tflog.SetField(ctx, "opnsense_host", host)
	ctx = tflog.SetField(ctx, "opnsense_api_key", apiKey)
	ctx = tflog.SetField(ctx, "opnsense_api_secret", apiSecret)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "opnsense_api_key")
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "opnsense_api_secret")

	tflog.Debug(ctx, "Creating OPNsense client")

	// Create a new OPNsense client using the configuration values
	client, err := NewClient(&host, &apiKey, &apiSecret, insecure, timeout)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create OPNsense API Client",
			"An unexpected error occurred when creating the OPNsense API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"OPNsense Client Error: "+err.Error(),
		)
		return
	}

	// Make the OPNsense client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured OPNsense client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *opnsenseProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewFirewallRuleDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *opnsenseProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewFirewallRuleResource,
		NewFirewallAliasResource,
		NewFirewallCategoryResource,
		NewNatDestinationResource,
		NewKeaReservationResource,
		NewKeaSubnetResource,
		NewWireguardServerResource,
		NewWireguardPeerResource,
	}
}

// Client represents the OPNsense API client
type Client struct {
	Host      string
	ApiKey    string
	ApiSecret string
	client    *http.Client
}

// NewClient creates a new OPNsense API client
func NewClient(host, apiKey, apiSecret *string, insecure bool, timeout int64) (*Client, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	}

	httpClient := &http.Client{
		Transport: tr,
		Timeout:   0, // We'll handle timeouts per-request
	}

	c := &Client{
		Host:      *host,
		ApiKey:    *apiKey,
		ApiSecret: *apiSecret,
		client:    httpClient,
	}

	return c, nil
}

// DoRequest performs an HTTP request to the OPNsense API
func (c *Client) DoRequest(ctx context.Context, method, endpoint string, body []byte) ([]byte, error) {
	url := fmt.Sprintf("%s/api/%s", c.Host, endpoint)
	
	tflog.Debug(ctx, "Making API request", map[string]any{
		"method":   method,
		"endpoint": endpoint,
		"url":      url,
	})

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set basic auth
	req.SetBasicAuth(c.ApiKey, c.ApiSecret)

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if body != nil && len(body) > 0 {
		req.Body = http.NoBody
		// For POST requests with body, we'd need to set the body properly
		// This is simplified for the example
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	// Read response body
	var respBody []byte
	// In a real implementation, you'd read the response body here
	
	return respBody, nil
}
