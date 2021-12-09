package routes

import "lab4/core"

type API struct {
	s core.Service
}

func NewAPI(s core.Service) *API {
	return &API{s: s}
}

func (a *API) Init() {
	a.s.Echo().GET("/book", a.GetBook)
	a.s.Echo().POST("/book", a.CreateBook)
}