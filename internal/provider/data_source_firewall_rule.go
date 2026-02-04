package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &FirewallRuleDataSource{}

func NewFirewallRuleDataSource() datasource.DataSource {
	return &FirewallRuleDataSource{}
}

type FirewallRuleDataSource struct {
	client *Client
}

type FirewallRuleDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Description types.String `tfsdk:"description"`
	Interface   types.String `tfsdk:"interface"`
	Protocol    types.String `tfsdk:"protocol"`
	SourceNet   types.String `tfsdk:"source_net"`
	DestNet     types.String `tfsdk:"destination_net"`
	Action      types.String `tfsdk:"action"`
	Enabled     types.Bool   `tfsdk:"enabled"`
}

func (d *FirewallRuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_rule"
}

func (d *FirewallRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches information about an OPNsense firewall rule",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Rule UUID",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the firewall rule",
				Computed:            true,
			},
			"interface": schema.StringAttribute{
				MarkdownDescription: "Interface name",
				Computed:            true,
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "Protocol",
				Computed:            true,
			},
			"source_net": schema.StringAttribute{
				MarkdownDescription: "Source network or IP address",
				Computed:            true,
			},
			"destination_net": schema.StringAttribute{
				MarkdownDescription: "Destination network or IP address",
				Computed:            true,
			},
			"action": schema.StringAttribute{
				MarkdownDescription: "Action to take",
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the rule is enabled",
				Computed:            true,
			},
		},
	}
}

func (d *FirewallRuleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Client, got: %T.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *FirewallRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FirewallRuleDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// In a real implementation, you would fetch the rule data from the API here
	// For now, this is a placeholder

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
