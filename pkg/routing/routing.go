package routing

import (
	"errors"
	"log"

	"net/http"
	"untypicalCompanyTestTask/pkg/repository"

	"github.com/labstack/echo/v4"
)

var (
	errKeyNotFound   = errors.New("Key is required but not provided\n")
	errValueNotFound = errors.New("Value is required but not provided\n")
)

//структура для внедрения зависимостей в обработчики
type Server struct {
	storage repository.Storager
	logger  *log.Logger
}

//запуск веб-сервера
func Run(s *repository.Storage, port string) {
	e := echo.New()

	server := NewServer(s)

	e.HideBanner = true

	//обработчкики
	e.GET("/get", server.getValue)
	e.GET("/list", server.list)
	e.POST("/upsert", server.upsert)
	e.DELETE("/delete", server.delete)
	e.Logger.Fatal(e.Start(":" + port))

}

//создание инстанса Server
func NewServer(s *repository.Storage) *Server {
	return &Server{
		storage: s,
		logger:  log.Default(),
	}

}

//метод для получения значения
func (s *Server) getValue(c echo.Context) error {
	key := c.QueryParam("key")

	//если ключ не был передан, выдаем ошибку
	if key == "" {
		c.String(http.StatusBadRequest, errKeyNotFound.Error())
		s.logger.Println(errKeyNotFound.Error())
		return errKeyNotFound
	}

	value, err := s.storage.GetValue(key)

	////если в хранилища такого значения нет - выдаем ошибку
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		s.logger.Println(err.Error())
	}

	c.String(http.StatusOK, value)

	return nil

}

//метод вывода списка всех объектов хранилища
func (s *Server) list(c echo.Context) error {
	list, err := s.storage.List()

	if err != nil {
		c.String(http.StatusOK, err.Error())
		s.logger.Println(err.Error())
		return err
	}

	c.JSONPretty(http.StatusOK, list, "	")
	return nil

}

//метод для добавления/изменения записи в хранилище
func (s *Server) upsert(c echo.Context) error {
	entry := &repository.Entry{}

	//получение данных из тела POST-запроса
	err := c.Bind(entry)

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		s.logger.Println(err.Error())
	}

	//проверка наличия в теле POST-запроса ключа
	if entry.Key == "" {
		c.String(http.StatusBadRequest, errKeyNotFound.Error())
		s.logger.Println(errKeyNotFound.Error())
		return errKeyNotFound
	}

	err = s.storage.Upsert(entry.Key, entry.Value)

	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, "Done")

	return nil
}

//метод удаления записи из хранилища
func (s *Server) delete(c echo.Context) error {
	key := c.QueryParam("key")

	//проверка наличия параметра ключа в запросе
	if key == "" {
		c.String(http.StatusBadRequest, errKeyNotFound.Error())
		s.logger.Println(errKeyNotFound.Error())
		return errKeyNotFound
	}

	err := s.storage.Delete(key)

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		s.logger.Println(err.Error())
		return err

	}

	c.String(http.StatusOK, "Done\n")

	return nil
}
