package core

import (
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/streadway/amqp"
)

type Service interface {
	DB() *gorm.DB
	Echo() *echo.Echo
	Redis() *redis.Client
	Broker() *amqp.Channel
}
