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
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &FirewallAliasResource{}
var _ resource.ResourceWithImportState = &FirewallAliasResource{}

func NewFirewallAliasResource() resource.Resource {
	return &FirewallAliasResource{}
}

type FirewallAliasResource struct {
	client *Client
}

type FirewallAliasResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Content     types.List   `tfsdk:"content"`
	Description types.String `tfsdk:"description"`
	Enabled     types.Bool   `tfsdk:"enabled"`
}

func (r *FirewallAliasResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_alias"
}

func (r *FirewallAliasResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages OPNsense firewall aliases",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Alias UUID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the alias",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of alias (host, network, port, url, urltable, geoip, networkgroup, mac, external, etc.)",
				Required:            true,
			},
			"content": schema.ListAttribute{
				MarkdownDescription: "List of alias entries (IPs, networks, ports, etc.)",
				Required:            true,
				ElementType:         types.StringType,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the alias",
				Optional:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the alias is enabled",
				Optional:            true,
			},
		},
	}
}

func (r *FirewallAliasResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FirewallAliasResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FirewallAliasResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert content list to newline-separated string (not comma-separated!)
	var contentItems []string
	resp.Diagnostics.Append(data.Content.ElementsAs(ctx, &contentItems, false)...)
	contentStr := strings.Join(contentItems, "\n")  // Changed from "," to "\n"

	aliasData := map[string]interface{}{
		"alias": map[string]interface{}{
			"name":    data.Name.ValueString(),
			"type":    data.Type.ValueString(),
			"content": contentStr,
		},
	}

	if !data.Description.IsNull() {
		aliasData["alias"].(map[string]interface{})["description"] = data.Description.ValueString()
	}
	if !data.Enabled.IsNull() {
		if data.Enabled.ValueBool() {
			aliasData["alias"].(map[string]interface{})["enabled"] = "1"
		} else {
			aliasData["alias"].(map[string]interface{})["enabled"] = "0"
		}
	} else {
		aliasData["alias"].(map[string]interface{})["enabled"] = "1"
	}

	jsonData, _ := json.Marshal(aliasData)

	url := fmt.Sprintf("%s/api/firewall/alias/addItem", r.client.Host)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonData)))
	httpReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create alias: %s", err))
		return
	}
	defer httpResp.Body.Close()

	body, _ := io.ReadAll(httpResp.Body)

	// Debug: Log raw response
	tflog.Debug(ctx, "Raw API Response", map[string]any{
		"status_code": httpResp.StatusCode,
		"body":        string(body),
	})

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to parse response: %s. Body: %s", err, string(body)))
		return
	}

	// Try multiple possible response formats
	var uuid string
	
	// Format 1: {"uuid": "..."}
	if val, ok := result["uuid"].(string); ok && val != "" {
		uuid = val
	}
	
	// Format 2: {"result": "saved", "uuid": "..."}
	if uuid == "" {
		if val, ok := result["uuid"].(string); ok && val != "" {
			uuid = val
		}
	}
	
	// Format 3: Check if there's a nested structure
	if uuid == "" {
		if alias, ok := result["alias"].(map[string]interface{}); ok {
			if val, ok := alias["uuid"].(string); ok && val != "" {
				uuid = val
			}
		}
	}

	// Format 4: Maybe it returns the data back with UUID embedded
	if uuid == "" {
		// Log what we got to help debug
		tflog.Warn(ctx, "UUID not found in expected locations", map[string]any{
			"response_keys": fmt.Sprintf("%v", getKeys(result)),
			"full_response": string(body),
		})
	}

	if uuid == "" {
		resp.Diagnostics.AddError(
			"API Error", 
			fmt.Sprintf("No UUID returned from API. Response: %s", string(body)),
		)
		return
	}

	data.ID = types.StringValue(uuid)

	// Apply configuration
	applyURL := fmt.Sprintf("%s/api/firewall/alias/reconfigure", r.client.Host)
	applyReq, _ := http.NewRequestWithContext(ctx, "POST", applyURL, nil)
	applyReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)
	r.client.client.Do(applyReq)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FirewallAliasResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FirewallAliasResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/api/firewall/alias/getItem/%s", r.client.Host, data.ID.ValueString())
	httpReq, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	httpReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read alias: %s", err))
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FirewallAliasResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data FirewallAliasResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var contentItems []string
	resp.Diagnostics.Append(data.Content.ElementsAs(ctx, &contentItems, false)...)
	contentStr := strings.Join(contentItems, "\n")  // Changed from "," to "\n"

	aliasData := map[string]interface{}{
		"alias": map[string]interface{}{
			"name":    data.Name.ValueString(),
			"type":    data.Type.ValueString(),
			"content": contentStr,
		},
	}

	if !data.Description.IsNull() {
		aliasData["alias"].(map[string]interface{})["description"] = data.Description.ValueString()
	}
	if !data.Enabled.IsNull() {
		if data.Enabled.ValueBool() {
			aliasData["alias"].(map[string]interface{})["enabled"] = "1"
		} else {
			aliasData["alias"].(map[string]interface{})["enabled"] = "0"
		}
	}

	jsonData, _ := json.Marshal(aliasData)

	url := fmt.Sprintf("%s/api/firewall/alias/setItem/%s", r.client.Host, data.ID.ValueString())
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonData)))
	httpReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update alias: %s", err))
		return
	}
	defer httpResp.Body.Close()

	// Apply configuration
	applyURL := fmt.Sprintf("%s/api/firewall/alias/reconfigure", r.client.Host)
	applyReq, _ := http.NewRequestWithContext(ctx, "POST", applyURL, nil)
	applyReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)
	r.client.client.Do(applyReq)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FirewallAliasResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data FirewallAliasResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/api/firewall/alias/delItem/%s", r.client.Host, data.ID.ValueString())
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, nil)
	httpReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete alias: %s", err))
		return
	}
	defer httpResp.Body.Close()

	// Apply configuration
	applyURL := fmt.Sprintf("%s/api/firewall/alias/reconfigure", r.client.Host)
	applyReq, _ := http.NewRequestWithContext(ctx, "POST", applyURL, nil)
	applyReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)
	r.client.client.Do(applyReq)
}

func (r *FirewallAliasResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper function to get keys from map for debugging
func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
