package kaniko

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BuildIDModifier returns a plan modifier set build id to unknown string to the planned value.
func BuildIDModifier() planmodifier.String {
	return buildIDModifier{}
}

// buildIDModifier implements the plan modifier.
type buildIDModifier struct{}

// Description returns a human-readable description of the plan modifier.
func (m buildIDModifier) Description(_ context.Context) string {
	return "Set build id to unknown string while need always run for every plan."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m buildIDModifier) MarkdownDescription(_ context.Context) string {
	return "Set build id to unknown string while need always run for every plan."
}

// PlanModifyString implements the plan modification logic.
func (m buildIDModifier) PlanModifyString(
	_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse,
) {
	var (
		ctx  = context.Background()
		plan imageResourceModel
	)

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.AlwaysRun.IsNull() && plan.AlwaysRun.ValueBool() {
		resp.PlanValue = types.StringUnknown()
	}
}
