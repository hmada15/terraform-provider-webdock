package provider

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hmada15/terraform-provider-webdock/api"
)

// implement resource interfaces.
var (
	_ resource.Resource                = &PublicKeyResource{}
	_ resource.ResourceWithConfigure   = &PublicKeyResource{}
	_ resource.ResourceWithImportState = &PublicKeyResource{}
)

// NewPublicKeyResource is a helper function to simplify the provider implementation.
func NewPublicKeyResource() resource.Resource {
	return &PublicKeyResource{}
}

// PublicKeyResource is the resource implementation.
type PublicKeyResource struct {
	client *api.Client
}

// PublicKeyResource is the model implementation.
type PublicKeyResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Key         types.String `tfsdk:"key"`
	Created     types.String `tfsdk:"created"`
	PublicKey   types.String `tfsdk:"public_key"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

// Metadata returns the resource type name.
func (s *PublicKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_key"
}

// Configure adds the provider configured client to the data source.
func (d *PublicKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected PublicKey Data Source Configure Type",
			fmt.Sprintf("Expected *api.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	d.client = client
}

// Schema defines the schema for the resource.
func (s *PublicKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"key": schema.StringAttribute{
				Computed: true,
			},
			"created": schema.StringAttribute{
				Computed: true,
			},
			"public_key": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Import using slug as the attribute
func (s *PublicKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create a new resource.
func (s *PublicKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "create publickey")
	// Retrieve values from plan
	var plan PublicKeyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	publicKeyRequest := api.PublicKeyRequest{
		Name:      plan.Name.ValueString(),
		PublicKey: plan.PublicKey.ValueString(),
	}

	publicKey, err := s.client.CreatePublicKey(ctx, publicKeyRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating publickey",
			"Could not create publickey, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = PublicKeyResourceModel{
		ID:          types.StringValue(strconv.Itoa(publicKey.ID)),
		Name:        types.StringValue(publicKey.Name),
		Key:         types.StringValue(publicKey.Key),
		Created:     types.StringValue(publicKey.Created),
		PublicKey:   types.StringValue(publicKey.Key),
		LastUpdated: types.StringValue(time.Now().Format(time.RFC850)),
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "finish createing publickey request")
}

// Read resource information.
func (s *PublicKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "read public key")

	// Get current state
	var state PublicKeyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "send get public key request")
	// Get refreshed public key value from Webdock
	publicKey, err := s.client.GetPublicKeyById(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Webdock publickey",
			"Could not read Webdock publickey id "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
	if (api.PublicKey{}) == publicKey {
		resp.Diagnostics.AddError(
			"Unable to get publickey",
			"Could not read Webdock publickey id "+state.ID.ValueString(),
		)
		return
	}

	// Overwrite items with refreshed state
	state = PublicKeyResourceModel{
		ID:        types.StringValue(strconv.Itoa(publicKey.ID)),
		Name:      types.StringValue(publicKey.Name),
		Key:       types.StringValue(publicKey.Key),
		PublicKey: types.StringValue(publicKey.Key),
		Created:   types.StringValue(publicKey.Created),
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "finish get publickey request")
}

func (s *PublicKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// updating resource is not supported
}

// Delete deletes the resource and removes the Terraform state on success.
func (s *PublicKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "delete publicKey")
	// Get current state
	var state PublicKeyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "send delete publicKey request")
	// delete publicKey
	err := s.client.DeletePublicKey(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleteing webdock publicKey",
			"Could not delete webdock publicKey "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
}
