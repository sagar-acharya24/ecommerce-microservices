package client

import (
	"testing"
)

func TestProductClientGetProduct(t *testing.T) {
	productClient := NewProductClient("http://localhost:8082")

	product, err := productClient.GetProduct(35)
	if err != nil {
		t.Fatalf("failed to get product: %v", err)
	}

	if product.ID != 35 {
		t.Fatalf(
			"expected product ID 34, got %d",
			product.ID,
		)
	}

	t.Logf(
		"Product: ID=%d Name=%s Price=%.2f Stock=%d",
		product.ID,
		product.Name,
		product.Price,
		product.Stock,
	)
}
