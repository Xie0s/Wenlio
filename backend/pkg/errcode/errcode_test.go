package errcode

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	appErr := New(400, 400001, "bad request")
	if appErr.HTTPCode != 400 {
		t.Errorf("HTTPCode = %d, want 400", appErr.HTTPCode)
	}
	if appErr.Code != 400001 {
		t.Errorf("Code = %d, want 400001", appErr.Code)
	}
	if appErr.Message != "bad request" {
		t.Errorf("Message = %q, want %q", appErr.Message, "bad request")
	}
	if appErr.Internal != nil {
		t.Errorf("Internal = %v, want nil", appErr.Internal)
	}
}

func TestAppErrorError(t *testing.T) {
	appErr := New(500, 500001, "服务器内部错误")
	got := appErr.Error()
	if got != "服务器内部错误" {
		t.Errorf("Error() = %q, want %q", got, "服务器内部错误")
	}
}

func TestAppErrorWrap(t *testing.T) {
	original := New(500, 500001, "server error")
	inner := errors.New("database connection failed")
	wrapped := original.Wrap(inner)

	if wrapped.HTTPCode != original.HTTPCode {
		t.Errorf("Wrapped HTTPCode = %d, want %d", wrapped.HTTPCode, original.HTTPCode)
	}
	if wrapped.Code != original.Code {
		t.Errorf("Wrapped Code = %d, want %d", wrapped.Code, original.Code)
	}
	if wrapped.Message != original.Message {
		t.Errorf("Wrapped Message = %q, want %q", wrapped.Message, original.Message)
	}
	if wrapped.Internal != inner {
		t.Errorf("Wrapped Internal = %v, want %v", wrapped.Internal, inner)
	}
}

func TestAppErrorWrapPreservesOriginal(t *testing.T) {
	original := New(404, 404001, "not found")
	inner := errors.New("mongo: no documents")
	wrapped := original.Wrap(inner)

	// Original should remain unchanged
	if original.Internal != nil {
		t.Error("Wrap() should not modify original error's Internal field")
	}
	// Wrapped is a new instance
	if wrapped == original {
		t.Error("Wrap() should return a new AppError instance")
	}
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      *AppError
		httpCode int
		code     int
	}{
		{"ErrInvalidParam", ErrInvalidParam, 400, 400001},
		{"ErrUnauthorized", ErrUnauthorized, 401, 401001},
		{"ErrTokenExpired", ErrTokenExpired, 401, 401002},
		{"ErrForbidden", ErrForbidden, 403, 403001},
		{"ErrResourceNotFound", ErrResourceNotFound, 404, 404001},
		{"ErrInternalServer", ErrInternalServer, 500, 500001},
		{"ErrDatabase", ErrDatabase, 500, 500002},
		{"ErrLoginFailed", ErrLoginFailed, 401, 401101},
		{"ErrAccountLocked", ErrAccountLocked, 401, 401102},
		{"ErrPasswordTooWeak", ErrPasswordTooWeak, 400, 400104},
		{"ErrTenantNotFound", ErrTenantNotFound, 404, 404201},
		{"ErrTenantIDExists", ErrTenantIDExists, 409, 409202},
		{"ErrThemeNotFound", ErrThemeNotFound, 404, 404401},
		{"ErrPageNotFound", ErrPageNotFound, 404, 404601},
		{"ErrPageSlugExists", ErrPageSlugExists, 409, 409602},
		{"ErrCommentNotFound", ErrCommentNotFound, 404, 404701},
		{"ErrUserNotFound", ErrUserNotFound, 404, 404801},
		{"ErrUsernameExists", ErrUsernameExists, 409, 409802},
		{"ErrUploadFileTooLarge", ErrUploadFileTooLarge, 400, 400901},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.HTTPCode != tt.httpCode {
				t.Errorf("%s.HTTPCode = %d, want %d", tt.name, tt.err.HTTPCode, tt.httpCode)
			}
			if tt.err.Code != tt.code {
				t.Errorf("%s.Code = %d, want %d", tt.name, tt.err.Code, tt.code)
			}
			if tt.err.Message == "" {
				t.Errorf("%s.Message is empty", tt.name)
			}
			if tt.err.Error() != tt.err.Message {
				t.Errorf("%s.Error() = %q, want %q", tt.name, tt.err.Error(), tt.err.Message)
			}
		})
	}
}

func TestAppErrorImplementsError(t *testing.T) {
	var err error = New(400, 400001, "test error")
	if err.Error() != "test error" {
		t.Errorf("AppError does not properly implement error interface")
	}
}
