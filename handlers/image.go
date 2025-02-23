package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"test-ai-api/stores"
	"test-ai-api/types"
	"test-ai-api/utils"
)

type ImageHandler struct {
	store *stores.ImageStore
}

func NewImageHandler(store *stores.ImageStore) *ImageHandler {
	return &ImageHandler{store: store}
}

func (h *ImageHandler) Create(w http.ResponseWriter, r *http.Request) {
	var image types.ImageCreate
	if err := json.NewDecoder(r.Body).Decode(&image); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	result, err := h.store.Create(image)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, result)
}

func (h *ImageHandler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid image ID")
		return
	}
	image, err := h.store.GetByID(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, image)
}

func (h *ImageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid image ID")
		return
	}

	if err := h.store.Delete(id); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Image deleted successfully"})
}
