package lantern_test

import (
	"testing"

	"spiritgen/internal/lantern"
)

func TestValidateAddress(t *testing.T) {
	tests := []struct {
		input    string
		wantErrs int
	}{
		{"서울특별시 강남구", 0},
		{" 어딘가 ", 0},
		{"", 1},
		{"   ", 1},
		{"\t\n", 1},
	}
	for _, tt := range tests {
		errs := lantern.ValidateAddress(tt.input)
		if len(errs) != tt.wantErrs {
			t.Errorf("ValidateAddress(%q) = %d errors, want %d", tt.input, len(errs), tt.wantErrs)
		}
	}
}

func TestValidatePerson(t *testing.T) {
	tests := []struct {
		name     string
		relation string
		wantErrs int
	}{
		{"홍길동", "부친", 0},
		{"홍길동", "", 1},
		{"", "부친", 1},
		{"", "", 2},
		{"  ", "  ", 2},
	}
	for _, tt := range tests {
		errs := lantern.ValidatePerson(tt.name, tt.relation)
		if len(errs) != tt.wantErrs {
			t.Errorf("ValidatePerson(%q, %q) = %d errors, want %d",
				tt.name, tt.relation, len(errs), tt.wantErrs)
		}
	}
}
