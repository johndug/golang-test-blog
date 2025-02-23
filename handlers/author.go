package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"test-ai-api/stores"
	"test-ai-api/types"
	"test-ai-api/utils"
)

type AuthorHandler struct {
	store *stores.AuthorStore
}

func NewAuthorHandler(store *stores.AuthorStore) *AuthorHandler {
	return &AuthorHandler{store: store}
}

func (h *AuthorHandler) Create(w http.ResponseWriter, r *http.Request) {
	var author types.AuthorCreate
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID := r.Context().Value("userID").(int64)
	result, err := h.store.Create(author, userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, result)
}

func (h *AuthorHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	authors, err := h.store.GetAll(10, 0) // TODO: Add pagination
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(authors) == 0 {
		utils.RespondWithJSON(w, http.StatusOK, []types.Author{})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, authors)
}

func (h *AuthorHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	author, err := h.store.GetBySlug(slug)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Author not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, author)
}

func (h *AuthorHandler) Update(w http.ResponseWriter, r *http.Request) {
	var author types.AuthorUpdate
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	slug := r.PathValue("slug")
	result, err := h.store.GetBySlug(slug)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Author not found")
		return
	}

	updated, err := h.store.Update(result.ID, author)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, updated)
}

func (h *AuthorHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid author ID")
		return
	}

	if err := h.store.Delete(id); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Author deleted successfully"})
}
