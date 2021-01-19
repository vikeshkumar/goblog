package api

import (
	"net/http"
)

func navigationHandler(writer http.ResponseWriter, request *http.Request) {
	type Nav struct {
		Link     string `json:"link"`
		Display  string `json:"display"`
		Children []*Nav `json:"children"`
	}
	m1 := &Nav{"/java", "Java", nil}
	m2 := &Nav{"/spring", "Spring", nil}
	m3 := &Nav{"/hibernate", "JPA", nil}
	m4 := &Nav{"/golang", "Golang", nil}
	m5 := &Nav{"/thoughts", "Thoughts", nil}
	m6 := &Nav{"/series", "Series", nil}
	m7 := &Nav{"/series/goblog", "Blog in Golang", nil}
	m8 := &Nav{"/series/spring-shop", "WebShop in Spring, Hibernate", nil}
	m9 := &Nav{"/series/hybris", "E-Commerce Site with Hybris", nil}
	m6.Children = []*Nav{m7, m8, m9}
	response := []*Nav{m1, m2, m3, m4, m5, m6}
	toJson(success(response), writer, nil)
}
