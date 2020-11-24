package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

//VARIABLES
var tpl *template.Template

const sessionLenght int = 30000

//FUNCTIONS
func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}
func mainPage(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "mainPage.gohtml", nil)
}
func profilePage(w http.ResponseWriter, req *http.Request) {

	u := getUser(w, req)
	fmt.Println("basliyor")
	fmt.Println(u.UserName)
	fmt.Println(u.Password)
	fmt.Println(u.FirstName)
	fmt.Println(u.Last)
	user1 := user{u.UserName, u.Password, u.FirstName, u.Last}
	/*logged := alreadyLoggedIn(w, req)
	if logged == false {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}*/
	tpl.ExecuteTemplate(w, "profilePage.gohtml", user1)
}
func signup(w http.ResponseWriter, req *http.Request) {
	//check if the user logged in
	/*if alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}*/
	var u user
	//check form submission
	if req.Method == http.MethodPost {
		//get from values
		un := req.FormValue("username")
		p := req.FormValue("password")
		f := req.FormValue("firstname")
		l := req.FormValue("lastname")

		//username taken?
		/*usernameCheck := checkUsername(un)
		if !usernameCheck {
			http.Error(w, "Username taken", http.StatusForbidden)
			return
		}*/
		//create session
		sID, err := uuid.NewUUID()
		if err != nil {
			panic(err)
		}
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(w, c)

		s := session{un, c.Value}
		insertSession(s)
		//store user in dbUsers
		bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		u = user{un, bs, f, l}
		insertUser(u)
		//Redirect
		http.Redirect(w, req, "/profilePage", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "signup.gohtml", u)

}

func login(w http.ResponseWriter, req *http.Request) {
	/*	if alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/display", http.StatusSeeOther)
		return
	}*/
	//process from submission
	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		p := req.FormValue("password")
		// username check

		//password check
		validatedUser := findUser(un, p)
		if !validatedUser {
			http.Error(w, "Username and password do not match", http.StatusForbidden)
			return
		}
		//create session
		sID, err := uuid.NewUUID()
		if err != nil {
			panic(err)
		}
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(w, c)
		s := session{un, c.Value}
		insertSession(s)
		http.Redirect(w, req, "/profilePage", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "login.gohtml", nil)
}

func logout(w http.ResponseWriter, req *http.Request) {
	/*if !alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/signup", http.StatusInternalServerError)
		return
	}*/
	c, _ := req.Cookie("session")
	// Delete session
	//deleteSession(c.Value)
	c = &http.Cookie{
		Name:   "session",
		Value:  "0",
		MaxAge: -1,
	}
	http.SetCookie(w, c)

	/*if time.Now().Sub(dbSessionsCleaned) > (time.Second * 30) {
		go cleanSessions()
	}*/
	http.Redirect(w, req, "/login", http.StatusSeeOther)
}

func productAdd(w http.ResponseWriter, req *http.Request) {
	//	if alreadyLoggedIn(w, req) {
	//		http.Redirect(w, req, "/", http.StatusSeeOther)
	//		return
	//	}
	var pro product
	if req.Method == http.MethodPost {
		n := req.FormValue("ProductName")
		c := req.FormValue("Condition")
		r := castBool(req.FormValue("Rentable"))
		p := req.FormValue("Price")
		d := req.FormValue("Description")

		pro = product{n, c, r, p, d}
		insertProduct(pro)
		//Redirect
		http.Redirect(w, req, "/productDisplay", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "productAdd.gohtml", pro)
}
func productDisplay(w http.ResponseWriter, req *http.Request) {
	/*if alreadyLoggedIn(w, req) {
		fmt.Println("adam burda")
	}
	*/

	var products []product

	products = getAllProducts()
	tpl.ExecuteTemplate(w, "productDisplay.gohtml", products)
	return

}
