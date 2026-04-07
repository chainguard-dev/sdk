/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1

import (
	"slices"
	"testing"

	cgannotations "chainguard.dev/sdk/proto/annotations"
	"chainguard.dev/sdk/uidp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

// getEventAttributes extracts EventAttributes from a proto method descriptor.
func getEventAttributes(t *testing.T, sd protoreflect.ServiceDescriptor, methodName string) *cgannotations.EventAttributes {
	t.Helper()
	md := sd.Methods().ByName(protoreflect.Name(methodName))
	if md == nil {
		t.Fatalf("method %s not found in service %s", methodName, sd.FullName())
	}
	opts, ok := md.Options().(*descriptorpb.MethodOptions)
	if !ok || opts == nil {
		t.Fatalf("method %s has no options", methodName)
	}
	ext, ok := proto.GetExtension(opts, cgannotations.E_Events).(*cgannotations.EventAttributes)
	if !ok || ext == nil {
		t.Fatalf("method %s is missing chainguard.annotations.events", methodName)
	}
	return ext
}

// annotationTest defines a test case for verifying proto event annotations.
type annotationTest struct {
	method   string
	wantType string
	wantExts []string
}

func runAnnotationTests(t *testing.T, sd protoreflect.ServiceDescriptor, tests []annotationTest) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			ea := getEventAttributes(t, sd, tt.method)
			if ea.GetType() != tt.wantType {
				t.Errorf("type = %q, want %q", ea.GetType(), tt.wantType)
			}
			if ea.GetAudience() != cgannotations.EventAttributes_CUSTOMER {
				t.Errorf("audience = %v, want CUSTOMER", ea.GetAudience())
			}
			if !slices.Equal(ea.GetExtensions(), tt.wantExts) {
				t.Errorf("extensions = %v, want %v", ea.GetExtensions(), tt.wantExts)
			}
		})
	}
}

func TestGroupsEventAnnotations(t *testing.T) {
	sd := File_chainguard_platform_iam_v2beta1_groups_proto.Services().ByName("GroupsService")
	if sd == nil {
		t.Fatal("GroupsService not found")
	}
	runAnnotationTests(t, sd, []annotationTest{
		{"CreateGroup", "dev.chainguard.api.iam.group.created.v1", []string{"group"}},
		{"UpdateGroup", "dev.chainguard.api.iam.group.updated.v1", []string{"group"}},
		{"DeleteGroup", "dev.chainguard.api.iam.group.deleted.v1", []string{"group"}},
	})
}

func TestGroupsEventInterfaces(t *testing.T) {
	groupUID := "abc123/def456"

	g := &Group{Uid: groupUID}
	if got := g.CloudEventsSubject(); got != groupUID {
		t.Errorf("Group.CloudEventsSubject() = %q, want %q", got, groupUID)
	}
	if got, ok := g.CloudEventsExtension("group"); !ok || got != groupUID {
		t.Errorf("Group.CloudEventsExtension(group) = (%q, %v), want (%q, true)", got, ok, groupUID)
	}
	if _, ok := g.CloudEventsExtension("unknown"); ok {
		t.Error("Group.CloudEventsExtension(unknown) returned true")
	}

	del := &DeleteGroupRequest{Uid: groupUID}
	if got := del.CloudEventsSubject(); got != groupUID {
		t.Errorf("DeleteGroupRequest.CloudEventsSubject() = %q, want %q", got, groupUID)
	}
	if got, ok := del.CloudEventsExtension("group"); !ok || got != uidp.Parent(groupUID) {
		t.Errorf("DeleteGroupRequest.CloudEventsExtension(group) = (%q, %v), want (%q, true)", got, ok, uidp.Parent(groupUID))
	}
}

func TestIdentitiesEventAnnotations(t *testing.T) {
	sd := File_chainguard_platform_iam_v2beta1_identities_proto.Services().ByName("IdentitiesService")
	if sd == nil {
		t.Fatal("IdentitiesService not found")
	}
	runAnnotationTests(t, sd, []annotationTest{
		{"CreateIdentity", "dev.chainguard.api.iam.identity.created.v1", []string{"group"}},
		{"UpdateIdentity", "dev.chainguard.api.iam.identity.updated.v1", []string{"group"}},
		{"DeleteIdentity", "dev.chainguard.api.iam.identity.deleted.v1", []string{"group"}},
	})
}

