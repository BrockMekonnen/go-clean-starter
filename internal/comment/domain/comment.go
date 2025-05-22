package domain

import (
	"time"
)

type Status string

const (
	StatusActive  Status = "ACTIVE"
	StatusDeleted Status = "DELETED"
)

type Comment struct {
	ID        string
	Body      string
	PostId    string
	UserId    string
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
	Version   int
}

type CommentProps struct {
	ID     string
	Body   string
	PostId string
	UserId string
}

func CreateComment(props CommentProps) Comment {
	return Comment{
		ID:        props.ID,
		Body:      props.Body,
		PostId:    props.PostId,
		UserId:    props.UserId,
		Status:    StatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   0,
	}
}

func MarkAsDeleted(c Comment) Comment {
	c.Status = StatusDeleted
	return c
}
