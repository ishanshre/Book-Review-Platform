package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ishanshre/Book-Review-Platform/internals/helpers"
	"github.com/ishanshre/Book-Review-Platform/internals/models"
)

func (h *Repository) FollowExistsApi(w http.ResponseWriter, r *http.Request) {
	author_id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	user_id := h.App.Session.GetInt(r.Context(), "user_id")
	follower := &models.Follower{
		AuthorID:   author_id,
		UserID:     user_id,
		FollowedAt: time.Now(),
	}
	exists, err := h.DB.FollowerExists(follower)
	if err != nil {
		helpers.StatusInternalServerError(w, "something went wrong")
		return
	}
	if !exists {
		helpers.WriteJson(w, http.StatusOK, map[string]bool{
			"exists": false,
		})
		return
	}
	helpers.WriteJson(w, http.StatusOK, map[string]bool{
		"exists": true,
	})

}

func (h *Repository) FollowApi(w http.ResponseWriter, r *http.Request) {
	author_id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	user_id := h.App.Session.GetInt(r.Context(), "user_id")
	follower := &models.Follower{
		AuthorID:   author_id,
		UserID:     user_id,
		FollowedAt: time.Now(),
	}
	exists, err := h.DB.FollowerExists(follower)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if exists {
		helpers.ServerError(w, errors.New("already exists"))
		return
	}
	if err := h.DB.InsertFollower(follower); err != nil {
		helpers.ServerError(w, err)
		return
	}
	helpers.ApiStatusOk(w, "follow success")
}
func (h *Repository) UnFollowApi(w http.ResponseWriter, r *http.Request) {
	author_id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	user_id := h.App.Session.GetInt(r.Context(), "user_id")
	follower := &models.Follower{
		AuthorID:   author_id,
		UserID:     user_id,
		FollowedAt: time.Now(),
	}
	exists, err := h.DB.FollowerExists(follower)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if !exists {
		helpers.ServerError(w, errors.New("follower does not exists"))
		return
	}
	if err := h.DB.DeleteFollower(user_id, author_id); err != nil {
		helpers.ServerError(w, err)
		return
	}
	helpers.ApiStatusOk(w, "unfollow success")
}

func (h *Repository) GetFollowingsListByUserIdApi(w http.ResponseWriter, r *http.Request) {
	user_id := h.App.Session.GetInt(r.Context(), "user_id")
	authors, err := h.DB.GetAllFollowingsByUserId(user_id)
	if err != nil {
		helpers.StatusInternalServerError(w, err.Error())
		return
	}
	helpers.ApiStatusOkData(w, authors)
}
