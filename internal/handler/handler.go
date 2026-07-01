package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Handler {
	return &Handler{pool: pool}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/", h.hello)
}

func (h *Handler) hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"hello": "world"})
}