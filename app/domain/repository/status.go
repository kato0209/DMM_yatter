package repository

import (
	"yatter-backend-go/app/domain/object"
)

type Status interface {
	PostStatus(s object.Status) error
}