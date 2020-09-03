package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type user struct {
	Uname string
	pwd   string
	Email string
	Amt   float32
}

func init() {

	dbUsers["bmerri"] = user{"bmerri", "pass", "bmerri@abc.com", 110.00}
	fmt.Println("Program started and init function is called")

}

//initiate database  connection
// func init() {
// 	db, err := sql.Open("mysql", "root:mpass@tcp(db-mysql:3306)/godb")
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	fmt.Println("db connection sucessful")
// 	defer db.Close()

// }

var dbUsers = map[string]user{}      // user ID, user
var dbSessions = map[string]string{} //// session ID, user ID

func createSessionID() string {

	var sess string
	var alph = [36]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', '9', '&', 'k', '6', 'm', '$', 'o', 'p',
		'q', 'r', 's', 't', 'v', '1', '2', '3', '0', '@', 'A', 'B', 'C', 'D', 'E', 'F', 'H', 'I', 'J', 'L',
	}
	fmt.Println(rand.Intn(36))

	for i := 0; i <= 50; i++ {
		sess += string(alph[rand.Intn(36)])
	}

	return sess

}

//Generic dbconn function
func dbConn() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:mpass@tcp(db-mysql:3306)/godb")
	if err != nil {
		panic(err.Error())
	}
	return db

}

func isLoggedIn(rw http.ResponseWriter, r *http.Request) (bool, string) {

	c := r.Cookies()
	leng := len(c)
	fmt.Println(c)
	if leng == 0 {
		return false, ""
	}
	var cval string
	for _, cook := range c {
		fmt.Println("cookie values are", cook, cook.Value, cook.Name)
		if cook.Name == "sessionid" {
			cval = cook.Value
			if cval == "" {
				return false, ""
			}
		}
	}

	db := dbConn()
	fmt.Println("db connection sucessful from isLoggedIn function")

	defer db.Close()
	var un string

	err1 := db.QueryRow("SELECT UNAME FROM sessions WHERE sessionID=?", cval).Scan(&un)
	if err1 != nil {
		panic(err1.Error())
	}
	return true, un

}

func login(rw http.ResponseWriter, r *http.Request) {
	log.Println(r.Method)

	if r.Method == http.MethodPost {
		ok, _ := isLoggedIn(rw, r)
		if ok {
			http.Redirect(rw, r, "/index", http.StatusSeeOther)

		}
		un := r.FormValue("Uname")
		pwd := r.FormValue("PWD")

		db := dbConn()
		log.Println("db connection sucessful for insertion in to session table")

		defer db.Close()
		//query := "INSERT INTO users VALUES (" + "\"" + un + "," + pwd + "\"" + ")"
		query := "SELECT UNAME FROM users WHERE UNAME =" + "\"" + un + "\"" + " and " + "PWD=" + "\"" + pwd + "\""
		var user string
		//fmt.Println(query)
		log.Println(query)

		err1 := db.QueryRow(query).Scan(&user)
		if err1 != nil {
			panic(err1.Error())
		}
		//fmt.Println(user)
		log.Println(user)

		sID := createSessionID()

		c := &http.Cookie{
			Name:  "sessionid",
			Value: sID,
		}

		//Set Cookie
		http.SetCookie(rw, c)

		//Setting the session
		iquery := "INSERT INTO sessions VALUES (" + "\"" + sID + "\"" + "," + "\"" + un + "\"" + ")"
		fmt.Println(query)

		iSess, err := db.Query(iquery)
		//fmt.Println(iSess)
		log.Println(iSess)

		// if there is an error inserting, handle it
		if err != nil {
			panic(err.Error())
		}
		// be careful deferring Queries if you are using transactions
		defer iSess.Close()

		//dbSessions[c.Value] = un
		http.Redirect(rw, r, "/index", http.StatusSeeOther)

	} else {

		ok, _ := isLoggedIn(rw, r)
		fmt.Println("cookie set? ", ok)
		if ok {
			http.Redirect(rw, r, "/index", http.StatusSeeOther)

		}

		tcl, err := template.ParseFiles("./login.html")

		if err != nil {
			log.Println(err)
		}
		ic := &http.Cookie{
			Name:  "sessionid",
			Value: "",
		}
		http.SetCookie(rw, ic)

		fmt.Println("Empty cookie is set")

		tcl.Execute(rw, nil)

	}

}

