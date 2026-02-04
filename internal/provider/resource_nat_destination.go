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

var _ resource.Resource = &NatDestinationResource{}
var _ resource.ResourceWithImportState = &NatDestinationResource{}

func NewNatDestinationResource() resource.Resource {
	return &NatDestinationResource{}
}

type NatDestinationResource struct {
	client *Client
}

type NatDestinationResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	Interface       types.String `tfsdk:"interface"`
	Protocol        types.String `tfsdk:"protocol"`
	SourceNet       types.String `tfsdk:"source_net"`
	SourcePort      types.String `tfsdk:"source_port"`
	DestinationNet  types.String `tfsdk:"destination_net"`
	DestinationPort types.String `tfsdk:"destination_port"`
	TargetIP        types.String `tfsdk:"target_ip"`
	TargetPort      types.String `tfsdk:"target_port"`
	Description     types.String `tfsdk:"description"`
	Log             types.Bool   `tfsdk:"log"`
}

func (r *NatDestinationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nat_destination"
}

func (r *NatDestinationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Destination NAT (Port Forward) rules in OPNsense 26.1",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "NAT rule UUID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable this NAT rule",
				Optional:            true,
			},
			"interface": schema.StringAttribute{
				MarkdownDescription: "Interface (e.g., 'wan')",
				Required:            true,
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "Protocol (tcp, udp, tcp/udp)",
				Required:            true,
			},
			"source_net": schema.StringAttribute{
				MarkdownDescription: "Source network (default: 'any')",
				Optional:            true,
			},
			"source_port": schema.StringAttribute{
				MarkdownDescription: "Source port",
				Optional:            true,
			},
			"destination_net": schema.StringAttribute{
				MarkdownDescription: "Destination network",
				Optional:            true,
			},
			"destination_port": schema.StringAttribute{
				MarkdownDescription: "Destination port",
				Required:            true,
			},
			"target_ip": schema.StringAttribute{
				MarkdownDescription: "Target IP address to forward to",
				Required:            true,
			},
			"target_port": schema.StringAttribute{
				MarkdownDescription: "Target port",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description",
				Optional:            true,
			},
			"log": schema.BoolAttribute{
				MarkdownDescription: "Log packets matching this rule",
				Optional:            true,
			},
		},
	}
}

func (r *NatDestinationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NatDestinationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NatDestinationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	natData := map[string]interface{}{
		"rule": map[string]interface{}{
			"interface":   data.Interface.ValueString(),
			"protocol":    data.Protocol.ValueString(),
			"dst_port":    data.DestinationPort.ValueString(),
			"target":      data.TargetIP.ValueString(),
			"local_port":  data.TargetPort.ValueString(),
		},
	}

	if !data.Enabled.IsNull() {
		if data.Enabled.ValueBool() {
			natData["rule"].(map[string]interface{})["enabled"] = "1"
		} else {
			natData["rule"].(map[string]interface{})["enabled"] = "0"
		}
	} else {
		natData["rule"].(map[string]interface{})["enabled"] = "1"
	}

	if !data.SourceNet.IsNull() {
		natData["rule"].(map[string]interface{})["source"] = data.SourceNet.ValueString()
	}

	if !data.SourcePort.IsNull() {
		natData["rule"].(map[string]interface{})["src_port"] = data.SourcePort.ValueString()
	}

	if !data.DestinationNet.IsNull() {
		natData["rule"].(map[string]interface{})["destination"] = data.DestinationNet.ValueString()
	}

	if !data.Description.IsNull() {
		natData["rule"].(map[string]interface{})["description"] = data.Description.ValueString()
	}

	if !data.Log.IsNull() && data.Log.ValueBool() {
		natData["rule"].(map[string]interface{})["log"] = "1"
	}

	jsonData, _ := json.Marshal(natData)

	url := fmt.Sprintf("%s/api/firewall/d_nat/add_rule", r.client.Host)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonData)))
	httpReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create NAT rule: %s", err))
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
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("No UUID returned from API: %s", string(body)))
		return
	}

	// Apply configuration
	applyURL := fmt.Sprintf("%s/api/firewall/apply", r.client.Host)
	applyReq, _ := http.NewRequestWithContext(ctx, "POST", applyURL, nil)
	applyReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)
	r.client.client.Do(applyReq)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NatDestinationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NatDestinationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/api/firewall/d_nat/get_rule/%s", r.client.Host, data.ID.ValueString())
	httpReq, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	httpReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read NAT rule: %s", err))
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NatDestinationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data NatDestinationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	natData := map[string]interface{}{
		"rule": map[string]interface{}{
			"interface":   data.Interface.ValueString(),
			"protocol":    data.Protocol.ValueString(),
			"dst_port":    data.DestinationPort.ValueString(),
			"target":      data.TargetIP.ValueString(),
			"local_port":  data.TargetPort.ValueString(),
		},
	}

	if !data.Enabled.IsNull() {
		if data.Enabled.ValueBool() {
			natData["rule"].(map[string]interface{})["enabled"] = "1"
		} else {
			natData["rule"].(map[string]interface{})["enabled"] = "0"
		}
	}

	if !data.SourceNet.IsNull() {
		natData["rule"].(map[string]interface{})["source"] = data.SourceNet.ValueString()
	}

	if !data.SourcePort.IsNull() {
		natData["rule"].(map[string]interface{})["src_port"] = data.SourcePort.ValueString()
	}

	if !data.DestinationNet.IsNull() {
		natData["rule"].(map[string]interface{})["destination"] = data.DestinationNet.ValueString()
	}

	if !data.Description.IsNull() {
		natData["rule"].(map[string]interface{})["description"] = data.Description.ValueString()
	}

	if !data.Log.IsNull() && data.Log.ValueBool() {
		natData["rule"].(map[string]interface{})["log"] = "1"
	}

	jsonData, _ := json.Marshal(natData)

	url := fmt.Sprintf("%s/api/firewall/d_nat/set_rule/%s", r.client.Host, data.ID.ValueString())
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonData)))
	httpReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update NAT rule: %s", err))
		return
	}
	defer httpResp.Body.Close()

	// Apply configuration
	applyURL := fmt.Sprintf("%s/api/firewall/apply", r.client.Host)
	applyReq, _ := http.NewRequestWithContext(ctx, "POST", applyURL, nil)
	applyReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)
	r.client.client.Do(applyReq)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NatDestinationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NatDestinationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/api/firewall/d_nat/del_rule/%s", r.client.Host, data.ID.ValueString())
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, nil)
	httpReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete NAT rule: %s", err))
		return
	}
	defer httpResp.Body.Close()

	// Apply configuration
	applyURL := fmt.Sprintf("%s/api/firewall/apply", r.client.Host)
	applyReq, _ := http.NewRequestWithContext(ctx, "POST", applyURL, nil)
	applyReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)
	r.client.client.Do(applyReq)
}

func (r *NatDestinationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
