package server

import (
	"backend/internal/database"
	"backend/internal/types"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("MABALLS")

type Claims struct {
	Username string `json:"Username"`
	jwt.RegisteredClaims
}

func newJwt(username string) (string, error) {
	expirationTime := time.Now().AddDate(0, 0, 7).Add(10 * time.Minute) // хайп ?
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func validateJwt(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("ИНВАЛИД ТОКЕН")
	}
	return claims, nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		http.Error(w, `{"error": "Ошибка парсинга формы"}`, http.StatusBadGateway)
		return
	}

	var user types.User
	parseLoginForm(&user, r)

	if err := database.CheckUser(&user); err != nil {
		http.Error(w, `{"error":"ОШИБКО ЧЕК ЮЗЕР"}`, http.StatusBadGateway)
		return
	}

	key, err := newJwt(user.Username)
	if err != nil {
		http.Error(w, `{"error": "Ошибка создания ключа"}`, http.StatusBadGateway)
		log.Panic("ошибка создания ключа", err)
		return
	}

	//заливаем в куки клиенту jwt ключ
	http.SetCookie(w, &http.Cookie{
		Name:     "key",
		Value:    key,
		Expires:  time.Now().AddDate(0, 0, 7),
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusOK)
}

func parseLoginForm(pUser *types.User, r *http.Request) error {
	if !(r.Form.Has("Username") && r.Form.Has("Password")) {
		return errors.New("invalid login_form")
	}
	pUser.Username = r.Form.Get("Username")
	pUser.Password = r.Form.Get("Password")
	return nil
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := r.Cookie("key")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		claims, err := validateJwt(tokenStr.Value)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:     "key",
				Value:    "",
				MaxAge:   -1,
				HttpOnly: true,
			})
			http.Redirect(w, r, "/login", http.StatusBadRequest)
		}
		form, err := database.GetForm(claims.Username)
		if err != nil {
			http.Error(w, `{"error": "Ошибка форма не найдена"}`, http.StatusBadGateway)
			return
		}
		form_json, _ := json.Marshal(form)
		setFormDataCookie(w, form_json)
		setSuccsessCookie(w)
		next.ServeHTTP(w, r)
	}
}
