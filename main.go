package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type studentinfo struct {
	Sid    string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Course string `json:"course,omitempty"`
}

func getMYsqlDB() *sql.DB {
	db, err := sql.Open("mysql", "root:@(127.0.0.1:3306)/studentinfo?parseTime=true")

	if err != nil {
		log.Fatal(err)
	}
	return db
}
func getstudent(w http.ResponseWriter, r *http.Request) {
	db := getMYsqlDB()
	defer db.Close()
	ss := []studentinfo{}
	s := studentinfo{}
	rows, err := db.Query("select * from student")
	if err != nil {
		fmt.Fprint(w, err)
	} else {
		for rows.Next() {
			rows.Scan(&s.Sid, &s.Name, &s.Course)
			ss = append(ss, s)
		}
		json.NewEncoder(w).Encode(ss)
	}

	fmt.Fprint(w, "GET the student")
}
func deletestudent(w http.ResponseWriter, r *http.Request) {
	db := getMYsqlDB()
	defer db.Close()
	parms := mux.Vars(r)
	sid, _ := strconv.Atoi(parms["sid"])
	result, err := db.Exec("delete from student where sid=?", sid)
	if err != nil {
		fmt.Fprint(w, err)
	} else {
		_, err := result.RowsAffected()
		if err != nil {
			fmt.Print(w, err)
		} else {
			fmt.Fprint(w, "Suscesssful")
		}
	}
	fmt.Fprint(w, "delete the student")
}
func updatestudent(w http.ResponseWriter, r *http.Request) {
	db := getMYsqlDB()
	defer db.Close()
	s := studentinfo{}
	json.NewDecoder(r.Body).Decode(&s)
	parms := mux.Vars(r)
	sid, _ := strconv.Atoi(parms["sid"])
	result, err := db.Exec("update student set name=?,course=?,where sid=?", sid, s.Name, s.Course)
	if err != nil {
		fmt.Fprint(w, err)
	} else {
		_, err := result.RowsAffected()
		if err != nil {
			json.NewEncoder(w).Encode("error")
		} else {
			json.NewEncoder(w).Encode(s)
		}
	}
	fmt.Fprint(w, "update the student")
}
func addstudent(w http.ResponseWriter, r *http.Request) {
	db := getMYsqlDB()
	defer db.Close()
	s := studentinfo{}
	json.NewDecoder(r.Body).Decode(s)
	sid, _ := strconv.Atoi(s.Sid)
	result, err := db.Exec("insert into student(sid, name ,course)value(?,?,?)", sid, s.Name, s.Course)
	if err != nil {
		fmt.Fprint(w, err)
	} else {
		_, err := result.LastInsertId()
		if err != nil {
			json.NewEncoder(w).Encode("{No record is inserted}")
		} else {
			json.NewEncoder(w).Encode(s)

		}
	}
	fmt.Fprint(w, "add the student")
}
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/students", getstudent).Methods("GET")
	r.HandleFunc("/students", addstudent).Methods("POST")
	r.HandleFunc("/students/{Sid}", updatestudent).Methods("PUT")
	r.HandleFunc("/students/{Sid}", deletestudent).Methods("DELETE")
	http.ListenAndServe(":8080", r)

}
