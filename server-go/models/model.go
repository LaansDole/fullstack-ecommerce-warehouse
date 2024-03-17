package models

import (
	"database/sql"
	"errors"
)

type LazadaUser struct {
	Username     string `json:"username" binding:"required"`
	RefreshToken string `json:"refresh_token"`
	PasswordHash string `json:"password_hash" binding:"required"`
}

type Buyer struct {
	Username string `json:"username" binding:"required"`
}

type Seller struct {
	Username string `json:"username" binding:"required"`
	ShopName string `json:"shop_name"`
	City     string `json:"city"`
}

type WHAdmin struct {
	Username     string `json:"username" binding:"required"`
	RefreshToken string `json:"refresh_token"`
	PasswordHash string `json:"password_hash" binding:"required"`
}

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

type Stockpile struct {
	ProductID   int `form:"product_id" json:"product_id"`
	WarehouseID int `form:"warehouse_id" json:"warehouse_id"`
	Quantity    int `form:"quantity" json:"quantity"`
}

type InboundOrder struct {
	ID            int    `form:"id" json:"id"`
	Quantity      int    `form:"quantity" json:"quantity" binding:"required"`
	ProductID     int    `form:"product_id" json:"product_id" binding:"required"`
	CreatedDate   string `form:"created_date" json:"created_date" binding:"required"`
	CreatedTime   string `form:"created_time" json:"created_time" binding:"required"`
	FulfilledDate string `form:"fulfilled_date" json:"fulfilled_date"`
	FulfilledTime string `form:"fulfilled_time" json:"fulfilled_time"`
	Seller        string `form:"seller" json:"seller" binding:"required"`
}

type BuyerOrder struct {
	ID            int    `form:"id" json:"id"`
	Quantity      int    `form:"quantity" json:"quantity" binding:"required"`
	ProductID     int    `form:"product_id" json:"product_id" binding:"required"`
	CreatedDate   string `form:"created_date" json:"created_date" binding:"required"`
	CreatedTime   string `form:"created_time" json:"created_time" binding:"required"`
	OrderStatus   string `form:"order_status" json:"order_status" binding:"required"`
	FulfilledDate string `form:"fulfilled_date" json:"fulfilled_date"`
	FulfilledTime string `form:"fulfilled_time" json:"fulfilled_time"`
	Buyer         string `form:"buyer" json:"buyer" binding:"required"`
}

type ProductCategory struct {
	CategoryName string `form:"category_name" json:"category_name"`
	Parent       string `form:"parent" json:"parent"`
}

// Endpoints for Buyers

