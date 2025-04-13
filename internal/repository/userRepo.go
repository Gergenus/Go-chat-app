package repository

import (
	"database/sql"
	"log"

	"github.com/Gergenus/StandardLib/pkg"
)

type UserRepo interface {
	AddUser(Name, password, email string) (int, error)
	DeleteUser(Name string) (int, error)
	GetUserExist(Name string) (bool, error)
	GetUser(Name string) (string, string, error)
}

type PostgresUserRepo struct {
	DB   pkg.DBwraper
	Hash pkg.Hasher
}

func (p *PostgresUserRepo) AddUser(Name, password, email string) (int, error) {
	var id int
	err := p.DB.DB.QueryRow("INSERT INTO users (username, password_hash, email) VALUES($1, $2, $3) RETURNING id", Name, password, email).Scan(&id)
	if err != nil {
		log.Println("AddUser err", err)
		return 0, err
	}
	return id, nil
}

func (p *PostgresUserRepo) DeleteUser(Name string) (int, error) {
	var id int
	err := p.DB.DB.QueryRow("DELETE FROM users WHERE username=$1 RETURNING id", Name).Scan(&id)
	if err != nil {
		log.Println("DeleteUser err", err)
		return 0, err
	}
	return id, nil
}

func (p *PostgresUserRepo) GetUserExist(Name string) (bool, error) {
	var id int
	row := p.DB.DB.QueryRow("SELECT id FROM users WHERE username = $1", Name)
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		log.Println("GetUserExist", err)
		return false, err
	}
	return true, nil
}

func (p *PostgresUserRepo) GetUser(Name string) (string, string, error) {
	var password_hash, email string
	row := p.DB.DB.QueryRow("SELECT password_hash, email FROM users WHERE username = $1", Name)
	err := row.Scan(&password_hash, &email)
	if err != nil {
		log.Println("GetUser err", err)
		return "", "", nil
	}
	return email, password_hash, nil
}
