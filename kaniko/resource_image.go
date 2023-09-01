package kaniko

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"k8s.io/client-go/rest"

	"github.com/seal-io/terraform-provider-kaniko/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &imageResource{}
	_ resource.ResourceWithConfigure = &imageResource{}
)

type imageResourceModel struct {
	ID          types.String `tfsdk:"id"`
	GitUsername types.String `tfsdk:"git_username"`
	GitPassword types.String `tfsdk:"git_password"`

	Context          types.String `tfsdk:"context"`
	Dockerfile       types.String `tfsdk:"dockerfile"`
	Destination      types.String `tfsdk:"destination"`
	BuildArg         types.Map    `tfsdk:"build_arg"`
	RegistryUsername types.String `tfsdk:"registry_username"`
	RegistryPassword types.String `tfsdk:"registry_password"`
	Cache            types.Bool   `tfsdk:"cache"`
	NoPush           types.Bool   `tfsdk:"no_push"`
	PushRetry        types.Int64  `tfsdk:"push_retry"`
	Reproducible     types.Bool   `tfsdk:"reproducible"`
	Verbosity        types.String `tfsdk:"verbosity"`
}

// NewImageResource is a helper function to simplify the provider implementation.
func NewImageResource() resource.Resource {
	return &imageResource{}
}

// imageResource is the resource implementation.
type imageResource struct {
	restConfig *rest.Config
}

// Metadata returns the resource type name.
func (r *imageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_image"
}

// Schema defines the schema for the resource.
func (r *imageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Specify the image to build.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					RandomModifier(),
				},
			},
			"git_username": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Username for the git clone",
			},
			"git_password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Password for the git clone",
			},
			"registry_username": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Username for the image registry",
			},
			"registry_password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Password for the image registry",
			},
			"context": schema.StringAttribute{
				Required:    true,
				Description: "Location of the build context",
			},
			"destination": schema.StringAttribute{
				Required:    true,
				Description: "Image name to be built and pushed.",
			},
			"dockerfile": schema.StringAttribute{
				Optional:    true,
				Description: "Path to the dockerfile to be built. (default \"Dockerfile\")",
			},
			"build_arg": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Arguments at build time.",
			},
			"cache": schema.BoolAttribute{
				Optional:    true,
				Description: "Set to true to opt in caching",
			},
			"no_push": schema.BoolAttribute{
				Optional:    true,
				Description: "Set to true if you only want to build the image, without pushing to a registry",
			},
			"push_retry": schema.Int64Attribute{
				Optional:    true,
				Description: "Number of retries for the push operation",
			},
			"reproducible": schema.BoolAttribute{
				Optional:    true,
				Description: "Set to true to strip timestamps out of the built image and make it reproducible.",
			},
			"verbosity": schema.StringAttribute{
				Optional:    true,
				Description: "Log level (trace, debug, info, warn, error, fatal, panic) (default info)",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *imageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Start Create")

	// Retrieve values from plan.
	var plan imageResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("kaniko-%s", utils.String(8)))

	state, err := r.build(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("kaniko build failed", err.Error())
		return
	}

	// Set state to fully populated data.
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *imageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state.
	var state imageResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state.
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *imageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Start Update")

	// Retrieve values from plan.
	var plan imageResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.build(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("kaniko build failed", err.Error())
		return
	}

	// Set state to fully populated data.
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *imageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

// Configure adds the provider configured client to the resource.
func (r *imageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	var ok bool
	r.restConfig, ok = req.ProviderData.(*rest.Config)
	if !ok {
		resp.Diagnostics.AddError("invalid provider data", "expected a rest config")
	}
}

func (r *imageResource) build(ctx context.Context, plan imageResourceModel) (*imageResourceModel, error) {
	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	gitUsername := os.Getenv("GIT_USERNAME")
	gitPassword := os.Getenv("GIT_PASSWORD")
	registryUsername := os.Getenv("REGISTRY_USERNAME")
	registryPassword := os.Getenv("REGISTRY_PASSWORD")
	var pushRetry int64 = 5
	verbosity := "debug"

	if !plan.GitUsername.IsNull() {
		gitUsername = plan.GitUsername.ValueString()
	}

	if !plan.GitPassword.IsNull() {
		gitPassword = plan.GitPassword.ValueString()
	}

	if !plan.RegistryUsername.IsNull() {
		registryUsername = plan.RegistryUsername.ValueString()
	}

	if !plan.RegistryPassword.IsNull() {
		registryPassword = plan.RegistryPassword.ValueString()
	}

	if !plan.PushRetry.IsNull() {
		pushRetry = plan.PushRetry.ValueInt64()
	}

	if !plan.Verbosity.IsNull() {
		verbosity = plan.Verbosity.ValueString()
	}

	options := &runOptions{
		ID:               plan.ID.ValueString(),
		GitPassword:      gitPassword,
		GitUsername:      gitUsername,
		RegistryUsername: registryUsername,
		RegistryPassword: registryPassword,
		Context:          plan.Context.ValueString(),
		Dockerfile:       plan.Dockerfile.ValueString(),
		Destination:      plan.Destination.ValueString(),
		Cache:            plan.Cache.ValueBool(),
		NoPush:           plan.NoPush.ValueBool(),
		PushRetry:        pushRetry,
		Reproducible:     plan.Reproducible.ValueBool(),
		Verbosity:        verbosity,
	}

	err := kanikoBuild(ctx, r.restConfig, options)
	if err != nil {
		return nil, err
	}

	return &plan, nil
}
