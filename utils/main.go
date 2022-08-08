package utils

import (
	"fmt"
	"github.com/gorilla/securecookie"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
)

var DbConn *gorm.DB
var err error

var CookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32),
)

func InitDB() {
	uri := "root:toor@tcp(127.0.0.1:3306)/GoAuth?charset=utf8mb4&parseTime=True&loc=Local"
	DbConn, err = gorm.Open(mysql.Open(uri), &gorm.Config{})
	if err != nil {
		fmt.Printf("[ERROR] DATABASE CONNECTION ERROR: %v", err)
	}
	fmt.Println("[INFO] Database Connection Successful")
}

func HashPassword(password string) ([]byte, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashPassword, nil
}

func VerifyPassword(hashPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}
func GetUsername(r *http.Request) (username string) {
	if cookie, err := r.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = CookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			username = cookieValue["name"]
		}
	}
	return username
}

func SetSession(w http.ResponseWriter, username string) {
	value := map[string]string{
		"name": username,
	}
	if encode, err := CookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encode,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
}

func ClearSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}
