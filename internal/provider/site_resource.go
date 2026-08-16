package provider

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdkmodels "github.com/NerdIT-Tech/tplink-omada-sdk-for-go/models"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &SiteResource{}
	_ resource.ResourceWithConfigure   = &SiteResource{}
	_ resource.ResourceWithImportState = &SiteResource{}
)

func NewSiteResource() resource.Resource {
	return &SiteResource{}
}

// SiteResource manages an Omada site, the top-level container a network
// administrator creates to group and manage the devices and clients at one physical
// location.
type SiteResource struct {
	client *omadaClient
}

// SiteResourceModel maps the omada_site resource schema to Go types.
type SiteResourceModel struct {
	ID        types.String  `tfsdk:"id"`
	Name      types.String  `tfsdk:"name"`
	Type      types.Int64   `tfsdk:"type"`
	Region    types.String  `tfsdk:"region"`
	TimeZone  types.String  `tfsdk:"timezone"`
	Address   types.String  `tfsdk:"address"`
	Latitude  types.Float64 `tfsdk:"latitude"`
	Longitude types.Float64 `tfsdk:"longitude"`
	Scenario  types.String  `tfsdk:"scenario"`
	SupportES types.Bool    `tfsdk:"support_es"`
	SupportL2 types.Bool    `tfsdk:"support_l2"`
	TagIDs    types.List    `tfsdk:"tag_ids"`
}

func (r *SiteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site"
}

