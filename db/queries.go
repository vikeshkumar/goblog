package db

const (
	CreateNewArticleReturningIdQuery    = `insert into articles(title, content) values('','') returning id`
	UpdateArticleQuery                  = `update articles set title=$1, content=$2, summary=$3, published=$4, modified_time=$5 where id=$6`
	UpdateArticleWithUrlQuery           = `update articles set title=$1, content=$2, summary=$3, published=$4, url=$5, modified_time=$6 where id=$7`
	GetArticlesCountQuery               = `select count(id) from articles`
	GetPublishedArticlesCountQuery      = `select count(id) from articles where published = true`
	GetArticlesAndContentQuery          = `select a.id, a.title, a.summary, a.creation_time, a.published, a.url from articles as a order by a.id desc`
	GetPublishedArticlesAndContentQuery = `select a.id, a.title, a.summary, a.creation_time, a.published, a.url from articles as a where a.published = true order by a.id desc`
	GetArticleByIdQuery                 = `select a.id, a.title, a.summary, a.content, a.published, a.creation_time, a.modified_time from articles as a where a.id = $1`
	GetArticleByTitleQuery              = `select a.id, a.title, a.summary, a.content, a.published, a.creation_time, a.modified_time from articles as a where a.url = $1 and a.published = true and a.deleted = false`
	DeleteArticleByIdQuery              = `update articles set deleted = true where id = $1`
	FindUserWithUserNameOrEmailQuery    = `select count(id) from users as u where u.user_name = lower($1) or lower(u.email) = lower($2)`
	FindTotalUserQuery                  = `select count(id) from users as u`
	CreateNewUserQuery                  = `insert into users(user_name, display_name, email, salt, hashed_password) values ($1, $2, $3, $4, $5) returning id`
	FindUserByUserNameQuery             = `select u.id, u.user_name, u.display_name from users as u where u.user_name = lower($1)`
	FindUserByIdQuery                   = `select u.user_name, u.display_name from users as u where u.id = $1`
	CreateTokenForUserQuery             = `insert into tokens(user_id, token_value) values( $1, $2) returning id`
	FindValidTokenQuery                 = `select t.id from tokens as t where t.id = $1 and t.user_id = $2 and t.token_value = $3 and revoked = false order by t.id desc limit 1`
	//Content Page Query
	CreateNewContentPageReturningIdQuery = `insert into pages(title, content) values('','') returning id`
	UpdateContentQuery                   = `update pages set title=$1, content=$2, summary=$3, published=$4, modified_time=$5 where id=$6`
	UpdateContentWithUrlQuery            = `update pages set title=$1, content=$2, summary=$3, published=$4, url=$5, modified_time=$6 where id=$7`
	GetContentPageByIdQuery              = `select a.id, a.title, a.summary, a.content, a.published, a.creation_time, a.modified_time from pages as a where a.id = $1`
	GetContentPageByUrlQuery             = `select a.id, a.title, a.summary, a.content, a.published, a.creation_time, a.modified_time from pages as a where a.url = $1 and a.published = true and a.deleted = false`
	DeleteContentPageByIdQuery           = `update pages set deleted = true where id = $1`
	GetContentPagesCountQuery            = `select count(id) from pages`
	GetPageAndContentQuery               = `select a.id, a.title, a.summary, a.creation_time, a.published, a.url from pages as a order by a.id desc`
)
