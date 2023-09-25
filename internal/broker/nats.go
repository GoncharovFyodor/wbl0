package broker

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"sync"
	"wbl0/internal/model"
	"wbl0/internal/service"
)

var Nconn *nats.Conn

// Соединение с сервером NATS
func InitNATS() {
	var err error
	Nconn, err = nats.Connect("http://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
}

// Подписка на указанный канал в сервере NATS и обработка полученных сообщений.
// Когда сообщение получено, оно декодируется из JSON в объект OrderInfo.
// Затем данные проходят валидацию, и если они проходят проверку, они сохраняются в кэш и БД через сервис s.
// Если происходит ошибка на любом из этапов, она регистрируется в журнале.
func SubscribeToNATS(s *service.Service, wg *sync.WaitGroup) {
	Nconn.Subscribe("order_info", func(m *nats.Msg) {
		defer wg.Done()
		var order model.OrderInfo
		err := json.Unmarshal(m.Data, &order)
		if err != nil {
			fmt.Println("Ошибка при разборе JSON:", err)
			return
		}

		fmt.Println("Получен заказ из NATS:", order)
		err = order.Validate()
		if err != nil {
			fmt.Println("Ошибка валидации данных:", err)
			return
		}
		fmt.Println("Валидация данных успешно завершена")

		err = s.SaveData(order)
		if err != nil {
			fmt.Println("Ошибка при сохранении данных:", err)
		} else {
			fmt.Println("Данные успешно сохранены")
		}
	})
}
