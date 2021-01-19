package dto

import (
	"time"
)

type Page struct {
	ID            int32     `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Published     bool      `json:"published"`
	Summary       string    `json:"summary"`
	ModifiedTime  time.Time `json:"modifiedTime"`
	CreationTime  time.Time `json:"creationTime"`
	PublishedTime time.Time `json:"publishedTime"`
	ModifiedDate  time.Time `json:"modifiedDate"`
}

type PageListing struct {
	ID           int32     `json:"id"`
	Title        string    `json:"title"`
	Summary      string    `json:"summary"`
	Url          string    `json:"url"`
	Published    bool      `json:"published"`
	CreationTime time.Time `json:"creationTime"`
}

type ListPageResult struct {
	Count    int           `json:"count"`
	Page     int           `json:"page"`
	Pages    int           `json:"pages"`
	Listings []PageListing `json:"listings"`
}

type Article struct {
	ID            int32     `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Published     bool      `json:"published"`
	Summary       string    `json:"summary"`
	ModifiedTime  time.Time `json:"modifiedTime"`
	CreationTime  time.Time `json:"creationTime"`
	PublishedTime time.Time `json:"publishedTime"`
	ModifiedDate  time.Time `json:"modifiedDate"`
}

type ArticleListing struct {
	ID           int32     `json:"id"`
	Title        string    `json:"title"`
	Summary      string    `json:"summary"`
	Url          string    `json:"url"`
	Published    bool      `json:"published"`
	CreationTime time.Time `json:"creationTime"`
}

type ListArticleResult struct {
	Count    int              `json:"count"`
	Page     int              `json:"page"`
	Pages    int              `json:"pages"`
	Listings []ArticleListing `json:"listings"`
}

type ErrorResponse struct {
	Error    error
	Status   int
	Response interface{} `json:"error"`
}

type SuccessResponse struct {
	Status int
	Result interface{} `json:"result"`
}

type User struct {
	UserName    string `json:"userName"`
	DisplayName string `json:"displayName"`
}
