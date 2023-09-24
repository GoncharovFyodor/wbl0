package database

import (
	"database/sql"
	"fmt"
)

// Создание таблиц в БД, если они еще не созданы.
func CreateTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS order_info (
	order_uid VARCHAR(50) PRIMARY KEY,
	track_number VARCHAR(50) NOT NULL,
	entry VARCHAR(50) NOT NULL,
	locale VARCHAR(10) NOT NULL,
	internal_signature VARCHAR(100),
	customer_id VARCHAR(50) NOT NULL,
	delivery_service VARCHAR(50) NOT NULL,
	shardkey VARCHAR(50) NOT NULL,
	sm_id INT NOT NULL,
	date_created VARCHAR(50) NOT NULL,
	oof_shard VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS deliveries (
	order_uid VARCHAR(50) PRIMARY KEY REFERENCES order_info(order_uid),
	name VARCHAR(100) NOT NULL,
	phone VARCHAR(20) NOT NULL,
	zip VARCHAR(10) NOT NULL,
	city VARCHAR(100) NOT NULL,
	address VARCHAR(100) NOT NULL,
	region VARCHAR(100) NOT NULL,
	email VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS payments (
	order_uid VARCHAR(50) PRIMARY KEY REFERENCES order_info(order_uid),
	transaction VARCHAR(100) NOT NULL,
	request_id VARCHAR(100),
	currency VARCHAR(10) NOT NULL,
	provider VARCHAR(50) NOT NULL,
	amount INT NOT NULL,
	payment_dt INT NOT NULL,
	bank VARCHAR(50) NOT NULL,
	delivery_cost INT NOT NULL,
	goods_total INT NOT NULL,
	custom_fee INT
);

CREATE TABLE IF NOT EXISTS items (
    order_uid VARCHAR(50) REFERENCES order_info(order_uid),
	chrt_id INT NOT NULL,
	track_number VARCHAR(50) NOT NULL,
	price INT NOT NULL,
	rid VARCHAR(100) NOT NULL,
	name VARCHAR(100) NOT NULL,
	sale INT NOT NULL,
	size VARCHAR(50) NOT NULL,
	total_price INT NOT NULL,
	nm_id INT NOT NULL,
	brand VARCHAR(100) NOT NULL,
	status INT NOT NULL,
    PRIMARY KEY (order_uid, chrt_id)
);`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("ошибка при создании таблиц в базе данных: %w", err)
	}

	return nil
}
