package handlers

import (
	//"fmt"

	"fmt"
	"main/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	//"gorm.io/driver/postgres"
	//"gorm.io/gorm"
	"main/db"
)

func clearHeader(c *gin.Context) {
	c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
}

// handlers

func RootHandler(c *gin.Context) {

	clearHeader(c)

	fmt.Println("root handler running")
	_, found := CookieFound(c)
	if found {
		c.Redirect(http.StatusSeeOther, "/home")
		return
	}
	fmt.Println("still  running ....")
	c.HTML(http.StatusOK, "login.html", nil)
}

func HomeHandler(c *gin.Context) {
	clearHeader(c)

	fmt.Println("homehandler running")

	claims, valid := CookieFound(c)
	if !valid {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	admin := (*claims)["admin"].(bool)
	if !admin {
		userid := uint((*claims)["userid"].(float64))
		username := db.Getusername(userid)
		c.HTML(http.StatusOK, "home.html", gin.H{
			"username": username,
		})
		return
	}
	c.Redirect(http.StatusSeeOther, "/adminhome")
}

func LoginGetHandler(c *gin.Context) {
	clearHeader(c)

	fmt.Println("login get running ")

	_, found := CookieFound(c)
	if found {
		c.Redirect(http.StatusSeeOther, "/home")
		return
	}
	c.HTML(http.StatusOK, "login.html", nil)
}

func LoginPostHandler(c *gin.Context) {
	clearHeader(c)

	fmt.Println("login post running ")

	username := c.PostForm("username")
	password := c.PostForm("password")

	isUser := authenticate(c, username, password)
	if !isUser {
		return
	}

	userid, err := db.GetUserid(username, password)
	if err != nil {
		fmt.Println("error while retriving id : ", err)
	}
	isadmin := db.Getrole(username, password)
	tokenString, err := getTokenString(userid, isadmin)
	if err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"tokenerr": "couldn't login please try again ... ",
		})
	}
	c.SetCookie("Authorise", tokenString, 3600*24, "/", "localhost", false, true)
	if isadmin {
		c.Redirect(http.StatusSeeOther, "/adminhome")
	} else {
		c.Redirect(http.StatusSeeOther, "/home")
	}

}

func SignupGetHandler(c *gin.Context) {
	clearHeader(c)
	c.HTML(http.StatusOK, "signup.html", nil)
}

func SignupPostHandler(c *gin.Context) {
	clearHeader(c)
	user := db.User{
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
		Email:    c.PostForm("email"),
	}

	userexist := db.CheckforUsername(user.Username)
	emailexist := db.CheckforEmail(user.Email)
	if userexist && emailexist {
		c.HTML(http.StatusAccepted, "signup.html", gin.H{
			"usernameErr": "username alredy taken !",
			"emailErr":    "Email already exist !",
		})
		return
	} else {
		if userexist {
			c.HTML(http.StatusAccepted, "signup.html", gin.H{
				"usernameErr": "username alredy taken !",
			})
			return
		}
		if emailexist {
			c.HTML(http.StatusAccepted, "signup.html", gin.H{
				"emailErr": "Email already exist !",
			})
			return
		}
	}
	fmt.Println(user.Username, user.Password)

	err := db.CreateUser(&user)
	if err != nil {
		fmt.Println("error while checking fior username :", err)
		c.HTML(http.StatusAccepted, "signup.html", gin.H{
			"signedup": "data base error !! , couldnt signup please try again...",
		})
		return
	}
	c.HTML(http.StatusAccepted, "signup.html", gin.H{ // can redirect into/ execute login page from here
		"signedup": "you have successfully signed up",
	})
}

func LogoutHandler(c *gin.Context) {
	clearHeader(c)

	fmt.Println("logout handler running")
	_, valid := CookieFound(c)
	if !valid {
		c.Redirect(http.StatusSeeOther, "/login")
	}
	c.SetCookie("Authorise", "", 3600, "/", "localhost", false, true)
	c.Redirect(http.StatusSeeOther, "/login")
}

