package service

import (
	"wbl0/internal/cache"
	"wbl0/internal/database"
	"wbl0/internal/model"
)

// Сервис для работы с данными заказов
type Service struct {
	db    *database.PostgresDB
	cache *cache.Cache
}

// Создание нового экземпляра сервиса и инициализация его зависимостей
func NewService(db *database.PostgresDB, cache *cache.Cache) *Service {
	return &Service{
		db:    db,
		cache: cache,
	}
}

// Получение всех данных заказов из БД
func (s *Service) GetAllDataFromDB() ([]model.OrderInfo, error) {
	return s.db.GetAllData()
}

// Получение данных заказа по идентификатору.
// Сначала поиск производится в кэше, если данные там есть - возвращаем их.
// Если данных нет в кэше, производится обращение к БД, полученные данные сохраняются в кэше и возвращаются
func (s *Service) GetDataById(id string) (model.OrderInfo, error) {
	data, ok := s.cache.GetById(id)
	if ok {
		return data, nil
	}

	data, err := s.db.GetDataById(id)
	if err != nil {
		return model.OrderInfo{}, err
	}

	s.cache.SetById(id, data)
	return data, nil
}

// Сохранение данных заказа.
// Сначала проверяется наличие данных в кэше. Если есть, они повторно не сохраняются.
// Если данных нет в кэше, производится сохранение этих данных в БД и обновление кэша.
func (s *Service) SaveData(data model.OrderInfo) error {
	_, ok := s.cache.GetById(data.OrderUid)
	if ok {
		//Данные уже есть в кэше, повторно их сохранять не нужно
		return nil
	}
	if err := s.db.SaveData(data); err != nil {
		return err
	}

	s.cache.SetById(data.OrderUid, data)
	return nil
}
