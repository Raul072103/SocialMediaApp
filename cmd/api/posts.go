package main

import (
	"SocialMediaApp/internal/store"
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
)

type postKey string

const postCtx postKey = "post"

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

// deletePostHandler godoc
//
//	@Summary		Deletes a post
//	@Description	Deletes a post by ID.
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int		true	"Post ID"
//	@Success		204		{string}	string	"Post deleted"
//	@Failure		400		{object}	error	"Bad request"
//	@Failure		404		{object}	error	"Post not found"
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID}/ [delete]
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	err := app.store.Comments.DeleteByPostID(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	err = app.store.Posts.DeleteById(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// getPostHandler godoc
//
//	@Summary		Retrieves a post
//	@Description	Retrieves a post by ID.
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int		true	"Post ID"
//	@Success		204		{string}	string	"Post retrieved successfully"
//	@Failure		400		{object}	error	"Bad request"
//	@Failure		404		{object}	error	"Post not found"
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID}/ [get]
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	comments, err := app.store.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = comments

	err = jsonResponse(w, http.StatusOK, &post)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// getPostHandler godoc
//
//	@Summary		Creates a post
//	@Description	Creates a post.
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreatePostPayload	true	"Post payload"
//	@Success		204		{string}	string				"Post created successfully"
//	@Failure		400		{object}	error				"Bad request"
//	@Failure		404		{object}	error				"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/posts/ [get]
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := getUserFromContext(r)

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  user.ID,
	}

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		log.Println(err.Error())
		app.internalServerError(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusCreated, post); err != nil {
		log.Println(err.Error())
		app.internalServerError(w, r, err)
		return
	}
}

type UpdatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

// updatePostHandler godoc
//
//	@Summary		Updates a post
//	@Description	Updates a post by ID.
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreatePostPayload	true	"Post payload"
//	@Success		204		{string}	string				"Post updated successfully"
//	@Failure		400		{object}	error				"Post not found"
//	@Failure		404		{object}	error				"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/posts/ [patch]
func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	var payload UpdatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if payload.Title != nil {
		post.Title = *payload.Title
	}

	if err := app.store.Posts.Update(r.Context(), post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postIDParam := chi.URLParam(r, "postID")
		postID, err := strconv.ParseInt(postIDParam, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()

		post, err := app.store.Posts.GetByID(ctx, postID)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
				return
			}
		}

		ctx = context.WithValue(ctx, postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *store.Post {
	post, _ := r.Context().Value(postCtx).(*store.Post)
	return post
}
