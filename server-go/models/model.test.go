package models

import "database/sql"

func GetAllBuyers() ([]*Buyer, error) {
	rows, err := DBBuyer.Query("SELECT * FROM buyer")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	var buyers []*Buyer
	for rows.Next() {
		var buyer Buyer
		err = rows.Scan(&buyer.Username)
		if err != nil {
			return nil, err
		}
		buyers = append(buyers, &buyer)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return buyers, nil
}
