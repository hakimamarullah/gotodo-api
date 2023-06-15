package authcontroller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gotodo/config"
	"github.com/gotodo/helpers"
	"github.com/gotodo/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var userCredential models.User

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userCredential); err != nil {
		helpers.ResponseJSON(w, models.ResponseBody{Message: err.Error(), Code: http.StatusBadRequest})
		return
	}

	defer r.Body.Close()

	var user models.User

	if err := models.DB.Where("Email = ?", userCredential.Email).First(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			helpers.ResponseJSON(w, models.ResponseBody{Message: "Email or password is invalid", Code: http.StatusUnauthorized})
			return
		default:
			helpers.ResponseJSON(w, models.ResponseBody{Message: "Something went wrong. Please try again later", Code: http.StatusInternalServerError})
			return
		}

	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userCredential.Password)); err != nil {
		helpers.ResponseJSON(w, models.ResponseBody{Message: "Wrong Password!", Code: http.StatusUnauthorized})
		return
	}

	expiryLimit := time.Now().Add(time.Minute * 15)
	tokenAlgorithm := jwt.NewWithClaims(jwt.SigningMethodHS256, &config.JWTClaim{
		Email: user.Email,
		ID:    user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "gotodo-api",
			ExpiresAt: jwt.NewNumericDate(expiryLimit),
		},
	})

	token, err := tokenAlgorithm.SignedString(config.JWT_KEY)
	if err != nil {
		helpers.ResponseJSON(w, models.ResponseBody{Message: "Something went wrong. Please try again later", Code: http.StatusInternalServerError})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    token,
		HttpOnly: true,
	})

	helpers.ResponseJSON(w, models.ResponseBody{Message: fmt.Sprintf("You are logged in as %s", user.Email)})
}

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&user); err != nil {
		helpers.ResponseJSON(w, models.ResponseBody{Message: "Invalid data", Code: http.StatusBadRequest})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	var checkUser models.User
	models.DB.Where("email = ?", user.Email).First(&checkUser)

	if checkUser.ID != 0 {
		helpers.ResponseJSON(w, models.ResponseBody{Message: "Email already used", Code: http.StatusConflict})
		return
	}

	if err := models.DB.Create(&user).Error; err != nil {
		helpers.ResponseJSON(w, models.ResponseBody{Message: "Something went wrong. Please try again", Code: http.StatusInternalServerError})
		return
	}

	response := map[string]interface{}{
		"email": user.Email,
		"id":    user.ID,
	}

	helpers.ResponseJSON(w, models.ResponseBody{Message: "User created", Code: http.StatusCreated, Data: response, Count: 1})
}
