package utils

import "testing"

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		name  string
		title string
		want  string
	}{
		{"lowercase ascii", "hello world", "hello-world"},
		{"mixed case", "Hello World", "hello-world"},
		{"special characters", "Hello, World! #123", "hello-world-123"},
		{"multiple spaces", "hello   world", "hello-world"},
		{"leading trailing special", "---hello---", "hello"},
		{"numbers", "chapter 1 section 2", "chapter-1-section-2"},
		{"already a slug", "my-slug", "my-slug"},
		{"empty string", "", ""},
		{"only special chars", "!@#$%", ""},
		{"unicode characters", "你好世界", ""},
		{"mixed unicode and ascii", "hello你好world", "hello-world"},
		{"consecutive dashes", "hello---world", "hello-world"},
		{"single word", "hello", "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateSlug(tt.title)
			if got != tt.want {
				t.Errorf("GenerateSlug(%q) = %q, want %q", tt.title, got, tt.want)
			}
		})
	}
}

func TestValidateSlug(t *testing.T) {
	tests := []struct {
		name string
		slug string
		want bool
	}{
		{"valid simple", "hello-world", true},
		{"valid numbers", "chapter-1", true},
		{"valid all numbers", "12345", true},
		{"valid min length", "ab", true},
		{"too short single char", "a", false},
		{"empty string", "", false},
		{"uppercase", "Hello-World", false},
		{"contains space", "hello world", false},
		{"starts with dash", "-hello", false},
		{"ends with dash", "hello-", false},
		{"consecutive dashes", "hello--world", true},
		{"special characters", "hello_world", false},
		{"single valid char pair", "a1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateSlug(tt.slug)
			if got != tt.want {
				t.Errorf("ValidateSlug(%q) = %v, want %v", tt.slug, got, tt.want)
			}
		})
	}
}

func TestValidateSlugMaxLength(t *testing.T) {
	// 128 characters - valid
	slug128 := "a"
	for i := 1; i < 128; i++ {
		slug128 += "b"
	}
	if !ValidateSlug(slug128) {
		t.Errorf("ValidateSlug() should accept slug of length 128")
	}

	// 129 characters - invalid
	slug129 := slug128 + "c"
	if ValidateSlug(slug129) {
		t.Errorf("ValidateSlug() should reject slug of length 129")
	}
}
