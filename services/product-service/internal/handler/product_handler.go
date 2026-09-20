package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/sagar-acharya24/ecommerce-microservices/services/product-service/internal/dto"
	"github.com/sagar-acharya24/ecommerce-microservices/services/product-service/internal/model"
	"github.com/sagar-acharya24/ecommerce-microservices/services/product-service/internal/service"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		service: productService,
	}
}

// CreateProduct handles POST /api/v1/products
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var request dto.CreateProductRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	product := &model.Product{
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Stock:       request.Stock,
	}

	if err := h.service.CreateProduct(product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetProduct handles GET /api/v1/products/:id
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product ID",
		})
		return
	}

	product, err := h.service.GetProductByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "product not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get product",
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetProducts handles GET /api/v1/products
func (h *ProductHandler) GetProducts(c *gin.Context) {
	products, err := h.service.GetProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get products",
		})
		return
	}

	c.JSON(http.StatusOK, products)
}

// UpdateProduct handles PUT /api/v1/products/:id
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product ID",
		})
		return
	}

	var request dto.UpdateProductRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	product := &model.Product{
		ID:          uint(id),
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Stock:       request.Stock,
	}

	if err := h.service.UpdateProduct(product); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "product not found",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	updatedProduct, err := h.service.GetProductByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch updated product",
		})
		return
	}

	c.JSON(http.StatusOK, updatedProduct)
}

// DeleteProduct handles DELETE /api/v1/products/:id
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product ID",
		})
		return
	}

	if err := h.service.DeleteProduct(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "product not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete product",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