func UserprofileHandler(c *gin.Context) {
	clearHeader(c)

	claims, valid := CookieFound(c)
	if !valid {
		c.Redirect(http.StatusSeeOther, "/login")
	}
	admin := (*claims)["admin"].(bool)
	if admin {
		c.Redirect(http.StatusSeeOther, "/")
	}
	userid := uint((*claims)["userid"].(float64))
	user := db.GetUserDetails(userid)
	c.HTML(http.StatusOK, "userprofile.html", gin.H{
		"username": user.Username,
		"email":    user.Email,
	})
}

// admin handlers

func AdminhomeHandler(c *gin.Context) {
	clearHeader(c)
	claims, valid := CookieFound(c)
	if !valid {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	admin := (*claims)["admin"].(bool)
	if !admin {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	userid := uint((*claims)["userid"].(float64))
	username := db.Getusername(userid)
	c.HTML(http.StatusOK, "adminhome.html", gin.H{
		"username": username,
	})

}

func AdminprofileHandler(c *gin.Context) {
	clearHeader(c)

	claims, valid := CookieFound(c)
	if !valid {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	admin := (*claims)["admin"].(bool)
	if !admin {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	userid := uint((*claims)["userid"].(float64))
	user := db.GetUserDetails(userid)
	c.HTML(http.StatusOK, "adminprofile.html", gin.H{
		"username": user.Username,
		"email":    user.Email,
	})
}

func NewuserGetHandler(c *gin.Context) {
	clearHeader(c)

	claims, valid := CookieFound(c)
	if !valid {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	admin := (*claims)["admin"].(bool)
	if !admin {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
	c.HTML(http.StatusOK, "newuser.html", nil)
}

func NewuserPostHandler(c *gin.Context) {
	clearHeader(c)

	claims, valid := CookieFound(c)
	if !valid {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	admin := (*claims)["admin"].(bool)
	if !admin {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
	user := db.User{
		Username: c.PostForm("username"),
		Email:    c.PostForm("email"),
		Password: c.PostForm("password"),
		Admin:    c.PostForm("isAdmin") == "true",
	}

	userexist := db.CheckforUsername(user.Username)
	emailexist := db.CheckforEmail(user.Email)

	if userexist && emailexist {
		c.HTML(http.StatusAccepted, "newuser.html", gin.H{
			"user":           user,
			"duplicateUser":  "username alredy taken !",
			"duplicateEmail": "Email already exist !",
		})
		return
	} else {
		if userexist {
			c.HTML(http.StatusAccepted, "newuser.html", gin.H{
				"user":          user,
				"duplicateUser": "username alredy taken !",
			})
			return
		}
		if emailexist {
			c.HTML(http.StatusAccepted, "newuser.html", gin.H{
				"user":           user,
				"duplicateEmail": "Email already exist !",
			})
			return
		}
	}

	err := db.CreateUser(&user)
	if err != nil {
		c.HTML(http.StatusOK, "newuser.html", gin.H{
			"added": "couldn't add the user ...database error",
		})
	} else {
		c.HTML(http.StatusOK, "newuser.html", gin.H{
			"added": "user added successfully",
		})
	}
}

func UserlistHandler(c *gin.Context) {
	clearHeader(c)

	claims, valid := CookieFound(c)
	if !valid {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	admin := (*claims)["admin"].(bool)
	if !admin {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	pagelimit := 8
	offset := 0
	var users []db.User
	var count int
	var err error

	users, count, err = db.GetUsers(pagelimit, offset)
	if err != nil {
		fmt.Println("error while getting userslist : ", err)
	}
	var userindex []int
	for i, _ := range users {
		userindex = append(userindex, i+1)
	}

	lastpage := count / pagelimit
	if count%pagelimit != 0 {
		lastpage++
	}

	c.HTML(http.StatusOK, "userlist.html", gin.H{
		"lastPage":    lastpage,
		"currentPage": 1,
		"nextPage":    2,
		"prvPage":     0,
		"userindex":   userindex,
		"search":      false,
		"data":        users,
	})

	// to dynamically show the list of users.
	// in the html there only two spans.
	// to learn how to show many list object with -
	//only one struct and dynamiclly put data.
}

func UserlistPostHandler(c *gin.Context) {
	clearHeader(c)

	claims, valid := CookieFound(c)
	if !valid {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	admin := (*claims)["admin"].(bool)
	if !admin {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	pagelimit := 8
	offset := 0
	var users []db.User
	var count int
	var err error
	searchWord := c.PostForm("searchword")

	fmt.Println("search eord in handler :", searchWord)
	users, count, err = db.GetSearchUsers(searchWord, pagelimit, offset)
	if err != nil {
		fmt.Println("error while searching users :", err)
	}
	var userindex []int
	for i, _ := range users {
		userindex = append(userindex, i+1)
		fmt.Println("appending index:", i)
	}
	fmt.Println(users, count)
	lastpage := count / pagelimit
	if count%pagelimit != 0 {
		lastpage++
	}

	c.HTML(http.StatusOK, "userlist.html", gin.H{
		"lastPage":    lastpage,
		"currentPage": 1,
		"nextPage":    2,
		"prvPage":     0,
		"search":      true,
		"userindex":   userindex,
		"searchWord":  searchWord,
		"data":        users,
	})
}

func NewpageHandler(c *gin.Context) {
	clearHeader(c)

	claims, valid := CookieFound(c)
	if !valid {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	admin := (*claims)["admin"].(bool)
	if !admin {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	newPage, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		fmt.Println("error while getting nextpage index :", err)
	}

	isSearch := c.Query("search") == "true"
	fmt.Println("boolean from userlisthandler", isSearch)

	pagelimit := 8 // number of users to list in single page
	offset := (newPage - 1) * pagelimit
	var users []db.User
	var count int
	searchWord := c.Query("word")

	if isSearch {
		users, count, err = db.GetSearchUsers(searchWord, pagelimit, offset)
		fmt.Println("user in new page when issearch :", users, "count", count)
		if err != nil {
			fmt.Println("error while searching users :", err)
		}
	} else {
		users, count, err = db.GetUsers(pagelimit, offset)
		fmt.Println("user in new page when is NOT  search :", users, "count", count)
		if err != nil {
			fmt.Println("error", err)
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}
	}
	var userindex []int
	for i, _ := range users {
		userindex = append(userindex, offset+i+1)

	}

	lastpage := count / pagelimit
	if count%pagelimit != 0 {
		lastpage = lastpage + 1
	}
	c.HTML(http.StatusOK, "userlist.html", gin.H{
		"lastPage":    lastpage,
		"currentPage": newPage,
		"nextPage":    newPage + 1,
		"prvPage":     newPage - 1,
		"userindex":   userindex,
		"search":      isSearch,
		"searchWord":  searchWord,
		"data":        users,
	})
}

func EdituserGetHandler(c *gin.Context) {
	clearHeader(c)

	claims, valid := CookieFound(c)
	if !valid {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	admin := (*claims)["admin"].(bool)
	if !admin {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	searchStatus := c.Query("search") == "true"
	searchWord := c.Query("word")
	userid, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		fmt.Println("error while converting userid in from string to uint : ", err)
	}
	user := db.GetUserDetails(uint(userid))
	fmt.Println("User in editgetusrhndlr", user)
	pagetoReturn, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		fmt.Println("error while getting index to retutn :", err)
	}

	c.HTML(http.StatusOK, "edituser.html", gin.H{
		"search":     searchStatus,
		"searchWord": searchWord,
		"page":       pagetoReturn,
		"user":       user,
	})
}

func EdituserPostHandler(c *gin.Context) {

	clearHeader(c)

	claims, valid := CookieFound(c)
	if !valid {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	admin := (*claims)["admin"].(bool)
	if !admin {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
	searchStatus := c.Query("search") == "true"
	searchWord := c.Query("word")
	userid, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		fmt.Println("error while converting userid in from string to uint : ", err)
	}
	fmt.Println("value of admin :", c.PostForm("admin"))
	fmt.Println(userid)
	userToUpdate := db.User{
		Username: c.PostForm("username"),
		Email:    c.PostForm("email"),
		Password: c.PostForm("password"),
		Admin:    c.PostForm("admin") == "true",
	}
	pagetoReturn, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		fmt.Println("error while getting index to retutn :", err)
	}

	userexist := db.CheckforUsername(userToUpdate.Username)
	emailexist := db.CheckforEmail(userToUpdate.Email)
	fmt.Println("userid :", userid)
	fmt.Println("unit of userid", uint(userid))
	user := db.GetUserDetails(uint(userid))
	if userexist && emailexist && user.Username != userToUpdate.Username && user.Email != userToUpdate.Email {
		c.HTML(http.StatusAccepted, "edituser.html", gin.H{
			"page":           pagetoReturn,
			"search":         searchStatus,
			"searchWord":     searchWord,
			"user":           userToUpdate,
			"duplicateUser":  "username alredy exist !",
			"duplicateEmail": "Email already exist !",
		})
		return
	} else {
		if userexist && user.Username != userToUpdate.Username {
			c.HTML(http.StatusAccepted, "edituser.html", gin.H{
				"page":          pagetoReturn,
				"search":        searchStatus,
				"searchWord":    searchWord,
				"user":          userToUpdate,
				"duplicateUser": "username alredy exist !",
			})
			return
		}
		if emailexist && user.Email != userToUpdate.Email {
			c.HTML(http.StatusAccepted, "edituser.html", gin.H{
				"page":           pagetoReturn,
				"search":         searchStatus,
				"searchWord":     searchWord,
				"user":           userToUpdate,
				"duplicateEmail": "Email already exist !",
			})
			return
		}
	}

	db.UpdateUser(userToUpdate, uint(userid))
	db.UpdateUserAdminStatus(uint(userid), userToUpdate.Admin)
	c.HTML(http.StatusOK, "edituser.html", gin.H{
		"page":       pagetoReturn,
		"search":     searchStatus,
		"searchWord": searchWord,
		"user":       userToUpdate,
		"updated":    "user successfully updated  :) ",
	})

	// below code to go back to userlist after updation
	/*updateduser, err := db.GetUsers()
	if err != nil {
		fmt.Println("error while getting user details :", err)
	}
	c.HTML(http.StatusOK, "userlist.html", gin.H{
		"data": updateduser,
	})*/
}

func DeleteuserHandler(c *gin.Context) {
	clearHeader(c)
	claims, valid := CookieFound(c)
	if !valid {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	admin := (*claims)["admin"].(bool)
	if !admin {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	userid, err := strconv.ParseUint(c.Query("id"), 10, 64)
	fmt.Println("id :", userid)
	if err != nil {
		fmt.Println("error while parasin userid from string to uint :", err)
	}
	db.DeletesUser(uint(userid)) // deletes the user
	//c.Redirect(http.StatusSeeOther, "/userlist") / and redirecting to the userlist handler
}

//functions

func getTokenString(userId uint, admin bool) (string, error) {

	fmt.Println("getTokenString  Running ")

	claims := jwt.MapClaims{
		"userid": userId,
		"admin":  admin,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(utils.SecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func CookieFound(c *gin.Context) (*jwt.MapClaims, bool) {

	fmt.Println("cookieFound running___> ")
	tokenString, err := c.Cookie("Authorise")

	if err != nil {
		fmt.Println(err)
		return nil, false
	}

	if claims, validated := validateTokenString(tokenString); validated {

		return claims, true
	}

	return nil, false
}

func validateTokenString(tokenString string) (*jwt.MapClaims, bool) {
	fmt.Println("validatetockenString  running")

	token, err := decodeTokenString(tokenString)
	if err != nil {

		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {

		return &claims, true
	}

	return nil, false
}

func decodeTokenString(tokenString string) (*jwt.Token, error) {

	fmt.Println("decodeTockenSring running ")

	token, err := jwt.Parse(tokenString, checkSigningMethod)
	if err != nil {

		return nil, err
	}

	return token, nil
}

func checkSigningMethod(token *jwt.Token) (interface{}, error) {

	fmt.Println("checksigning method running")

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return utils.SecretKey, nil
}

func authenticate(c *gin.Context, username, password string) bool {

	fmt.Println("authenticate running")

	userexist := db.CheckforUsername(username)

	if !userexist {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"usernameErr": "Username Not Found",
			"passwordErr": "Invalid Password",
		})
		return false
	}
	passwordmatched := db.Verifypassword(username, password)
	if passwordmatched {
		return true
	} else {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"passwordErr": "Incorrect Password",
		})
		return false
	}
}
