package main

import (
    "github.com/titmuscody/bank/db"
	"fmt"
	"time"
	"net/http" 
	"io/ioutil"
	"encoding/json"
    "strings"
    )

type Page struct{
	Title string
	Body []byte
}

type Operation struct {
	Add []int 
}

func (p *Page) save() error {
	filename := p.Title + ".html"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error){
	filename := title + ".html"
	body, _ := ioutil.ReadFile(filename)
	return &Page{Title:title, Body:body}, nil
}

func handler(w http.ResponseWriter, r *http.Request){
	fmt.Println("printing headers")
	fmt.Println(r.Header)
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Println("printing body")
	fmt.Println(string(body))
	store := &Operation{}
	if string(body) != ""{
		m := string(body)
		json.Unmarshal([]byte(m), &store)
		sum := 0
		for _, e := range store.Add{
			sum += e
		}
		fmt.Fprintf(w, "{\"result\":%d}", sum)
	}
    fmt.Fprintf(w, "%d", db.GetData())
}

func viewHandler(w http.ResponseWriter, r *http.Request){
	//expires := time.Now().Add(365 * 24 * time.Hour)
	//cookie := http.Cookie{Name:"username", Value:"me", Expires:expires,
	//	HttpOnly:false}
	//http.SetCookie(w, &cookie)
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	fmt.Println("starting cookie")
	fmt.Fprintf(w, "%s", p.Body)
	//cookie, _ := r.Cookie("username")
	//fmt.Fprint(w, cookie)
	for _, c := range r.Cookies(){
		fmt.Println(c.Name)
	}
	//fmt.Println(w, r.Cookies())

	
}
func scriptHandler(w http.ResponseWriter, r *http.Request){
	title := r.URL.Path[len("/scripts/"):]
	body, _ := ioutil.ReadFile(title + ".js")
	//p, _ := loadPage(title)
	fmt.Fprintf(w, "%s", body)
}

func cssHandler(w http.ResponseWriter, r *http.Request){
	title := r.URL.Path[len("/css/"):]
	body, _ := ioutil.ReadFile(title + ".css")
	//p, _ := loadPage(title)
    w.Header().Set("Content-Type", "text/css")
	fmt.Fprintf(w, "%s", body)
}

func visitedHandler(w http.ResponseWriter, r *http.Request){
	//title := r.URL.Path[len("/visited/"):]
	
	//c := database.Get()
	//defer c.Close()
	//items, err := redis.Values(c.Do("smembers", "visited"))
	//if err != nil {
	//fmt.Fprintf(w, "%s", err)
//}
	
	//fmt.Fprintf(w, "%s", items)
}

func loginHandler(w http.ResponseWriter, r *http.Request){
    auth := r.Header["Authorization"][0]
    if strings.Contains(auth, ":") {
        userPass := strings.Split(auth, ":")
        username := userPass[0]
        pass := userPass[1]
        userHash := db.GetUserHash(username)
        if userHash == pass {
            fmt.Println("user authenticated")
            expires := time.Now().Add(24 * time.Hour)
            cookie := http.Cookie{Name:"Id", Value:db.CreateSessionId(username), Expires:expires, Path:"/"}
            http.SetCookie(w, &cookie)
        } else {
            fmt.Fprintf(w, "%s", "no log in for you")

        }
    } else if auth != "" {
        key := db.GetUserKey(auth)
        fmt.Fprintf(w, "%s", key)
    } else {
        fmt.Println("in else")
        fmt.Fprintf(w, "%s", "unable to determine intentions")
    }
	
}

func secureHandler(w http.ResponseWriter, r *http.Request){
    id, _ := r.Cookie("Id")
    username := db.Validate(id.Value)
    if username == "" {
        fmt.Fprintf(w, "%s", "please re-authenticate")
        return
    }
        cookie := http.Cookie{Name:"Id", Value:db.CreateSessionId(username), Expires:time.Now().Add(time.Duration(15)*time.Minute), Path:"/"}
        http.SetCookie(w, &cookie)
    fmt.Println("opening page for " + username)
	title := r.URL.Path[len("/secure/"):]
    body, _ := ioutil.ReadFile(title)
    fmt.Fprintf(w, "%s", body)
    

}

func main(){

	http.HandleFunc("/visited/", visitedHandler)
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/scripts/", scriptHandler)
	http.HandleFunc("/css/", cssHandler)
    http.HandleFunc("/secure/", secureHandler)
	http.HandleFunc("/login/", loginHandler)
	
	http.ListenAndServe(":8090", nil)

}
