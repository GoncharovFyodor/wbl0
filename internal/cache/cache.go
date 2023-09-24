package cache

import (
	"sync"
	"wbl0/internal/model"
)

// Структура кэша для хранения заказов in-memory
type Cache struct {
	mu   sync.Mutex
	data map[string]model.OrderInfo
}

// Создать новый экземпляр кэша и вернуть указатель на него
func NewCache() *Cache {
	return &Cache{
		data: make(map[string]model.OrderInfo),
	}
}

// Получение данных заказа по его ID.
// Если данные для указанного ID отсутствуют в кэше, возвращается второй аргумент со значением false.
func (c *Cache) GetById(id string) (model.OrderInfo, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, ok := c.data[id]
	return data, ok
}

// Сохранение данных заказа в кэше по его ID.
// Если данные для указанного ID уже присутствуют в кэше, они будут заменены новыми данными.
func (c *Cache) SetById(id string, data model.OrderInfo) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[id] = data
}
