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
	_ datasource.DataSource              = &ServerDataSource{}
	_ datasource.DataSourceWithConfigure = &ServerDataSource{}
)

type ServerDataSource struct {
	client *api.Client
}

type ServerDataSourceModel struct {
	Slug                   types.String `tfsdk:"slug"`
	Name                   types.String `tfsdk:"name"`
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
}

func NewServerDataSource() datasource.DataSource {
	return &ServerDataSource{}
}

func (*ServerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server"
}

// Schema defines the schema for the data source.
func (d *ServerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"slug": schema.StringAttribute{
				Required:    true,
				Description: "Server slug",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Server name",
			},
			"date": schema.StringAttribute{
				Computed:    true,
				Description: "Creation date/time",
			},
			"location": schema.StringAttribute{
				Computed:    true,
				Description: "Location ID of the server",
			},
			"image": schema.StringAttribute{
				Computed:    true,
				Description: "Server image",
			},
			"profile": schema.StringAttribute{
				Computed:    true,
				Description: "Server profile",
			},
			"ipv4": schema.StringAttribute{
				Computed:    true,
				Description: "IPv4 address",
			},
			"ipv6": schema.StringAttribute{
				Computed:    true,
				Description: "IPv6 address",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Server status Enum: provisioning, running, stopped, error, rebooting, starting, stopping, reinstalling",
			},
			"virtualization": schema.StringAttribute{
				Computed:    true,
				Description: "Server virtualization type indicating whether it's a Webdock LXD VPS or a KVM Virtual Machine Enum: container, kvm",
			},
			"web_server": schema.StringAttribute{
				Computed:    true,
				Description: "Webserver type Enum: Apache, Nginx, None",
			},
			"snapshot_run_time": schema.Int64Attribute{
				Computed:    true,
				Description: "Last known snapshot runtime (seconds)",
			},
			"word_press_lock_down": schema.BoolAttribute{
				Computed:    true,
				Description: "Wordpress lockdown status",
			},
			"ssh_password_auth_enabled": schema.BoolAttribute{
				Computed:    true,
				Description: "SSH Password Authentication Enabled for this Server",
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *ServerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read refreshes the Terraform state with the latest data.
func (d *ServerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Preparing to read `item` data source")
	var state ServerDataSourceModel

	// get the user supplied data from the tf datasoruce block
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	server, err := d.client.GetServerBYSlug(ctx, state.Slug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get `server` by slug",
			err.Error(),
		)
		return
	}

	// Map response body to model
	state = ServerDataSourceModel{
		Slug:                   types.StringValue(server.Slug),
		Name:                   types.StringValue(server.Name),
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

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Finished reading `server` data source", map[string]any{"success": true})
}
