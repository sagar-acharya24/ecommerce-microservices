package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Product struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

type ProductClient interface {
	GetProduct(productID uint) (*Product, error)
}

type productClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewProductClient(baseURL string) ProductClient {
	return &productClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *productClient) GetProduct(productID uint) (*Product, error) {
	url := fmt.Sprintf(
		"%s/api/v1/products/%d",
		c.baseURL,
		productID,
	)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to call product service: %w",
			err,
		)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("product not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"product service returned status %d",
			resp.StatusCode,
		)
	}

	var product Product

	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, fmt.Errorf(
			"failed to decode product response: %w",
			err,
		)
	}

	return &product, nil
}
