package main

import (
	//"html/template"

	"main/db"
	"main/handlers"

	//"net/http"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	// no  need of parsing with html/templates when using using gin as there is
	// a method in the gin that can do that .we only do this where we usse net/http
	//as it lack the capablity
	//templates := template.Must(template.New("").ParseGlob("templates/*.html"))
	// ____and we we use code below to connect parsed html to the gin as we doing it seperately.
	//router.SetHTMLTemplate(templates)

	db.ConnectToDb()

	router.LoadHTMLGlob("templates/*")

	router.GET("/", handlers.RootHandler)

	router.GET("/home", handlers.HomeHandler)

	router.GET("/login", handlers.LoginGetHandler)

	router.POST("/login", handlers.LoginPostHandler)

	router.GET("/signup", handlers.SignupGetHandler)

	router.POST("/signup", handlers.SignupPostHandler)

	router.GET("/userprofile", handlers.UserprofileHandler)

	router.GET("/logout", handlers.LogoutHandler)

	router.GET("/adminhome", handlers.AdminhomeHandler)

	router.GET("/adminprofile", handlers.AdminprofileHandler)

	router.GET("/userlist", handlers.UserlistHandler)

	router.POST("/userlist", handlers.UserlistPostHandler)

	router.GET("/newuser", handlers.NewuserGetHandler)

	router.POST("/newuser", handlers.NewuserPostHandler)

	router.GET("/edit", handlers.EdituserGetHandler)

	router.POST("/update", handlers.EdituserPostHandler)

	router.GET("/delete", handlers.DeleteuserHandler, handlers.NewpageHandler)

	router.GET("/Newpage", handlers.NewpageHandler)

	router.Run(":6600")
}
