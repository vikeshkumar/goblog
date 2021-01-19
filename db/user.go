package db

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"golang.org/x/crypto/bcrypt"
	"net.vikesh/goshop/config"
	"net.vikesh/goshop/dto"
	"strconv"
	"strings"
)

type IUserService interface {
	FindUserForCookieValue(token string) (dto.User, error)
	RegisterNewUser(form *dto.RegistrationForm) error
	FindUserByUserName(uname string) (int, dto.User, error)
	CreateTokenForUser(userId int) (string, error)
	IsValidToken(token string) bool
}

type UserService struct {
}

func (u *UserService) FindUserForCookieValue(token string) (dto.User, error) {
	ctx := context.Background()
	nilUser := dto.User{}
	isValidUser := u.IsValidToken(token)
	if !isValidUser {
		return nilUser, nil
	}
	decodedString, _ := base64.StdEncoding.DecodeString(token)
	parts := strings.Split(string(decodedString), ":")
	userId, _ := strconv.Atoi(parts[1])
	name := &pgtype.Text{}
	dName := &pgtype.Text{}
	err := db.QueryRow(ctx, FindUserByIdQuery, userId).Scan(name, dName)
	if err != nil {
		return dto.User{}, errors.New("user not valid")
	}
	var userName string
	var displayName string
	assign([]pgtype.Value{name, dName}, []interface{}{&userName, &displayName})
	return dto.User{UserName: userName, DisplayName: displayName}, nil
}

func (u *UserService) RegisterNewUser(form *dto.RegistrationForm) error {
	var count int
	ctx := context.Background()
	totalUserError := db.QueryRow(ctx, FindTotalUserQuery).Scan(&count)
	if totalUserError != nil || count > 0 {
		return errors.New("sorry, new users not allowed at this time")
	}
	queryError := db.QueryRow(ctx, FindUserWithUserNameOrEmailQuery, form.Username, form.Email).Scan(&count)
	if queryError != nil || count > 0 {
		return errors.New("username already taken")
	}
	tx, txBeginError := db.BeginTx(ctx, txOptions)
	defer tx.Rollback(ctx)
	if txBeginError != nil {
		return errors.New(txBeginError.Error())
	}
	userName := form.Username
	randomSalt := uuid.New().String()
	saltedUserName := strings.Join([]string{userName, randomSalt}, "")
	hashedPassword, hashError := bcrypt.GenerateFromPassword([]byte(saltedUserName), config.Get().GetInt(config.BcryptCost))
	if hashError != nil {
		return errors.New("error creating a new user at this time")
	}
	var userId int
	row := tx.QueryRow(ctx, CreateNewUserQuery, strings.ToLower(form.Username), form.DisplayName, strings.ToLower(form.Email), randomSalt, hashedPassword).Scan(&userId)
	if row != nil || userId == 0 {
		return errors.New("cannot create user with given username")
	}
	tx.Commit(ctx)
	return nil
}

func (u *UserService) FindUserByUserName(uname string) (int, dto.User, error) {
	ctx := context.Background()
	id := &pgtype.Int4{}
	name := &pgtype.Text{}
	dName := &pgtype.Text{}
	err := db.QueryRow(ctx, FindUserByUserNameQuery, uname).Scan(id, name, dName)
	if err != nil {
		return -1, dto.User{}, errors.New("cannot find user with id")
	}
	var userName string
	var userId int
	var displayName string
	assign([]pgtype.Value{id, name, dName}, []interface{}{&userId, &userName, &displayName})
	return userId, dto.User{UserName: userName, DisplayName: displayName}, nil
}

func (u *UserService) CreateTokenForUser(userId int) (string, error) {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, txOptions)
	defer tx.Rollback(ctx)
	if err != nil {
		return "", errors.New("cannot login at this moment")
	}
	tokenValue := uuid.New().String()
	defer tx.Commit(ctx)
	var tokenId int
	tx.QueryRow(ctx, CreateTokenForUserQuery, userId, tokenValue).Scan(&tokenId)
	if tokenId == 0 {
		return "", errors.New("cannot login at this moment")
	}
	return base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(tokenId) + ":" + strconv.Itoa(userId) + ":" + tokenValue)), nil
}

func (u *UserService) IsValidToken(token string) bool {
	ctx := context.Background()
	decodedString, decodeErr := base64.StdEncoding.DecodeString(token)
	if decodeErr != nil {
		return false
	}
	parts := strings.Split(string(decodedString), ":")
	if len(parts) < 3 {
		return false
	}
	var count int
	err := db.QueryRow(ctx, FindValidTokenQuery, parts[0], parts[1], parts[2]).Scan(&count)
	if err != nil || count == 0 {
		return false
	}
	return true
}
