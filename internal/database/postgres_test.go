package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"wbl0/testutils"
)

func TestPostgresDB_GetDataById(t *testing.T) {
	//Подготовка к тестированию БД
	db := testutils.InitTestDatabase(t)
	defer db.Close()

	//Создание экземпляра PostgresDB
	pgDB := NewDB(db)

	//Сохранение заказа в БД
	err := pgDB.SaveData(testutils.TestOrder)
	assert.NoError(t, err)

	//Получение заказа по его ИД
	data, err := pgDB.GetDataById(testutils.TestOrder.OrderUid)
	assert.NoError(t, err)
	assert.Equal(t, testutils.TestOrder, data)
}
