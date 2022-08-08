package handlers

import (
	"github.com/gorilla/securecookie"
	"go-authentication/models"
	"go-authentication/utils"
	"gorm.io/gorm"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

type Message struct {
	Level string
	Value string
}
type payload struct {
	Title    string
	Message  Message
	AuthUser string
}


// LoginHandler - handles Login Logic
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("templates/Login.html")
	if err != nil {
		log.Printf("\n\n Template Error: %v ", err)
	}
	if strings.Compare(r.Method, "GET") == 0 {
		tpl.Execute(w, payload{Title: "Login"})
		return
	} else {
		//	Get form Data
		email, password := r.FormValue("email"), r.FormValue("password")
		var user models.User
		utils.DbConn.Where("email=?", email).Find(&user)
		if user.Email == "" {
			tpl.Execute(w, payload{Title: "Login", Message: Message{
				Level: "Error",
				Value: "Invalid Credentials",
			}})
			return
		}
		if ok := utils.VerifyPassword(user.Password, password); !ok {
			tpl.Execute(w, payload{Title: "Login", Message: Message{
				Level: "Error",
				Value: "Invalid Credentials",
			}})
			return
		}
		utils.SetSession(w, email)
		time.Sleep(time.Second * 2)
		http.Redirect(w, r, "/", 302)

	}

}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("templates/SignUp.html")
	if err != nil {
		log.Printf("template Error: %v", err)
	}
	if strings.Compare(r.Method, "GET") == 0 {
		tpl.Execute(w, payload{Title: "SignUp"})
		return
	} else {
		// Get Form Data
		username, password, email := r.FormValue("username"), r.FormValue("password"), r.FormValue("email")
		//Hash plain Password
		if len(username) < 1 && len(email) < 1 && len(password) < 1 {
			tpl.Execute(w, payload{Title: "SignUp", Message: Message{
				Level: "Error",
				Value: "Fill the Required Fields",
			}})
			return
		}
		hash, _ := utils.HashPassword(password)
		utils.DbConn.Create(&models.User{
			Model:    gorm.Model{},
			Username: username,
			Email:    email,
			Password: string(hash),
			Role:     "User",
		})

		tpl.Execute(w, payload{Title: "SignUp", Message: Message{
			Level: "Success",
			Value: "User Created, You can Now Login",
		}})

	}

}

//Auth Middleware Function
func Auth(next func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
	})
}

// HomeHandler - Handles Home route logic
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	username := utils.GetUsername(r)
	if username == "" {
		http.Redirect(w, r, "/sign-in", 302)
		return
	}
	mainfile := "templates/index.html"
	headerfile := "templates/partials/header.html"
	tpl, err := template.ParseFiles(mainfile, headerfile)
	if err != nil {
		log.Printf("Template Error: %v", err)
	}
	tpl.Execute(w, payload{Title: "Home", AuthUser: username})
}