func logout(rw http.ResponseWriter, r *http.Request) {
	log.Println("In Logout Function")

	ok, _ := isLoggedIn(rw, r)

	if ok {

		c, err := r.Cookie("sessionid")
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(c.Value)
		cval := c.Value

		db := dbConn()
		//fmt.Println("db connection sucessful for deletion")
		log.Println("db connection sucessful for deletion")

		defer db.Close()
		var un string

		err1 := db.QueryRow("SELECT UNAME FROM sessions WHERE sessionID=?", cval).Scan(&un)
		if err1 != nil {
			panic(err1.Error())
		}

		delForm, err := db.Prepare("DELETE FROM sessions WHERE UNAME=?")
		if err != nil {
			panic(err.Error())
		}
		delForm.Exec(un)
		log.Println("Session DELETED")

		ec := &http.Cookie{
			Name:  "sessionid",
			Value: "",
		}

		//Empty Cookie
		http.SetCookie(rw, ec)
		fmt.Println("cookie set to Nil")

		http.Redirect(rw, r, "/login", 303)
	} else {
		http.Redirect(rw, r, "/login", 303)
	}

}

func register(rw http.ResponseWriter, r *http.Request) {
	//fmt.Println("In the register func and method is ", r.Method)

	log.Println("In the register func and method is ", r.Method)

	ok, _ := isLoggedIn(rw, r)
	if ok {
		http.Redirect(rw, r, "/index", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {

		un := r.FormValue("Uname")
		pwd := r.FormValue("PWD")
		email := r.FormValue("email")
		amt := r.FormValue("amt")
		fmt.Println(un, pwd, email, amt)

		db, err := sql.Open("mysql", "root:mpass@tcp(db-mysql:3306)/godb")
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("db connection sucessful", db.Stats())

		defer db.Close()
		query := "INSERT INTO users VALUES (" + "\"" + un + "\"" + "," + "\"" + pwd + "\"" + "," + "\"" + email + "\"" + "," + amt + ")"
		fmt.Println(query)

		insert, err := db.Query(query)
		fmt.Println(insert)

		// if there is an error inserting, handle it
		if err != nil {
			panic(err.Error())
		}
		// be careful deferring Queries if you are using transactions
		defer insert.Close()
		http.Redirect(rw, r, "/login", http.StatusSeeOther)

	}
	regtmp, err := template.ParseFiles("./register.html")
	if err != nil {
		fmt.Println(err)
	}
	regtmp.Execute(rw, nil)
	http.Redirect(rw, r, "/login", http.StatusSeeOther)
}

func index(rw http.ResponseWriter, r *http.Request) {

	//fmt.Println(r.Method)
	log.Println(r.Method)

	ok, loggedUSer := isLoggedIn(rw, r)
	fmt.Println("Cookie set? ", ok)
	if !ok {
		http.Redirect(rw, r, "/login", http.StatusSeeOther)
		return
	}

	db := dbConn()
	//fmt.Println("INDEX Fn: db connection sucessful for retrival of data")
	var eml string
	var amt float32

	defer db.Close()

	err1 := db.QueryRow("SELECT email, amt FROM users WHERE UNAME=?", loggedUSer).Scan(&eml, &amt)
	if err1 != nil {
		panic(err1.Error())
	}

	var u = user{loggedUSer, "", eml, amt}
	fmt.Println("In the index page ")

	tcl, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Println(err)
	}
	err = tcl.Execute(rw, u)
	if err != nil {
		fmt.Println(err)
	}

}

func redirindex(rw http.ResponseWriter, r *http.Request) {
	http.Redirect(rw, r, "/login", http.StatusSeeOther)

}

func main() {

	http.HandleFunc("/", redirindex)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/index", index)
	http.HandleFunc("/register", register)
	log.Fatal(http.ListenAndServe(":8081", nil))

}
