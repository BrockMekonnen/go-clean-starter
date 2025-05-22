package domain

import (
	"github.com/BrockMekonnen/go-clean-starter/internal/_shared/domain"
	"time"
)

type PostState string

const (
	StateDraft     PostState = "DRAFT"
	StatePublished PostState = "PUBLISHED"
	StateDeleted   PostState = "DELETED"
)

type Post struct {
	ID          string
	Title       string
	Content     string
	UserId      string
	State       PostState
	PublishedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Version     int
}

type PostProps struct {
	ID      string
	Title   string
	Content string
	UserId  string
}

func NewPost(props PostProps) (*Post, error) {
	a := &Post{
		ID:        props.ID,
		Title:     props.Title,
		Content:   props.Content,
		UserId:    props.UserId,
		State:     StateDraft,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   0,
	}
	if err := a.validate(); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *Post) PublishPost() (*Post, error) {
	if err := a.validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	a.State = StatePublished
	a.PublishedAt = &now
	a.UpdatedAt = now
	return a, nil
}

func (a *Post) MarkAsDeleted() (*Post, error) {
	if err := a.validate(); err != nil {
		return nil, err
	}
	a.State = StateDeleted
	a.UpdatedAt = time.Now()
	return a, nil
}

func (a *Post) ChangeTitle(title string) (*Post, error) {
	a.Title = title
	a.UpdatedAt = time.Now()
	if err := a.validate(); err != nil {
		return nil, err
	}
	return a, nil
}
func (a *Post) ChangeContent(content string) (*Post, error) {
	a.Content = content
	a.UpdatedAt = time.Now()
	if err := a.validate(); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *Post) IsPublished() bool {
	return a.State == StatePublished
}

func (a *Post) validate() error {
	if len(a.Title) == 0 {
		return domain.NewBusinessError("title cannot be empty", "")
	}
	if len(a.Content) == 0 {
		return domain.NewBusinessError("content cannot be empty", "")
	}
	return nil
}
