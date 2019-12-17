// Package classification Petstore API.
//
// the purpose of this application is to provide an application
// that is using plain go code to define an API
//
// This should demonstrate all the possible comment annotations
// that are available to turn go code into a fully compliant swagger 2.0 spec
//
//     Schemes: http, https
//     Host: localhost:12345
//
//     Consumes:
//     - application/json
//     - application/xml
//
//     Produces:
//     - application/json
//     - application/xml
//
//
// swagger:meta
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	session, err := mgo.Dial("localhost/muxgocrud")
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB("muxgocrud")
}

// Employee represents body of employee response.
type Employee struct {
	ID        bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string        `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string        `json:"lastname,omitempty" bson:"lastname,omitempty"`
	EmpID     int           `json:"empid,omitempty" bson:"empid,omitempty"`
	Salary    float64       `json:"salary,omitempty" bson:"salary,omitempty"`
	Practice  string        `json:"practice,omitempty" bson:"practice,omitempty"`
}

// EmployeeCollection holds collection of emp records and total count.
type EmployeeCollection struct {
	AllEmployees []Employee `json:"employees"`
	Count        int        `json:"count"`
}

var mockMarshal = json.Marshal

// CreateEmployeeEndpoint creates an employee record.
func CreateEmployeeEndpoint(response http.ResponseWriter, request *http.Request) {

	// swagger:operation POST /employees CreateEmployeeEndpoint
	//
	// Creates an employee record.
	// ---
	// produces:
	// - application/json
	// consumes:
	// - application/json
	// parameters:
	// - in: body
	//   name: employee
	//   description: The employee to create.
	//   schema:
	//    type: object
	//    required:
	//     - firstname
	//     - lastname
	//     - empid
	//     - salary
	//     - practice
	//   properties:
	//     firstname:
	//	    type: string
	//     lasttname:
	//	    type: string
	//     empid:
	//	    type: integer
	//     salary:
	//	    type: number
	//     practice:
	//	    type: string
	// responses:
	//   '201':
	//     description: employee response
	//   '500':
	//     description: internal server error
	//   default:
	//     description: unexpected error

	setResponseHeader(response)
	var employee Employee
	json.NewDecoder(request.Body).Decode(&employee)
	collection := db.C("employee")
	err := collection.Insert(employee)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	result, err := mockMarshal(&employee)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
	}
	response.WriteHeader(http.StatusCreated)
	response.Write(result)
}

// GetEmployeesEndpoint returns list of employees.
func GetEmployeesEndpoint(response http.ResponseWriter, request *http.Request) {

	// swagger:operation GET /employees GetEmployeesEndpoint
	//
	//  Get specific employee record.
	//	Set response headers.
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     description: employee response
	//   '500':
	//     description: internal server error
	//   default:
	//     description: unexpected error

	setResponseHeader(response)
	var employees []Employee
	collection := db.C("employee")
	limit, _ := strconv.Atoi(request.FormValue("limit"))
	page, _ := strconv.Atoi(request.FormValue("page"))
	skips := limit * (page - 1)
	collection.Find(nil).Limit(limit).Skip(skips).All(&employees)
	employeeCollection := EmployeeCollection{AllEmployees: employees, Count: len(employees)}
	result, err := mockMarshal(employeeCollection)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return

	}
	response.Write(result)
}

//GetEmployeeEndpoint returns single employee record.
func GetEmployeeEndpoint(response http.ResponseWriter, request *http.Request) {

	// swagger:operation GET /employee/{id} GetEmployeeEndpoint
	//
	//  Get specific employee record.
	//	Set response headers.
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: path
	//   description: primitive id
	//   required: true
	// responses:
	//   '200':
	//     description: employee response
	//   '404':
	//     description: not found
	//   '500':
	//     description: internal server error
	//   default:
	//     description: unexpected error

	setResponseHeader(response)
	params := mux.Vars(request)
	id := bson.ObjectIdHex(params["id"])
	var employee Employee
	collection := db.C("employee")
	err := collection.Find(bson.M{"_id": id}).One(&employee)
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	result, er := mockMarshal(&employee)
	if er != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + er.Error() + `" }`))
		return
	}
	response.Write(result)
}

// UpdateEmployeeEndpoint updates given employee record.
func UpdateEmployeeEndpoint(response http.ResponseWriter, request *http.Request) {

	// swagger:operation PUT /employee/{id} UpdateEmployeeEndpoint
	//
	//  Update specific employee record.
	//	Set response headers.
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: path
	//   description: primitive id
	//   required: true
	// - in: body
	//   name: employee
	//   description: The employee to create.
	//   schema:
	//    type: object
	//    required:
	//     - firstname
	//     - lastname
	//     - empid
	//     - salary
	//     - practice
	//   properties:
	//     firstname:
	//	    type: string
	//     lasttname:
	//	    type: string
	//     empid:
	//	    type: integer
	//     salary:
	//	    type: number
	//     practice:
	//	    type: string
	// responses:
	//   '200':
	//     description: employee response
	//   default:
	//     description: unexpected error

	setResponseHeader(response)
	params := mux.Vars(request)
	id := bson.ObjectIdHex(params["id"])
	var employee Employee
	err := json.NewDecoder(request.Body).Decode(&employee)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	collection := db.C("employee")
	err = collection.Update(bson.M{"_id": id}, bson.M{"$set": employee})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	result, err := mockMarshal(&employee)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	response.Write(result)
}

// DeleteEmployeeEndpoint deletes given employee record.
func DeleteEmployeeEndpoint(response http.ResponseWriter, request *http.Request) {

	// swagger:operation DELETE /employee/{id} UpdateEmployeeEndpoint
	//
	//  Delete specific employee record.
	//	Set response headers.
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: path
	//   description: primitive id
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//     description: employee response
	//   default:
	//     description: unexpected error

	setResponseHeader(response)
	params := mux.Vars(request)
	id := bson.ObjectIdHex(params["id"])
	collection := db.C("employee")
	err := collection.Remove(bson.M{"_id": id})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	response.Write([]byte("Employee deleted successfully."))
}

var headers = handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
var methods = handlers.AllowedMethods([]string{"GET", "PUT", "POST", "DELETE", "OPTIONS", "HEAD"})
var origins = handlers.AllowedOrigins([]string{"*"})

// DefineRoute : collection of all routes.
func DefineRoute() {
	fmt.Println("Starting the application...")
	router.HandleFunc("/employees", CreateEmployeeEndpoint).Methods("POST")
	router.HandleFunc("/employees", GetEmployeesEndpoint).Methods("GET")
	router.HandleFunc("/employee/{id}", GetEmployeeEndpoint).Methods("GET")
	router.HandleFunc("/employee/{id}", UpdateEmployeeEndpoint).Methods("PUT")
	router.HandleFunc("/employee/{id}", DeleteEmployeeEndpoint).Methods("DELETE")
}

// The main function.
func main() {
	DefineRoute()
	http.ListenAndServe(":12345", handlers.CORS(headers, methods, origins)(router))
}
