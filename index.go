package main

import (
	"appengine"
	"appengine/datastore"
	"html/template"
	"log"
	"net/http"
)

type Person struct {
	ID     *datastore.Key
	Name   string
	Animal string
}
type People []Person

func init() {

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/enroll", enroll)
	http.HandleFunc("/", index)

	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("main").ParseFiles("templates/main.tmpl", "templates/index.tmpl")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	data := make(map[string]interface{})
	data["people"], err = getAllPeople(r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = t.ExecuteTemplate(w, "main", data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func signup(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("main").ParseFiles("templates/main.tmpl", "templates/signup.tmpl")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = t.ExecuteTemplate(w, "main", nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func enroll(w http.ResponseWriter, r *http.Request) {
	log.Print("GO")
	name := r.FormValue("name")
	animal := r.FormValue("animal")
	var p Person
	p.Name = name
	p.Animal = animal
	c := appengine.NewContext(r)
	c.Infof("%v doo ", p)
	key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Person", nil), &p)
	if err != nil {
		c.Errorf("%v doo ", err)
		return
	}
	p.ID = key

	http.Redirect(w, r, "/", 301)

}

func getAllPeople(r *http.Request) (People, error) {
	var err error
	var ps People

	c := appengine.NewContext(r)
	_, err = datastore.NewQuery("Person").GetAll(c, &ps)
	if err != nil {
		return ps, err
	}

	return ps, err
}
