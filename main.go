package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Course struct {
	Coursename string  `json:"course name"`
	CourseId   string  `json:"course id"`
	Author     *Author `json:"Author"`
}

type Author struct {
	Fullname string `json:"full name"`
}

var courses []Course

func (c *Course) IsEmpty() bool {
	return c.Coursename == ""

}

func main() {
	r := mux.NewRouter()
	courses = append(courses, Course{CourseId: "2", Coursename: "react", Author: &Author{Fullname: "pukar"}})
	r.HandleFunc("/", serverhome).Methods("GET")
	r.HandleFunc("/all", getallcourse).Methods("GET")
	r.HandleFunc("/addp", addcourse).Methods("POST")
	r.HandleFunc("/update/{id}", getonecourse).Methods("GET")
	r.HandleFunc("/update/{id}", update).Methods("PUT")
	r.HandleFunc("/delete/{id}", deleteonecourse).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":4000", r))

}

func serverhome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>hello world<h1>"))
}

func getallcourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contant-type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func getonecourse(w http.ResponseWriter, r *http.Request) {
	fmt.Print("get one course")
	w.Header().Set("Contant-type", "application/json")
	params := mux.Vars(r)
	fmt.Println(params)
	for _, course := range courses {
		if course.CourseId == params["id"] {
			json.NewEncoder(w).Encode(course)
			return
		}
	}
	json.NewEncoder(w).Encode("NO course is found in that id")
	return
}
func addcourse(w http.ResponseWriter, r *http.Request) {
	fmt.Print("create one course")
	w.Header().Set("Contant-type", "application/json")
	if r.Body == nil {
		json.NewEncoder(w).Encode("please insert a value")
	}
	var course Course
	_ = json.NewDecoder(r.Body).Decode(&course)
	if course.IsEmpty() {
		json.NewEncoder(w).Encode("please insert a value")
		return
	}
	rand.Seed(time.Now().UnixNano())
	course.CourseId = strconv.Itoa(rand.Intn(100))

	courses = append(courses, course)
	json.NewEncoder(w).Encode(course)

}
func update(w http.ResponseWriter, r *http.Request) {
	fmt.Print("update one course")
	w.Header().Set("Contant-type", "application/json")
	params := mux.Vars(r)
	for index, course := range courses {
		if course.CourseId == params["id"] {
			courses = append(courses[:index], courses[index+1:]...)
			var course Course
			_ = json.NewDecoder(r.Body).Decode(&course)
			course.CourseId = params["id"]
			courses = append(courses, course)
			json.NewEncoder(w).Encode(course)
			return
		}
	}
}

func deleteonecourse(w http.ResponseWriter, r *http.Request) {
	fmt.Print("delete one course")
	w.Header().Set("Contant-type", "application/json")
	params := mux.Vars(r)
	for index, course := range courses {
		if course.CourseId == params["id"] {
			courses = append(courses[:index], courses[index+1:]...)
			break

		}
	}
}
