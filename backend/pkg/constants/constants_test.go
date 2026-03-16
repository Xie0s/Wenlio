package constants

import "testing"

func TestRoleConstants(t *testing.T) {
	if RoleSuperAdmin != "super_admin" {
		t.Errorf("RoleSuperAdmin = %q, want %q", RoleSuperAdmin, "super_admin")
	}
	if RoleTenantAdmin != "tenant_admin" {
		t.Errorf("RoleTenantAdmin = %q, want %q", RoleTenantAdmin, "tenant_admin")
	}
}

func TestVersionStatusConstants(t *testing.T) {
	if VersionStatusDraft != "draft" {
		t.Errorf("VersionStatusDraft = %q, want %q", VersionStatusDraft, "draft")
	}
	if VersionStatusPublished != "published" {
		t.Errorf("VersionStatusPublished = %q, want %q", VersionStatusPublished, "published")
	}
	if VersionStatusArchived != "archived" {
		t.Errorf("VersionStatusArchived = %q, want %q", VersionStatusArchived, "archived")
	}
}

func TestPageStatusConstants(t *testing.T) {
	if PageStatusDraft != "draft" {
		t.Errorf("PageStatusDraft = %q, want %q", PageStatusDraft, "draft")
	}
	if PageStatusPublished != "published" {
		t.Errorf("PageStatusPublished = %q, want %q", PageStatusPublished, "published")
	}
}

func TestCommentStatusConstants(t *testing.T) {
	if CommentStatusPending != "pending" {
		t.Errorf("CommentStatusPending = %q, want %q", CommentStatusPending, "pending")
	}
	if CommentStatusApproved != "approved" {
		t.Errorf("CommentStatusApproved = %q, want %q", CommentStatusApproved, "approved")
	}
	if CommentStatusRejected != "rejected" {
		t.Errorf("CommentStatusRejected = %q, want %q", CommentStatusRejected, "rejected")
	}
}

func TestTenantStatusConstants(t *testing.T) {
	if TenantStatusActive != "active" {
		t.Errorf("TenantStatusActive = %q, want %q", TenantStatusActive, "active")
	}
	if TenantStatusSuspended != "suspended" {
		t.Errorf("TenantStatusSuspended = %q, want %q", TenantStatusSuspended, "suspended")
	}
	if TenantStatusDeleting != "deleting" {
		t.Errorf("TenantStatusDeleting = %q, want %q", TenantStatusDeleting, "deleting")
	}
}

func TestUserStatusConstants(t *testing.T) {
	if UserStatusActive != "active" {
		t.Errorf("UserStatusActive = %q, want %q", UserStatusActive, "active")
	}
	if UserStatusInactive != "inactive" {
		t.Errorf("UserStatusInactive = %q, want %q", UserStatusInactive, "inactive")
	}
}

func TestReservedTenantIDs(t *testing.T) {
	expectedReserved := []string{
		"admin", "api", "assets", "static",
		"health", "favicon.ico", "robots.txt", "sitemap",
	}

	for _, id := range expectedReserved {
		if _, ok := ReservedTenantIDs[id]; !ok {
			t.Errorf("ReservedTenantIDs missing expected entry %q", id)
		}
	}
}

func TestReservedTenantIDsNotContainNormal(t *testing.T) {
	normalIDs := []string{
		"mycompany", "docs", "blog", "team",
	}
	for _, id := range normalIDs {
		if _, ok := ReservedTenantIDs[id]; ok {
			t.Errorf("ReservedTenantIDs should not contain %q", id)
		}
	}
}
