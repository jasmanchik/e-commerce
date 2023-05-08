package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jasmanchik/garage-sale/cmd/sales-api/internal/handlers"
	"github.com/jasmanchik/garage-sale/internal/platform/database/databasetest"
	"github.com/jasmanchik/garage-sale/internal/schema"
)

func TestProducts(t *testing.T) {
	db, teardown := databasetest.Setup(t)
	defer teardown()

	if err := schema.Seed(db); err != nil {
		t.Fatalf("could not seed database: %v", err)
	}
	l := log.New(os.Stdout, "TEST : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	tests := ProductTests{
		app: handlers.Routes(l, db),
	}

	t.Run("List", tests.List)
	t.Run("ProductCRUD", tests.ProductCRUD)
}

type ProductTests struct {
	app http.Handler
}

func (pt *ProductTests) ProductCRUD(t *testing.T) {
	var created map[string]interface{}

	// Create
	{
		body := strings.NewReader(`{"name": "product0", "cost": 55, "quantity": 6}`)

		req := httptest.NewRequest(http.MethodPost, "/v1/products", body)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		pt.app.ServeHTTP(resp, req)
		if http.StatusCreated != resp.Code {
			t.Fatalf("post: expecting status code: %d, got: %d", http.StatusCreated, resp.Code)
		}

		if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
			t.Fatalf("can not decode post request: %v", err)
		}

		if created["id"] == "" || created["id"] == nil {
			t.Fatalf("can not find product_id in created product")
		}
		if created["date_created"] == "" || created["date_created"] == nil {
			t.Fatalf("can not find date_created in created product")
		}
		if created["date_updated"] == "" || created["date_updated"] == nil {
			t.Fatalf("can not find date_updated in created product")
		}

		want := map[string]interface{}{
			"id":           created["id"],
			"date_created": created["date_created"],
			"date_updated": created["date_updated"],
			"name":         "product0",
			"cost":         float64(55),
			"quantity":     float64(6),
		}

		if diff := cmp.Diff(created, want); diff != "" {
			t.Fatalf("Created object doesnt the same as i want to see: %s\n", diff)
		}
	}

	// Get by id
	{
		endpoint := fmt.Sprintf("/v1/products/%s", created["id"])
		req := httptest.NewRequest(http.MethodGet, endpoint, nil)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		pt.app.ServeHTTP(resp, req)

		if http.StatusOK != resp.Code {
			t.Fatalf("grtting item http status is not correct. Exp: %d, Got: %d", http.StatusOK, resp.Code)
		}

		var itemFromDB map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&itemFromDB); err != nil {
			t.Fatalf("can not to detch data from db: %v", err)
		}

		if diff := cmp.Diff(itemFromDB, created); diff != "" {
			t.Fatalf("item from DB is not the same as I created: %s", diff)
		}
	}

	// Delete by id
	{
		endpoint := fmt.Sprintf("/v1/products/%s", created["id"])
		req := httptest.NewRequest(http.MethodDelete, endpoint, nil)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		pt.app.ServeHTTP(resp, req)

		if http.StatusOK != resp.Code {
			t.Fatalf("deleting item http status is not correct. Exp: %d, Got: %d", http.StatusOK, resp.Code)
		}

		req = httptest.NewRequest(http.MethodGet, endpoint, nil)
		resp = httptest.NewRecorder()
		pt.app.ServeHTTP(resp, req)
		if resp.Code != http.StatusNotFound {
			t.Fatalf("getting item after delete http status is not correct. Exp: %d, Got: %d", http.StatusOK, resp.Code)
		}
	}
}

func (pt *ProductTests) List(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v1/products", nil)
	resp := httptest.NewRecorder()

	pt.app.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("getting: expexting status code %v, got %v", http.StatusOK, resp.Code)
	}

	var list []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		t.Fatalf("Decode error: %v", err)
	}

	want := []map[string]interface{}{
		{
			"id":           "4eabde0b-3331-4927-9091-701f829b0262",
			"name":         "Comic Books",
			"cost":         float64(50),
			"quantity":     float64(42),
			"date_created": "2019-01-01T00:00:01Z",
			"date_updated": "2019-01-01T00:00:01Z",
		},
		{
			"id":           "b0de7d30-42e4-4ee2-8f1a-a382be080c32",
			"name":         "McDonalds Toys",
			"cost":         float64(75),
			"quantity":     float64(120),
			"date_created": "2019-01-01T00:00:02Z",
			"date_updated": "2019-01-01T00:00:02Z",
		},
		{
			"id":           "2378aa21-db61-4d71-b7c8-3ee573df000a",
			"name":         "Big Wheels",
			"cost":         float64(500),
			"quantity":     float64(2),
			"date_created": "2019-01-01T00:00:03Z",
			"date_updated": "2019-01-01T00:00:03Z",
		},
	}

	if diff := cmp.Diff(want, list); diff != "" {
		t.Fatalf("Response did not match expected. Diff:\n %s", diff)
	}

}
