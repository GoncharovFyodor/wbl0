package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"wbl0/testutils"
)

func TestCache_GetAndSetById(t *testing.T) {
	//Создание нового экземпляра кэша
	cache := NewCache()

	//Сохранение заказа в кэш
	cache.SetById(testutils.TestOrder.OrderUid, testutils.TestOrder)

	//Получение заказа по его идентификатору из кэша
	data, ok := cache.GetById(testutils.TestOrder.OrderUid)
	assert.True(t, ok)
	assert.Equal(t, testutils.TestOrder, data)
}
