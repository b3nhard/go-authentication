package main

import (
	"fmt"
	"go-authentication/handlers"
	"go-authentication/models"
	"go-authentication/utils"
	"net/http"
	"os"
)

func main() {
	os.Setenv("SESSION_KEY", "_)+(*YU1GV@!}{:L:M@B!)W*(&YGF@!()IOJ@!~09iu3 ")
	utils.InitDB()
	utils.DbConn.AutoMigrate(&models.User{}, models.Role{})
	fmt.Println("[INFO] Migrated Database")
	fs := http.FileServer(http.Dir("static"))
	http.HandleFunc("/", handlers.Auth(handlers.HomeHandler))
	http.HandleFunc("/sign-in", handlers.LoginHandler)
	http.HandleFunc("/sign-up", handlers.SignUpHandler)
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	fmt.Println("[INFO] Server running on 127.0.0.1:8000")
	http.ListenAndServe(":8000", nil)
}
