package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hmada15/terraform-provider-webdock/api"
)

var (
	_ datasource.DataSource              = &ProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &ProfileDataSource{}
)

type ProfileDataSource struct {
	client *api.Client
}

type (
	ProfileDataSourceModel struct {
		LocationId types.String   `tfsdk:"location_id"`
		Profile    []ProfileModel `tfsdk:"profiles"`
	}
	ProfileModel struct {
		Slug  types.String `tfsdk:"slug"`
		Name  types.String `tfsdk:"name"`
		RAM   types.Int64  `tfsdk:"ram"`
		Disk  types.Int64  `tfsdk:"disk"`
		CPU   CPU          `tfsdk:"cpu"`
		Price Price        `tfsdk:"price"`
	}
	CPU struct {
		Cores   types.Int64 `tfsdk:"cores"`
		Threads types.Int64 `tfsdk:"threads"`
	}
	Price struct {
		Amount   types.Int64  `tfsdk:"amount"`
		Currency types.String `tfsdk:"currency"`
	}
)

func NewProfileDataSource() datasource.DataSource {
	return &ProfileDataSource{}
}

func (*ProfileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_profile"
}

// Schema defines the schema for the data source.
func (d *ProfileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"location_id": schema.StringAttribute{
				Required:    true,
				Description: "Location of the profile",
			},
			"profiles": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{

						"slug": schema.StringAttribute{
							Computed:    true,
							Description: "Profile slug",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Profile name",
						},
						"ram": schema.Int64Attribute{
							Computed:    true,
							Description: "RAM memory (in MiB)",
						},
						"disk": schema.Int64Attribute{
							Computed:    true,
							Description: "Disk size (in MiB)",
						},
						"cpu": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "CPU model",
							Attributes: map[string]schema.Attribute{
								"cores": schema.Int64Attribute{
									Computed:    true,
									Description: "cpu cores",
								},
								"threads": schema.Int64Attribute{
									Computed:    true,
									Description: "cpu threads",
								},
							},
						},
						"price": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Price model",
							Attributes: map[string]schema.Attribute{
								"amount": schema.Int64Attribute{
									Computed:    true,
									Description: "Price amount",
								},
								"currency": schema.StringAttribute{
									Computed:    true,
									Description: "Price currency",
								},
							},
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *ProfileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Profile Data Source Configure Type",
			fmt.Sprintf("Expected *api.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	d.client = client
}

// Read refreshes the Terraform state with the latest data
func (d *ProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Preparing to read `item` data source")

	// get location_id from config
	var locationId types.String
	diags := req.Config.GetAttribute(ctx, path.Root("location_id"), &locationId)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// list profiles
	profiles, err := d.client.ListProfiles(ctx, locationId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list `profile`",
			err.Error(),
		)
		return
	}
	// Map response body to model
	var state ProfileDataSourceModel
	for _, profile := range profiles {
		profilestate := ProfileModel{
			Name: types.StringValue(profile.Name),
			Slug: types.StringValue(profile.Slug),
			RAM:  types.Int64Value(int64(profile.RAM)),
			Disk: types.Int64Value(int64(profile.Disk)),
		}
		profilestate.CPU.Cores = types.Int64Value(int64(profile.CPU.Cores))
		profilestate.CPU.Threads = types.Int64Value(int64(profile.CPU.Threads))
		profilestate.Price.Amount = types.Int64Value(int64(profile.CPU.Threads))
		profilestate.Price.Currency = types.StringValue(profile.Price.Currency)

		state.Profile = append(state.Profile, profilestate)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Finished reading `profile` data source", map[string]any{"success": true})
}
