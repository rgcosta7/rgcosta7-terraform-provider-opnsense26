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

var _ resource.Resource = &WireguardPeerResource{}
var _ resource.ResourceWithImportState = &WireguardPeerResource{}

func NewWireguardPeerResource() resource.Resource {
	return &WireguardPeerResource{}
}

type WireguardPeerResource struct {
	client *Client
}

type WireguardPeerResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Enabled          types.Bool   `tfsdk:"enabled"`
	PublicKey        types.String `tfsdk:"public_key"`
	AllowedIPs       types.String `tfsdk:"allowed_ips"`
	Endpoint         types.String `tfsdk:"endpoint"`
	EndpointPort     types.Int64  `tfsdk:"endpoint_port"`
	PresharedKey     types.String `tfsdk:"preshared_key"`
	Keepalive        types.Int64  `tfsdk:"keepalive"`
}

func (r *WireguardPeerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_wireguard_peer"
}

func (r *WireguardPeerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages WireGuard peers in OPNsense",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Peer UUID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the peer",
				Required:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the peer is enabled",
				Optional:            true,
			},
			"public_key": schema.StringAttribute{
				MarkdownDescription: "Peer's public key",
				Required:            true,
			},
			"allowed_ips": schema.StringAttribute{
				MarkdownDescription: "Comma-separated list of allowed IP addresses/networks",
				Required:            true,
			},
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Endpoint hostname or IP address",
				Optional:            true,
			},
			"endpoint_port": schema.Int64Attribute{
				MarkdownDescription: "Endpoint port",
				Optional:            true,
			},
			"preshared_key": schema.StringAttribute{
				MarkdownDescription: "Pre-shared key for additional security",
				Optional:            true,
				Sensitive:           true,
			},
			"keepalive": schema.Int64Attribute{
				MarkdownDescription: "Persistent keepalive interval in seconds",
				Optional:            true,
			},
		},
	}
}

func (r *WireguardPeerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *WireguardPeerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data WireguardPeerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	peerData := map[string]interface{}{
		"client": map[string]interface{}{
			"name":       data.Name.ValueString(),
			"pubkey":     data.PublicKey.ValueString(),
			"tunneladdress": data.AllowedIPs.ValueString(),
		},
	}

	if !data.Enabled.IsNull() {
		if data.Enabled.ValueBool() {
			peerData["client"].(map[string]interface{})["enabled"] = "1"
		} else {
			peerData["client"].(map[string]interface{})["enabled"] = "0"
		}
	} else {
		peerData["client"].(map[string]interface{})["enabled"] = "1"
	}

	if !data.Endpoint.IsNull() {
		peerData["client"].(map[string]interface{})["serveraddress"] = data.Endpoint.ValueString()
	}

	if !data.EndpointPort.IsNull() {
		peerData["client"].(map[string]interface{})["serverport"] = fmt.Sprintf("%d", data.EndpointPort.ValueInt64())
	}

	if !data.PresharedKey.IsNull() {
		peerData["client"].(map[string]interface{})["psk"] = data.PresharedKey.ValueString()
	}

	if !data.Keepalive.IsNull() {
		peerData["client"].(map[string]interface{})["keepalive"] = fmt.Sprintf("%d", data.Keepalive.ValueInt64())
	}

	jsonData, _ := json.Marshal(peerData)

	url := fmt.Sprintf("%s/api/wireguard/client/add_client", r.client.Host)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonData)))
	httpReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create peer: %s", err))
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
	applyURL := fmt.Sprintf("%s/api/wireguard/service/reconfigure", r.client.Host)
	applyReq, _ := http.NewRequestWithContext(ctx, "POST", applyURL, nil)
	applyReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)
	r.client.client.Do(applyReq)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WireguardPeerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data WireguardPeerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/api/wireguard/client/get_client/%s", r.client.Host, data.ID.ValueString())
	httpReq, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	httpReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read peer: %s", err))
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WireguardPeerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data WireguardPeerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	peerData := map[string]interface{}{
		"client": map[string]interface{}{
			"name":       data.Name.ValueString(),
			"pubkey":     data.PublicKey.ValueString(),
			"tunneladdress": data.AllowedIPs.ValueString(),
		},
	}

	if !data.Enabled.IsNull() {
		if data.Enabled.ValueBool() {
			peerData["client"].(map[string]interface{})["enabled"] = "1"
		} else {
			peerData["client"].(map[string]interface{})["enabled"] = "0"
		}
	}

	if !data.Endpoint.IsNull() {
		peerData["client"].(map[string]interface{})["serveraddress"] = data.Endpoint.ValueString()
	}

	if !data.EndpointPort.IsNull() {
		peerData["client"].(map[string]interface{})["serverport"] = fmt.Sprintf("%d", data.EndpointPort.ValueInt64())
	}

	if !data.PresharedKey.IsNull() {
		peerData["client"].(map[string]interface{})["psk"] = data.PresharedKey.ValueString()
	}

	if !data.Keepalive.IsNull() {
		peerData["client"].(map[string]interface{})["keepalive"] = fmt.Sprintf("%d", data.Keepalive.ValueInt64())
	}

	jsonData, _ := json.Marshal(peerData)

	url := fmt.Sprintf("%s/api/wireguard/client/set_client/%s", r.client.Host, data.ID.ValueString())
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonData)))
	httpReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update peer: %s", err))
		return
	}
	defer httpResp.Body.Close()

	// Apply configuration
	applyURL := fmt.Sprintf("%s/api/wireguard/service/reconfigure", r.client.Host)
	applyReq, _ := http.NewRequestWithContext(ctx, "POST", applyURL, nil)
	applyReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)
	r.client.client.Do(applyReq)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WireguardPeerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data WireguardPeerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/api/wireguard/client/del_client/%s", r.client.Host, data.ID.ValueString())
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, nil)
	httpReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete peer: %s", err))
		return
	}
	defer httpResp.Body.Close()

	// Apply configuration
	applyURL := fmt.Sprintf("%s/api/wireguard/service/reconfigure", r.client.Host)
	applyReq, _ := http.NewRequestWithContext(ctx, "POST", applyURL, nil)
	applyReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)
	r.client.client.Do(applyReq)
}

func (r *WireguardPeerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
