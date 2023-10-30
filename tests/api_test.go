package main

import (
	"bytes"
	"cloudflare-takehome/routes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var rt = routes.Routes()

func TestPingEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	rt.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "handler returned wrong status code")
	assert.Contains(t, "pong", rec.Body.String(), "returned unexpected body")
}

func TestCreateURLEndpoint(t *testing.T) {
	jsonReq := []byte(`{
		"original_url": "cloudflare.com"
	}`)

	req, err := http.NewRequest("POST", "/create", bytes.NewBuffer(jsonReq))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	rt.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "handler returned wrong status code")
	assert.Contains(t, rec.Body.String(), "data", "returned unexpected body")
}

func TestCreateURLEndpointWithExpiration(t *testing.T) {

	jsonReq := []byte(`{
		"original_url": "cloudflare.com",
		"expiration_date": "2023-10-30T12:04:05Z"
	}`)

	req, err := http.NewRequest("POST", "/create", bytes.NewBuffer(jsonReq))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	rt.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "handler returned wrong status code")
	assert.Contains(t, rec.Body.String(), "data", "returned unexpected body")
}

func TestCreateURLEndpointNoURL(t *testing.T) {
	jsonReq := []byte(`{
		"original_url": ""
	}`)

	req, err := http.NewRequest("POST", "/create", bytes.NewBuffer(jsonReq))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	rt.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code, "handler returned wrong status code")
}

func TestGoToURLEndpoint(t *testing.T) {
	//Create new url
	jsonReq := []byte(`{
		"original_url": "cloudflare.com"
	}`)

	req, err := http.NewRequest("POST", "/create", bytes.NewBuffer(jsonReq))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	rt.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "handler returned wrong status code")
	assert.Contains(t, rec.Body.String(), "data", "returned unexpected body")

	var jsonResp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &jsonResp)
	if err != nil {
		t.Fatal(err)
	}

	//Attempt to go to newly created url
	data := jsonResp["data"].(map[string]interface{})
	req, err = http.NewRequest("GET", data["short_url"].(string), nil)
	if err != nil {
		t.Fatal(err)
	}

	rec = httptest.NewRecorder()
	rt.ServeHTTP(rec, req)

	assert.Equal(t, "http://cloudflare.com", rec.Result().Header.Get("Location"))
}

func TestDeleteURLEndpoint(t *testing.T) {
	//Create new url
	jsonReq := []byte(`{
		"original_url": "cloudflare.com"
	}`)

	req, err := http.NewRequest("POST", "/create", bytes.NewBuffer(jsonReq))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	rt.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "handler returned wrong status code")
	assert.Contains(t, rec.Body.String(), "data", "returned unexpected body")

	var jsonResp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &jsonResp)
	if err != nil {
		t.Fatal(err)
	}

	//Delete url we just created
	data := jsonResp["data"].(map[string]interface{})
	req, err = http.NewRequest("DELETE", "/delete/"+data["key"].(string), nil)
	if err != nil {
		t.Fatal(err)
	}

	rec = httptest.NewRecorder()
	rt.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "handler returned wrong status code")
	assert.Contains(t, rec.Body.String(), "deleted", "returned unexpected body")
}

func TestURLStatsEndpoint(t *testing.T) {
	//Create new url
	jsonReq := []byte(`{
		"original_url": "cloudflare.com"
	}`)

	req, err := http.NewRequest("POST", "/create", bytes.NewBuffer(jsonReq))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	rt.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "handler returned wrong status code")
	assert.Contains(t, rec.Body.String(), "data", "returned unexpected body")

	var jsonResp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &jsonResp)
	if err != nil {
		t.Fatal(err)
	}

	//Access new short url
	data := jsonResp["data"].(map[string]interface{})
	req, err = http.NewRequest("GET", data["short_url"].(string), nil)
	if err != nil {
		t.Fatal(err)
	}

	rec = httptest.NewRecorder()
	rt.ServeHTTP(rec, req)

	//Get stats
	req, err = http.NewRequest("GET", "/stats/"+data["key"].(string), nil)
	if err != nil {
		t.Fatal(err)
	}

	rec = httptest.NewRecorder()
	rt.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "handler returned wrong status code")
	assert.Contains(t, rec.Body.String(), "\"all_time\":1", "unexpected all time count")
	assert.Contains(t, rec.Body.String(), "\"past_day\":1", "unexpected past day count")
	assert.Contains(t, rec.Body.String(), "\"past_week\":1", "unexpected past week count")
}
