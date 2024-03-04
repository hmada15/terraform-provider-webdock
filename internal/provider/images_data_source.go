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
	_ datasource.DataSource              = &ImagesDataSource{}
	_ datasource.DataSourceWithConfigure = &ImagesDataSource{}
)

type ImagesDataSource struct {
	client *api.Client
}

type (
	ImagesDataSourceModel struct {
		Images []ImagesModel `tfsdk:"images"`
	}
	ImagesModel struct {
		Slug       types.String `tfsdk:"slug"`
		Name       types.String `tfsdk:"name"`
		WebServer  types.String `tfsdk:"web_server"`
		PhpVersion types.String `tfsdk:"php_version"`
	}
)

func NewImagesDataSource() datasource.DataSource {
	return &ImagesDataSource{}
}

func (*ImagesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_images"
}

// Schema defines the schema for the data source.
func (d *ImagesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"images": schema.ListNestedAttribute{
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
						"web_server": schema.StringAttribute{
							Computed:    true,
							Description: "RAM memory (in MiB)",
						},
						"php_version": schema.StringAttribute{
							Computed:    true,
							Description: "Disk size (in MiB)",
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *ImagesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Image Data Source Configure Type",
			fmt.Sprintf("Expected *api.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	d.client = client
}

// Read refreshes the Terraform state with the latest data
func (d *ImagesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Preparing to read `item` data source")
	var state ImagesDataSourceModel

	images, err := d.client.ListImages(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list `image`",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, image := range images {
		imageState := ImagesModel{
			Slug:       types.StringValue(image.Slug),
			Name:       types.StringValue(image.Name),
			WebServer:  types.StringValue(image.WebServer),
			PhpVersion: types.StringValue(image.PhpVersion),
		}
		state.Images = append(state.Images, imageState)
	}
	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Finished reading `image` data source", map[string]any{"success": true})
}
