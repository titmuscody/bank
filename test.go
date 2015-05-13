package main

import (
	"fmt"
	"net/http" 
	"io/ioutil"
	"encoding/json"
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
	filename := title + ".txt"
	body, _ := ioutil.ReadFile(filename)
	return &Page{Title:title, Body:body}, nil
}

func handler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "hi there %s", r.URL.Path[1:])
	r.ParseForm()
	//fmt.Println(r.Form)
	//fmt.Println(r.Body)
	m := r.Form
	fmt.Println(m)
	store := &Operation{}
	for k := range m {
		json.Unmarshal([]byte(k), &store)
		fmt.Println(k)
		fmt.Println(store)
		
		sum := 0
		for _, e := range store.Add{
			sum += e
			fmt.Println(e)
		}
		fmt.Println(sum)
		fmt.Fprintf(w, "the result was %d", sum)
	}

	//fmt.Println(m)
}

func viewHandler(w http.ResponseWriter, r *http.Request){
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	//fmt.Fprintf(w, "header=%s body=%s", p.Title, p.Body)
	fmt.Fprintf(w, "%s", p.Body)
}

func main(){
	p1 := &Page{Title:"testpage", Body: []byte("this is the test")}
	p1.save()
	p2, _  := loadPage("testpage")
	fmt.Println(string(p2.Body))
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
	http.ListenAndServe(":8090", nil)
	fmt.Println("hello world")

}
