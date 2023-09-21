package main

import "github.com/nats-io/stan.go"

var validJson = `{
"order_uid": "testtest",
"track_number": "TRRRRRRRACK",
"entry": "WBIL",
"delivery": {
"name": "Test2 Testov2",
"phone": "+9720000000",
"zip": "2639809",
"city": "Kiryat Mozkin",
"address": "Ploshad Mira 15",
"region": "Kraiot",
"email": "test@gmail.com"
},
"payment": {
"transaction": "testtest",
"request_id": "",
"currency": "USD",
"provider": "wbpay",
"amount": 1817,
"payment_dt": 1637907727,
"bank": "vtb",
"delivery_cost": 1500,
"goods_total": 317,
"custom_fee": 0
},
"items": [
{
"chrt_id": 9934930,
"track_number": "TRRRRRRRACK",
"price": 453,
"rid": "ab4219087a764ae0btest",
"name": "Mascaras",
"sale": 30,
"size": "0",
"total_price": 317,
"nm_id": 2389212,
"brand": "Vivienne Sabo",
"status": 202
}
],
"locale": "en",
"internal_signature": "",
"customer_id": "test",
"delivery_service": "meest",
"shardkey": "9",
"sm_id": 99,
"date_created": "2021-11-26T06:22:19Z",
"oof_shard": "1"
}`

var test = `{
  "order_uid": "b563feb7b2b84b6test",
  "track_number": "WBILMTESTTRACK",
  "entry": "WBIL",
  "delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
  },
  "payment": {
    "transaction": "b563feb7b2b84b6test",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 1817,
    "payment_dt": 1637907727,
    "bank": "alpha",
    "delivery_cost": 1500,
    "goods_total": 317,
    "custom_fee": 0
  },
  "items": [
    {
      "chrt_id": 9934930,
      "track_number": "WBILMTESTTRACK",
      "price": 453,
      "rid": "ab4219087a764ae0btest",
      "name": "Mascaras",
      "sale": 30,
      "size": "0",
      "total_price": 317,
      "nm_id": 2389212,
      "brand": "Vivienne Sabo",
      "status": 202
    }
  ],
  "locale": "en",
  "internal_signature": "",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}`

var manyItems = `{
"order_uid": "manyitems",
"track_number": "WBILMTESTTRACK",
"entry": "WBIL",
"delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
},
"payment": {
    "transaction": "manyitems",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 1817,
    "payment_dt": 1637907727,
    "bank": "alpha",
    "delivery_cost": 1500,
    "goods_total": 317,
    "custom_fee": 0
},
"items": [
    {
    "chrt_id": 1,
    "track_number": "WBILMTESTTRACK",
    "price": 453,
    "rid": "ab4219087a764ae0btest",
    "name": "whatever",
    "sale": 30,
    "size": "0",
    "total_price": 317,
    "nm_id": 2389212,
    "brand": "Vivienne Sabo",
    "status": 202
    },
    {
    "chrt_id": 22,
    "track_number": "WBILMTESTTRACK",
    "price": 453,
    "rid": "ab4219087a764ae0btest",
    "name": "WhiteHat",
    "sale": 30,
    "size": "0",
    "total_price": 317,
    "nm_id": 2389212,
    "brand": "vnnn",
    "status": 202
    },
    {
    "chrt_id": 23,
    "track_number": "WBILMTESTTRACK",
    "price": 453,
    "rid": "ab4219087a764ae0btest",
    "name": "RedHat",
    "sale": 30,
    "size": "0",
    "total_price": 317,
    "nm_id": 2389212,
    "brand": "Vivienne Sabo",
    "status": 202
    },
    {
    "chrt_id": 88,
    "track_number": "WBILMTESTTRACK",
    "price": 453,
    "rid": "ab4219087a764ae0btest",
    "name": "BlackHat",
    "sale": 30,
    "size": "0",
    "total_price": 317,
    "nm_id": 2389212,
    "brand": "Vivienne Sabo",
    "status": 202
    }
],
"locale": "en",
"internal_signature": "",
"customer_id": "test",
"delivery_service": "meest",
"shardkey": "9",
"sm_id": 99,
"date_created": "2021-11-26T06:22:19Z",
"oof_shard": "1"
}`
var validJson2 = `{
"order_uid": "827347feoih",
"track_number": "test_track_number",
"entry": "WBIL",
"delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
},
"payment": {
    "transaction": "827347feoih",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 1817,
    "payment_dt": 1637907727,
    "bank": "alpha",
    "delivery_cost": 1500,
    "goods_total": 317,
    "custom_fee": 0
},
"items": [
    {
    "chrt_id": 9934930,
    "track_number": "test_track_number",
    "price": 453,
    "rid": "ab4219087a764ae0btest",
    "name": "Mascaras",
    "sale": 30,
    "size": "0",
    "total_price": 317,
    "nm_id": 2389212,
    "brand": "Vivienne Sabo",
    "status": 202
    }
],
"locale": "en",
"internal_signature": "",
"customer_id": "test",
"delivery_service": "meest",
"shardkey": "9",
"sm_id": 99,
"date_created": "2021-11-26T06:22:19Z",
"oof_shard": "1"
}`

func main() {
	sc, _ := stan.Connect("test-cluster", "wbl0")
	defer sc.Close()
	sc.Publish("service", []byte(validJson))
	sc.Publish("service", []byte("{invalid}"))
	sc.Publish("service", []byte(test))
	sc.Publish("service", []byte(manyItems))
	sc.Publish("service", []byte(validJson2))
	sc.Publish("service", []byte("notjson"))
}
