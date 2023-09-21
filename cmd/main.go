package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/joho/godotenv"
	"github.com/nats-io/stan.go"
	"github.com/patrickmn/go-cache"
	"net/http"
	"os"
)

type OrderInfo struct {
	OrderUid          string   `json:"order_uid"`
	TrackNumber       string   `json:"track_number"`
	Entry             string   `json:"entry"`
	Delivery          Delivery `json:"delivery"`
	Payment           Payment  `json:"payment"`
	Items             []Item   `json:"items"`
	Locale            string   `json:"locale"`
	InternalSignature string   `json:"internal_signature"`
	CustomerId        string   `json:"customer_id"`
	DeliveryService   string   `json:"delivery_service"`
	Shardkey          string   `json:"shardkey"`
	SmId              int      `json:"sm_id"`
	DateCreated       string   `json:"date_created"`
	OofShard          string   `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestId    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ChrtId      int64  `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmId        int64  `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func main() {
	fmt.Println("Инициализация приложения...")
	fmt.Println("Загрузка переменных окружения...")
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Ошибка загрузки .env файла")
		os.Exit(1)
	}
	fmt.Println("Переменные окружения загружены")
	fmt.Println("Подключение к БД:", os.Getenv("DB_NAME"), "...")

	//Формирование полного адреса для подключения к БД
	urlDb := fmt.Sprintf(
		"%v://%v:%v@%v:%v/%v",
		os.Getenv("DRIVER"),
		os.Getenv("USER_NAME"),
		os.Getenv("PASSWORD"),
		os.Getenv("HOST"),
		os.Getenv("PORT"),
		os.Getenv("DB_NAME"),
	)
	connection, err := pgx.Connect(context.Background(), urlDb)
	if err != nil {
		_, printErr := fmt.Fprintf(os.Stderr, "Не удалось подключиться к БД: %v\n", err)
		PrintError(printErr)
		os.Exit(1)
	}
	defer connection.Close(context.Background())
	fmt.Println("Подключение установлено")

	//Получение из БД всех order_uid
	orderUids := GetOrderUids(connection)

	//Создание кэша
	Cache := cache.New(cache.NoExpiration, -1)

	//Запись полученных по order_uid данных в кэш
	for i := range orderUids {
		data, _ := GetDataByUid(connection, orderUids[i])
		Cache.Set(orderUids[i], data, cache.NoExpiration)
	}
	fmt.Printf("Восстановление кэша из БД ... записей найдено (%v)\n", len(Cache.Items()))

	//Подключение и подписка на канал nats-streaming
	sc, _ := stan.Connect("test-cluster", "wbl0-pub")
	defer sc.Close()
	var order OrderInfo
	_, subscriptionErr := sc.Subscribe("service", func(msg *stan.Msg) {
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			fmt.Println("Got: Invalid json")
			insertionErr := InsertInvalidData(connection, string(msg.Data))
			PrintError(insertionErr)
		} else {
			fmt.Println("Got: Valid json")
			Cache.Set(order.OrderUid, string(msg.Data), cache.NoExpiration)
			InsertData(connection, order)
		}
	})
	if subscriptionErr != nil {
		fmt.Println("Subscription error:", subscriptionErr)
	}

	//http-сервер для получения информации о заказе по order_uid
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	fmt.Println("http-сервер запущен на http://localhost:8080")
	router.LoadHTMLGlob("templates/*.html")
	router.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", gin.H{
			"content": "index page",
		})
	})
	router.POST("/result", func(context *gin.Context) {
		result, _ := Cache.Get(context.PostForm("order_uid"))
		if result == nil {
			context.IndentedJSON(http.StatusOK, "Записей по данному order_uid не найдено")
		} else {
			context.IndentedJSON(http.StatusOK, order)
		}
	})
	err = router.Run(":8080")
	PrintError(err)
}

func GetOrderUids(connection *pgx.Conn) (sliceUid []string) {
	query := `select array_agg(order_uid) from order_info`
	err := connection.QueryRow(context.Background(), query).Scan(&sliceUid)
	if err != nil {
		_, printErr := fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		PrintError(printErr)
		os.Exit(1)
	}
	return sliceUid
}

func GetDataByUid(connection *pgx.Conn, uid string) (string, error) {
	var order OrderInfo
	query := `
		select order_info.*, to_jsonb(p.*) as "payment", to_jsonb(d.*) as "delivery",
			(select jsonb_agg(to_jsonb(i.*)) from items i) as "items"
		from order_info
		left join payments p on order_info.order_uid = p.transaction
		left join items i on order_info.track_number = i.track_number
		join 
			(select d.name, d.phone, d.zip, d.city, d.address, d.region, d.email
			 from deliveries d
			 where d.id = (
			     select od.delivery_id
			     from order_delivery od 
			     where od.order_uid = $1
			 )) as d on true
		where order_info.order_uid = $1
		limit 1`
	if err := connection.QueryRow(context.Background(), query, uid).Scan(
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
		&order.Payment,
		&order.Delivery,
		&order.Items,
	); err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			fmt.Println(pgError)
			return "", pgError
		}
	}
	res, _ := json.Marshal(&order)
	return string(res), nil
}

// Вставка невалидных данных
func InsertInvalidData(connection *pgx.Conn, data string) (err error) {
	query := `insert into invalid_data(data) values ($1)`
	if err = connection.QueryRow(context.Background(), query, data).Scan(); err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			fmt.Println(pgError)
			return pgError
		}
	}
	return nil
}

// Парсинг JSON и вставка данных
func InsertData(connection *pgx.Conn, order OrderInfo) {
	uid, err := InsertDataOrder(connection, order)
	PrintError(err)
	id, err := InsertDataDelivery(connection, order)
	PrintError(err)
	err = InsertDataPayment(connection, order)
	PrintError(err)
	err = InsertDataItems(connection, order)
	PrintError(err)
	err = InsertDataOrderDelivery(connection, uid, id)
	PrintError(err)
}

// Вставка данных о заказе в таблицу order_info
func InsertDataOrder(connection *pgx.Conn, order OrderInfo) (orderUid string, err error) {
	query := `insert into order_info
		(order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		returning order_uid
		`
	if err = connection.QueryRow(context.Background(), query,
		order.OrderUid,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerId,
		order.DeliveryService,
		order.Shardkey,
		order.SmId,
		order.DateCreated,
		order.OofShard,
	).Scan(&orderUid); err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			fmt.Println(pgError)
			return "", pgError
		}
	}
	return orderUid, nil
}

// Вставка данных о доставке в таблицу deliveries
func InsertDataDelivery(connection *pgx.Conn, order OrderInfo) (id int, err error) {
	query := `insert into deliveries
		(name, phone, zip, city, address, region, email)
		values ($1, $2, $3, $4, $5, $6, $7)
		returning id
		`
	if err = connection.QueryRow(context.Background(), query,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	).Scan(&id); err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			fmt.Println(pgError)
			return 0, pgError
		}
	}
	return id, nil
}

// Вставка данных о платеже в таблицу payments
func InsertDataPayment(connection *pgx.Conn, order OrderInfo) (err error) {
	query := `insert into payments
		(transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_totals, custom_fee)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`
	if err = connection.QueryRow(context.Background(), query,
		order.Payment.Transaction,
		order.Payment.RequestId,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDt,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	).Scan(); err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			fmt.Println(pgError)
			return pgError
		}
	}
	return nil
}

// Вставка данных о товарах в таблицу items
func InsertDataItems(connection *pgx.Conn, order OrderInfo) (err error) {
	for i := 0; i < len(order.Items); i++ {
		query := `insert into items
		(chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`
		if err = connection.QueryRow(context.Background(), query,
			order.Items[i].ChrtId,
			order.Items[i].TrackNumber,
			order.Items[i].Price,
			order.Items[i].Rid,
			order.Items[i].Name,
			order.Items[i].Sale,
			order.Items[i].Size,
			order.Items[i].TotalPrice,
			order.Items[i].NmId,
			order.Items[i].Brand,
			order.Items[i].Status,
		).Scan(); err != nil {
			var pgError *pgconn.PgError
			if errors.As(err, &pgError) {
				fmt.Println(pgError)
				return pgError
			}
		}
	}
	return nil
}

// Вставка данных в таблицу order_delivery для связи между таблицами order_info и deliveries
func InsertDataOrderDelivery(connection *pgx.Conn, orderUid string, id int) (err error) {
	query := `insert into order_delivery
		(order_uid, delivery_id)
		values ($1, $2)
		`
	if err = connection.QueryRow(context.Background(), query, orderUid, id).Scan(); err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			fmt.Println(pgError)
			return pgError
		}
	}
	return nil
}

func PrintError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
