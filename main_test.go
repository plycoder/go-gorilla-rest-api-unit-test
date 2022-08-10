package main

import (
    "os"
    "testing"   
    "log"

    "net/http"
    "net/http/httptest"
    "bytes"
    "encoding/json"
    "strconv"
)

var a App

func TestMain(m *testing.M) {
    a.Initialize("root","","127.0.0.1","3306","go-gorilla-rest-api-swagger")
    ensureTableExists()
    code := m.Run()
    clearTable()
    os.Exit(code)
}

func ensureTableExists() {
    if _, err := a.DB.Exec(tableCreationQuery); err != nil {
        log.Fatal(err)
    }
}

func clearTable() {
    a.DB.Exec("truncate articles")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS articles (
  id bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  title text NOT NULL,
  authors varchar(300) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY id (id)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci`


func TestEmptyTable(t *testing.T) {
    clearTable()

    req, _ := http.NewRequest("GET", "/articles", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)

    if body := response.Body.String(); body != "[]" {
        t.Errorf("Expected an empty array. Got %s", body)
    }
}


func executeRequest(req *http.Request) *httptest.ResponseRecorder {
    rr := httptest.NewRecorder()
    a.Router.ServeHTTP(rr, req)

    return rr
}


func checkResponseCode(t *testing.T, expected, actual int) {
    if expected != actual {
        t.Errorf("Expected response code %d. Got %d\n", expected, actual)
    }
}


func TestGetNonExistentArticle(t *testing.T) {
    clearTable()

    req, _ := http.NewRequest("GET", "/article/11", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusNotFound, response.Code)

    var m map[string]string
    json.Unmarshal(response.Body.Bytes(), &m)
    if m["error"] != "Article not found" {
        t.Errorf("Expected the 'error' key of the response to be set to 'Article not found'. Got '%s'", m["error"])
    }
}


func TestCreateArticle(t *testing.T) {

    clearTable()

    var jsonStr = []byte(`{"title":"test article", "authors": "Muhammad Rahman,Md Abd Ar Rahman"}`)
    req, _ := http.NewRequest("POST", "/article", bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")

    response := executeRequest(req)
    checkResponseCode(t, http.StatusCreated, response.Code)

    var m map[string]interface{}
    json.Unmarshal(response.Body.Bytes(), &m)

    if m["title"] != "test article" {
        t.Errorf("Expected article title to be 'test article'. Got '%v'", m["title"])
    }

    if m["authors"] != "Muhammad Rahman,Md Abd Ar Rahman" {
        t.Errorf("Expected article authors to be 'Muhammad Rahman,Md Abd Ar Rahman'. Got '%v'", m["authors"])
    }

    // the id is compared to 1.0 because JSON unmarshaling converts numbers to
    // floats, when the target is a map[string]interface{}
    if m["id"] != 1.0 {
        t.Errorf("Expected article ID to be '1'. Got '%v'", m["id"])
    }
}


func TestGetArticle(t *testing.T) {
    clearTable()
    addArticles(1)

    req, _ := http.NewRequest("GET", "/article/1", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)
}

// main_test.go

func addArticles(count int) {
    if count < 1 {
        count = 1
    }

    for i := 0; i < count; i++ {
        a.DB.Exec("INSERT INTO articles(title, authors) VALUES(?, ?)", "Article "+strconv.Itoa(i), "Authors "+strconv.Itoa(i))
    }
}


func TestUpdateArticle(t *testing.T) {

    clearTable()
    addArticles(1)

    req, _ := http.NewRequest("GET", "/article/1", nil)
    response := executeRequest(req)
    var originalArticle map[string]interface{}
    json.Unmarshal(response.Body.Bytes(), &originalArticle)

    var jsonStr = []byte(`{"title":"test article - updated title", "authors": "test author, author 2"}`)
    req, _ = http.NewRequest("PUT", "/article/1", bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")

    response = executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)

    var m map[string]interface{}
    json.Unmarshal(response.Body.Bytes(), &m)

    if m["id"] != originalArticle["id"] {
        t.Errorf("Expected the id to remain the same (%v). Got %v", originalArticle["id"], m["id"])
    }

    if m["title"] == originalArticle["title"] {
        t.Errorf("Expected the title to change from '%v' to '%v'. Got '%v'", originalArticle["title"], m["title"], m["title"])
    }

    if m["authors"] == originalArticle["authors"] {
        t.Errorf("Expected the authors to change from '%v' to '%v'. Got '%v'", originalArticle["authors"], m["authors"], m["authors"])
    }
}


func TestDeleteArticle(t *testing.T) {
    clearTable()
    addArticles(1)

    req, _ := http.NewRequest("GET", "/article/1", nil)
    response := executeRequest(req)
    checkResponseCode(t, http.StatusOK, response.Code)

    req, _ = http.NewRequest("DELETE", "/article/1", nil)
    response = executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)

    req, _ = http.NewRequest("GET", "/article/1", nil)
    response = executeRequest(req)
    checkResponseCode(t, http.StatusNotFound, response.Code)
}