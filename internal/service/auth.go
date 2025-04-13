package service

import (
	"errors"
	"fmt"
	"log"

	"github.com/Gergenus/StandardLib/internal/repository"
	"github.com/Gergenus/StandardLib/pkg"
)

var (
	ErrUserAlreadyExists = errors.New("user alredy exists")
	ErrUnauthorized      = errors.New("incorrect password")
)

type Auth interface {
	SignUp(name, email, password string) (int, error)
	SignIn(name, password string) (string, error)
}

type JWTauth struct {
	UserRepo repository.UserRepo
	Hasher   pkg.Hasher
	Auther   pkg.JWTpkg
}

func (j *JWTauth) SignUp(name, email, password string) (int, error) {
	hashpassword := j.Hasher.Hash(password)
	ok, err := j.UserRepo.GetUserExist(name)
	if err != nil {
		log.Println("Getting user error", err)
		return 0, err
	}
	if ok {
		return 0, ErrUserAlreadyExists
	}
	id, err := j.UserRepo.AddUser(name, hashpassword, email)
	if err != nil {
		log.Println("Adding user error")
		return 0, err
	}
	return id, nil
}

func (j *JWTauth) SignIn(name, password string) (string, error) {
	_, hash, err := j.UserRepo.GetUser(name)
	fmt.Println(hash, "SIGN")
	fmt.Println(password, "Sign password")
	if err != nil {
		log.Println("SignIn err", err)
		return "", err
	}
	if !j.Hasher.Check(hash, password) {
		log.Println("incorrect password")
		return "", ErrUnauthorized
	}
	tkn, err := j.Auther.GenerateToken(name)
	if err != nil {
		log.Println("SignIn GenerateToken err", err)
		return "", err
	}
	return tkn, nil
}
