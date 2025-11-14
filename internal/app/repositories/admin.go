package repositories

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

type AdminRepo struct {
	db *sql.DB
}

func NewAdminRepo(db *sql.DB) AdminRepo {
	return AdminRepo{db: db}
}

func (r AdminRepo) CreateAdmin(email string, hashPassword string) error {
	_, err := r.db.Exec(`INSERT INTO admin (email, hash_pass) VALUES ($1, $2)`, email, hashPassword)
	if err != nil {
		return err
	}
	return nil
}

func (r AdminRepo) GetAdminByEmail(email string) (string, error) {
	var hashPassword string
	err := r.db.QueryRow(`SELECT hash_pass FROM admin WHERE email = $1`, email).Scan(&hashPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("admin not found")
		}
		return "", err
	}
	return hashPassword, nil
}

func (r AdminRepo) AdminExists(email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM admin WHERE email = $1)`, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
