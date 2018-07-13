package env_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/nmeji/tstr/env"
)

func TestEnvStrVal_Present(t *testing.T) {
	test := "test"

	os.Setenv("envvar", test)

	val := env.StrVal("envvar")

	if val != test {
		t.Errorf("expected %s, got %s", test, val)
	}
}

func TestEnvStrVal_Nil(t *testing.T) {
	val := env.StrVal("notexistingvar")

	if len(val) > 0 {
		t.Errorf("expected \"\" (empty string), got %s", val)
	}
}

func TestEnvIntVal_Error(t *testing.T) {
	os.Setenv("number", "number")

	val, err := env.IntVal("number")

	if err == nil {
		t.Error("expected error, got nil")
	}

	if val != 0 {
		t.Errorf("expected 0, got %d", val)
	}
}

func TestEnvIntVal_OK(t *testing.T) {
	test := 0xFF

	os.Setenv("number", fmt.Sprintf("%d", test))

	val, err := env.IntVal("number")

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	if val != test {
		t.Errorf("expected %d, got %d", test, val)
	}
}

func TestEnvIntVal_Nil(t *testing.T) {
	val, err := env.IntVal("notexistingvar")

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	if val != 0 {
		t.Errorf("expected 0, got %d", val)
	}
}

func TestEnvBoolVal_Error(t *testing.T) {
	os.Setenv("flag", "flag")

	val, err := env.BoolVal("flag")

	if err == nil {
		t.Error("expected error, got nil")
	}

	if val == true {
		t.Error("expected false, got true")
	}
}

func TestEnvBoolVal_OK(t *testing.T) {
	test := true

	os.Setenv("flag", strings.ToUpper(fmt.Sprintf("%t", test)))

	val, err := env.BoolVal("flag")

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	if val != test {
		t.Errorf("expected %t, got %t", test, val)
	}
}

func TestEnvBoolVal_Nil(t *testing.T) {
	val, err := env.BoolVal("notexistingvar")

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	if val == true {
		t.Error("expected false, got true")
	}
}
