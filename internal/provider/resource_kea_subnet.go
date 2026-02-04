package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &KeaSubnetResource{}
var _ resource.ResourceWithImportState = &KeaSubnetResource{}

func NewKeaSubnetResource() resource.Resource {
	return &KeaSubnetResource{}
}

type KeaSubnetResource struct {
	client *Client
}

type KeaSubnetResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Subnet      types.String `tfsdk:"subnet"`
	Pools       types.String `tfsdk:"pools"`
	Option      types.String `tfsdk:"option_data"`
	Description types.String `tfsdk:"description"`
}

func (r *KeaSubnetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kea_subnet"
}

func (r *KeaSubnetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Kea DHCP subnets in OPNsense",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Subnet UUID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"subnet": schema.StringAttribute{
				MarkdownDescription: "Subnet in CIDR notation (e.g., 192.168.1.0/24)",
				Required:            true,
			},
			"pools": schema.StringAttribute{
				MarkdownDescription: "IP address pools (comma-separated ranges, e.g., '192.168.1.100-192.168.1.200')",
				Optional:            true,
			},
			"option_data": schema.StringAttribute{
				MarkdownDescription: "DHCP options data",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the subnet",
				Optional:            true,
			},
		},
	}
}

func (r *KeaSubnetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Client, got: %T.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *KeaSubnetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data KeaSubnetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subnetData := map[string]interface{}{
		"subnet": map[string]interface{}{
			"subnet": data.Subnet.ValueString(),
		},
	}

	if !data.Pools.IsNull() {
		subnetData["subnet"].(map[string]interface{})["pools"] = data.Pools.ValueString()
	}
	if !data.Option.IsNull() {
		subnetData["subnet"].(map[string]interface{})["option_data"] = data.Option.ValueString()
	}
	if !data.Description.IsNull() {
		subnetData["subnet"].(map[string]interface{})["description"] = data.Description.ValueString()
	}

	jsonData, _ := json.Marshal(subnetData)

	url := fmt.Sprintf("%s/api/kea/dhcpv4/addSubnet", r.client.Host)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonData)))
	httpReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create subnet: %s", err))
		return
	}
	defer httpResp.Body.Close()

	body, _ := io.ReadAll(httpResp.Body)

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to parse response: %s", err))
		return
	}

	if uuid, ok := result["uuid"].(string); ok {
		data.ID = types.StringValue(uuid)
	} else {
		resp.Diagnostics.AddError("API Error", "No UUID returned from API")
		return
	}

	// Apply configuration
	applyURL := fmt.Sprintf("%s/api/kea/service/reconfigure", r.client.Host)
	applyReq, _ := http.NewRequestWithContext(ctx, "POST", applyURL, nil)
	applyReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)
	r.client.client.Do(applyReq)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *KeaSubnetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data KeaSubnetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/api/kea/dhcpv4/getSubnet/%s", r.client.Host, data.ID.ValueString())
	httpReq, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	httpReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read subnet: %s", err))
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *KeaSubnetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data KeaSubnetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subnetData := map[string]interface{}{
		"subnet": map[string]interface{}{
			"subnet": data.Subnet.ValueString(),
		},
	}

	if !data.Pools.IsNull() {
		subnetData["subnet"].(map[string]interface{})["pools"] = data.Pools.ValueString()
	}
	if !data.Option.IsNull() {
		subnetData["subnet"].(map[string]interface{})["option_data"] = data.Option.ValueString()
	}
	if !data.Description.IsNull() {
		subnetData["subnet"].(map[string]interface{})["description"] = data.Description.ValueString()
	}

	jsonData, _ := json.Marshal(subnetData)

	url := fmt.Sprintf("%s/api/kea/dhcpv4/setSubnet/%s", r.client.Host, data.ID.ValueString())
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonData)))
	httpReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update subnet: %s", err))
		return
	}
	defer httpResp.Body.Close()

	// Apply configuration
	applyURL := fmt.Sprintf("%s/api/kea/service/reconfigure", r.client.Host)
	applyReq, _ := http.NewRequestWithContext(ctx, "POST", applyURL, nil)
	applyReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)
	r.client.client.Do(applyReq)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *KeaSubnetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data KeaSubnetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/api/kea/dhcpv4/delSubnet/%s", r.client.Host, data.ID.ValueString())
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, nil)
	httpReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete subnet: %s", err))
		return
	}
	defer httpResp.Body.Close()

	// Apply configuration
	applyURL := fmt.Sprintf("%s/api/kea/service/reconfigure", r.client.Host)
	applyReq, _ := http.NewRequestWithContext(ctx, "POST", applyURL, nil)
	applyReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)
	r.client.client.Do(applyReq)
}

func (r *KeaSubnetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
