package main

import (
	"SocialMediaApp/internal/store"
	"net/http"
)

// getUserFeedHandler godoc
//
//	@Summary		Gets the user's feed.
//	@Description	Retrieves the user's feed, using pagination and filtering, with a maximum of 20 posts per request.
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//
//	@Param			limit	path		int			false	"Length of the response"
//	@Param			offset	path		int			false	"Offset of the response"
//	@Param			sort	path		string		false	"Method of sorting the posts"
//	@Param			tags	path		[]string	false	"The tags which the posts must contain"
//	@Param			search	path		string		false	"Keyword that must appear in the posts"
//	@Success		200		{string}	string		"Feed retrieved successfully!"
//	@Failure		400		{object}	error		"Bad request"
//	@Failure		404		{object}	error		"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	// pagination, filters, sort
	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
		Tags:   []string{},
	}

	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(fq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	feed, err := app.store.Posts.GetUserFeed(ctx, int64(13), fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
