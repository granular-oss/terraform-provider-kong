package provider

import (
	"context"
	"fmt"
	"strings"

	kongModels "github.com/granular-oss/terraform-provider-kong/internal/provider/models"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kong/go-kong/kong"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &kongConsumerOauth2Resource{}
	_ resource.ResourceWithConfigure   = &kongConsumerOauth2Resource{}
	_ resource.ResourceWithImportState = &kongConsumerOauth2Resource{}
)

// NewkongConsumerOauth2Resource is a helper function to simplify the provider implementation.
func KongConsumerOauth2Resource() resource.Resource {
	return &kongConsumerOauth2Resource{}
}

// kongConsumerOauth2Resource is the resource implementation.
type kongConsumerOauth2Resource struct {
	client *kong.Client
}

// Metadata returns the resource type name.
func (r *kongConsumerOauth2Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "kong_consumer_oauth2"
}

// Schema defines the schema for the resource.
func (r *kongConsumerOauth2Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Optional: false, Required: false, PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			}},
			"kong_id": schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			}},
			"consumer_id": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			}},
			"name":          schema.StringAttribute{Required: true},
			"client_id":     schema.StringAttribute{Required: true},
			"client_secret": schema.StringAttribute{Required: true, Sensitive: true},
			"hash_secret":   schema.BoolAttribute{Optional: true, Computed: true, Default: booldefault.StaticBool(false)},
			"redirect_uris": schema.ListAttribute{Optional: true, ElementType: types.StringType},
			"tags":          schema.ListAttribute{Optional: true, ElementType: types.StringType},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *kongConsumerOauth2Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*kong.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *kong.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *kongConsumerOauth2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan kongModels.KongConsumerOauth2Model
	client := r.client.Oauth2Credentials
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	reqBody := kongModels.RequestFromConsumerOauth2Model(ctx, &plan)
	tflog.Debug(ctx, "Creating kong_consumer_oauth2", map[string]any{"body": reqBody, "model": plan})

	oauth, err := client.Create(ctx, plan.ConsumerID.ValueStringPointer(), reqBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create kong_consumer_oauth2", err.Error())
		return
	}
	newState := kongModels.ConsumerOauth2ModelFromResponse(oauth)
	if plan.HashSecret.ValueBool() {
		newState.ClientSecret = plan.ClientSecret
	}
	diags = resp.State.Set(ctx, newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *kongConsumerOauth2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state kongModels.KongConsumerOauth2Model
	diags := req.State.Get(ctx, &state)
	idParts := strings.Split(state.ID.ValueString(), ":")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	oauth, err := r.client.Oauth2Credentials.Get(ctx, &idParts[0], &idParts[1])
	if err != nil {
		resp.Diagnostics.AddError("Error Reading kong_consumer_oauth2", "Could not read kong_consumer_oauth2 with ID "+state.ID.ValueString()+": "+err.Error())
		return
	}
	tflog.Debug(ctx, "Reading kong_consumer_oauth2 state", map[string]any{"state_file": state, "upstream": oauth})
	newState := kongModels.ConsumerOauth2ModelFromResponse(oauth)
	if state.HashSecret.ValueBool() {
		newState.ClientSecret = state.ClientSecret
	}
	diags = resp.State.Set(ctx, newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *kongConsumerOauth2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan kongModels.KongConsumerOauth2Model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// This will actually call PUT because the ID property is set. If we use the .Update(), it will
	// use the PATCH verb and exclude and nulls. Using PUT allows us to update values to null
	reqBody := kongModels.RequestFromConsumerOauth2Model(ctx, &plan)
	tflog.Debug(ctx, "Updating kong_consumer_oauth2", map[string]any{"body": reqBody, "model": plan})
	oauth, err := r.client.Oauth2Credentials.Create(ctx, plan.ConsumerID.ValueStringPointer(), reqBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create kong_consumer_oauth2", err.Error())
		return
	}
	newState := kongModels.ConsumerOauth2ModelFromResponse(oauth)
	if plan.HashSecret.ValueBool() {
		newState.ClientSecret = plan.ClientSecret
	}
	diags = resp.State.Set(ctx, newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *kongConsumerOauth2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state kongModels.KongConsumerOauth2Model
	diags := req.State.Get(ctx, &state)
	idParts := strings.Split(state.ID.ValueString(), ":")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.Oauth2Credentials.Delete(ctx, &idParts[0], &idParts[1])
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete kong_consumer_oauth2", err.Error())
		return
	}
}

func (r *kongConsumerOauth2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")
	key, err := r.client.Oauth2Credentials.Get(ctx, &idParts[0], &idParts[1])
	if err != nil {
		resp.Diagnostics.AddError("Cannot find kong_consumer_oauth2 for the given name/id:group", err.Error())
	}
	id := *key.Consumer.ID + ":" + *key.ID
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
