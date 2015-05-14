package main

import (
	"fmt"
	"net/http" 
	"io/ioutil"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	)

type Page struct{
	Title string
	Body []byte
}

type Operation struct {
	Add []int 
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
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
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	fmt.Fprintf(w, "%s", p.Body)
}

func visitedHandler(w http.ResponseWriter, r *http.Request){
	//title := r.URL.Path[len("/visited/"):]
	c, err := redis.Dial("tcp", ":6379") 
	if err != nil {
	panic(err)
}
	defer c.Close()
	
	items, err := redis.Values(c.Do("smembers", "visited"))
	if err != nil {
	fmt.Fprintf(w, "%s", err)
}
	
	fmt.Fprintf(w, "%s", items)
}


func main(){
	
	

	http.HandleFunc("/visited/", visitedHandler)
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
	http.ListenAndServe(":8090", nil)

}
