package controllers

import (
	"net/http"

	"github.com/LaansDole/fullstack-ecommerce-warehouse/server-go/models"
	"github.com/gin-gonic/gin"
)

type Stockpile struct {
	ProductID   int `form:"product_id" json:"product_id"`
	WarehouseID int `form:"warehouse_id" json:"warehouse_id"`
	Quantity    int `form:"quantity" json:"quantity"`
}

func GetAllStockpile(c *gin.Context) {
	rows, err := models.DBAdmin.Query("SELECT * FROM stockpile")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	defer rows.Close()

	var results []Stockpile // assuming you have a Stockpile struct
	for rows.Next() {
		var stockpile Stockpile
		if err := rows.Scan(&stockpile.ProductID, &stockpile.WarehouseID, &stockpile.Quantity); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		results = append(results, stockpile)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, results)
}