func (r *SiteResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an Omada site: the top-level container a network administrator " +
			"creates to manage the devices and clients at a specific physical location.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Site ID assigned by the controller.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the site. Must be 1 to 64 characters.",
				Required:            true,
			},
			"type": schema.Int64Attribute{
				MarkdownDescription: "Type of the site: `0` for a basic site, `1` for a pro site. " +
					"Defaults to `0`. Cannot be changed after creation; changing it forces recreation.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "Country/region of the site, as an ISO country code abbreviation " +
					"(e.g. `United States` for the United States of America).",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"timezone": schema.StringAttribute{
				MarkdownDescription: "Timezone of the site. See section 5.1 of the Omada Open API Access Guide " +
					"for accepted values.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"address": schema.StringAttribute{
				MarkdownDescription: "Street address of the site.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"latitude": schema.Float64Attribute{
				MarkdownDescription: "Latitude of the site, within the range -90 to 90.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
				},
			},
			"longitude": schema.Float64Attribute{
				MarkdownDescription: "Longitude of the site, within the range -180 to 180.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
				},
			},
			"scenario": schema.StringAttribute{
				MarkdownDescription: "Deployment scenario of the site. See the controller's " +
					"\"Get scenario list\" interface for accepted values.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"support_es": schema.BoolAttribute{
				MarkdownDescription: "Whether the site supports adopting Agile (ES) series switches.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"support_l2": schema.BoolAttribute{
				MarkdownDescription: "Whether the site supports adopting non-Agile (L2) series switches.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"tag_ids": schema.ListAttribute{
				MarkdownDescription: "IDs of the site tags applied to this site. Tag IDs are created via " +
					"the controller's site tag endpoints.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *SiteResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*omadaClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *omadaClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *SiteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SiteResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	site := sdkmodels.NewCreateSiteEntity()
	name := plan.Name.ValueString()
	site.SetName(&name)

	if !plan.Type.IsNull() {
		t, diags := int32FromInt64(plan.Type.ValueInt64(), path.Root("type"))
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		site.SetTypeEscaped(&t)
	}
	if !plan.Region.IsNull() {
		v := plan.Region.ValueString()
		site.SetRegion(&v)
	}
	if !plan.TimeZone.IsNull() {
		v := plan.TimeZone.ValueString()
		site.SetTimeZone(&v)
	}
	if !plan.Address.IsNull() {
		v := plan.Address.ValueString()
		site.SetAddress(&v)
	}
	if !plan.Latitude.IsNull() {
		v := plan.Latitude.ValueFloat64()
		site.SetLatitude(&v)
	}
	if !plan.Longitude.IsNull() {
		v := plan.Longitude.ValueFloat64()
		site.SetLongitude(&v)
	}
	if !plan.Scenario.IsNull() {
		v := plan.Scenario.ValueString()
		site.SetScenario(&v)
	}
	if !plan.SupportES.IsNull() {
		v := plan.SupportES.ValueBool()
		site.SetSupportES(&v)
	}
	if !plan.SupportL2.IsNull() {
		v := plan.SupportL2.ValueBool()
		site.SetSupportL2(&v)
	}
	if !plan.TagIDs.IsNull() {
		var tagIDs []string
		resp.Diagnostics.Append(plan.TagIDs.ElementsAs(ctx, &tagIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		site.SetTagIds(tagIDs)
	}

	tflog.Debug(ctx, "Creating Omada site", map[string]any{"name": name})

	apiResp, err := r.client.api.Openapi().V1().ByOmadacId(r.client.omadacID).Sites().Post(ctx, site, nil)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create Omada Site", err.Error())
		return
	}
	if diags := checkOperationResponse(apiResp.GetErrorCode(), apiResp.GetMsg()); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	siteID, ok := extractCreatedSiteID(apiResp.GetResult())
	if !ok || siteID == "" {
		resp.Diagnostics.AddError(
			"Unexpected Omada API Response",
			"The controller accepted the site creation request but did not return a site ID.",
		)
		return
	}

	var state SiteResourceModel
	resp.Diagnostics.Append(r.readSite(ctx, siteID, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SiteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SiteResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	siteID := state.ID.ValueString()
	var refreshed SiteResourceModel
	diags := r.readSite(ctx, siteID, &refreshed)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	if refreshed.ID.IsNull() {
		// The site no longer exists on the controller; drop it from state so
		// Terraform plans to recreate it instead of erroring.
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &refreshed)...)
}

func (r *SiteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SiteResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	siteID := plan.ID.ValueString()

	site := sdkmodels.NewUpdateSiteEntity()
	name := plan.Name.ValueString()
	site.SetName(&name)

	if !plan.Region.IsNull() {
		v := plan.Region.ValueString()
		site.SetRegion(&v)
	}
	if !plan.TimeZone.IsNull() {
		v := plan.TimeZone.ValueString()
		site.SetTimeZone(&v)
	}
	if !plan.Address.IsNull() {
		v := plan.Address.ValueString()
		site.SetAddress(&v)
	}
	if !plan.Latitude.IsNull() {
		v := plan.Latitude.ValueFloat64()
		site.SetLatitude(&v)
	}
	if !plan.Longitude.IsNull() {
		v := plan.Longitude.ValueFloat64()
		site.SetLongitude(&v)
	}
	if !plan.Scenario.IsNull() {
		v := plan.Scenario.ValueString()
		site.SetScenario(&v)
	}
	if !plan.SupportES.IsNull() {
		v := plan.SupportES.ValueBool()
		site.SetSupportES(&v)
	}
	if !plan.SupportL2.IsNull() {
		v := plan.SupportL2.ValueBool()
		site.SetSupportL2(&v)
	}
	if !plan.TagIDs.IsNull() {
		var tagIDs []string
		resp.Diagnostics.Append(plan.TagIDs.ElementsAs(ctx, &tagIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		site.SetTagIds(tagIDs)
	}

	tflog.Debug(ctx, "Updating Omada site", map[string]any{"site_id": siteID})

	apiResp, err := r.client.api.Openapi().V1().ByOmadacId(r.client.omadacID).Sites().BySiteId(siteID).Put(ctx, site, nil)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Update Omada Site", err.Error())
		return
	}
	if diags := checkOperationResponse(apiResp.GetErrorCode(), apiResp.GetMsg()); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	var state SiteResourceModel
	resp.Diagnostics.Append(r.readSite(ctx, siteID, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SiteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SiteResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	siteID := state.ID.ValueString()
	tflog.Debug(ctx, "Deleting Omada site", map[string]any{"site_id": siteID})

	apiResp, err := r.client.api.Openapi().V1().ByOmadacId(r.client.omadacID).Sites().BySiteId(siteID).Delete(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Delete Omada Site", err.Error())
		return
	}
	if diags := checkOperationResponse(apiResp.GetErrorCode(), apiResp.GetMsg()); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
}

func (r *SiteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// readSite fetches siteID from the controller and populates model. If the controller
// reports that the site no longer exists, model.ID is left null (unset) rather than
// an error being added to the returned diagnostics, so callers can distinguish
// "gone" from a real failure and react accordingly (Read removes it from state;
// Create/Update treat it as an unexpected condition).
func (r *SiteResource) readSite(ctx context.Context, siteID string, model *SiteResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	apiResp, err := r.client.api.Openapi().V1().ByOmadacId(r.client.omadacID).Sites().BySiteId(siteID).Get(ctx, nil)
	if err != nil {
		diags.AddError("Unable to Read Omada Site", err.Error())
		return diags
	}

	code := apiResp.GetErrorCode()
	if code != nil && *code != 0 {
		msg := ""
		if apiResp.GetMsg() != nil {
			msg = *apiResp.GetMsg()
		}
		if isNotFoundMessage(msg) {
			return diags
		}
		diags.AddError("Omada Controller Rejected Site Read", fmt.Sprintf("controller error %d: %s", *code, msg))
		return diags
	}

	site := apiResp.GetResult()
	if site == nil {
		return diags
	}

	model.ID = types.StringValue(siteID)
	model.Name = stringOrNull(site.GetName())
	model.Region = stringOrNull(site.GetRegion())
	model.TimeZone = stringOrNull(site.GetTimeZone())
	model.Address = stringOrNull(site.GetAddress())
	model.Latitude = float64OrNull(site.GetLatitude())
	model.Longitude = float64OrNull(site.GetLongitude())
	model.Scenario = stringOrNull(site.GetScenario())
	model.SupportES = boolOrNull(site.GetSupportES())
	model.SupportL2 = boolOrNull(site.GetSupportL2())

	if site.GetTypeEscaped() != nil {
		model.Type = types.Int64Value(int64(*site.GetTypeEscaped()))
	} else {
		model.Type = types.Int64Null()
	}

	tagIDs, tagDiags := types.ListValueFrom(ctx, types.StringType, site.GetTagIds())
	diags.Append(tagDiags...)
	model.TagIDs = tagIDs

	return diags
}

// extractCreatedSiteID pulls the new site's ID out of a create-site response's
// result payload. The Omada Open API description doesn't type this payload's
// properties, so the SDK exposes it only as untyped additional data; "siteId" is
// the key the controller has been observed to send.
func extractCreatedSiteID(result sdkmodels.OperationResponse_resultable) (string, bool) {
	if result == nil {
		return "", false
	}
	raw, ok := result.GetAdditionalData()["siteId"]
	if !ok {
		return "", false
	}
	id, ok := raw.(string)
	return id, ok
}

// isNotFoundMessage reports whether msg looks like the controller's way of saying a
// site ID doesn't exist, as opposed to some other controller-side error. The Omada
// Open API doesn't document a stable error code for this case, so a Read that hits
// it falls back to matching the message text.
func isNotFoundMessage(msg string) bool {
	lower := strings.ToLower(msg)
	return strings.Contains(lower, "not exist") || strings.Contains(lower, "not found")
}

// checkOperationResponse returns error diagnostics if code is missing or nonzero,
// the Omada Open API's convention for a controller-side failure.
func checkOperationResponse(code *int32, msg *string) diag.Diagnostics {
	var diags diag.Diagnostics
	if code != nil && *code == 0 {
		return diags
	}
	m := ""
	if msg != nil {
		m = *msg
	}
	errCode := int32(0)
	if code != nil {
		errCode = *code
	}
	diags.AddError("Omada Controller Rejected Request", fmt.Sprintf("controller error %d: %s", errCode, m))
	return diags
}

func stringOrNull(v *string) types.String {
	if v == nil {
		return types.StringNull()
	}
	return types.StringValue(*v)
}

func float64OrNull(v *float64) types.Float64 {
	if v == nil {
		return types.Float64Null()
	}
	return types.Float64Value(*v)
}

func boolOrNull(v *bool) types.Bool {
	if v == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*v)
}

// int32FromInt64 narrows v to an int32, returning an attribute-scoped error
// diagnostic instead of silently overflowing if v is out of int32 range.
func int32FromInt64(v int64, attr path.Path) (int32, diag.Diagnostics) {
	var diags diag.Diagnostics
	if v < math.MinInt32 || v > math.MaxInt32 {
		diags.AddAttributeError(
			attr,
			"Value Out Of Range",
			fmt.Sprintf("value %d does not fit in a 32-bit integer (must be between %d and %d).", v, math.MinInt32, math.MaxInt32),
		)
		return 0, diags
	}
	return int32(v), diags
}
