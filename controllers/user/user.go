package user

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sing3demons/gin-backend-api/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type handler struct {
	db *gorm.DB
}

func New(db *gorm.DB) *handler {
	return &handler{db}
}

func (h handler) GetAll(c *gin.Context) {
	var users []models.User
	if err := h.db.Find(&users).Error; err != nil {
		log.Panicln(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	responses := []responseUser{}
	for _, user := range users {
		responses = append(responses, responseUser{
			ID:        user.ID,
			FullName:  user.Fullname,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.String(),
			UpdatedAt: user.UpdatedAt.String(),
		})
	}

	getResponseJson(c, responses)
}

func (h handler) SearchByName(c *gin.Context) {
	fullName := c.Param("name")
	if fullName != "" {
		c.JSON(http.StatusOK, gin.H{
			"message": "users with fullname: " + fullName,
		})
		return
	}
}

func (h handler) GetById(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := h.db.First(&user, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})
		return
	}

	getResponseJson(c, user)
}

func (h *handler) Register(c *gin.Context) {
	var form FormRegister

	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), 8)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Fullname: form.FullName,
		Email:    form.Email,
		Password: string(hashedPassword),
	}

	// err = bcrypt.CompareHashAndPassword(hashedPassword, password)

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	getResponseJson(c, user)
}

func (h handler) Login(c *gin.Context) {
	var user FormLogin

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	getResponseJson(c, user)
}

func getResponseJson(c *gin.Context, data any) {
	json, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	c.Writer.Header().Add("Response-Json", string(json))
	c.JSON(http.StatusOK, gin.H{
		"message":    "success",
		"resultDate": data,
	})
}
