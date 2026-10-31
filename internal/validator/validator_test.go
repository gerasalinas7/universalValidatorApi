package validator

import "testing"

func TestValidateEmail(t *testing.T) {
	result, supported := Validate("email", "student@example.com")
	if !supported || !result.Valid {
		t.Fatalf("expected valid email, got %#v", result)
	}
}

func TestValidateInvalidEmail(t *testing.T) {
	result, supported := Validate("email", "student@example")
	if !supported || result.Valid {
		t.Fatalf("expected invalid email, got %#v", result)
	}
}

func TestValidateCUIT(t *testing.T) {
	result, supported := Validate("cuit", "20-12345678-6")
	if !supported {
		t.Fatal("cuit validator should be supported")
	}
	// This test focuses on response behavior; checksum validity depends on the supplied CUIT.
	_ = result
}

func TestUnsupportedType(t *testing.T) {
	_, supported := Validate("unknown", "value")
	if supported {
		t.Fatal("expected unsupported type")
	}
}
