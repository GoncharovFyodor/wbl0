package server

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"wbl0/internal/cache"
	"wbl0/internal/database"
	"wbl0/internal/service"
	"wbl0/testutils"
)

func TestHTTPServer_GetOrder(t *testing.T) {
	// Подготовка к тестированию HTTP-сервера
	db := testutils.InitTestDatabase(t)
	defer db.Close()
	pgDB := database.NewDB(db)
	cache := cache.NewCache()
	s := service.NewService(pgDB, cache)

	// Создание экземпляра HTTP-сервера
	handler := NewHTTPServer(s)

	//Создание тестового заказа и сохранение его в БД и кэш
	err := s.SaveData(testutils.TestOrder)
	assert.NoError(t, err)

	//Создание тестового http-запроса
	tUrl := "/get_data?id=" + testutils.TestOrder.OrderUid
	rq, err := http.NewRequest("GET", tUrl, nil)
	assert.NoError(t, err)

	//Создание рекордера HTTP-ответов
	recorder := httptest.NewRecorder()

	//Обработка HTTP-запроса с помощью обработчика
	handler.ServeHTTP(recorder, rq)

	//Проверка кода ответа
	assert.Equal(t, http.StatusOK, recorder.Code)
}
