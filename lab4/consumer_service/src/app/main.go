package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

type Service struct {
	db     *gorm.DB
	cache  *redis.Client
	broker *amqp.Channel
}

type Book struct {
	ID        int       `json:"id"`
	Name      string     `json:"name"`
	Author    string     `json:"author"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	db, err := initSQLConnection()
	if err != nil {
		log.Fatalf("unable to establish Mysql connection: %s", err.Error())
	}
	if err := db.DB().Ping(); err != nil {
		log.Fatalf("mysql ping failed: %s", err.Error())
	}

	rds := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	broker, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatalf("unable to establish connection with rabbitmq broker: %s", err.Error())
	}
	ch, err := broker.Channel()
	if err != nil {
		log.Fatalf("unable to establish connection with rabbitmq broker: %s", err.Error())
	}

	s := &Service{db: db, cache: rds, broker: ch}
	if err := initConsumer("books", "key", s.saveHandler, ch); err != nil {
		log.Fatal(err)
	}
}

func initConsumer(queueName, routingKey string, handler func(d amqp.Delivery) bool, broker *amqp.Channel) error {
	log.Println("Starting a consumer!")
	_, err := broker.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	if err := broker.QueueBind(queueName, routingKey, "events", false, nil); err != nil {
		return err
	}

	msgs, err := broker.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	for msg := range msgs {
		if handler(msg) {
			_ = msg.Ack(false)
		} else {
			_ = msg.Nack(false, true)
		}
	}
	fmt.Println("Rabbit consumer closed - critical Error")

	return nil
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

func (s *Service) saveHandler(d amqp.Delivery) bool {
	if d.Body == nil {
		return false
	}

	book := new(Book)
	err := json.Unmarshal(d.Body, book)
	if err := s.db.Model(book).Save(book).Error; err != nil {
		return false
	}

	b, err := json.Marshal(book)
	if err != nil {
		return false
	}

	s.cache.Set(context.Background(), strconv.Itoa(book.ID), string(b), 0)
	if err != nil {
		return false
	}

	return true
}
