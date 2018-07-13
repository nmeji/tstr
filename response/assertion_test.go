package response_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nmeji/tstr/request"
)

func TestAssertJSON_WithPath(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{
			"books": [ 
					{"id": 1, "price": 88.94},
					{"id": 2, "price": 97.10},
					{"id": 3, "price": 11.87}
				],
			"author": "niko"
		}`)
	}))
	defer ts.Close()

	checker, err := request.Get(ts.URL)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	books := []struct {
		ID    int
		Price float64
	}{
		{ID: 2, Price: 97.10},
		{ID: 3, Price: 11.87},
	}
	checker.
		ExpectStatus(200).
		ExpectBody.ToHaveInJson("$.books[1,2]", books).
		ExpectBody.ToHaveInJson("$.author", "niko").
		ExpectBody.ToHaveInJson("$.books[0].price", 88.94).
		MakeAssertion(t)
}
