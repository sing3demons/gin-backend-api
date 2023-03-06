package blog

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sing3demons/gin-backend-api/models"
	"github.com/sing3demons/gin-backend-api/utils"
	"gorm.io/gorm"
)

type handler struct {
	db *gorm.DB
}

func New(db *gorm.DB) *handler {
	return &handler{db}
}

func (h *handler) GetAll(c *gin.Context) {
	blogs := []models.Blog{}
	if err := h.db.Find(&blogs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}


	utils.ResponseJsonWithLogger(c, 200, blogs)

}

func (h *handler) GetById(c *gin.Context) {}

type formBlog struct {
	Topic string `json:"topic"`
}

func (h *handler) Create(c *gin.Context) {
	var body formBlog
	sub, _ := c.Get("sub")
	userId, _ := strconv.Atoi(sub.(string))

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	blog := models.Blog{
		Topic:  body.Topic,
		UserID: uint(userId),
	}

	if err := h.db.Create(&blog).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success create blog",
	})

}

func (h *handler) Update(c *gin.Context) {

}

func (h *handler) Delete(c *gin.Context) {

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
		"statusCode": 200,
		"resultDate": data,
	})
}
