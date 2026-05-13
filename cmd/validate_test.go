package cmd_test

import (
	"testing"

	"github.com/yourorg/envchain-cli/internal/validate"
)

// TestValidate_ChainName_Integration exercises the validate package as it
// would be called by CLI command handlers before touching the store.
func TestValidate_ChainName_Integration(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid simple", "myproject", false},
		{"valid hyphen", "my-project", false},
		{"empty", "", true},
		{"spaces", "my project", true},
		{"slash", "a/b", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validate.ChainName(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("ChainName(%q) error = %v, wantErr %v", tc.input, err, tc.wantErr)
			}
		})
	}
}

// TestValidate_Key_Integration exercises key validation as called by add/set
// command handlers.
func TestValidate_Key_Integration(t *testing.T) {
	cases := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{"valid upper", "DATABASE_URL", false},
		{"valid underscore start", "_INTERNAL", false},
		{"empty key", "", true},
		{"digit start", "1VAR", true},
		{"hyphen", "MY-VAR", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validate.Key(tc.key)
			if (err != nil) != tc.wantErr {
				t.Errorf("Key(%q) error = %v, wantErr %v", tc.key, err, tc.wantErr)
			}
		})
	}
}
