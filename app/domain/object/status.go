package object

import (
	"time"
)


type MediaAttachment struct {
	ID int64 `json:"id,omitempty"`

	Type int64 `json:"type,omitempty"`

	Url string `json:"url,omitempty"`

	Description string `json:"description,omitempty"`
}

type Status struct {

	ID int64 `json:"id,omitempty"`

	Account Account `json:"account,omitempty"`

	Content string `json:"content,omitempty"`
	
	CreateAt time.Time `json:"create_at,omitempty" db:"create_at"`

	MediaAttachment MediaAttachment `json:"media_attachement,omitempty"`
}