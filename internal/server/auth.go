package server

import (
	"backend/internal/database"
	"backend/internal/types"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("4ijk")

type Claims struct {
	Username string `json:"Username"`
	jwt.RegisteredClaims
}

func newJwt(username string) (string, error) {
	expirationTime := time.Now().AddDate(0, 0, 7).Add(10 * time.Minute)
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
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

func login(w http.ResponseWriter, user types.User) error {
	prevPref := log.Prefix()
	log.SetPrefix(prevPref + "LOGIN ")
	defer log.SetPrefix(prevPref)

	if err := database.CheckUser(&user); err != nil {
		//http.Error(w, `{"error":"ОШИБКО ЧЕК ЮЗЕР: `+err.Error()+`"}`, http.StatusBadGateway)
		setLoginErrorCookie(w, "Не верный логин или пароль \n (invalid username or password) \n (登录名或密码无效)")
		return err
	}

	key, err := newJwt(user.Username)
	if err != nil {
		http.Error(w, `{"error": "Ошибка создания ключа"}`, http.StatusBadGateway)
		return err
	}

	//заливаем в куки клиенту jwt ключ
	setJwtCookie(w, key)
	return nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	prevPref := log.Prefix()
	log.SetPrefix(prevPref + "LoginHandler ")
	defer log.SetPrefix(prevPref)
	var user types.User

	if err := parseLoginRequest(r, &user); err != nil {
		setLoginErrorCookie(w, err.Error())
		RedirectToProfile(w, r)
		return
	}

	if err := login(w, user); err != nil {
		//установка LoginErrorCookie реализована в login()
		RedirectToProfile(w, r)
		return
	}

	form, err := database.GetForm(user.Username)
	if err != nil {
		setLoginErrorCookie(w, fmt.Errorf(`ошибка форма не найдена`).Error())
		RedirectToProfile(w, r)
		return
	}
	form_json, _ := json.Marshal(form)
	setUsernameCookie(w, user.Username)
	setFormDataCookie(w, form_json)
	setSuccessCookie(w)
	RedirectToProfile(w, r)
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
		pUser.Password = strings.TrimSpace(r.Form["Password"][0])
		return validateLoginData(pUser)
	}

	return errors.New("unsupported content type")
}

func validateLoginData(pUser *types.User) error {
	prevPref := log.Prefix()
	log.SetPrefix(prevPref + "ValidateLoginData")
	defer log.SetPrefix(prevPref)
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

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prevPref := log.Prefix()
		log.SetPrefix(prevPref + "AuthMiddleware ")
		defer log.SetPrefix(prevPref)
		defer next.ServeHTTP(w, r)

		tokenStr, err := getJwtFromCookies(r)
		if err != nil {
			log.Println(err)
			log.SetPrefix(prevPref)
			return
		}

		claims, err := validateJwt(tokenStr)
		if err != nil {
			log.Println(err)
			removeJwtFromCookies(w)
			log.SetPrefix(prevPref)
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			return
		}
		form, err := database.GetForm(claims.Username)
		if err != nil {
			log.Println(err)
			log.SetPrefix(prevPref)
			http.Error(w, `{"error": "Ошибка форма не найдена: `+err.Error()+`"}`, http.StatusBadGateway)
			return
		}
		form_json, err := json.Marshal(form)
		if err != nil {
			log.Println(err)
			log.SetPrefix(prevPref)
			http.Error(w, `{"error": "Ошибка преобразования формы в json: `+err.Error()+`"}`, http.StatusBadGateway)
			return
		}
		_, err = getUsernameFromCookies(r)
		setUsernameCookie(w, claims.Username)
		setFormDataCookie(w, form_json)
		setSuccessCookie(w)
		if err != nil {
			if strings.Contains(r.URL.Path, "/profile") && strings.Contains(r.URL.Path, "/process") {
				RedirectToProfile(w, r)
			} else {
				http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
			}
		}
	})
}

func exitHandler(w http.ResponseWriter, r *http.Request) {
	prevPref := log.Prefix()
	log.SetPrefix(prevPref + "ExitHandler ")
	defer log.SetPrefix(prevPref)

	clearCookies(w)
	removeJwtFromCookies(w)
	removeUsernameFromCookies(w)
	removePasswordFromCookies(w)
	http.Redirect(w, r, "/", http.StatusFound)
}
