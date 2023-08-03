package dao

import (
	//"database/sql"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type (
	status struct {
		db *sqlx.DB
	}
)

func NewStatus(db *sqlx.DB) repository.Status {
	return &status{db: db}
}

func (r *status) PostStatus(s object.Status) error {

	accountID := s.Account.ID
	query := `
		INSERT INTO Status (account_id, content)
		VALUES (:account_id, :content);
	`

	_, err := r.db.NamedExec(query, map[string]interface{}{
		"account_id": accountID,
		"content":    s.Content,
	})
	if err != nil {
		return err
	}
	return nil
}