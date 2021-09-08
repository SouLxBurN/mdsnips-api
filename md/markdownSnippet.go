package md

import "time"

// MarkdownSnippet
type MarkdownSnippet struct {
	// Markdown snippet guid.
	ID string `json:"id,omitempty" bson:"id" format:"uuid"`
    // Markdown snippet title.
    Title string `json:"title" bson:"title" example:"SouLxBurN Is Awesome!"`
	// Markdown body to save.
	Body string `json:"body" bson:"body" example:"# Markdown Snippet\nSome Text"`
	// Date markdown snippet was created.
	CreateDate time.Time `json:"createDate,omitempty" bson:"createDate" format:"date-time"`
}

// MDListItem
type MDListItem struct {
	// Markdown snippet guid.
	ID string `json:"id,omitempty" bson:"id" format:"uuid"`
    // Markdown snippet title.
    Title string `json:"title" bson:"title" example:"SouLxBurN Is Awesome!"`
	// Date markdown snippet was created
	CreateDate time.Time `json:"createDate,omitempty" bson:"createDate" format:"date-time"`
}

// CreateMDReq
type CreateMDReq struct {
    // Markdown snippet title.
    Title string `json:"title" validate:"required,min=1,max=64" minLength:"1" maxLength:"64" example:"SouLxBurN Is Awesome!"`
	// Markdown body to save.
	Body string `json:"body" validate:"required,min=1,max=1024" minLength:"1" maxLength:"1024" example:"# Markdown Snippet\nSome Text"`
}

// UpdateMarkdownSnippet
type UpdateMDReq struct {
	// Markdown snippet guid.
	ID string `json:"id,omitempty" format:"uuid" validate:"required"`
    // Markdown snippet title.
    Title string `json:"title" validate:"required,min=1,max=64" minLength:"1" maxLength:"64" example:"SouLxBurN Is Awesome!"`
	// Markdown body to save.
	Body string `json:"body" validate:"required,min=1,max=1024" minLength:"1" maxLength:"1024" example:"# Markdown Snippet\nSome Text"`
}
