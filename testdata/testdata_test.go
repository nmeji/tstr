package testdata_test

import (
	"reflect"
	"testing"

	"github.com/nmeji/tstr/testdata"
)

func TestNew_CSV_NoError(t *testing.T) {
	_, err := testdata.New("testdata/1.csv")

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	if t == nil {
		t.Error("expected non-nil TestData, got nil")
	}
}

func TestUnmarshal_CSV_2DSlices(t *testing.T) {
	testdata, _ := testdata.New("testdata/1.csv")

	result := [][]string{}
	err := testdata.Unmarshal(&result)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	expected := [][]string{
		[]string{"color", "hex"},
		[]string{"RED", "#FF0000"},
		[]string{"GREEN", "#00FF00"},
		[]string{"BLUE", "#0000FF"},
	}
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestUnmarshal_CSV_SliceOfStructs(t *testing.T) {
	testdata, _ := testdata.New("testdata/1.csv")

	result := []struct {
		Color string
		Hex   string
	}{}
	err := testdata.Unmarshal(&result)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	if len(result) != 3 {
		t.Errorf("expected 3, got %d", len(result))
	}
}

func TestUnmarshal_TOML_NoError(t *testing.T) {
	testdata, _ := testdata.New("testdata/1.toml")

	type env struct {
		ID        string
		UserEmail string `toml:"user_email"`
		Vertical  string
	}
	var result struct{ QA env }
	err := testdata.Unmarshal(&result)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	if result.QA.ID != "180b3347-210d-443f-92dc-369f51752188" {
		t.Errorf("expected 180b3347-210d-443f-92dc-369f51752188, got %s", result.QA.ID)
	}

	if result.QA.UserEmail != "sxaykvtiwco@mailinator.com" {
		t.Errorf("expected sxaykvtiwco@mailinator.com, got %s", result.QA.UserEmail)
	}

	if result.QA.Vertical != "mobile" {
		t.Errorf("expected mobile, got %s", result.QA.Vertical)
	}
}

func TestUnmarshal_JSON_NoError(t *testing.T) {
	testdata, _ := testdata.New("testdata/1.json")

	type Address struct {
		Street string `json:"streetAddress"`
		City   string
	}
	type Profile struct {
		FirstName string
		LastName  string
		Age       int
		City      Address `json:"address"`
	}

	var data Profile
	if err := testdata.Unmarshal(&data); err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	expected := Profile{
		FirstName: "John",
		LastName:  "doe",
		Age:       26,
		City: Address{
			Street: "naist street",
			City:   "Nara",
		},
	}
	if !reflect.DeepEqual(expected, data) {
		t.Errorf("expected %v, got %v", expected, data)
	}
}
