package db

import (
	"context"
	"errors"
	"github.com/jackc/pgtype"
	log "github.com/sirupsen/logrus"
	"net.vikesh/goshop/dto"
	"strconv"
	"strings"
	"time"
)

type IArticleService interface {
	NewArticle() (*dto.Article, error)
	SaveArticle(article *dto.Article) error
	GetArticleById(id string) (*dto.Article, error)
	GetArticleByTitle(articleTitle string) (*dto.Article, error)
	DeleteArticleById(id string) error
	ListArticles(searchParams *dto.SearchParams) (*dto.ListArticleResult, error)
	ListPublishedArticles(searchParams *dto.SearchParams) (*dto.ListArticleResult, error)
}

type ArticleService struct {
}

func (a *ArticleService) NewArticle() (*dto.Article, error) {
	var lastInsertedId = &pgtype.Int4{}
	queryError := db.QueryRow(context.Background(), CreateNewArticleReturningIdQuery).Scan(lastInsertedId)
	if queryError != nil {
		log.Errorf("error in executing [sql=%v] [params=%v]", CreateNewArticleReturningIdQuery, queryError)
		return nil, queryError
	}
	if present(lastInsertedId) {
		return &dto.Article{ID: lastInsertedId.Int}, nil
	} else {
		return &dto.Article{}, errors.New("failed to create a new article")
	}
}

func (a *ArticleService) SaveArticle(article *dto.Article) error {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, txOptions)
	defer tx.Rollback(ctx)
	if err != nil {
		panic(err)
	}
	existingArticle, err := a.GetArticleById(strconv.Itoa(int(article.ID)))
	var titleFromUrl string
	if err != nil || existingArticle.ID == 0 {
		return errors.New("article not found")
	} else {
		published := existingArticle.Published
		if !published && article.Published {
			titleFromUrl = createUrlFromTitle(article.Title)
		}
	}
	log.Debugf("[sql=%v] [params=%v]", UpdateArticleQuery, []interface{}{article.Title, article.Content, article.Summary, time.Now(), article.ID})
	if len(titleFromUrl) == 0 {
		exec, execError := tx.Exec(ctx, UpdateArticleQuery, article.Title, article.Content, article.Summary, article.Published, time.Now(), article.ID)
		defer tx.Commit(ctx)
		if execError != nil {
			log.Errorf("error in executing [sql=%v] [params=%v], %v", UpdateArticleQuery, []interface{}{article.Title, article.Content, article.Summary, article.Published, article.ID}, exec)
			return execError
		}
	} else {
		exec, execError := tx.Exec(ctx, UpdateArticleWithUrlQuery, article.Title, article.Content, article.Summary, article.Published, titleFromUrl, time.Now(), article.ID)
		defer tx.Commit(ctx)
		if execError != nil {
			log.Errorf("error in executing [sql=%v] [params=%v], %v", UpdateArticleWithUrlQuery, []interface{}{article.Title, article.Content, article.Summary, article.Published, titleFromUrl, article.ID}, exec)
			return execError
		}
	}
	return nil
}

func (a *ArticleService) GetArticleById(id string) (*dto.Article, error) {
	article := &dto.Article{}
	row := db.QueryRow(context.Background(), GetArticleByIdQuery, id)
	articleId := &pgtype.Int4{}
	title := &pgtype.Text{}
	summary := &pgtype.Text{}
	content := &pgtype.Text{}
	published := &pgtype.Bool{}
	creationTime := &pgtype.Timestamp{}
	modifiedTime := &pgtype.Timestamp{}
	scanError := row.Scan(articleId, title, summary, content, published, creationTime, modifiedTime)
	if scanError != nil {
		return nil, scanError
	}
	assign([]pgtype.Value{articleId, title, summary, content, published, creationTime, modifiedTime},
		[]interface{}{&article.ID, &article.Title, &article.Summary, &article.Content, &article.Published, &article.CreationTime, &article.ModifiedTime})
	return article, nil
}

func (a *ArticleService) GetArticleByTitle(articleTitle string) (*dto.Article, error) {
	article := &dto.Article{}
	row := db.QueryRow(context.Background(), GetArticleByTitleQuery, strings.ToLower(articleTitle))
	articleId := &pgtype.Int4{}
	title := &pgtype.Text{}
	summary := &pgtype.Text{}
	content := &pgtype.Text{}
	published := &pgtype.Bool{}
	creationTime := &pgtype.Timestamp{}
	modifiedTime := &pgtype.Timestamp{}
	scanError := row.Scan(articleId, title, summary, content, published, creationTime, modifiedTime)
	if scanError != nil {
		return nil, scanError
	}
	assign([]pgtype.Value{articleId, title, summary, content, published, creationTime, modifiedTime},
		[]interface{}{&article.ID, &article.Title, &article.Summary, &article.Content, &article.Published, &article.CreationTime, &article.ModifiedTime})
	return article, nil
}
func (a *ArticleService) DeleteArticleById(id string) error {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, txOptions)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	exec, err := tx.Exec(context.Background(), DeleteArticleByIdQuery, id)
	if exec != nil {
		return err
	}
	defer tx.Commit(ctx)
	return nil
}

func (a *ArticleService) ListArticles(searchParams *dto.SearchParams) (*dto.ListArticleResult, error) {
	var count int
	countError := db.QueryRow(context.Background(), GetArticlesCountQuery).Scan(&count)
	if countError != nil {
		panic(count)
	}
	rows, findError := db.Query(context.Background(), GetArticlesAndContentQuery)
	if findError != nil {
		return nil, findError
	}
	defer rows.Close()
	listing := &dto.ListArticleResult{Count: count}
	listing.Listings = make([]dto.ArticleListing, listing.Count)
	for i := 0; i < listing.Count && rows.Next(); i++ {
		id := &pgtype.Int4{}
		title := &pgtype.Text{}
		publishedTime := &pgtype.Timestamp{}
		summary := &pgtype.Text{}
		published := &pgtype.Bool{}
		url := &pgtype.Text{}
		scanError := rows.Scan(id, title, summary, publishedTime, published, url)
		l := &listing.Listings[i]
		if scanError != nil {
			return nil, scanError
		}
		assign([]pgtype.Value{id, title, publishedTime, summary, published, url}, []interface{}{&l.ID, &l.Title, &l.CreationTime, &l.Summary, &l.Published, &l.Url})
	}
	return listing, nil
}

func (a *ArticleService) ListPublishedArticles(searchParams *dto.SearchParams) (*dto.ListArticleResult, error) {
	var count int
	countError := db.QueryRow(context.Background(), GetPublishedArticlesCountQuery).Scan(&count)
	if countError != nil {
		panic(count)
	}
	rows, findError := db.Query(context.Background(), GetPublishedArticlesAndContentQuery)
	if findError != nil {
		return nil, findError
	}
	defer rows.Close()
	listing := &dto.ListArticleResult{Count: count}
	listing.Listings = make([]dto.ArticleListing, listing.Count)
	for i := 0; i < listing.Count && rows.Next(); i++ {
		id := &pgtype.Int4{}
		title := &pgtype.Text{}
		publishedTime := &pgtype.Timestamp{}
		summary := &pgtype.Text{}
		published := &pgtype.Bool{}
		url := &pgtype.Text{}
		scanError := rows.Scan(id, title, summary, publishedTime, published, url)
		l := &listing.Listings[i]
		if scanError != nil {
			return nil, scanError
		}
		assign([]pgtype.Value{id, title, publishedTime, summary, published, url}, []interface{}{&l.ID, &l.Title, &l.CreationTime, &l.Summary, &l.Published, &l.Url})
	}
	return listing, nil
}
