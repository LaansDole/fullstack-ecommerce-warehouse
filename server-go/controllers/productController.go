// controllers/productController.go
package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/LaansDole/fullstack-ecommerce-warehouse/server-go/models"
	"github.com/gin-gonic/gin"
)

type Product struct {
	ID                 int     `form:"id" json:"id"`
	Image              string  `form:"image" json:"image"`
	Title              string  `form:"title" json:"title"`
	ProductDescription string  `form:"product_description" json:"product_description"`
	Category           string  `form:"category" json:"category"`
	Price              float64 `form:"price" json:"price"`
	Width              int     `form:"width" json:"width"`
	Length             int     `form:"length" json:"length"`
	Height             int     `form:"height" json:"height"`
	Seller             string  `form:"seller" json:"seller"`
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

// GET product endpoints

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

// CREATE product endpoint

func CreateProduct(c *gin.Context) {
	var product Product
	if err := c.ShouldBind(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error at bind": err.Error()})
		return
	}

	fmt.Printf("\nProduct after binding: %v \n", product)

	product.Seller = c.MustGet("username").(string)
	product.Image = c.MustGet("imagePath").(string)

	fmt.Println("Product after setting image:", product)

	query := `INSERT INTO product (title, image, product_description, category, price, width, length, height, seller) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := models.DBAdmin.Exec(query, product.Title, product.Image, product.ProductDescription, product.Category, product.Price, product.Width, product.Length, product.Height, product.Seller)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"id":      id,
		"product": product,
	})
}
