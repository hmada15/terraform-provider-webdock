package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hmada15/terraform-provider-webdock/api"
)

var (
	_ datasource.DataSource              = &LocationDataSource{}
	_ datasource.DataSourceWithConfigure = &LocationDataSource{}
)

type LocationDataSource struct {
	client *api.Client
}

type LocationDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	City        types.String `tfsdk:"city"`
	Country     types.String `tfsdk:"country"`
	Description types.String `tfsdk:"description"`
	Icon        types.String `tfsdk:"icon"`
}

func NewLocationDataSource() datasource.DataSource {
	return &LocationDataSource{}
}

func (*LocationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_location"
}

// Schema defines the schema for the data source.
func (d *LocationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Location name",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Location name",
			},
			"city": schema.StringAttribute{
				Computed:    true,
				Description: "Location name",
			},
			"country": schema.StringAttribute{
				Computed:    true,
				Description: "Location name",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Location name",
			},
			"icon": schema.StringAttribute{
				Computed:    true,
				Description: "Location name",
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *LocationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Location Data Source Configure Type",
			fmt.Sprintf("Expected *api.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	d.client = client
}

// Read refreshes the Terraform state with the latest data
func (d *LocationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Preparing to read `item` data source")
	var state LocationDataSourceModel

	// get the user supplied data from the tf datasoruce block
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	locations, err := d.client.ListLocations(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list `location`",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, location := range locations {
		state = LocationDataSourceModel{
			ID:          types.StringValue(location.ID),
			Name:        types.StringValue(location.Name),
			City:        types.StringValue(location.City),
			Country:     types.StringValue(location.Country),
			Description: types.StringValue(location.Description),
			Icon:        types.StringValue(location.Icon),
		}
	}
	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Finished reading `location` data source", map[string]any{"success": true})
}
