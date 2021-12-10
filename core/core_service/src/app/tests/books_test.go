package tests

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/require"
	"lab4/routes"
	"testing"
	"time"
)

type ServiceMock struct {
	Dd *gorm.DB
}

func (s *ServiceMock) DB() *gorm.DB {
	return s.Dd
}
func (s *ServiceMock) Echo() *echo.Echo {
	return nil
}
func (s *ServiceMock)Redis() *redis.Client {
	return nil
}
func (s *ServiceMock) Broker() *amqp.Channel {
	return nil
}

func TestGetZoeFromDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	resp := sqlmock.NewRows([]string{"id", "name", "author", "createdAt", "updatedAt", "deletedAt"}).AddRow(1, "testName", "testAuthor", time.Now(), time.Now(), time.Now())
	mock.ExpectQuery("^SELECT (.+)$").WillReturnRows(resp)
	mock.ExpectCommit()
	require.Nil(t, err)
	defer db.Close()

	conn, err := gorm.Open("mysql", db)
	serviceMock := &ServiceMock{Dd: conn}
	conn.LogMode(true)
	book, err := routes.GetBookFromDB(1, serviceMock)
	require.Nil(t, err)
	require.Equal(t, "testName", book.Name)
}