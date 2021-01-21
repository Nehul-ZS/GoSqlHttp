package main

import (
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http/httptest"
	"testing"
	//"net/http"
	"bytes"
	//"fmt"
)
type httpreq struct{
reqMethod string
reqPath string
reqBody []byte
expBody string
expCode int
}
func TestAccessData(t *testing.T){
	testcases:=[] httpreq{
		{"GET","/hello",nil,"404 page not found\n",404},
		{"GET","/employee",nil,`[{"id":1,"name":"Nehul","age":21,"gender":"M","jobid":1,"jobrole":"SDE Intern"},{"id":2,"name":"Vipul","age":22,"gender":"M","jobid":2,"jobrole":"GoLang Intern"},{"id":3,"name":"Shiva","age":22,"gender":"M","jobid":1,"jobrole":"SDE Intern"},{"id":4,"name":"Vaishnav","age":23,"gender":"M","jobid":3,"jobrole":"Backend Developer"}]`,200},
		{"GET","/employee/1",nil,`[{"id":1,"name":"Nehul","age":21,"gender":"M","jobid":1,"jobrole":"SDE Intern"}]`,200},
		{"GET","/employee/2",nil,`[{"id":2,"name":"Vipul","age":22,"gender":"M","jobid":2,"jobrole":"GoLang Intern"}]`,200},
		{"GET","/employee/8",nil,"Employee details not found",400},
		{"POST","/hello",[]byte(`{"name":"Roy","age":20,"gender":"M","jobid":2}`),"404 page not found\n",404},
		{"POST","/employee",[]byte(`{"name":"Sid","age":23,"gender":"M","jobid":3}`),"Employee added to database",201},
		{"POST","/employee/1",nil,"Employee not created",400},
		{"PUT","/employee/1",[]byte(`{"id":1,"name":"Nehul","age":22,"gender":"M","jobid":3}`),"Employee data updated in database",202},
		//{"PUT","/employee/8",[]byte(`{"id":1,"name":"Nehul","age":22,"gender":"M","jobid":3}`),"Employee not found in database",404},
		{"PUT","/hello",[]byte(`{"name":"Cap","age":21,"gender":"M","jobid":1}`),"404 page not found\n",404},
		{"PUT","/employee",[]byte(`{"name":"Rio","age":24,"gender":"M","jobid":2}`),"Employee Id Missing",400},
		{"DELETE","/employee",nil,"Employee Id Missing",400},
		{"DELETE","/hello",nil,"404 page not found\n",404},
		{"DELETE","/employee/1",nil,"Employee data deleted in database",202},
		{"DELETE","/employee/8",nil,"Employee not found in database",404},

	}
	for i,testcase:=range testcases{
		w := httptest.NewRecorder()
		r:= httptest.NewRequest(testcase.reqMethod,testcase.reqPath,bytes.NewBuffer(testcase.reqBody))
		myRouter := mux.NewRouter().StrictSlash(true)
		myRouter.HandleFunc("/employee/{id}", fetchEmp)
		myRouter.HandleFunc("/employee", fetchAll)
		myRouter.ServeHTTP(w, r)
		//handler := http.HandlerFunc(fetchAll)
		//handler.ServeHTTP(w, r)
		res:=w.Result()
		resBody,resErr:=ioutil.ReadAll(res.Body)
		resCode:=w.Code
		//fmt.Println(resBody)
		if resCode!=testcase.expCode{
			t.Errorf("Test %v has failed, Expected status code: %v but got %v",i,testcase.expCode,resCode)
		}
		if string(resBody)!=testcase.expBody{
			t.Errorf("Test %v has failed, Expected: %v but got %v",i,testcase.expBody,string(resBody))
		}
		if resErr!=nil{
			t.Errorf("Expected nil error but got %v",resErr)
		}
	}
}