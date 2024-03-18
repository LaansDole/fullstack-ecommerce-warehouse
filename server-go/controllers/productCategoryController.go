package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/LaansDole/fullstack-ecommerce-warehouse/server-go/models"
	"github.com/gin-gonic/gin"
)

type ProductCategory struct {
	CategoryName string  `form:"category_name" json:"category_name"`
	Parent       *string `form:"parent" json:"parent"`
}

func queryProductCategories(query string, args ...interface{}) ([]ProductCategory, error) {
	rows, err := models.DBSeller.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var categories []ProductCategory
	for rows.Next() {
		var category ProductCategory
		if err := rows.Scan(&category.CategoryName, &category.Parent); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows reported an error: %w", err)
	}

	return categories, nil
}

// GET category endpoint

func GetAllProductCategories(c *gin.Context) {
	categories, err := queryProductCategories("SELECT * FROM product_category")
	fmt.Printf("\nCategories: %v\n", categories)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func GetProductCategoryByName(c *gin.Context) {
	categoryName := c.Param("category_name")
	categories, err := queryProductCategories("SELECT * FROM product_category WHERE category_name = ?", categoryName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if len(categories) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Product category with name: %s not found", categoryName)})
		return
	}

	c.JSON(http.StatusOK, categories[0])
}

// CREATE category endpoint

func CreateProductCategory(c *gin.Context) {
	var category ProductCategory
	if err := c.ShouldBind(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error at bind": err.Error()})
		return
	}

	if category.Parent != nil && category.CategoryName == *category.Parent {
		c.JSON(http.StatusConflict, gin.H{"message": "Product Category and parent cannot have same name"})
		return
	}

	var query string
	var result sql.Result
	var err error
	if category.Parent == nil {
		query = "INSERT INTO product_category (category_name) VALUES (?)"
		result, err = models.DBAdmin.Exec(query, category.CategoryName)
	} else {
		parentCategory, err := models.GetProductCategoryByName(*category.Parent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Get category by name error": err.Error()})
			return
		}
		if parentCategory == nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Category %s not found", *category.Parent)})
			return
		}
		query = "INSERT INTO product_category (category_name, parent) VALUES (?, ?)"
		result, err = models.DBAdmin.Exec(query, category.CategoryName, *category.Parent)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Insert category error": err.Error()})
		return
	}

	fmt.Printf("\nResult: %v\n", result)

	c.JSON(http.StatusCreated, gin.H{
		"message":       fmt.Sprintf("Product category with name: %s created", category.CategoryName),
		"category_name": category.CategoryName,
		"parent":        category.Parent,
	})
}

// UPDATE category endpoint

func UpdateProductCategory(c *gin.Context) {
	categoryName := c.Param("category_name")
	var category ProductCategory
	if err := c.ShouldBind(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error at bind": err.Error()})
		return
	}

	if category.Parent != nil && categoryName == *category.Parent {
		c.JSON(http.StatusConflict, gin.H{"error": "Product Category and parent cannot have same name"})
		return
	}

	transaction, err := models.DBAdmin.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Transaction error": "Internal server error"})
		return
	}

	if category.Parent == nil {
		_, err = transaction.Exec("UPDATE product_category SET parent = NULL WHERE category_name = ?", categoryName)
	} else {
		_, err = transaction.Exec("UPDATE product_category SET parent = ? WHERE category_name = ?", *category.Parent, categoryName)
	}
	if err != nil {
		transaction.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	_, err = transaction.Exec("UPDATE product_category child JOIN product_category parent ON child.parent = parent.category_name SET child.parent = ? WHERE parent.category_name = ?", category.Parent, categoryName)
	if err != nil {
		transaction.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	err = transaction.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       fmt.Sprintf("Product category with name: %s updated", categoryName),
		"category_name": categoryName,
		"parent":        category.Parent,
	})
}

// DELETE category endpoint

func DeleteProductCategory(c *gin.Context) {
	categoryName := c.Param("category_name")

	tx, err := models.DBAdmin.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	_, err = tx.Exec("UPDATE product_category SET parent = NULL WHERE parent = ?", categoryName)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while updating product categories"})
		return
	}

	_, err = tx.Exec("DELETE FROM product_category WHERE category_name = ?", categoryName)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while deleting a product category"})
		return
	}

	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       fmt.Sprintf("Product category with name: %s deleted", categoryName),
		"category_name": categoryName,
	})
}
