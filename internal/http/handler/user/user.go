package user

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"effective_mobile/internal/exception"
	"effective_mobile/internal/http/handler/user/request"
	"effective_mobile/internal/http/handler/user/response"
	resp "effective_mobile/internal/lib/api/response"
	"effective_mobile/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	log     *slog.Logger
	service Service
}

func New(log *slog.Logger, service Service) *Handler {
	return &Handler{log: log, service: service}
}

// Index
// @tags users
// @param page query int false "page"
// @param name query string false "name"
// @param surname query string false "surname"
// @success 200 {object} response.Index
// @response 500 {object} response.Error
// @router /users [get]
func (h *Handler) Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.user.Index"

		log := h.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		queryParams := r.URL.Query()

		name := queryParams.Get("name")
		surname := queryParams.Get("surname")
		page, err := strconv.Atoi(queryParams.Get("page"))
		if err != nil {
			page = 1
		}

		users, err := h.service.Index(page, name, surname)
		if err != nil {
			log.Error("failed to get list of users", sl.Err(err))

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Err("failed to get list"))

			return
		}

		log.Info("lit got")

		render.JSON(w, r, response.Index{
			Success: resp.OK(),
			Users:   users,
		})
	}
}

// Store
// @tags users
// @param request body request.Store true "request"
// @success 200 {object} response.Store
// @response 400 {object} response.Error
// @response 500 {object} response.Error
// @router /users [post]
func (h *Handler) Store() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.user.Store"

		log := h.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req request.Store

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Err("empty request"))

			return
		}

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Err("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			var validatorErr validator.ValidationErrors

			errors.As(err, &validatorErr)

			log.Error("invalid request", sl.Err(err))

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.ValidationErr(validatorErr))

			return
		}

		user, err := h.service.Store(req.Name, req.Surname, req.Patronymic)
		if err != nil {
			log.Error("failed to store user", sl.Err(err))

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Err("failed to store user"))

			return
		}

		log.Info("user stored", slog.String("id", user.Id))

		render.JSON(w, r, response.Store{
			Success: resp.OK(),
			User:    user,
		})
	}
}

// Show
// @tags users
// @param id path string true "id"
// @success 200 {object} response.Show
// @response 404 {object} response.Error
// @response 400 {object} response.Error
// @response 500 {object} response.Error
// @router /users/{id} [get]
func (h *Handler) Show() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.user.Show"

		log := h.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id := chi.URLParam(r, "id")
		if id == "" {
			log.Info("id is empty")

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Err("invalid request"))

			return
		}

		user, err := h.service.Show(id)
		if errors.Is(err, exception.ErrUserNotFound) {
			log.Debug("user not found id: " + id)

			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, resp.Err("user not found"))

			return
		}

		if err != nil {
			log.Error("failed to get user", sl.Err(err))

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Err("internal error"))

			return
		}

		log.Info("user got", slog.String("id", user.Id))

		render.JSON(w, r, response.Show{
			Success: resp.OK(),
			User:    user,
		})
	}
}

// Update
// @tags users
// @param id path string true "id"
// @param request body request.Update true "request"
// @success 200 {object} response.Update
// @response 404 {object} response.Error
// @response 400 {object} response.Error
// @response 500 {object} response.Error
// @router /users/{id} [patch]
func (h *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.user.Update"

		log := h.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id := chi.URLParam(r, "id")
		if id == "" {
			log.Info("id is empty")

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Err("invalid request"))

			return
		}

		var req request.Update

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Err("empty request"))

			return
		}

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Err("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			var validatorErr validator.ValidationErrors

			errors.As(err, &validatorErr)

			log.Error("invalid request", sl.Err(err))

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.ValidationErr(validatorErr))

			return
		}

		user, err := h.service.Update(id, req.Name, req.Surname, req.Patronymic, req.Gender, req.CountryId, req.Age)
		if errors.Is(err, exception.ErrUserNotFound) {
			log.Debug("user not found id: " + id)

			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, resp.Err("user not found"))

			return
		}

		if err != nil {
			log.Error("failed to update user", sl.Err(err))

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Err("failed to store user"))

			return
		}

		log.Info("user updated", slog.String("id", user.Id))

		render.JSON(w, r, response.Update{
			Success: resp.OK(),
			User:    user,
		})
	}
}

// Delete
// @tags users
// @param id path string true "id"
// @success 200 {object} response.Delete
// @response 404 {object} response.Error
// @response 400 {object} response.Error
// @response 500 {object} response.Error
// @router /users/{id} [delete]
func (h *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.user.Delete"

		log := h.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id := chi.URLParam(r, "id")
		if id == "" {
			log.Info("id is empty")

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Err("invalid request"))

			return
		}

		id, err := h.service.Delete(id)
		if errors.Is(err, exception.ErrUserNotFound) {
			log.Debug("user not found id: " + id)

			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, resp.Err("user not found"))

			return
		}

		if err != nil {
			log.Error("failed to delete user", sl.Err(err))

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Err("internal err"))

			return
		}

		log.Info("user deleted", slog.String("id", id))

		render.JSON(w, r, response.Delete{
			Success: resp.OK(),
			Id:      id,
		})
	}
}
