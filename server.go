package main

import (
	"fmt"
	"time"
	"net/http" 
	"io/ioutil"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"crypto/sha512"
    "crypto/hmac"
	"math/rand"
	"encoding/hex"
    //"strconv"
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
}

func viewHandler(w http.ResponseWriter, r *http.Request){
	expires := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name:"username", Value:"me", Expires:expires,
		HttpOnly:false}
	http.SetCookie(w, &cookie)
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

func visitedHandler(w http.ResponseWriter, r *http.Request){
	//title := r.URL.Path[len("/visited/"):]
	
	c := database.Get()
	defer c.Close()
	items, err := redis.Values(c.Do("smembers", "visited"))
	if err != nil {
	fmt.Fprintf(w, "%s", err)
}
	
	fmt.Fprintf(w, "%s", items)
}

func loginHandler(w http.ResponseWriter, r *http.Request){
	//mess := "my secret code"
	//pass := r.Header["Authorization"]
	//cookie := http.Cookie{Name:"username", Value:"me", Expires:time.Now().Add(364),}
	//http.SetCookie(w, &cookie)
    auth := r.Header["Authorization"][0]
	conn := database.Get()
	defer conn.Close()
    fmt.Printf("auth=%s\n", auth)
    if strings.Contains(auth, ":") {
        fmt.Println("in primary")

        userPass := strings.Split(auth, ":")
        username := userPass[0]
        pass := userPass[1]
        user := strings.Join([]string{"user", username}, ":")
        myKey, _ := conn.Do("hget", user, "key")
        myPass, _ := conn.Do("hget", user, "password")
        fmt.Println(string(myKey.([]uint8)))
        fmt.Println(string(myPass.([]uint8)))
        hash := hmac.New(sha512.New, []byte(string(myKey.([]uint8))))
        hash.Write([]byte(myPass.([]uint8)))
        myHash := hex.EncodeToString(hash.Sum(nil))
        fmt.Println(myHash)
        conn.Do("hset", user, "key", rand.Int())
        if myHash == pass {
            fmt.Fprintf(w, "%s", "you have logged in")

        } else {
            fmt.Fprintf(w, "%s", "failure no log in for you")

        }
    } else if auth != "" {
        //get number and send
        fmt.Println("in auth")
        user := strings.Join([]string{"user", auth}, ":")
        key, _ := conn.Do("hget", user, "key")
        //fmt.Println(string(key.([]uint8)))
        //key_val := strconv.Itoa(int(key.(uint8[])[0]))
        //key_val := key.(string)
        fmt.Fprintf(w, "%s", string(key.([]uint8)))
        
    } else {
        fmt.Println("in else")
        fmt.Fprintf(w, "%s", "unable to determine intentions")
    }
    //io.WriteString(hash, "password")
	//hash.Write(data)
    //fmt.Println(hex.EncodeToString(hash))
    
    // correct= fmt.Println(hex.EncodeToString(hash.Sum(nil)))
	
	
}

func newPool() *redis.Pool {

	return &redis.Pool{
		MaxIdle:3,
		IdleTimeout:20* time.Second,
		Dial: func() (redis.Conn, error){
		c, err := redis.Dial("tcp", ":6379")
		if err != nil{panic(err.Error())}
		return c, err	
			
		},

	}
}

var database = newPool()
func main(){

	http.HandleFunc("/visited/", visitedHandler)
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/scripts/", scriptHandler)
	http.HandleFunc("/login/", loginHandler)
	
	http.ListenAndServe(":8090", nil)

}
