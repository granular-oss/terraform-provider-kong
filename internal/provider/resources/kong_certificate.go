package provider

import (
	"context"
	"fmt"

	kongModels "github.com/granular-oss/terraform-provider-kong/internal/provider/models"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kong/go-kong/kong"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &kongCertificateResource{}
	_ resource.ResourceWithConfigure   = &kongCertificateResource{}
	_ resource.ResourceWithImportState = &kongCertificateResource{}
)

func KongCertificateResource() resource.Resource {
	return &kongCertificateResource{}
}

// kongCertificateResource is the resource implementation.
type kongCertificateResource struct {
	client *kong.Client
}

// Metadata returns the resource type name.
func (r *kongCertificateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "kong_certificate"
}

// Schema defines the schema for the resource.
func (r *kongCertificateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, Optional: false, Required: false, PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			}},
			"cert": schema.StringAttribute{Required: true, Sensitive: true},
			"key":  schema.StringAttribute{Required: true, Sensitive: true},
			"snis": schema.ListAttribute{Optional: true, Computed: true, ElementType: types.StringType},
			"tags": schema.ListAttribute{Optional: true, ElementType: types.StringType},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *kongCertificateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *kongCertificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan kongModels.KongCertificateModel
	client := r.client.Certificates
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	reqBody := kongModels.RequestFromCertificateModel(ctx, &plan)
	tflog.Debug(ctx, "Creating kong_certificate", map[string]any{"body": reqBody, "model": plan})
	cert, err := client.Create(ctx, reqBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create kong_certificate", err.Error())
		return
	}
	result := kongModels.CertificateModelFromResponse(cert)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *kongCertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state kongModels.KongCertificateModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	cert, err := r.client.Certificates.Get(ctx, state.ID.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError("Error Reading kong_certificate", "Could not read kong_certificate with ID "+state.ID.ValueString()+": "+err.Error())
		return
	}
	tflog.Debug(ctx, "Reading kong_certificate state", map[string]any{"state_file": state, "upstream": cert})
	diags = resp.State.Set(ctx, kongModels.CertificateModelFromResponse(cert))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *kongCertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan kongModels.KongCertificateModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// This will actually call PUT because the ID property is set. If we use the .Update(), it will
	// use the PATCH verb and exclude and nulls. Using PUT allows us to update values to null
	reqBody := kongModels.RequestFromCertificateModel(ctx, &plan)
	tflog.Debug(ctx, "Updating kong_certificate", map[string]any{"body": reqBody, "model": plan})
	cert, err := r.client.Certificates.Create(ctx, reqBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to update kong_certificate", err.Error())
		return
	}
	result := kongModels.CertificateModelFromResponse(cert)
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *kongCertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state kongModels.KongCertificateModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.Certificates.Delete(ctx, state.ID.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete kong certificate", err.Error())
		return
	}
}

func (r *kongCertificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	cert, err := r.client.Certificates.Get(ctx, kong.String(req.ID))
	if err != nil {
		resp.Diagnostics.AddError("Cannot find kong certificate for the given id", err.Error())
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), cert.ID)...)
}
