package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"wbl0/internal/model"
)

// PostgresDB представляет БД PostgreSQL
type PostgresDB struct {
	db *sql.DB
}

// Создает новый экземпляр БД и возвращает указатель на него
func NewDB(db *sql.DB) *PostgresDB {
	return &PostgresDB{db: db}
}

// Сохранение данных заказа в БД. Реализована в service.go
func (pg *PostgresDB) SaveData(data model.OrderInfo) error {
	return SaveToDB(pg.db, data)
}

// Инициализация и подключение к БД PostgreSQL. Возвращает указатель на созданное подключение.
func InitDB(host string, port string, user string, password string, dbname string) *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// Сохранение данных заказа в БД PostgreSQL с использованием транзакций.
func SaveToDB(db *sql.DB, data model.OrderInfo) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("insert into order_info (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		data.OrderUid, data.TrackNumber, data.Entry, data.Locale, data.InternalSignature, data.CustomerId, data.DeliveryService, data.Shardkey, data.SmId, data.DateCreated, data.OofShard)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		data.OrderUid, data.Delivery.Name, data.Delivery.Phone, data.Delivery.Zip, data.Delivery.City, data.Delivery.Address, data.Delivery.Region, data.Delivery.Email)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("INSERT INTO payments (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		data.OrderUid, data.Payment.Transaction, data.Payment.RequestId, data.Payment.Currency, data.Payment.Provider, data.Payment.Amount, data.Payment.PaymentDt, data.Payment.Bank, data.Payment.DeliveryCost, data.Payment.GoodsTotal, data.Payment.CustomFee)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, item := range data.Items {
		_, err = tx.Exec("INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
			data.OrderUid, item.ChrtId, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmId, item.Brand, item.Status)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// Получение данных заказа по идентификатору из БД PostgreSQL
func (pg *PostgresDB) GetDataById(id string) (model.OrderInfo, error) {
	row := pg.db.QueryRow("SELECT track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM order_info WHERE order_uid=$1", id)
	var order model.OrderInfo
	var trackNumber, entry, locale, internalSignature, customerId, deliveryService, shardkey, dateCreated, oofShard string
	var smId int
	err := row.Scan(&trackNumber, &entry, &locale, &internalSignature, &customerId, &deliveryService, &shardkey, &smId, &dateCreated, &oofShard)
	if err != nil {
		return model.OrderInfo{}, err
	}

	order.OrderUid = id
	order.TrackNumber = trackNumber
	order.Entry = entry
	order.Locale = locale
	order.InternalSignature = internalSignature
	order.CustomerId = customerId
	order.DeliveryService = deliveryService
	order.Shardkey = shardkey
	order.SmId = smId
	order.DateCreated = dateCreated
	order.OofShard = oofShard

	row = pg.db.QueryRow("SELECT name, phone, zip, city, address, region, email FROM deliveries WHERE order_uid=$1", id)
	var delivery model.Delivery
	err = row.Scan(&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address, &delivery.Region, &delivery.Email)
	if err != nil {
		return model.OrderInfo{}, err
	}

	order.Delivery = delivery

	row = pg.db.QueryRow("SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payments WHERE order_uid=$1", id)
	var payment model.Payment
	err = row.Scan(&payment.Transaction, &payment.RequestId, &payment.Currency, &payment.Provider, &payment.Amount, &payment.PaymentDt, &payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee)
	if err != nil {
		return model.OrderInfo{}, err
	}

	order.Payment = payment

	rows, err := pg.db.Query("SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items WHERE order_uid=$1", id)
	if err != nil {
		return model.OrderInfo{}, err
	}
	defer rows.Close()

	var items []model.Item
	for rows.Next() {
		var item model.Item
		err := rows.Scan(&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmId, &item.Brand, &item.Status)
		if err != nil {
			return model.OrderInfo{}, err
		}
		items = append(items, item)
	}

	order.Items = items

	return order, nil
}

// Получение всех данных заказов их PostgreSQL
func (pg *PostgresDB) GetAllData() ([]model.OrderInfo, error) {
	rows, err := pg.db.Query(`
		SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard 
		FROM order_info`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []model.OrderInfo
	for rows.Next() {
		var order model.OrderInfo
		if err := rows.Scan(
			&order.OrderUid,
			&order.TrackNumber,
			&order.Entry,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerId,
			&order.DeliveryService,
			&order.Shardkey,
			&order.SmId,
			&order.DateCreated,
			&order.OofShard,
		); err != nil {
			log.Printf("Ошибка при сканировании данных из базы данных: %v\n", err)
			return nil, err
		}

		deliveryRow := pg.db.QueryRow(
			`SELECT name, phone, zip, city, address, region, email FROM deliveries WHERE order_uid = $1`,
			order.OrderUid)
		var delivery model.Delivery
		if err := deliveryRow.Scan(
			&delivery.Name,
			&delivery.Phone,
			&delivery.Zip,
			&delivery.City,
			&delivery.Address,
			&delivery.Region,
			&delivery.Email,
		); err != nil {
			log.Printf("Ошибка при сканировании данных из таблицы delivery: %v\n", err)
			return nil, err
		}
		order.Delivery = delivery

		paymentRow := pg.db.QueryRow(`
			SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee 
			FROM payments WHERE order_uid = $1`,
			order.OrderUid)
		var payment model.Payment
		if err := paymentRow.Scan(
			&payment.Transaction,
			&payment.RequestId,
			&payment.Currency,
			&payment.Provider,
			&payment.Amount,
			&payment.PaymentDt,
			&payment.Bank,
			&payment.DeliveryCost,
			&payment.GoodsTotal,
			&payment.CustomFee,
		); err != nil {
			log.Printf("Ошибка при сканировании данных из таблицы payment: %v\n", err)
			return nil, err
		}
		order.Payment = payment

		itemsRows, err := pg.db.Query("SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items WHERE order_uid = $1", order.OrderUid)
		if err != nil {
			log.Printf("Ошибка при выполнении запроса к таблице item: %v\n", err)
			return nil, err
		}
		defer itemsRows.Close()

		var items []model.Item
		for itemsRows.Next() {
			var itemData model.Item
			if err := itemsRows.Scan(
				&itemData.ChrtId,
				&itemData.TrackNumber,
				&itemData.Price,
				&itemData.Rid,
				&itemData.Name,
				&itemData.Sale,
				&itemData.Size,
				&itemData.TotalPrice,
				&itemData.NmId,
				&itemData.Brand,
				&itemData.Status,
			); err != nil {
				log.Printf("Ошибка при сканировании данных из таблицы item: %v\n", err)
				return nil, err
			}
			items = append(items, itemData)
		}
		order.Items = items

		data = append(data, order)
	}

	return data, nil
}
