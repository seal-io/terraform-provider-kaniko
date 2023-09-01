package kaniko

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/seal-io/terraform-provider-kaniko/utils"
)

// RandomModifier returns a plan modifier set a random string to the planned value.
// Use this when need to generate different plan.
func RandomModifier() planmodifier.String {
	return randomModifier{}
}

// randomModifier implements the plan modifier.
type randomModifier struct{}

// Description returns a human-readable description of the plan modifier.
func (m randomModifier) Description(_ context.Context) string {
	return "Generate random string value for every plan."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m randomModifier) MarkdownDescription(_ context.Context) string {
	return "Generate random string value for every plan."
}

// PlanModifyString implements the plan modification logic.
func (m randomModifier) PlanModifyString(
	_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse,
) {
	switch {
	case !req.PlanValue.IsNull() && req.StateValue.IsNull():
		// For create, skip if there is a plan value and state value is null.
		return
	case !req.PlanValue.IsNull() && !req.StateValue.IsNull() && req.PlanValue != req.StateValue:
		// For update, skip if there is a plan value, and it is not equal to state value.
		return
	default:
		// For update, set a random string to the planned value to change the plan.
		id := types.StringValue(fmt.Sprintf("kaniko-%s", utils.String(8)))
		req.PlanValue = id
		req.StateValue = id

		resp.PlanValue = id
		resp.RequiresReplace = true
	}
}
