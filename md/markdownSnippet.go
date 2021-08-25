package md

import "time"

// MarkdownSnippet
type MarkdownSnippet struct {
	// Markdown snippet guid.
	ID string `json:"id,omitempty" bson:"id" format:"uuid"`
	// Markdown body to save.
	Body string `json:"body" bson:"body" example:"# Markdown Snippet\nSome Text"`
	// Date markdown snippet was created
	CreateDate time.Time `json:"createDate,omitempty" bson:"createDate" format:"date-time"`
}

// MDListItem
type MDListItem struct {
	// Markdown snippet guid.
	ID string `json:"id,omitempty" bson:"id" format:"uuid"`
	// Date markdown snippet was created
	CreateDate time.Time `json:"createDate,omitempty" bson:"createDate" format:"date-time"`
}

// CreateMDReq
type CreateMDReq struct {
	// Markdown body to save.
	Body string `json:"body" validate:"required,min=1,max=1024" minLength:"1" maxLength:"1024" example:"# Markdown Snippet\nSome Text"`
}

// UpdateMarkdownSnippet
type UpdateMDReq struct {
	// Markdown snippet guid.
	ID string `json:"id,omitempty" format:"uuid" validate:"required"`
	// Markdown body to save.
	Body string `json:"body" validate:"required,min=1,max=1024" minLength:"1" maxLength:"1024" example:"# Markdown Snippet\nSome Text"`
}
