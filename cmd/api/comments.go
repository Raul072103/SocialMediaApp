package main

import (
	"SocialMediaApp/internal/store"
	"net/http"
	"strconv"
)

type CommentPayload struct {
	UserID  string `json:"user_id" validate:"required,max=100"`
	Content string `json:"content" validate:"required,max=100"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payload CommentPayload
	err := readJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	post := getPostFromCtx(r)
	if post == nil {
		app.notFoundResponse(w, r, err)
		return
	}

	userID, err := strconv.ParseInt(payload.UserID, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.store.Users.GetByID(r.Context(), userID)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	comment := &store.Comment{
		PostID:  post.ID,
		UserID:  user.ID,
		Content: payload.Content,
		User:    *user,
	}

	err = app.store.Comments.Create(r.Context(), comment)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}