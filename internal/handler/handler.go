package handler

import (
	"fmt"
	"net/http"

	"github.com/Stipulations/goAuthy/internal/jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	pool  *pgxpool.Pool
	token *jwt.Service
}

func New(pool *pgxpool.Pool, token *jwt.Service) *Handler {
	return &Handler{pool: pool, token: token}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/", h.hello)

	admin := r.Group("/admin")
	{
		admin.POST("/login", h.login)
		admin.GET("/me", h.me)
	}

}

func (h *Handler) hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"hello": "world"})
}

func (h *Handler) login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	var (
		id   uuid.UUID
		hash string
	)

	err := h.pool.QueryRow(c.Request.Context(), "SELECT id, password_hash FROM admins WHERE username = $1", req.Username).Scan(&id, &hash)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid details"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid details"})
		return
	}

	gennedtoken, err := h.token.Generate(req.Username, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
		return
	}

	c.SetCookie("authToken", gennedtoken, 15*60, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged in"})
}

func (h *Handler) me(c *gin.Context) {
	hToken, err := c.Cookie("authToken")
	if err != nil {
		fmt.Printf("%s", hToken)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized1"})
		return
	}

	verified, err := h.token.Verify(hToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized2"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"username": verified.Subject, "id": verified.ID, "issuedAt": verified.IssuedAt})
}
