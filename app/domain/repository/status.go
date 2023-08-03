package repository

import (
	"yatter-backend-go/app/domain/object"
	"context"
)

type Status interface {
	PostStatus(s object.Status) (int64, error)

	FindById(ctx context.Context, id int64) (*object.Status, int64, error)

	FindMediaById(ctx context.Context, id int64) (*object.MediaAttachment,  error)
}