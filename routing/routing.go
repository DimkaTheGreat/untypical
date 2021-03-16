package routing

import (
	"errors"

	"net/http"
	"untypicalCompanyTestTask/repository"

	"github.com/labstack/echo/v4"
)

var (
	errKeyNotFound   = errors.New("Key is required but not provided")
	errValueNotFound = errors.New("Value is required but not provided")
)

type Server struct {
	storage *repository.Storage
}

func Run(s *repository.Storage, port string) {
	server := NewServer(s)

	e := echo.New()

	e.HideBanner = true
	e.GET("/get", server.getValue)
	e.GET("/list", server.list)
	e.POST("/upsert", server.upsert)
	e.DELETE("/delete", server.delete)
	e.Start(":" + port)

}

func NewServer(s *repository.Storage) *Server {
	return &Server{
		storage: s,
	}

}

func (s *Server) getValue(c echo.Context) error {
	key := c.QueryParam("key")
	value, err := s.storage.GetValue(key)

	if err != nil {
		return err
	}

	c.String(http.StatusOK, value)
	if err != nil {
		return err
	}
	return nil

}

func (s *Server) list(c echo.Context) error {
	list, err := s.storage.List()

	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, list)
	return nil

}

func (s *Server) upsert(c echo.Context) error {
	entry := &repository.Entry{}

	err := c.Bind(entry)

	if err != nil {
		return err
	}

	err = s.storage.Upsert(entry.Key, entry.Value)

	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, "Done")

	return nil
}

func (s *Server) delete(c echo.Context) error {
	key := c.QueryParam("key")

	if key == "" {
		c.String(http.StatusBadRequest, "Key parameter is required")
	}

	err := s.storage.Delete(key)

	if err != nil {
		c.String(http.StatusBadGateway, err.Error())
		return err

	}

	c.String(http.StatusOK, "Done")

	return nil
}
