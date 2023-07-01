package product_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jasmanchik/garage-sale/internal/platform/database/databasetest"
	"github.com/jasmanchik/garage-sale/internal/product"
	"github.com/jasmanchik/garage-sale/internal/schema"
)

func TestProducts(t *testing.T) {
	db, cleanup := databasetest.Setup(t)
	defer cleanup()

	ctx := context.Background()

	np := product.NewProduct{
		Name:     "Comic Book",
		Cost:     10,
		Quantity: 55,
	}

	d := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)

	p, err := product.Create(ctx, db, &np, d)
	if err != nil {
		t.Fatalf("can not create product: %+v", err)
	}

	saved, err := product.Retrieve(ctx, db, p.ID)
	if err != nil {
		t.Fatalf("can not retrieve product: %+v", err)
	}

	if str := cmp.Diff(p, saved); str != "" {
		t.Fatalf("saved product does not match: %v != %v", p, saved)
	}

	pl, err := product.List(ctx, db)
	if err != nil {
		t.Fatalf("can not get list of products: %+v", err)
	}
	npFound := false
	for _, pil := range pl {
		if str2 := cmp.Diff(p, &pil); str2 == "" {
			npFound = true
			break
		}
	}
	if !npFound {
		t.Fatalf("can not find new product in list of products")
	}
}

func TestProductList(t *testing.T) {
	db, cleanup := databasetest.Setup(t)
	defer cleanup()
	ctx := context.Background()

	if err := schema.Seed(db); err != nil {
		t.Fatalf("can not seed database: %+v", err)
	}

	pl, err := product.List(ctx, db)
	if err != nil {
		t.Fatalf("can not get list of products: %+v", err)
	}
	if exp, got := 3, len(pl); exp != got {
		t.Fatalf("expected 3 products, got %d", len(pl))
	}
}
