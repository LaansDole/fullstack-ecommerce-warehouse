package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/LaansDole/fullstack-ecommerce-warehouse/server-go/models"
	"github.com/gin-gonic/gin"
)

type InboundOrder struct {
	ID            int     `form:"id" json:"id"`
	Quantity      int     `form:"quantity" json:"quantity"`
	ProductID     int     `form:"product_id" json:"product_id"`
	CreatedDate   string  `form:"created_date" json:"created_date"`
	CreatedTime   string  `form:"created_time" json:"created_time"`
	FulfilledDate *string `form:"fulfilled_date" json:"fulfilled_date"`
	FulfilledTime *string `form:"fulfilled_time" json:"fulfilled_time"`
	Seller        string  `form:"seller" json:"seller"`
}

func queryInboundOrders(query string, args ...interface{}) ([]InboundOrder, error) {
	rows, err := models.DBSeller.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var orders []InboundOrder
	for rows.Next() {
		var order InboundOrder
		if err := rows.Scan(&order.ID, &order.Quantity, &order.ProductID, &order.CreatedDate, &order.CreatedTime, &order.FulfilledDate, &order.FulfilledTime, &order.Seller); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows reported an error: %w", err)
	}

	return orders, nil
}

// GET inbound orders endpoint

func GetAllInboundOrders(c *gin.Context) {
	orders, err := queryInboundOrders("SELECT * FROM inbound_order")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func GetInboundOrderById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid inbound order ID"})
		return
	}

	products, err := queryInboundOrders("SELECT * FROM inbound_order WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "inbound order not found"})
		return
	}

	c.JSON(http.StatusOK, products[0])
}

// CREATE inbound order endpoint

func CreateInboundOrder(c *gin.Context) {
	var order InboundOrder
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	seller := c.MustGet("username").(string)

	product, err := models.GetProduct(order.ProductID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Get product error": err.Error()})
		return
	}
	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	result, err := models.DBSeller.Exec(
		"INSERT INTO view_inbound_order_noid (quantity, product_id, created_date, created_time, seller) VALUES (?, ?, CURDATE(), CURTIME(), ?)",
		order.Quantity,
		order.ProductID,
		seller,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	order.ID = int(id)
	order.CreatedDate = time.Now().Format("2006-01-02")
	order.CreatedTime = time.Now().Format("15:04:05")
	order.Seller = seller

	c.JSON(http.StatusCreated, order)
}

// Update inbound order endpoint

func UpdateInboundOrder(c *gin.Context) {
	inboundOrderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid inbound order ID"})
		return
	}

	var inboundOrder InboundOrder
	if err := c.ShouldBindJSON(&inboundOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	seller := c.MustGet("username").(string)

	result, err := models.DBSeller.Exec(
		"UPDATE inbound_order SET quantity = ?, created_date = CURDATE(), created_time = CURTIME(), seller = ? WHERE id = ?",
		inboundOrder.Quantity,
		seller,
		inboundOrderID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("\n Row affected %v\n", rowsAffected)

	inboundOrder.CreatedDate = time.Now().Format("2006-01-02")
	inboundOrder.CreatedTime = time.Now().Format("15:04:05")
	inboundOrder.Seller = seller
	inboundOrder.ID = inboundOrderID

	c.JSON(http.StatusOK, inboundOrder)
}

// DELETE inbound order endpoint

func DeleteInboundOrder(c *gin.Context) {
	id := c.Param("id")

	_, err := models.DBSeller.Exec("DELETE FROM inbound_order WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while deleting an inbound order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Inbound order with ID: %s deleted", id),
		"id":      id,
	})
}

func FulfillInboundOrder(c *gin.Context) {
	id := c.Param("id")
	seller := c.MustGet("username").(string)

	_, err := models.DBSeller.Exec("CALL sp_fulfill_inbound_order(?, ?, @result)", id, seller)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var resultCode int
	err = models.DBSeller.QueryRow("SELECT @result as result").Scan(&resultCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	switch resultCode {
	case 0:
		row := models.DBSeller.QueryRow("SELECT * FROM inbound_order WHERE id = ?", id)
		var order InboundOrder
		err := row.Scan(&order.ID, &order.Quantity, &order.ProductID, &order.CreatedDate, &order.CreatedTime, &order.FulfilledDate, &order.FulfilledTime, &order.Seller)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":        "Inbound order successfully committed",
			"result":         resultCode,
			"id":             order.ID,
			"fulfilled_date": order.FulfilledDate,
			"fulfilled_time": order.FulfilledTime,
		})
	case 1:
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "No available warehouses or inbound order already fulfilled",
			"result": resultCode,
		})
	case 2:
		c.JSON(http.StatusNotFound, gin.H{
			"error":  fmt.Sprintf("Inbound order ID %s does not exist", id),
			"result": resultCode,
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "An error occurred while fulfilling an inbound order",
			"result": resultCode,
		})
	}
}
