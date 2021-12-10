package service

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"lab4/cmd"
	"lab4/routes"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/streadway/amqp"
)

type Service struct {
	e      *echo.Echo
	db     *gorm.DB
	cache  *redis.Client
	broker *amqp.Channel
}

type Book struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func New() (*Service, error) {
	e := echo.New()
	db, err := initSQLConnection()
	if err != nil {
		return nil, fmt.Errorf("unable to establish Mysql connection: %w", err)
	}
	if err := db.DB().Ping(); err != nil {
		return nil, fmt.Errorf("mysql ping failed: %w", err)
	}

	if err := cmd.Migrate(db); err != nil {
		return nil, fmt.Errorf("mysql migration failed: %w", err)
	}

	rds := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	broker, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		return nil, fmt.Errorf("unable to establish connection with rabbitmq broker: %w", err)
	}
	ch, err := broker.Channel()
	if err != nil {
		return nil, fmt.Errorf("unable to establish connection with rabbitmq broker: %w", err)
	}
	if err := ch.ExchangeDeclare("events", "topic", true, false, false, false, nil); err != nil {
		return nil, fmt.Errorf("unable to declare an exchange: %w", err)
	}

	s := &Service{
		e:      e,
		db:     db,
		cache:  rds,
		broker: ch,
	}

	api := routes.NewAPI(s)
	api.Init()

	return s, nil
}

func (s *Service) Echo() *echo.Echo {
	return s.e
}

func (s *Service) DB() *gorm.DB {
	return s.db
}

func (s *Service) Redis() *redis.Client {
	return s.cache
}

func (s *Service) Broker() *amqp.Channel {
	return s.broker
}

func initSQLConnection() (*gorm.DB, error) {
	connString := "%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local"

	return gorm.Open("mysql", fmt.Sprintf(connString,
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME")))
}