func TestIdentitiesEventInterfaces(t *testing.T) {
	identityUID := "abc123/def456"
	parentUID := uidp.Parent(identityUID)

	id := &Identity{Uid: identityUID}
	if got := id.CloudEventsSubject(); got != identityUID {
		t.Errorf("Identity.CloudEventsSubject() = %q, want %q", got, identityUID)
	}
	if got, ok := id.CloudEventsExtension("group"); !ok || got != parentUID {
		t.Errorf("Identity.CloudEventsExtension(group) = (%q, %v), want (%q, true)", got, ok, parentUID)
	}

	del := &DeleteIdentityRequest{Uid: identityUID}
	if got, ok := del.CloudEventsExtension("group"); !ok || got != parentUID {
		t.Errorf("DeleteIdentityRequest.CloudEventsExtension(group) = (%q, %v), want (%q, true)", got, ok, parentUID)
	}
}

func TestRoleBindingsEventAnnotations(t *testing.T) {
	sd := File_chainguard_platform_iam_v2beta1_role_bindings_proto.Services().ByName("RoleBindingsService")
	if sd == nil {
		t.Fatal("RoleBindingsService not found")
	}
	runAnnotationTests(t, sd, []annotationTest{
		{"CreateRoleBinding", "dev.chainguard.api.iam.rolebindings.created.v1", []string{"group"}},
		{"UpdateRoleBinding", "dev.chainguard.api.iam.rolebindings.updated.v1", []string{"group"}},
		{"DeleteRoleBinding", "dev.chainguard.api.iam.rolebindings.deleted.v1", []string{"group"}},
	})
}

func TestRoleBindingsEventInterfaces(t *testing.T) {
	rbUID := "abc123/def456"

	rb := &RoleBinding{Uid: rbUID}
	if got := rb.CloudEventsSubject(); got != rbUID {
		t.Errorf("RoleBinding.CloudEventsSubject() = %q, want %q", got, rbUID)
	}
	if got, ok := rb.CloudEventsExtension("group"); !ok || got != uidp.Parent(rbUID) {
		t.Errorf("RoleBinding.CloudEventsExtension(group) = (%q, %v), want (%q, true)", got, ok, uidp.Parent(rbUID))
	}

	del := &DeleteRoleBindingRequest{Uid: rbUID}
	if got, ok := del.CloudEventsExtension("group"); !ok || got != uidp.Parent(rbUID) {
		t.Errorf("DeleteRoleBindingRequest.CloudEventsExtension(group) = (%q, %v), want (%q, true)", got, ok, uidp.Parent(rbUID))
	}
}

func TestIdentityProvidersEventAnnotations(t *testing.T) {
	sd := File_chainguard_platform_iam_v2beta1_identity_providers_proto.Services().ByName("IdentityProvidersService")
	if sd == nil {
		t.Fatal("IdentityProvidersService not found")
	}
	runAnnotationTests(t, sd, []annotationTest{
		{"CreateIdentityProvider", "dev.chainguard.api.iam.identity_providers.created.v1", []string{"group"}},
	})
}

func TestIdentityProvidersEventInterfaces(t *testing.T) {
	idpUID := "abc123/def456"

	idp := &IdentityProvider{
		Uid:  idpUID,
		Name: "test-idp",
		Configuration: &IdentityProvider_Oidc{
			Oidc: &IdentityProvider_OIDC{
				Issuer:       "https://accounts.google.com",
				ClientId:     "client-id",
				ClientSecret: "super-secret",
			},
		},
	}
	if got := idp.CloudEventsSubject(); got != idpUID {
		t.Errorf("IdentityProvider.CloudEventsSubject() = %q, want %q", got, idpUID)
	}
	if got, ok := idp.CloudEventsExtension("group"); !ok || got != uidp.Parent(idpUID) {
		t.Errorf("IdentityProvider.CloudEventsExtension(group) = (%q, %v), want (%q, true)", got, ok, uidp.Parent(idpUID))
	}

	// Verify redaction strips the client secret.
	redacted := idp.CloudEventsRedact().(*IdentityProvider)
	if redacted.GetName() != "test-idp" {
		t.Errorf("redacted name = %q, want %q", redacted.GetName(), "test-idp")
	}
	oidc := redacted.GetOidc()
	if oidc.GetClientSecret() != "" {
		t.Error("redacted IdentityProvider still contains client secret")
	}
	if oidc.GetIssuer() != "https://accounts.google.com" {
		t.Errorf("redacted issuer = %q, want preserved", oidc.GetIssuer())
	}
	if oidc.GetClientId() != "client-id" {
		t.Errorf("redacted client_id = %q, want preserved", oidc.GetClientId())
	}
}

