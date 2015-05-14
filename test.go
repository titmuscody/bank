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
	fmt.Fprintf(w, "hi there %s", r.URL.Path[1:])
	fmt.Println("printing headers")
	fmt.Println(r.Header)
	//r.ParseForm()
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Println("printing body")
	fmt.Println(string(body))
	//fmt.Println(r.Form)
	//fmt.Println(r.Body)
	fmt.Println("printing map")
	fmt.Println(r.Form)
	store := &Operation{}
	if string(body) != ""{
		m := string(body)
		json.Unmarshal([]byte(m), &store)
		//fmt.Println(k)
		//fmt.Println(store)
		
		sum := 0
		for _, e := range store.Add{
			sum += e
			//fmt.Println(e)
		}
		//fmt.Println(sum)
		fmt.Fprintf(w, "the result was %d", sum)
	}else{
		m := r.Form
			for k := range m {
				json.Unmarshal([]byte(k), &store)
				//fmt.Println(k)
				//fmt.Println(store)
				
				sum := 0
				for _, e := range store.Add{
					sum += e
					//fmt.Println(e)
				}
				//fmt.Println(sum)
				fmt.Fprintf(w, "the result was %d", sum)
			}
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request){
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	//fmt.Fprintf(w, "header=%s body=%s", p.Title, p.Body)
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
