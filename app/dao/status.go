package dao

import (
	//"database/sql"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/domain/repository"
	"yatter-backend-go/app/handler/timelines"
	"fmt"
	"context"
	"errors"
	"database/sql"
	"time"

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

type StatusModel struct {
	ID        int64     `json:"id,omitempty" db:"id"`
	AccountId int64     `json:"account_id,omitempty" db:"account_id"`
	Content   string    `json:"content,omitempty" db:"content"`
	CreateAt  time.Time `json:"create_at,omitempty" db:"create_at"`
}

func (r *status) PostStatus(s object.Status) (int64, error) {

	tx,  _ := r.db.Beginx()
	var err error
	
	defer func() {
		switch r := recover(); {
		case r != nil:
			tx.Rollback()
			panic(r)
		case err != nil:
			tx.Rollback()
		}
	}()

	status := new(StatusModel)
	status.AccountId = s.Account.ID
	status.Content = s.Content
	query := `
		INSERT INTO Status (account_id, content)
		VALUES (:account_id, :content);
	`

	result, err := tx.NamedExec(query, status)
	fmt.Println(result)
	if err != nil {
		return 0, err
	}
	

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	if s.MediaAttachment.Url != "" {
		query := `
		INSERT INTO Media (type, url, description, status_id)
		VALUES (:type, :url, :description, :status_id);
	`
		fmt.Println(lastInsertID)
		_, err := tx.NamedExec(query, map[string]interface{}{
			"type": s.MediaAttachment.Type,
			"url":    s.MediaAttachment.Url,
			"description": s.MediaAttachment.Description,
			"status_id": lastInsertID,
		})
		if err != nil {
			fmt.Println(err)
			return 0, err
		}

		if err = tx.Commit(); err != nil {
			return 0, err
		}
	}

	return lastInsertID, nil
}

func (r *status) FindById(ctx context.Context, id int64) (*object.Status, int64, error) {
	entity := new(StatusModel)
	err := r.db.QueryRowxContext(ctx, "select * from status where id = ?", id).StructScan(entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, 0, nil
		}

		return nil, 0, fmt.Errorf("failed to find status from db: %w", err)
	}

	media, err := r.FindMediaById(ctx, entity.ID)
	if err != nil {
		return nil, 0, err
	}

	status := new(object.Status)
	if media != nil {
		status.MediaAttachment = *media
	}
	status.ID = entity.ID
	status.Content = entity.Content
	status.CreateAt = entity.CreateAt
	

	AccountId := entity.AccountId

	return status, AccountId, nil
}

func (r *status) FindMediaById(ctx context.Context, id int64) (*object.MediaAttachment,  error) {
	entity := new(object.MediaAttachment)
	err := r.db.QueryRowxContext(ctx, "select * from Media where status_id = ?", id).StructScan(entity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to find media from db: %w", err)
	}

	return entity, nil
}


func (r *status) GetTimelines(ctx context.Context) ([]object.Status, error) {

	statuses := []object.Status{}
	var args []interface{}

	maxID, _ := ctx.Value(timelines.ContextMaxID).(int64)
	
    sinceID, _ := ctx.Value(timelines.ContextSinceID).(int64)

    limit, _ := ctx.Value(timelines.ContextLimit).(int64)

	onlyMedia, onlyMediaOk := ctx.Value(timelines.ContextOnlyMedia).(bool)

	query := `SELECT * FROM Status WHERE id > ? `
	args = append(args, sinceID)

	if maxID > 0 {
		query += " AND id < ?"
		args = append(args, maxID)
	}

	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}
	

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()


	entity := new(StatusModel)
	status := new(object.Status)
	Ar := ctx.Value(timelines.ContextAr).(repository.Account)
	
	if !onlyMediaOk || !onlyMedia{
		for rows.Next() {
			
			if err := rows.StructScan(entity); err != nil {
				return nil, fmt.Errorf("failed to scan row: %w", err)
			}
			AccountId := entity.AccountId
			account, err := Ar.FindById(ctx, AccountId)
			if err != nil {
				return nil, fmt.Errorf("failed to find account from db: %w", err)
			}
			status.Account = *account
			status.ID = entity.ID
			status.Content = entity.Content
			status.CreateAt = entity.CreateAt
			
			media, err := r.FindMediaById(ctx, status.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to find media from db: %w", err)
			}
			if media != nil {
				status.MediaAttachment = *media
			}
			statuses = append(statuses, *status)
		}
	} else {
		for rows.Next() {
			
			if err := rows.StructScan(&entity); err != nil {
				return nil, fmt.Errorf("failed to scan row: %w", err)
			}

			AccountId := entity.AccountId
			account, err := Ar.FindById(ctx, AccountId)
			if err != nil {
				return nil, fmt.Errorf("failed to find account from db: %w", err)
			}
			status.Account = *account
			status.ID = entity.ID
			status.Content = entity.Content
			status.CreateAt = entity.CreateAt

			media, err := r.FindMediaById(ctx, status.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to find media from db: %w", err)
			}
			if media != nil {
				status.MediaAttachment = *media
				statuses = append(statuses, *status)
			}
		}
	}
	
	return statuses, nil

}