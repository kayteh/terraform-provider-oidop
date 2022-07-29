package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kayteh/podio-go"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.ResourceType = appFieldResourceType{}
var _ tfsdk.Resource = appField{}

type appFieldResourceType struct{}

func (t appFieldResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "A field within an app",

		Attributes: map[string]tfsdk.Attribute{
			"field_id": {
				Type:     types.Int64Type,
				Computed: true,
			},
			"app_id": {
				Type:     types.Int64Type,
				Required: true,
			},
			"type": {
				Type:     types.StringType,
				Required: true,
			},
			"label": {
				Type:     types.StringType,
				Required: true,
			},
			"description": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
			},
			"required": {
				Type:     types.BoolType,
				Optional: true,
				Computed: true,
			},
		},
	}, nil
}

func (t appFieldResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return appField{
		provider: provider,
	}, diags
}

type appFieldData struct {
	FieldID types.Int64 `tfsdk:"field_id"`
	// "type": The type of the field (see area for more information),
	AppID types.Int64 `tfsdk:"app_id"`
	// "field_id": The id of the field,
	Type types.String `tfsdk:"type"`
	// "label": The label of the field, which is what the users will see,
	Label types.String `tfsdk:"label"`
	// "description": The description of the field, shown to the user when inserting and editing,
	Description types.String `tfsdk:"description"`
	// "required": True if the field is required when creating and editing items, false otherwise
	Required types.Bool `tfsdk:"required"`
}

type appField struct {
	provider provider
}

func (r appField) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var plan appFieldData

	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var required bool
	if plan.Required.Unknown || plan.Required.Null {
		required = false
	} else {
		required = plan.Required.Value
	}

	createdAppField, err := r.provider.client.CreateField(strconv.Itoa(int(plan.AppID.Value)), podio.CreateFieldParams{
		Type: plan.Type.Value,
		Config: podio.FieldConfig{
			Label:       plan.Label.Value,
			Description: plan.Description.Value,
			Required:    required,
		},
	})

	if err != nil {
		resp.Diagnostics.AddError("Error creating app_field resource", err.Error())
		return
	}

	tflog.Trace(ctx, "created a field in Podio")

	result := appFieldData{
		FieldID:     types.Int64{Value: int64(createdAppField.FieldID)},
		AppID:       types.Int64{Value: plan.AppID.Value},
		Type:        types.String{Value: createdAppField.Type},
		Label:       types.String{Value: createdAppField.Config.Label},
		Description: types.String{Value: createdAppField.Config.Description},
		Required:    types.Bool{Value: createdAppField.Config.Required},
	}

	diags = resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags...)
}

func (r appField) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state appFieldData

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	readAppField, err := r.provider.client.GetField(strconv.Itoa(int(state.AppID.Value)), strconv.Itoa(int(state.FieldID.Value)))

	if err != nil {
		resp.Diagnostics.AddError("Error refreshing internal state for an app field", err.Error())
	}

	latestState := appFieldData{
		FieldID:     types.Int64{Value: int64(readAppField.FieldID)},
		AppID:       types.Int64{Value: state.AppID.Value},
		Type:        types.String{Value: readAppField.Type},
		Label:       types.String{Value: readAppField.Config.Label},
		Description: types.String{Value: readAppField.Config.Description},
		Required:    types.Bool{Value: readAppField.Config.Required},
	}

	diags = resp.State.Set(ctx, &latestState)
	resp.Diagnostics.Append(diags...)
}

func (r appField) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var plan appFieldData
	var state appFieldData

	planReadDiags := req.Config.Get(ctx, &plan)
	stateReadDiags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(planReadDiags...)
	resp.Diagnostics.Append(stateReadDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var required bool
	if plan.Required.Unknown || plan.Required.Null {
		required = false
	} else {
		required = plan.Required.Value
	}

	tflog.Info(ctx, "FieldID: "+strconv.Itoa(int(state.FieldID.Value)))

	updatedAppField, err := r.provider.client.UpdateField(strconv.Itoa(int(state.AppID.Value)), strconv.Itoa(int(state.FieldID.Value)), podio.FieldConfig{
		Label:       plan.Label.Value,
		Description: plan.Description.Value,
		Required:    required,
	})

	if err != nil {
		resp.Diagnostics.AddError("Error updating the app_field resource", err.Error())
		return
	}

	tflog.Trace(ctx, "updated the field in Podio")

	result := appFieldData{
		FieldID:     types.Int64{Value: state.FieldID.Value},
		AppID:       types.Int64{Value: state.AppID.Value},
		Label:       types.String{Value: updatedAppField.Config.Label},
		Description: types.String{Value: updatedAppField.Config.Description},
		Required:    types.Bool{Value: updatedAppField.Config.Required},
		Type:        types.String{Value: updatedAppField.Type},
	}

	planReadDiags = resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(planReadDiags...)
}

func (r appField) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state appFieldData

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.provider.client.DeleteField(strconv.Itoa(int(state.AppID.Value)), strconv.Itoa(int(state.FieldID.Value)), false)

	if err != nil {
		resp.Diagnostics.AddError("Error occurred while attempting to delete app_field resource", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r appField) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: <app_id>,<field_id>. Got: %q", req.ID),
		)
		return
	}

	app_id, err := strconv.Atoi(idParts[0])
	if err != nil {
		resp.Diagnostics.AddError("Error parsing app_id while importing podio_app_field", err.Error())
	}

	field_id, err := strconv.Atoi(idParts[1])
	if err != nil {
		resp.Diagnostics.AddError("Error parsing field_id while importing podio_app_field", err.Error())
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("app_id"), types.Int64{Value: int64(app_id)})...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("field_id"), types.Int64{Value: int64(field_id)})...)
}
