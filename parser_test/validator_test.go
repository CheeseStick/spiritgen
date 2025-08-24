package parser_test

import (
	"spiritgen/internal/parser"
	"testing"
)

func TestValidatePresenter(t *testing.T) {
	tests := []struct {
		input    string
		wantErrs int
	}{
		{"누군가", 0},
		{"누 군가", 0},
		{"   ", 1},
		{"", 1},
		{"\t", 1},
		{"\t\t", 1},
		{"\t\n", 1},
		{" \t", 1},
	}

	for _, tt := range tests {
		errs := parser.ValidatePresenter(tt.input)
		if len(errs) != tt.wantErrs {
			t.Errorf("ValidatePresenter(%q) = %d errors, want %d", tt.input, len(errs), tt.wantErrs)
		}
	}
}

func TestValidateDeceasedInput(t *testing.T) {
	tests := []struct {
		name     string
		dharma   string
		clan     string
		relation string
		wantErrs int
	}{
		// 정상 - 모든 데이터
		{"누군가", "법명", "본관", "관계", 0},
		{"누군가", "법명", "", "관계", 0},

		// 정상 - 이름 사용
		{"누군가", "", "본관", "관계", 0},
		{"누군가", "", "", "관계", 0},

		// 정상 - 법명 사용
		{"", "법명", "본관", "관계", 0},
		{"", "법명", "", "관계", 0},

		// 오류 - 이름, 법명이 비었을 때
		{"", "", "본관", "관계", 1},
		{"", "", "", "관계", 1},

		// 오류 - 관계 필드가 비어있을 때
		{"누군가", "", "본관", "", 1},
		{"누군가", "", "", "", 1},
		{"", "법명", "본관", "", 1},
		{"", "법명", "", "", 1},

		// 오류 - 모두 비었을 때
		{"", "", "", "", 2},
	}

	for _, tt := range tests {
		errs := parser.ValidateDeceasedInput(tt.name, tt.dharma, tt.clan, tt.relation)
		if len(errs) != tt.wantErrs {
			t.Errorf("ValidateDeceasedInput(%q, %q, %q, %q) = %d errors, want %d",
				tt.name, tt.dharma, tt.clan, tt.relation, len(errs), tt.wantErrs)
		}
	}
}
