package server

import (
	"backend/internal/database"
	"backend/internal/types"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
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

func login(w http.ResponseWriter, user types.User) error {
	log.SetPrefix("LOGIN ")

	if err := database.CheckUser(&user); err != nil {
		//http.Error(w, `{"error":"ОШИБКО ЧЕК ЮЗЕР: `+err.Error()+`"}`, http.StatusBadGateway)
		setLoginErrorCookie(w, err.Error())
		return err
	}

	key, err := newJwt(user.Username)
	if err != nil {
		http.Error(w, `{"error": "Ошибка создания ключа"}`, http.StatusBadGateway)
		return err
	}

	//заливаем в куки клиенту jwt ключ
	http.SetCookie(w, &http.Cookie{
		Name:     "key",
		Value:    key,
		Expires:  time.Now().AddDate(0, 0, 7),
		HttpOnly: true,
	})
	return nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	var user types.User

	if err := parseLoginRequest(r, &user); err != nil {
		http.Error(w, `{"error": "Ошибка парсинга формы "`+err.Error()+`}`, http.StatusBadGateway)
		return
	}

	if err := login(w, user); err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	form, err := database.GetForm(user.Username)
	if err != nil {
		http.Error(w, `{"error": "Ошибка форма не найдена: `+err.Error()+`"}`, http.StatusBadGateway)
		return
	}
	form_json, _ := json.Marshal(form)
	setUsernameCookie(w, user.Username)
	setFormDataCookie(w, form_json)
	setSuccessCookie(w)
	http.Redirect(w, r, "/profile", http.StatusFound)
}

func parseLoginRequest(r *http.Request, pUser *types.User) error {
	// Определяем Content-Type
	contentType := r.Header.Get("Content-Type")

	// Парсим JSON
	if strings.Contains(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(pUser); err != nil {
			return errors.New("invalid JSON format")
		}
		if len(pUser.Username) == 0 || len(pUser.Password) == 0 {
			return errors.New("username and password are required")
		}
		return validateLoginData(pUser)
	}

	// Парсим форму
	if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		if err := r.ParseForm(); err != nil {
			return errors.New("failed to parse form")
		}

		if len(r.Form["Username"]) == 0 || len(r.Form["Password"]) == 0 {
			return errors.New("username and password are required")
		}

		pUser.Username = strings.TrimSpace(r.Form["Username"][0])
		pUser.Password = os.Getenv("SALT") + strings.TrimSpace(r.Form["Password"][0])
		return validateLoginData(pUser)
	}

	return errors.New("unsupported content type")
}

func validateLoginData(pUser *types.User) error {
	l, err := regexp.Compile(`^FormUser_[0-9]{1,}$`)
	if err != nil {
		log.Print(err)
		return err
	}
	p, err := regexp.Compile(`^[a-zA-z0-9_]{0,}$`)
	if err != nil {
		log.Print(err)
		return err
	}
	if !l.MatchString(pUser.Username) || !p.MatchString(pUser.Password) {
		return errors.New("username or password invalid")
	}

	return nil
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := getJWtFromCookies(r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		claims, err := validateJwt(tokenStr)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:     "key",
				Value:    "",
				MaxAge:   -1,
				HttpOnly: true,
			})
			http.Redirect(w, r, "/", http.StatusBadRequest)
		}
		form, err := database.GetForm(claims.Username)
		if err != nil {
			http.Error(w, `{"error": "Ошибка форма не найдена: `+err.Error()+`"}`, http.StatusBadGateway)
			next.ServeHTTP(w, r)
			return
		}
		form_json, _ := json.Marshal(form)
		setUsernameCookie(w, claims.Username)
		setFormDataCookie(w, form_json)
		setSuccessCookie(w)
		next.ServeHTTP(w, r)
	}
}

func exitHandler(w http.ResponseWriter, r *http.Request) {
	clearCookies(w)
	removeJwtFromCookies(w)
	removeUsernameFromCookies(w)
	removePasswordFromCookies(w)
	http.Redirect(w, r, "/", http.StatusFound)
}