func TestAccountAssociationsEventAnnotations(t *testing.T) {
	sd := File_chainguard_platform_iam_v2beta1_account_associations_proto.Services().ByName("AccountAssociationsService")
	if sd == nil {
		t.Fatal("AccountAssociationsService not found")
	}
	runAnnotationTests(t, sd, []annotationTest{
		{"CreateAccountAssociation", "dev.chainguard.api.iam.account_associations.created.v1", []string{"group"}},
		{"UpdateAccountAssociation", "dev.chainguard.api.iam.account_associations.updated.v1", []string{"group"}},
		{"DeleteAccountAssociation", "dev.chainguard.api.iam.account_associations.deleted.v1", []string{"group"}},
	})
}

func TestAccountAssociationsEventInterfaces(t *testing.T) {
	// AccountAssociation UID is the group UIDP itself.
	groupUID := "abc123"

	aa := &AccountAssociation{Uid: groupUID}
	if got := aa.CloudEventsSubject(); got != groupUID {
		t.Errorf("AccountAssociation.CloudEventsSubject() = %q, want %q", got, groupUID)
	}
	if got, ok := aa.CloudEventsExtension("group"); !ok || got != groupUID {
		t.Errorf("AccountAssociation.CloudEventsExtension(group) = (%q, %v), want (%q, true)", got, ok, groupUID)
	}

	del := &DeleteAccountAssociationRequest{Uid: groupUID}
	if got := del.CloudEventsSubject(); got != groupUID {
		t.Errorf("DeleteAccountAssociationRequest.CloudEventsSubject() = %q, want %q", got, groupUID)
	}
	if got, ok := del.CloudEventsExtension("group"); !ok || got != groupUID {
		t.Errorf("DeleteAccountAssociationRequest.CloudEventsExtension(group) = (%q, %v), want (%q, true)", got, ok, groupUID)
	}
	if redacted := del.CloudEventsRedact(); redacted != nil {
		t.Errorf("DeleteAccountAssociationRequest.CloudEventsRedact() = %v, want nil", redacted)
	}
}

func TestGroupInvitesEventAnnotations(t *testing.T) {
	sd := File_chainguard_platform_iam_v2beta1_group_invites_proto.Services().ByName("GroupInvitesService")
	if sd == nil {
		t.Fatal("GroupInvitesService not found")
	}
	runAnnotationTests(t, sd, []annotationTest{
		{"CreateGroupInvite", "dev.chainguard.api.iam.group_invite.created.v1", []string{"group"}},
		{"DeleteGroupInvite", "dev.chainguard.api.iam.group_invite.deleted.v1", []string{"group"}},
	})
}

