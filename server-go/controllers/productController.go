// controllers/productController.go
package controllers

import (
	"net/http"
	"strconv"

	"github.com/LaansDole/fullstack-ecommerce-warehouse/server-go/models"
	"github.com/gin-gonic/gin"
)

type Product struct {
	ID                 int     `json:"id"`
	Image              string  `json:"image"`
	Title              string  `json:"title"`
	ProductDescription string  `json:"product_description"`
	Category           string  `json:"category"`
	Price              float64 `json:"price"`
	Width              int     `json:"width"`
	Length             int     `json:"length"`
	Height             int     `json:"height"`
	Seller             string  `json:"seller"`
}

func queryProducts(query string, args ...interface{}) ([]Product, error) {
	rows, err := models.DBAdmin.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Image, &product.Title, &product.ProductDescription, &product.Category, &product.Price, &product.Width, &product.Length, &product.Height, &product.Seller); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func GetAllProducts(c *gin.Context) {
	products, err := queryProducts("SELECT * FROM product")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func GetAllProductsASC(c *gin.Context) {
	products, err := queryProducts("SELECT * FROM product ORDER BY price ASC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func GetAllProductsDSC(c *gin.Context) {
	products, err := queryProducts("SELECT * FROM product ORDER BY price DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func GetProductById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	products, err := queryProducts("SELECT * FROM product WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, products[0])
}

func GetProductByTitle(c *gin.Context) {
	title := c.Param("title")

	products, err := queryProducts("SELECT * FROM product WHERE title LIKE CONCAT('%', ?, '%')", title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, products)
}
