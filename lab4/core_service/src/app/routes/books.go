package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"lab4/core"
	"net/http"
	"strconv"

	"lab4/models"
	"lab4/routes/requests"

	"github.com/labstack/echo/v4"
)

func (a *API) CreateBook(ctx echo.Context) error {
	req := new(requests.CreateBookRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Wrong request body")
	}

	book := &models.Book{
		Name:   req.Name,
		Author: req.Author,
	}

	b, err := json.Marshal(book)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("error while marshaling a book"))
	}

	if err := a.Publish(b); err != nil {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("error while saving a book"))
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (a *API) Publish(data []byte) error {
	return a.s.Broker().Publish("events", "key", false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         data,
		DeliveryMode: amqp.Persistent,
	})
}

func (a *API) GetBook(ctx echo.Context) error {
	bookStr := ctx.QueryParam("id")
	bookId, err := strconv.Atoi(bookStr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("<%s> is an invalid book id", bookStr))
	}

	bookJson, err := a.s.Redis().Get(context.Background(), bookStr).Result()
	if err != nil {
		book, err := GetBookFromDB(bookId, a.s)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, err.Error())
		}

		return ctx.JSON(http.StatusOK, book)
	}

	book := new(models.Book)
	if err := json.Unmarshal([]byte(bookJson), book); err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("Error unmarshalling book with id <%d>", bookId))
	}

	return ctx.JSON(http.StatusOK, book)
}

func GetBookFromDB(bookId int, s core.Service) (*models.Book, error) {
	b := new(models.Book)
	if err := s.DB().Model(b).Where("id = ?", bookId).First(b).Error; err != nil {
		return nil, fmt.Errorf("book <%d> not found", bookId)
	}

	return b, nil
}
