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
