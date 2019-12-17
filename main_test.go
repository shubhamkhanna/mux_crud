package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

type VarMock struct {
	err error
}

func (m *VarMock) dummyMarshal(v interface{}) ([]byte, error) {
	return nil, m.err
}

func TestMuxCrudAPI(t *testing.T) {
	t.Run("TestUpdateEmployeeEndpoint", func(t *testing.T) {
		testUpdateEmployeeEndpoint(t)
	})

	t.Run("TestDeleteEmployeeEndpoint", func(t *testing.T) {
		testDeleteEmployeeEndpoint(t)
	})

	t.Run("TestGetEmployeeEndpoint", func(t *testing.T) {
		testGetEmployeeEndpoint(t)
	})

	t.Run("TestGetEmployeesEndpoint", func(t *testing.T) {
		testGetEmployeesEndpoint(t)
	})

	t.Run("TestCreateEmployeesEndpoint", func(t *testing.T) {
		testCreateEmployeesEndpoint(t)
	})

}
func testGetEmployeeEndpoint(t *testing.T) {
	collection := db.C("employee")
	var existingEmployee Employee
	collection.Find(bson.M{}).One(&existingEmployee)
	req, err := http.NewRequest("GET", "/employee/"+existingEmployee.ID.Hex(), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	router.HandleFunc("/employee/{id}", GetEmployeeEndpoint)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected, _ := json.Marshal(existingEmployee)
	if rr.Body.String() != string(expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), string(expected))
	}

	t.Run("It returns 404", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/employee/000000000000000000000000", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusNotFound)
		}
	})

	t.Run("it mocks marshal error", func(t *testing.T) {
		m := &VarMock{err: errors.New("failed")}
		mockMarshal = m.dummyMarshal
		req, _ := http.NewRequest("GET", "/employee/"+existingEmployee.ID.Hex(), nil)
		req.Header.Set("Content-Type", "application/json")
		router.HandleFunc("/employee/{id}", GetEmployeeEndpoint)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})

	defer func() {
		mockMarshal = json.Marshal
	}()
}

func testGetEmployeesEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/employees", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetEmployeesEndpoint)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var employeeCollection EmployeeCollection
	json.Unmarshal([]byte(rr.Body.String()), &employeeCollection)

	if employeeCollection.Count == 0 {
		t.Errorf("Data collection not found!")
	}

	t.Run("it tests pagination", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/employees?limit=1&page=1", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetEmployeesEndpoint)
		handler.ServeHTTP(rr, req)
		var employeeCollection EmployeeCollection
		json.Unmarshal([]byte(rr.Body.String()), &employeeCollection)
		if employeeCollection.Count != 1 {
			t.Errorf("Data collection not found!")
		}

	})

	t.Run("it mocks marshal error", func(t *testing.T) {
		m := &VarMock{err: errors.New("failed")}
		mockMarshal = m.dummyMarshal
		req, _ := http.NewRequest("GET", "/employees", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetEmployeesEndpoint)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})

	defer func() {
		mockMarshal = json.Marshal
	}()
}

func testCreateEmployeesEndpoint(t *testing.T) {

	payload := []byte(`{
	    "firstname": "aditi",
	    "lastname": "patil",
	    "empid": 1200,
	    "salary": 20000,
	    "practice": "IBM"
	}`)
	collection := db.C("employee")
	count, _ := collection.Count()
	req, err := http.NewRequest("POST", "/employees", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateEmployeeEndpoint)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
	createdCount, _ := collection.Count()
	if createdCount == count {
		t.Errorf("New record is not created.")
	}

	t.Run("It returns internal server error", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/employees", bytes.NewBuffer(payload))
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CreateEmployeeEndpoint)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})

	t.Run("it mocks marshal error", func(t *testing.T) {
		payload := []byte(`{
			"firstname": "aditi",
			"lastname": "patil",
			"empid": 1500,
			"salary": 20000,
			"practice": "IBM"
		}`)
		m := &VarMock{err: errors.New("failed")}
		mockMarshal = m.dummyMarshal
		req, _ := http.NewRequest("POST", "/employees", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CreateEmployeeEndpoint)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})

	defer func() {
		collection.Remove(bson.M{"empid": 1200})
		collection.Remove(bson.M{"empid": 1500})
		mockMarshal = json.Marshal
	}()
}

func testUpdateEmployeeEndpoint(t *testing.T) {
	payload := []byte(`{
	    "firstname": "updated_firstname",
	    "lastname": "updated_lastname",
	    "empid": 1400,
	    "salary": 20000,
	    "practice": "IBM"
	}`)

	newRecord := Employee{
		bson.NewObjectId(),
		"new_firstname",
		"new_lastname",
		1000,
		20000,
		"IBM",
	}

	collection := db.C("employee")
	collection.Insert(newRecord)
	var existingEmployee Employee
	count, _ := collection.Count()
	collection.Find(nil).Skip(count - 1).One(&existingEmployee)
	req, err := http.NewRequest("PUT", "/employee/"+existingEmployee.ID.Hex(), bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	router.HandleFunc("/employee/{id}", UpdateEmployeeEndpoint).Methods("PUT")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	var updatedEmployee Employee
	collection.Find(nil).Skip(count - 1).One(&updatedEmployee)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	assert.Equal(t, updatedEmployee.Firstname, "updated_firstname")

	t.Run("It returns internal server error for invalid payload", func(t *testing.T) {
		req.Header.Set("Content-Type", "application/json")
		router.HandleFunc("/employee/000000000000000000000000", UpdateEmployeeEndpoint).Methods("PUT")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})

	t.Run("It returns internal server error in records update", func(t *testing.T) {
		req.Header.Set("Content-Type", "application/json")
		req, _ := http.NewRequest("PUT", "/employee/000000000000000000000000", bytes.NewBuffer(payload))
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})

	t.Run("it mocks marshal error", func(t *testing.T) {
		m := &VarMock{err: errors.New("failed")}
		mockMarshal = m.dummyMarshal
		req, _ := http.NewRequest("PUT", "/employee/"+existingEmployee.ID.Hex(), bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		router.HandleFunc("/employee/{id}", UpdateEmployeeEndpoint).Methods("PUT")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})

	defer func() {
		collection.Remove(bson.M{"empid": 1400})
		mockMarshal = json.Marshal
	}()
}

func testDeleteEmployeeEndpoint(t *testing.T) {
	newRecord := Employee{
		bson.NewObjectId(),
		"new_firstname",
		"new_lastname",
		2000,
		20000,
		"IBM",
	}
	collection := db.C("employee")
	collection.Insert(newRecord)
	var existingEmployee Employee
	count, _ := collection.Count()
	collection.Find(nil).Skip(count - 1).One(&existingEmployee)
	req, err := http.NewRequest("DELETE", "/employee/"+existingEmployee.ID.Hex(), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	router.HandleFunc("/employee/{id}", DeleteEmployeeEndpoint).Methods("DELETE")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	deletedCount, _ := collection.Count()
	if deletedCount == count {
		t.Errorf("Record is not deleted.")
	}

	t.Run("It returns internal server error in records delete", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/employee/000000000000000000000000", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})

	defer func() {
		collection.Remove(bson.M{"empid": 2000})
	}()
}

func TestDefineRoute(t *testing.T) {
	DefineRoute()
	if router == nil {
		t.Errorf("Error in defining router!!!")
	}
}
