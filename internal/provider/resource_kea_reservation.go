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

var _ resource.Resource = &KeaReservationResource{}
var _ resource.ResourceWithImportState = &KeaReservationResource{}

func NewKeaReservationResource() resource.Resource {
	return &KeaReservationResource{}
}

type KeaReservationResource struct {
	client *Client
}

type KeaReservationResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Subnet      types.String `tfsdk:"subnet"`
	IPAddress   types.String `tfsdk:"ip_address"`
	HWAddress   types.String `tfsdk:"hw_address"`
	Hostname    types.String `tfsdk:"hostname"`
	Description types.String `tfsdk:"description"`
}

func (r *KeaReservationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kea_reservation"
}

func (r *KeaReservationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Kea DHCP reservations in OPNsense",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Reservation UUID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"subnet": schema.StringAttribute{
				MarkdownDescription: "Subnet UUID this reservation belongs to",
				Required:            true,
			},
			"ip_address": schema.StringAttribute{
				MarkdownDescription: "Reserved IP address",
				Required:            true,
			},
			"hw_address": schema.StringAttribute{
				MarkdownDescription: "Hardware (MAC) address",
				Required:            true,
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "Hostname for this reservation",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the reservation",
				Optional:            true,
			},
		},
	}
}

func (r *KeaReservationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *KeaReservationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data KeaReservationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reservationData := map[string]interface{}{
		"reservation": map[string]interface{}{
			"subnet":     data.Subnet.ValueString(),
			"ip_address": data.IPAddress.ValueString(),
			"hw_address": data.HWAddress.ValueString(),
		},
	}

	if !data.Hostname.IsNull() {
		reservationData["reservation"].(map[string]interface{})["hostname"] = data.Hostname.ValueString()
	}
	if !data.Description.IsNull() {
		reservationData["reservation"].(map[string]interface{})["description"] = data.Description.ValueString()
	}

	jsonData, _ := json.Marshal(reservationData)

	url := fmt.Sprintf("%s/api/kea/dhcpv4/add_reservation", r.client.Host)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonData)))
	httpReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create reservation: %s", err))
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

func (r *KeaReservationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data KeaReservationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/api/kea/dhcpv4/get_reservation/%s", r.client.Host, data.ID.ValueString())
	httpReq, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	httpReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read reservation: %s", err))
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *KeaReservationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data KeaReservationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reservationData := map[string]interface{}{
		"reservation": map[string]interface{}{
			"subnet":     data.Subnet.ValueString(),
			"ip_address": data.IPAddress.ValueString(),
			"hw_address": data.HWAddress.ValueString(),
		},
	}

	if !data.Hostname.IsNull() {
		reservationData["reservation"].(map[string]interface{})["hostname"] = data.Hostname.ValueString()
	}
	if !data.Description.IsNull() {
		reservationData["reservation"].(map[string]interface{})["description"] = data.Description.ValueString()
	}

	jsonData, _ := json.Marshal(reservationData)

	url := fmt.Sprintf("%s/api/kea/dhcpv4/set_reservation/%s", r.client.Host, data.ID.ValueString())
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonData)))
	httpReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update reservation: %s", err))
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

func (r *KeaReservationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data KeaReservationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/api/kea/dhcpv4/del_reservation/%s", r.client.Host, data.ID.ValueString())
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, nil)
	httpReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete reservation: %s", err))
		return
	}
	defer httpResp.Body.Close()

	// Apply configuration
	applyURL := fmt.Sprintf("%s/api/kea/service/reconfigure", r.client.Host)
	applyReq, _ := http.NewRequestWithContext(ctx, "POST", applyURL, nil)
	applyReq.SetBasicAuth(r.client.ApiKey, r.client.ApiSecret)
	r.client.client.Do(applyReq)
}

func (r *KeaReservationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