func GetBuyer(username string) (*Buyer, error) {
	var buyer Buyer
	err := DBBuyer.QueryRow("SELECT * FROM buyer WHERE username = ?", username).Scan(&buyer.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &buyer, nil
}

func InsertBuyer(username string) (sql.Result, error) {
	result, err := DBBuyer.Exec("INSERT INTO buyer (username) VALUES (?)", username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

// Endpoints for Sellers

func GetSeller(username string) (*Seller, error) {
	var seller Seller
	err := DBSeller.QueryRow("SELECT * FROM seller WHERE username = ?", username).Scan(&seller.Username, &seller.ShopName, &seller.City)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &seller, nil
}

func InsertSeller(username, shopName, city string) (sql.Result, error) {
	result, err := DBSeller.Exec("INSERT INTO seller (username, shop_name, city) VALUES (?, ?, ?)", username, shopName, city)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func GetShopName(shopName string) (*Seller, error) {
	var seller Seller
	err := DBSeller.QueryRow("SELECT * FROM seller WHERE shop_name = ?", shopName).Scan(&seller.Username, &seller.ShopName, &seller.City)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &seller, nil
}

// Endpoints for WH Admins

func GetWHAdmin(username string) (*WHAdmin, error) {
	var admin WHAdmin
	var refreshToken sql.NullString // use sql.NullString for nullable columns

	err := DBSeller.QueryRow("SELECT * FROM lazada_user WHERE username = ?", username).Scan(&admin.Username, &refreshToken, &admin.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			// No admin with the provided username was found
			return nil, nil
		}
		return nil, err
	}

	admin.RefreshToken = refreshToken.String // convert sql.NullString to string
	return &admin, nil
}

func InsertWHAdmin(username, passwordHash string) (sql.Result, error) {
	result, err := DBAdmin.Exec("INSERT INTO wh_admin (username, password_hash) VALUES (?, ?)", username, passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func DeleteWHAdminToken(username string) error {
	_, err := DBAdmin.Exec("UPDATE wh_admin SET refresh_token = NULL WHERE username = ?", username)
	return err
}

// Endpoints for Lazada Users

func GetLazadaUser(username string) (*LazadaUser, error) {
	var user LazadaUser
	var refreshToken sql.NullString // use sql.NullString for nullable columns

	err := DBSeller.QueryRow("SELECT * FROM lazada_user WHERE username = ?", username).Scan(&user.Username, &refreshToken, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			// No user with the provided username was found
			return nil, nil
		}
		return nil, err
	}

	user.RefreshToken = refreshToken.String // convert sql.NullString to string
	return &user, nil
}

func GetLazadaUserByRole(role, username string) (interface{}, error) {
	switch role {
	case "seller":
		return GetSeller(username)
	case "buyer":
		return GetBuyer(username)
	case "admin":
		return GetWHAdmin(username)
	case "lazada_user":
		return GetLazadaUser(username)
	default:
		return nil, errors.New("invalid role")
	}
}

func InsertLazadaUserByRole(role, username, hashedPassword, shopName, city string) error {
	_, err := DBSeller.Exec("INSERT INTO lazada_user (username, password_hash) VALUES (?, ?)", username, hashedPassword)
	if err != nil {
		return err
	}

	switch role {
	case "buyer":
		_, err = InsertBuyer(username)
	case "seller":
		_, err = InsertSeller(username, shopName, city)
	case "admin":
		_, err = InsertWHAdmin(username, hashedPassword)
	default:
		err = errors.New("invalid role")
	}

	return err
}

func DeleteLazadaUserToken(username string) error {
	_, err := DBSeller.Exec("UPDATE lazada_user SET refresh_token = NULL WHERE username = ?", username)
	return err
}

// Endpoints for Products

// GetProduct retrieves a product by its ID.
func GetProduct(id int) (*Product, error) {
	var product Product
	err := DBAdmin.QueryRow("SELECT * FROM product WHERE id = ?", id).Scan(&product.ID, &product.Image, &product.Title, &product.ProductDescription, &product.Category, &product.Price, &product.Width, &product.Length, &product.Height, &product.Seller)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

// DeleteProduct deletes a product by its ID.
func DeleteProduct(id int) error {
	_, err := DBAdmin.Exec("DELETE FROM product WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

// Endpoints for Inbound Order

// GetInboundOrderByProduct retrieves an inbound order by product ID and seller.
func GetInboundOrderByProduct(productID int, seller string) (*InboundOrder, error) {
	var inboundOrder InboundOrder
	err := DBSeller.QueryRow("SELECT * FROM inbound_order WHERE product_id = ?", productID).Scan(&inboundOrder.ID, &inboundOrder.Quantity, &inboundOrder.ProductID, &inboundOrder.CreatedDate, &inboundOrder.CreatedTime, &inboundOrder.FulfilledDate, &inboundOrder.FulfilledTime, &inboundOrder.Seller)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &inboundOrder, nil
}

func GetInboundOrder(inboundOrderID int, seller string) (*InboundOrder, error) {
	var inboundOrder InboundOrder
	err := DBSeller.QueryRow("SELECT * FROM inbound_order WHERE id = ?", inboundOrderID).Scan(&inboundOrder.ID, &inboundOrder.Quantity, &inboundOrder.ProductID, &inboundOrder.CreatedDate, &inboundOrder.CreatedTime, &inboundOrder.FulfilledDate, &inboundOrder.FulfilledTime, &inboundOrder.Seller)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &inboundOrder, nil
}

// Endpoints for Buyer Order

// GetBuyerOrderByProduct retrieves a buyer order by product ID.
func GetBuyerOrderByProduct(productID int) (*BuyerOrder, error) {
	var buyerOrder BuyerOrder
	err := DBBuyer.QueryRow("SELECT * FROM buyer_order WHERE product_id = ?", productID).Scan(&buyerOrder.ProductID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &buyerOrder, nil
}

// Endpoints for Stockpile

// GetStockPileByProduct retrieves a stockpile by product ID.
func GetStockPileByProduct(productID int) (*Stockpile, error) {
	var stockPile Stockpile
	err := DBAdmin.QueryRow("SELECT * FROM stockpile WHERE product_id = ?", productID).Scan(&stockPile.ProductID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &stockPile, nil
}

// Endpoints for Product Category

// GetProductCategoryByName retrieves a product category by its name.
func GetProductCategoryByName(categoryName string) (*ProductCategory, error) {
	var category ProductCategory
	err := DBAdmin.QueryRow("SELECT * FROM product_category WHERE category_name = ?", categoryName).Scan(&category.CategoryName, &category.Parent)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}
