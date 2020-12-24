package oauth2_test

import (
	"testing"

	"github.com/go-oauth2/oauth2/v4"
)

func TestValidatePlain(t *testing.T) {
	cc := oauth2.CodeChallengePlain
	if !cc.Validate("plaintest", "plaintest") {
		t.Fatal("not valid")
	}
}

func TestValidateS256(t *testing.T) {
	cc := oauth2.CodeChallengeS256
	if !cc.Validate("W6YWc_4yHwYN-cGDgGmOMHF3l7KDy7VcRjf7q2FVF-o=", "s256test") {
		t.Fatal("not valid")
	}
}

func TestValidateS256NoPadding(t *testing.T) {
	cc := oauth2.CodeChallengeS256
	if !cc.Validate("W6YWc_4yHwYN-cGDgGmOMHF3l7KDy7VcRjf7q2FVF-o", "s256test") {
		t.Fatal("not valid")
	}
}
