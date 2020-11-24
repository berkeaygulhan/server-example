package main

import "net/http"

func main() {

	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./css/"))))
	http.Handle("/favicon.co", http.NotFoundHandler())
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/", mainPage)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/profilePage", profilePage)
	http.HandleFunc("/productAdd", productAdd)
	http.HandleFunc("/productDisplay", productDisplay)

	http.ListenAndServe(":8080", nil)

}
