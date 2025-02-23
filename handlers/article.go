package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"test-ai-api/stores"
	"test-ai-api/types"
	"test-ai-api/utils"
)

type ArticleHandler struct {
	store       *stores.ArticleStore
	authorStore *stores.AuthorStore
}

func NewArticleHandler(store *stores.ArticleStore, authorStore *stores.AuthorStore) *ArticleHandler {
	return &ArticleHandler{store: store, authorStore: authorStore}
}

func (h *ArticleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var article types.ArticleCreate
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := r.Context().Value("userID").(int64)

	// Check if user is an author
	author, err := h.authorStore.GetByUserID(userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusForbidden, "Only authors can create articles")
		return
	}

	result, err := h.store.Create(article, author.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, result)
}

func (h *ArticleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	articles, err := h.store.GetAll(10, 0) // TODO: Add pagination
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(articles) == 0 {
		utils.RespondWithJSON(w, http.StatusOK, []types.Article{})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, articles)
}

func (h *ArticleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid article ID")
		return
	}

	article, err := h.store.GetByID(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Article not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, article)
}

func (h *ArticleHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid article ID")
		return
	}

	var article types.ArticleUpdate
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Get user ID from context
	userID := r.Context().Value("userID").(int64)

	// Check if user is the author of this article
	existingArticle, err := h.store.GetByID(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Article not found")
		return
	}

	author, err := h.authorStore.GetByUserID(userID)
	if err != nil || author.ID != existingArticle.AuthorID {
		utils.RespondWithError(w, http.StatusForbidden, "Not authorized to update this article")
		return
	}

	updated, err := h.store.Update(id, article)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, updated)
}

func (h *ArticleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid article ID")
		return
	}

	// Get user ID from context
	userID := r.Context().Value("userID").(int64)

	// Check if user is the author of this article
	existingArticle, err := h.store.GetByID(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Article not found")
		return
	}

	author, err := h.authorStore.GetByUserID(userID)
	if err != nil || author.ID != existingArticle.AuthorID {
		utils.RespondWithError(w, http.StatusForbidden, "Not authorized to delete this article")
		return
	}

	if err := h.store.Delete(id); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Article deleted successfully"})
}
