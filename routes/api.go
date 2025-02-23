package routes

import (
	"database/sql"
	"net/http"
	"test-ai-api/handlers"
	"test-ai-api/middleware"
	"test-ai-api/stores"
)

func SetupRoutes(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	userStore := stores.NewUserStore(db)
	authHandler := handlers.NewAuthHandler(userStore)

	// Public routes
	mux.HandleFunc("POST /api/login", authHandler.Login)
	mux.HandleFunc("POST /api/register", authHandler.Register)

	authorStore := stores.NewAuthorStore(db)
	authorHandler := handlers.NewAuthorHandler(authorStore)

	// Protected routes
	mux.HandleFunc("GET /api/authors", authorHandler.GetAll)
	mux.HandleFunc("GET /api/authors/{slug}", authorHandler.GetBySlug)
	mux.HandleFunc("POST /api/authors", middleware.AuthMiddleware(authorHandler.Create))
	mux.HandleFunc("PUT /api/authors/{slug}", middleware.AuthMiddleware(authorHandler.Update))
	mux.HandleFunc("DELETE /api/authors/{slug}", middleware.AuthMiddleware(authorHandler.Delete))

	imageStore := stores.NewImageStore(db)
	imageHandler := handlers.NewImageHandler(imageStore)

	// public routes
	mux.HandleFunc("GET /api/images/{id}", imageHandler.GetById)

	// Protected routes
	mux.HandleFunc("POST /api/images", middleware.AuthMiddleware(imageHandler.Create))
	mux.HandleFunc("DELETE /api/images/{id}", middleware.AuthMiddleware(imageHandler.Delete))

	articleStore := stores.NewArticleStore(db)
	articleHandler := handlers.NewArticleHandler(articleStore, authorStore)

	// Public routes
	mux.HandleFunc("GET /api/articles", articleHandler.GetAll)
	mux.HandleFunc("GET /api/articles/{id}", articleHandler.GetByID)

	// Protected routes
	mux.HandleFunc("POST /api/articles", middleware.AuthMiddleware(articleHandler.Create))
	mux.HandleFunc("PUT /api/articles/{id}", middleware.AuthMiddleware(articleHandler.Update))
	mux.HandleFunc("DELETE /api/articles/{id}", middleware.AuthMiddleware(articleHandler.Delete))

	// Protected routes
	mux.HandleFunc("GET /api/me", middleware.AuthMiddleware(authHandler.GetCurrentUser))

	// Wrap the mux with CORS middleware
	handler := middleware.CorsMiddleware(mux)

	return handler
}
