package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)
type emp struct{
	Id int `json:"id"`
	Name string `json:"name"`
	Age int `json:"age"`
	Gender string `json:"gender"`
	Role int `json:"jobid"`
	Job string `json:"jobrole"`
}
var Emp emp

func getAll(db *sql.DB,w http.ResponseWriter){
	var res *sql.Rows
	var err error
	empList:=make([]emp,0)
	res, err = db.Query("select E.EmpId,E.name,E.age,E.gender,E.role,J.JobRole from employee as E inner join jobrole as J on E.role=J.JobId")
	defer res.Close()
	if err != nil {
		log.Println(err)
	}
	resCount:=0
	for res.Next() {
		resCount=resCount+1
		err := res.Scan(&Emp.Id, &Emp.Name, &Emp.Age,&Emp.Gender,&Emp.Role,&Emp.Job)
		if err != nil {
			log.Println(err)
		}
		//fmt.Printf("%v\n", Emp)
		//res,_:=json.Marshal(Emp)
		//fmt.Println(string(res))
		empList=append(empList,Emp)
	}
	if resCount==0{
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No Employee details available"))
		return
	}
	jsonRes,_:=json.Marshal(empList)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonRes)
}
func getEmp(db *sql.DB,w http.ResponseWriter,key int){
	var res *sql.Rows
	var err error
	empList:=make([]emp,0)
	res, err = db.Query(fmt.Sprintf("select E.EmpId,E.name,E.age,E.gender,E.role,J.JobRole from employee as E inner join jobrole as J on E.role=J.JobId where E.EmpId=%v",key))
	defer res.Close()
	if err != nil {
		log.Println(err)
	}
	resCount:=0
	for res.Next() {
		resCount=resCount+1
		err := res.Scan(&Emp.Id, &Emp.Name, &Emp.Age,&Emp.Gender,&Emp.Role,&Emp.Job)
		if err != nil {
			log.Println(err)
		}
		//fmt.Printf("%v\n", Emp)
		//res,_:=json.Marshal(Emp)
		//fmt.Println(string(res))
		empList=append(empList,Emp)
	}
	if resCount==0{
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Employee details not found"))
		return
	}
	jsonRes,_:=json.Marshal(empList)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonRes)
}
func postData(db *sql.DB,r *http.Request,w http.ResponseWriter){
	body:=r.Body
	err:=json.NewDecoder(body).Decode(&Emp)
	if err!=nil{
		log.Println(err)
	}
	//fmt.Println(Emp)
	execQuery:=fmt.Sprintf("insert into employee(name,age,gender,role) values('%v',%v,'%v',%v)",Emp.Name,Emp.Age,Emp.Gender,Emp.Role)
	res, err := db.Exec(execQuery)
	if err != nil {
		//panic(err.Error())
		log.Println(err)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Employee added to database"))
	fmt.Printf("The last inserted row id: %d\n", lastId)
	//log.Println(Emp)
}
func updateData(db *sql.DB,r *http.Request,w http.ResponseWriter,key int){
	body:=r.Body
	err:=json.NewDecoder(body).Decode(&Emp)
	//fmt.Println(key)
	if err!=nil{
		log.Println(err)
	}

	//fmt.Println(Emp.Role)
	execQuery:=fmt.Sprintf("update employee set name='%v',age=%v,gender='%v',role=%v where EmpId=%v",Emp.Name,Emp.Age,Emp.Gender,Emp.Role,key)
	_, err = db.Exec(execQuery)
	if err != nil {
		//panic(err.Error())
		log.Println(err)
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Employee data updated in database"))
	//log.Println(Emp)
}
func deleteData(db *sql.DB,r *http.Request,w http.ResponseWriter,key int){
	body:=r.Body
	err:=json.NewDecoder(body).Decode(&Emp)
	if err!=nil{
		log.Println(err)
	}
	//fmt.Println(Emp)
	execQuery:=fmt.Sprintf("delete from employee where EmpId=%v",key)
	res, err:= db.Exec(execQuery)
	if err != nil {
		//panic(err.Error())
		log.Println(err)
	}
	changes,err:=res.RowsAffected()
	if err != nil {
		//panic(err.Error())
		log.Println(err)
	}
	if !(changes>0){
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Employee not found in database"))
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Employee data deleted in database"))
	//log.Println(Emp)
}
func fetchAll(w http.ResponseWriter, r *http.Request){
	db, err := sql.Open("mysql", "nehul:9618181838@tcp(127.0.0.1)/company")
	defer db.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server Error"))
		log.Println(err)
	}
	if r.Method=="GET" {
		getAll(db, w)
	}
	if r.Method=="POST"{
		postData(db,r,w)
	}
	if r.Method=="PUT"{
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Employee Id Missing"))
	}
	if r.Method=="DELETE"{
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Employee Id Missing"))
	}
}
func fetchEmp(w http.ResponseWriter, r *http.Request){
	db, err := sql.Open("mysql", "nehul:9618181838@tcp(127.0.0.1)/company")
	defer db.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server Error"))
		log.Println(err)
	}

	vars := mux.Vars(r)
	key,err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Employee Id Missing"))
		log.Println(err)
	}
	if r.Method=="GET" {
		getEmp(db, w, key)
	}
	if r.Method=="PUT"{
		updateData(db,r,w,key)
	}
	if r.Method=="DELETE"{
		deleteData(db,r,w,key)
	}
	if r.Method=="POST" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Employee not created"))
	}
	//fmt.Fprintf(w, "Key: " + string(key))
}

func main(){
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/employee",fetchAll)
	myRouter.HandleFunc("/employee/{id}", fetchEmp)
	http.ListenAndServe(":8080",myRouter)
}