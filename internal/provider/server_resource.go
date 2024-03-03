package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hmada15/terraform-provider-webdock/api"
	"github.com/hmada15/terraform-provider-webdock/helper"
)

// implement resource interfaces.
var (
	_ resource.Resource                = &ServerResource{}
	_ resource.ResourceWithConfigure   = &ServerResource{}
	_ resource.ResourceWithModifyPlan  = &ServerResource{}
	_ resource.ResourceWithImportState = &ServerResource{}
)

// NewServerResource is a helper function to simplify the provider implementation.
func NewServerResource() resource.Resource {
	return &ServerResource{}
}

// ServerResource is the resource implementation.
type ServerResource struct {
	client *api.Client
}

// ServerResource is the model implementation.
type ServerResourceModel struct {
	Slug                   types.String `tfsdk:"slug"`
	Name                   types.String `tfsdk:"name"`
	LocationID             types.String `tfsdk:"location_id"`
	ProfileSlug            types.String `tfsdk:"profile_slug"`
	ImageSlug              types.String `tfsdk:"image_slug"`
	Date                   types.String `tfsdk:"date"`
	Location               types.String `tfsdk:"location"`
	Image                  types.String `tfsdk:"image"`
	Profile                types.String `tfsdk:"profile"`
	Ipv4                   types.String `tfsdk:"ipv4"`
	Ipv6                   types.String `tfsdk:"ipv6"`
	Status                 types.String `tfsdk:"status"`
	Virtualization         types.String `tfsdk:"virtualization"`
	WebServer              types.String `tfsdk:"web_server"`
	SnapshotRunTime        types.Int64  `tfsdk:"snapshot_run_time"`
	WordPressLockDown      types.Bool   `tfsdk:"word_press_lock_down"`
	SSHPasswordAuthEnabled types.Bool   `tfsdk:"ssh_password_auth_enabled"`
	LastUpdated            types.String `tfsdk:"last_updated"`
}

// Metadata returns the resource type name.
func (s *ServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server"
}

