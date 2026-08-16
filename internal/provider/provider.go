package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var _ provider.Provider = &OmadaProvider{}

// OmadaProvider defines the provider implementation.
type OmadaProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and run locally, and "test" when running acceptance
	// testing.
	version string
}

// OmadaProviderModel describes the provider data model, i.e. the fields
// exposed in a root `provider "omada" { ... }` block.
type OmadaProviderModel struct {
	Host         types.String `tfsdk:"host"`
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	OmadacID     types.String `tfsdk:"omadac_id"`
	Insecure     types.Bool   `tfsdk:"insecure"`
}

// Environment variable fallbacks for provider configuration, mirroring the
// convention used by most Terraform providers so credentials can be kept out
// of configuration files (e.g. in CI).
const (
	envHost         = "OMADA_HOST"
	envClientID     = "OMADA_CLIENT_ID"
	envClientSecret = "OMADA_CLIENT_SECRET" //nolint:gosec // this is an env var name, not a credential
	envOmadacID     = "OMADA_OMADAC_ID"
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &OmadaProvider{
			version: version,
		}
	}
}

func (p *OmadaProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "omada"
	resp.Version = p.version
}

func (p *OmadaProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Omada provider is used to manage a TP-Link Omada SDN Controller. " +
			"It must be configured with credentials for an Omada OpenAPI application before it can be used.",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "URL of the Omada Controller, e.g. `https://omada.example.com:8043`. " +
					"May also be set via the `OMADA_HOST` environment variable.",
				Optional: true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "Client ID of the Omada OpenAPI application. " +
					"May also be set via the `OMADA_CLIENT_ID` environment variable.",
				Optional: true,
			},
			"client_secret": schema.StringAttribute{
				MarkdownDescription: "Client secret of the Omada OpenAPI application. " +
					"May also be set via the `OMADA_CLIENT_SECRET` environment variable.",
				Optional:  true,
				Sensitive: true,
			},
			"omadac_id": schema.StringAttribute{
				MarkdownDescription: "ID of the Omada Controller to manage. " +
					"May also be set via the `OMADA_OMADAC_ID` environment variable.",
				Optional: true,
			},
			"insecure": schema.BoolAttribute{
				MarkdownDescription: "Skip TLS certificate verification when connecting to the controller. " +
					"Useful for controllers using a self-signed certificate. Defaults to `false`.",
				Optional: true,
			},
		},
	}
}

func (p *OmadaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Omada provider")

	var data OmadaProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Terraform config values take precedence over environment variables.
	host := firstNonEmpty(data.Host.ValueString(), os.Getenv(envHost))
	clientID := firstNonEmpty(data.ClientID.ValueString(), os.Getenv(envClientID))
	clientSecret := firstNonEmpty(data.ClientSecret.ValueString(), os.Getenv(envClientSecret))
	omadacID := firstNonEmpty(data.OmadacID.ValueString(), os.Getenv(envOmadacID))

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Omada Controller Host",
			"The provider cannot create the Omada API client because the controller host is unknown. "+
				"Set the host value in the provider configuration or use the OMADA_HOST environment variable.",
		)
	}
	if clientID == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Missing Omada Client ID",
			"The provider cannot create the Omada API client because the client_id is unknown. "+
				"Set the client_id value in the provider configuration or use the OMADA_CLIENT_ID environment variable.",
		)
	}
	if clientSecret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Missing Omada Client Secret",
			"The provider cannot create the Omada API client because the client_secret is unknown. "+
				"Set the client_secret value in the provider configuration or use the OMADA_CLIENT_SECRET environment variable.",
		)
	}
	if omadacID == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("omadac_id"),
			"Missing Omada Controller ID",
			"The provider cannot create the Omada API client because the omadac_id is unknown. "+
				"Set the omadac_id value in the provider configuration or use the OMADA_OMADAC_ID environment variable.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "omada_host", host)
	ctx = tflog.SetField(ctx, "omada_omadac_id", omadacID)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "omada_client_secret")
	tflog.Debug(ctx, "Creating Omada API client")

	// TODO: once github.com/NerdIT-Tech/tplink-omada-sdk-for-go publishes a
	// client, construct and authenticate it here, then hand it to resources
	// and data sources via resp.ResourceData / resp.DataSourceData, e.g.:
	//
	//   client, err := omadasdk.NewClient(host, clientID, clientSecret, omadacID,
	//       omadasdk.WithInsecureSkipVerify(data.Insecure.ValueBool()))
	//   if err != nil {
	//       resp.Diagnostics.AddError("Unable to Create Omada API Client", err.Error())
	//       return
	//   }
	//   resp.ResourceData = client
	//   resp.DataSourceData = client

	tflog.Info(ctx, "Configured Omada provider", map[string]any{"success": true})
}

func (p *OmadaProvider) Resources(_ context.Context) []func() resource.Resource {
	// Resources are registered here as they are implemented, e.g.:
	//   return []func() resource.Resource{ NewSiteResource }
	return []func() resource.Resource{}
}

func (p *OmadaProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	// Data sources are registered here as they are implemented, e.g.:
	//   return []func() datasource.DataSource{ NewSiteDataSource }
	return []func() datasource.DataSource{}
}
