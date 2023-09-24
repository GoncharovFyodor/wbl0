package main

import (
	"fmt"
	"log"
	"net/http"
	"wbl0/internal/broker"
	"wbl0/internal/cache"
	"wbl0/internal/database"
	"wbl0/internal/server"
	"wbl0/internal/service"
)

func main() {
	//Инициализация БД postgres
	db := database.InitDB("localhost", "5432", "postgres", "postgres", "wbl0")
	defer db.Close()

	//Создание таблиц в БД, если они еще не созданы
	if err := database.CreateTables(db); err != nil {
		log.Fatal("Ошибка при создании таблиц в базе данных")
	}

	//Инициализация кэша
	cache := cache.NewCache()

	//Создание сервиса для работы с БД
	dataService := service.NewService(database.NewDB(db), cache)

	//При старте сервиса восстанавливаются данные из БД в кэш.
	data, err := dataService.GetAllDataFromDB()
	if err != nil {
		log.Fatal("Ошибка при чтении данных из базы данных:", err)
	}

	for _, order := range data {
		cache.SetById(order.OrderUid, order)
	}

	//Инициализация NATS
	broker.InitNATS()

	//Подписка на канал NATS для получения данных
	broker.SubscribeToNATS(dataService)

	//Обработчик http для корневого пути
	http.Handle("/", server.NewHTTPServer(dataService))
	addr := ":8080"
	fmt.Println("Сервер запущен и слушает на", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
