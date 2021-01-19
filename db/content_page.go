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


type IContentPageService interface {
	NewContentPage() (*dto.Page, error)
	GetPageById(id string) (*dto.Page, error)
	SaveContentPage(page *dto.Page) error
	GetContentPageById(id string) (*dto.Page, error)
	GetContentPageByUrl(pageTitle string) (*dto.Page, error)
	DeleteContentPageById(id string) error
	ListContentPages(searchParams *dto.SearchParams) (*dto.ListPageResult, error)
}

type ContentPageService struct {
}

func (c *ContentPageService) NewContentPage() (*dto.Page, error) {
	var lastInsertedId = &pgtype.Int4{}
	queryError := db.QueryRow(context.Background(), CreateNewContentPageReturningIdQuery).Scan(lastInsertedId)
	if queryError != nil {
		log.Errorf("error in executing [sql=%v] [params=%v]", CreateNewContentPageReturningIdQuery, queryError)
		return nil, queryError
	}
	if present(lastInsertedId) {
		return &dto.Page{ID: lastInsertedId.Int}, nil
	} else {
		return &dto.Page{}, errors.New("failed to create a new page")
	}
}

func (c *ContentPageService) GetPageById(id string) (*dto.Page, error) {
	page := &dto.Page{}
	row := db.QueryRow(context.Background(), GetContentPageByIdQuery, id)
	pageId := &pgtype.Int4{}
	title := &pgtype.Text{}
	summary := &pgtype.Text{}
	content := &pgtype.Text{}
	published := &pgtype.Bool{}
	creationTime := &pgtype.Timestamp{}
	modifiedTime := &pgtype.Timestamp{}
	scanError := row.Scan(pageId, title, summary, content, published, creationTime, modifiedTime)
	if scanError != nil {
		return nil, scanError
	}
	assign([]pgtype.Value{pageId, title, summary, content, published, creationTime, modifiedTime},
		[]interface{}{&page.ID, &page.Title, &page.Summary, &page.Content, &page.Published, &page.CreationTime, &page.ModifiedTime})
	return page, nil
}

func (c *ContentPageService) SaveContentPage(page *dto.Page) error {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, txOptions)
	defer tx.Rollback(ctx)
	if err != nil {
		panic(err)
	}
	existingPage, err := c.GetPageById(strconv.Itoa(int(page.ID)))
	var titleFromUrl string
	if err != nil || existingPage.ID == 0 {
		return errors.New("page not found")
	} else {
		published := existingPage.Published
		if !published && page.Published {
			titleFromUrl = createUrlFromTitle(page.Title)
		}
	}
	log.Info("[sql=%v] [params=%v]", UpdateContentQuery, []interface{}{page.Title, page.Content, page.Summary, time.Now(), page.ID})
	if len(titleFromUrl) == 0 {
		exec, execError := tx.Exec(ctx, UpdateContentQuery, page.Title, page.Content, page.Summary, page.Published, time.Now(), page.ID)
		defer tx.Commit(ctx)
		if execError != nil {
			log.Errorf("error in executing [sql=%v] [params=%v], %v", UpdateContentQuery, []interface{}{page.Title, page.Content, page.Summary, page.Published, page.ID}, exec)
			return execError
		}
	} else {
		exec, execError := tx.Exec(ctx, UpdateContentWithUrlQuery, page.Title, page.Content, page.Summary, page.Published, titleFromUrl, time.Now(), page.ID)
		defer tx.Commit(ctx)
		if execError != nil {
			log.Errorf("error in executing [sql=%v] [params=%v], %v", UpdateContentWithUrlQuery, []interface{}{page.Title, page.Content, page.Summary, page.Published, titleFromUrl, page.ID}, exec)
			return execError
		}
	}
	return nil
}

func (c *ContentPageService) GetContentPageById(id string) (*dto.Page, error) {
	page := &dto.Page{}
	row := db.QueryRow(context.Background(), GetContentPageByIdQuery, id)
	pageId := &pgtype.Int4{}
	title := &pgtype.Text{}
	summary := &pgtype.Text{}
	content := &pgtype.Text{}
	published := &pgtype.Bool{}
	creationTime := &pgtype.Timestamp{}
	modifiedTime := &pgtype.Timestamp{}
	scanError := row.Scan(pageId, title, summary, content, published, creationTime, modifiedTime)
	if scanError != nil {
		return nil, scanError
	}
	assign([]pgtype.Value{pageId, title, summary, content, published, creationTime, modifiedTime},
		[]interface{}{&page.ID, &page.Title, &page.Summary, &page.Content, &page.Published, &page.CreationTime, &page.ModifiedTime})
	return page, nil
}

func (c *ContentPageService) GetContentPageByUrl(pageTitle string) (*dto.Page, error) {
	page := &dto.Page{}
	row := db.QueryRow(context.Background(), GetContentPageByUrlQuery, strings.ToLower(pageTitle))
	pageId := &pgtype.Int4{}
	title := &pgtype.Text{}
	summary := &pgtype.Text{}
	content := &pgtype.Text{}
	published := &pgtype.Bool{}
	creationTime := &pgtype.Timestamp{}
	modifiedTime := &pgtype.Timestamp{}
	scanError := row.Scan(pageId, title, summary, content, published, creationTime, modifiedTime)
	if scanError != nil {
		return nil, scanError
	}
	assign([]pgtype.Value{pageId, title, summary, content, published, creationTime, modifiedTime},
		[]interface{}{&page.ID, &page.Title, &page.Summary, &page.Content, &page.Published, &page.CreationTime, &page.ModifiedTime})
	return page, nil
}
func (c *ContentPageService) DeleteContentPageById(id string) error {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, txOptions)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	exec, err := tx.Exec(context.Background(), DeleteContentPageByIdQuery, id)
	if exec != nil {
		return err
	}
	defer tx.Commit(ctx)
	return nil
}

func (c *ContentPageService) ListContentPages(searchParams *dto.SearchParams) (*dto.ListPageResult, error) {
	var count int
	countError := db.QueryRow(context.Background(), GetContentPagesCountQuery).Scan(&count)
	if countError != nil {
		panic(count)
	}
	rows, findError := db.Query(context.Background(), GetPageAndContentQuery)
	if findError != nil {
		return nil, findError
	}
	defer rows.Close()
	listing := &dto.ListPageResult{Count: count}
	listing.Listings = make([]dto.PageListing, listing.Count)
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
