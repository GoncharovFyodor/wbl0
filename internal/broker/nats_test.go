package broker

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"wbl0/internal/cache"
	"wbl0/internal/database"
	"wbl0/internal/model"
	"wbl0/internal/service"
	"wbl0/testutils"
)

func TestSubscribeToNATS(t *testing.T) {
	//Подготовка к тестированию брокера NATS
	//Создание экземпляра сервиса и кэша
	db := testutils.InitTestDatabase(t)
	defer db.Close()

	pgDB := database.NewDB(db)
	cache := cache.NewCache()
	s := service.NewService(pgDB, cache)

	//Инициализация брокера
	InitNATS()
	SubscribeToNATS(s)

	//Публикация тестового заказа в брокер
	PublishOrderToNATS(testutils.TestOrder)

	//Проверка сохранения данных заказа в БД и кэше
	data, err := s.GetDataById(testutils.TestOrder.OrderUid)
	assert.NoError(t, err)
	assert.Equal(t, testutils.TestOrder, data)
}

// Вспомогательная функция для публикации заказа в NATS
func PublishOrderToNATS(order model.OrderInfo) {
	payload, err := json.Marshal(order)
	if err != nil {
		fmt.Println("Ошибка при преобразовании в JSON:", err)
		return
	}

	fmt.Println("Отправка заказа в NATS:", payload)
	err = Nconn.Publish("order_info", payload)
	if err != nil {
		fmt.Println("Ошибка при отправке данных в канал:", err)
	} else {
		fmt.Println("Данные успешно отправлены в канал")
	}
}