// Configure adds the provider configured client to the data source.
func (d *ServerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Server Data Source Configure Type",
			fmt.Sprintf("Expected *api.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	d.client = client
}

// Schema defines the schema for the resource.
func (s *ServerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"slug": schema.StringAttribute{
				Computed: true,
				Optional: true,
				// Requires Replace if the value change
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "Must be unique",
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"location_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the location. Get this from the /locations endpoint.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"profile_slug": schema.StringAttribute{
				Required:    true,
				Description: "Slug of the server profile. Get this from the /profiles endpoint.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"virtualization": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"image_slug": schema.StringAttribute{
				Required:    true,
				Description: "Slug of the server image. Get this from the /images endpoint. You must pass either this parameter or snapshotId",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"date": schema.StringAttribute{
				Computed: true,
			},
			"location": schema.StringAttribute{
				Computed: true,
			},
			"image": schema.StringAttribute{
				Computed: true,
			},
			"profile": schema.StringAttribute{
				Computed: true,
			},
			"ipv4": schema.StringAttribute{
				Computed: true,
			},
			"ipv6": schema.StringAttribute{
				Computed: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"web_server": schema.StringAttribute{
				Computed: true,
			},
			"snapshot_run_time": schema.Int64Attribute{
				Computed: true,
			},
			"word_press_lock_down": schema.BoolAttribute{
				Computed: true,
			},
			"ssh_password_auth_enabled": schema.BoolAttribute{
				Computed: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// ModifyPlan tailor the plan to match the expected end state.
func (s *ServerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Check if the resource is being created.
	if req.State.Raw.IsNull() {
		tflog.Debug(ctx, "start ModifyPlan")

		var state ServerResourceModel
		diags := req.Plan.Get(ctx, &state)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		tflog.Debug(ctx, "check if server exist")
		// check if a server with the slug exist
		exist, err := s.client.ServerExist(ctx, state.Slug.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error checking Webdock server exist",
				"Error: "+err.Error(),
			)
			return
		}
		if exist == helper.YES {
			resp.Diagnostics.AddError(
				"Server with the same slug exist",
				"Webdock require a unique slug for each server and a server with the slug"+
					state.Slug.ValueString()+" already exists please choose new slug",
			)
		}
	}
	// Check if the resource is being destroyed.
	if req.Plan.Raw.IsNull() {
		resp.Diagnostics.AddWarning(
			"Server Deletion requires special privileges which cannot be obtained in the Webdock dashboard without first contacting Webdock Support!",
			"This will nuke the server from orbit including all data and server snapshots. Use with care.",
		)
	}
}

// Import using slug as the attribute
func (s *ServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("slug"), req, resp)
}

// Create a new resource.
func (s *ServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "create server")
	// Retrieve values from plan
	var plan ServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	serverRequest := api.ServerRequest{
		Name:           plan.Name.ValueString(),
		Slug:           plan.Slug.ValueString(),
		LocationID:     plan.LocationID.ValueString(),
		ProfileSlug:    plan.ProfileSlug.ValueString(),
		Virtualization: plan.Virtualization.ValueString(),
		ImageSlug:      plan.ImageSlug.ValueString(),
	}
	tflog.Debug(ctx, "send create server request")
	server, err := s.client.CreateServer(ctx, serverRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating server",
			"Could not create server, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = ServerResourceModel{
		Slug:                   types.StringValue(server.Slug),
		Name:                   types.StringValue(server.Name),
		LocationID:             types.StringValue(server.Location),
		ProfileSlug:            types.StringValue(server.Profile),
		ImageSlug:              types.StringValue(server.Image),
		Date:                   types.StringValue(server.Date),
		Location:               types.StringValue(server.Location),
		Image:                  types.StringValue(server.Image),
		Profile:                types.StringValue(server.Profile),
		Ipv4:                   types.StringValue(server.Ipv4),
		Ipv6:                   types.StringValue(server.Ipv6),
		Status:                 types.StringValue(server.Status),
		Virtualization:         types.StringValue(server.Virtualization),
		WebServer:              types.StringValue(server.WebServer),
		SnapshotRunTime:        types.Int64Value(server.SnapshotRunTime),
		WordPressLockDown:      types.BoolValue(server.WordPressLockDown),
		SSHPasswordAuthEnabled: types.BoolValue(server.SSHPasswordAuthEnabled),
		LastUpdated:            types.StringValue(time.Now().Format(time.RFC850)),
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "finish create server request")
}

// Read resource information.
func (s *ServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "read server")

	// Get current state
	var state ServerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "send get server request")
	// Get refreshed server value from Webdock
	server, err := s.client.GetServerBYSlug(ctx, state.Slug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Webdock server",
			"Could not read Webdock server Slug "+state.Slug.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state = ServerResourceModel{
		Slug:                   types.StringValue(server.Slug),
		Name:                   types.StringValue(server.Name),
		LocationID:             types.StringValue(server.Location),
		ProfileSlug:            types.StringValue(server.Profile),
		ImageSlug:              types.StringValue(server.Image),
		Date:                   types.StringValue(server.Date),
		Location:               types.StringValue(server.Location),
		Image:                  types.StringValue(server.Image),
		Profile:                types.StringValue(server.Profile),
		Ipv4:                   types.StringValue(server.Ipv4),
		Ipv6:                   types.StringValue(server.Ipv6),
		Status:                 types.StringValue(server.Status),
		Virtualization:         types.StringValue(server.Virtualization),
		WebServer:              types.StringValue(server.WebServer),
		SnapshotRunTime:        types.Int64Value(server.SnapshotRunTime),
		WordPressLockDown:      types.BoolValue(server.WordPressLockDown),
		SSHPasswordAuthEnabled: types.BoolValue(server.SSHPasswordAuthEnabled),
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "finish get server request")
}

func (s *ServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// updating resource is not supported
}

// Delete deletes the resource and removes the Terraform state on success.
func (s *ServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "delete server")
	// Get current state
	var state ServerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "send delete server request")
	// delete server
	err := s.client.DeleteServer(ctx, state.Slug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleteing webdock server",
			"Could not delete webdock server "+state.Slug.ValueString()+": "+err.Error(),
		)
		return
	}
}
