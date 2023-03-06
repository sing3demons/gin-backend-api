package user

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/matthewhartstonge/argon2"
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
	var users []models.User
	if err := h.db.Preload("Blogs").Find(&users).Error; err != nil {
		log.Panicln(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	utils.ResponseJsonWithLogger(c, 200, users)
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

	utils.ResponseJsonWithLogger(c, 200, user)
}

func (h *handler) Register(c *gin.Context) {
	var form FormRegister

	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	argon := argon2.DefaultConfig()
	hashedPassword, err := argon.HashEncoded([]byte(form.Password))
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), 8)
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

	utils.ResponseJsonWithLogger(c, 201, user)
}

func (h handler) Login(c *gin.Context) {
	var body FormLogin

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.db.Where("email = ?", body.Email).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})
		return
	}

	ok, err := argon2.VerifyEncoded([]byte(body.Password), []byte(user.Password))
	if !ok || err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid password",
		})
		return
	}

	type CustomClaims struct {
		Name string `json:"name"`
		jwt.RegisteredClaims
	}

	claims := &CustomClaims{
		Name: user.Fullname,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(int(user.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	utils.ResponseJsonWithLogger(c, 200, gin.H{
		"access_token": accessToken,
	})
}

func (h *handler) GetProfile(c *gin.Context) {
	id, ok := c.Get("sub")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	var user models.User
	if err := h.db.First(&user, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})
		return
	}

	utils.ResponseJsonWithLogger(c, 200, user)
}