func TestGroupInvitesEventInterfaces(t *testing.T) {
	inviteUID := "abc123/def456"
	parentUID := uidp.Parent(inviteUID)

	gi := &GroupInvite{
		Uid:     inviteUID,
		RoleUid: "role-uid",
		Email:   "test@example.com",
		Code:    "secret-code",
		KeyId:   "key-123",
	}
	if got := gi.CloudEventsSubject(); got != inviteUID {
		t.Errorf("GroupInvite.CloudEventsSubject() = %q, want %q", got, inviteUID)
	}
	if got, ok := gi.CloudEventsExtension("group"); !ok || got != parentUID {
		t.Errorf("GroupInvite.CloudEventsExtension(group) = (%q, %v), want (%q, true)", got, ok, parentUID)
	}
	if _, ok := gi.CloudEventsExtension("unknown"); ok {
		t.Error("GroupInvite.CloudEventsExtension(unknown) returned true")
	}

	// Verify redaction keeps only uid and expiration_time (matching v1 behavior).
	redacted := gi.CloudEventsRedact().(*GroupInvite)
	if redacted.GetUid() != inviteUID {
		t.Errorf("redacted uid = %q, want %q", redacted.GetUid(), inviteUID)
	}
	if redacted.GetCode() != "" {
		t.Error("redacted GroupInvite still contains code")
	}
	if redacted.GetKeyId() != "" {
		t.Error("redacted GroupInvite still contains key_id")
	}
	if redacted.GetRoleUid() != "" {
		t.Error("redacted GroupInvite still contains role_uid")
	}
	if redacted.GetEmail() != "" {
		t.Error("redacted GroupInvite still contains email")
	}

	del := &DeleteGroupInviteRequest{Uid: inviteUID}
	if got := del.CloudEventsSubject(); got != inviteUID {
		t.Errorf("DeleteGroupInviteRequest.CloudEventsSubject() = %q, want %q", got, inviteUID)
	}
	if got, ok := del.CloudEventsExtension("group"); !ok || got != parentUID {
		t.Errorf("DeleteGroupInviteRequest.CloudEventsExtension(group) = (%q, %v), want (%q, true)", got, ok, parentUID)
	}
}

// TestReadOnlyMethodsHaveNoEvents verifies List/Get methods don't have event annotations.
func TestReadOnlyMethodsHaveNoEvents(t *testing.T) {
	readOnlyMethods := []struct {
		file    protoreflect.FileDescriptor
		service string
		method  string
	}{
		{File_chainguard_platform_iam_v2beta1_group_invites_proto, "GroupInvitesService", "GetGroupInvite"},
		{File_chainguard_platform_iam_v2beta1_group_invites_proto, "GroupInvitesService", "ListGroupInvites"},
		{File_chainguard_platform_iam_v2beta1_groups_proto, "GroupsService", "GetGroup"},
		{File_chainguard_platform_iam_v2beta1_groups_proto, "GroupsService", "ListGroups"},
		{File_chainguard_platform_iam_v2beta1_identities_proto, "IdentitiesService", "GetIdentity"},
		{File_chainguard_platform_iam_v2beta1_identities_proto, "IdentitiesService", "ListIdentities"},
		{File_chainguard_platform_iam_v2beta1_role_bindings_proto, "RoleBindingsService", "GetRoleBinding"},
		{File_chainguard_platform_iam_v2beta1_role_bindings_proto, "RoleBindingsService", "ListRoleBindings"},
		{File_chainguard_platform_iam_v2beta1_identity_providers_proto, "IdentityProvidersService", "GetIdentityProvider"},
		{File_chainguard_platform_iam_v2beta1_identity_providers_proto, "IdentityProvidersService", "ListIdentityProviders"},
		{File_chainguard_platform_iam_v2beta1_account_associations_proto, "AccountAssociationsService", "GetAccountAssociation"},
		{File_chainguard_platform_iam_v2beta1_account_associations_proto, "AccountAssociationsService", "ListAccountAssociations"},
	}

	for _, tt := range readOnlyMethods {
		t.Run(tt.service+"/"+tt.method, func(t *testing.T) {
			sd := tt.file.Services().ByName(protoreflect.Name(tt.service))
			if sd == nil {
				t.Fatalf("service %s not found", tt.service)
			}
			md := sd.Methods().ByName(protoreflect.Name(tt.method))
			if md == nil {
				t.Fatalf("method %s not found", tt.method)
			}
			opts, ok := md.Options().(*descriptorpb.MethodOptions)
			if !ok || opts == nil {
				return // No options means no events — correct.
			}
			ext, ok := proto.GetExtension(opts, cgannotations.E_Events).(*cgannotations.EventAttributes)
			if ok && ext != nil {
				t.Errorf("read-only method %s.%s should NOT have event annotations, but has type=%q", tt.service, tt.method, ext.GetType())
			}
		})
	}
}
