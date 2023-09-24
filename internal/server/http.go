package server

import (
	"encoding/json"
	"html/template"
	"net/http"
	"wbl0/internal/service"
)

// HTTP-сервер для обработки запросов
type HTTPServer struct {
	srv *service.Service
}

// Создание нового HTTP-сервера c заданных сервисом
func NewHTTPServer(srv *service.Service) *HTTPServer {
	return &HTTPServer{srv: srv}
}

// Обработка входящего HTTP-запроса и его маршрутизация в соответствии с URL
func (hs *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		hs.handleIndexPage(w, r)
	case "/get_data":
		hs.handleGetDataById(w, r)
	default:
		http.NotFound(w, r)
	}
}

// Обработка запроса на главную страницу и отображение шаблона для ввода данных заказа
func (hs *HTTPServer) handleIndexPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Ошибка при загрузке шаблона", 500)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Ошибка при выполнении шаблона", 500)
		return
	}
}

// Обработка запроса на получение данныхзаказа по идентификатору. Результат возвращается в формате JSON
func (hs *HTTPServer) handleGetDataById(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	data, err := hs.srv.GetDataById(id)
	if err != nil {
		http.Error(w, "Ошибка при получении данных", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
