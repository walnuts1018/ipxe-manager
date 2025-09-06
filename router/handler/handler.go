package handler

import (
	"github.com/walnuts1018/ipxe-manager/usecase"
	"github.com/walnuts1018/ipxe-manager/util/random"
)

type Handler struct {
	usecase *usecase.Usecase
	random  random.Random
}

func NewHandler(usecase *usecase.Usecase, random random.Random) *Handler {
	return &Handler{usecase: usecase, random: random}
}

