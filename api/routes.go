// api/routes.go
package api

import (
	"net/http"
	"os"

	"github.com/arcedo/financial-ai-backend/api/handlers"
	"github.com/arcedo/financial-ai-backend/api/helpers"
	"github.com/arcedo/financial-ai-backend/api/middlewares"
)

func (s *Server) setupRoutes() {
	router := http.NewServeMux()

	authMiddleware := middlewares.JWTAuthMiddleware([]byte(os.Getenv("SECRET")), s.store)

	router.HandleFunc("/login", helpers.MakeHTTPHandleFunc(handlers.Login, s.store, []string{"POST"}))
	router.HandleFunc("/register", helpers.MakeHTTPHandleFunc(handlers.CreateUser, s.store, []string{"POST"}))
	/*router.HandleFunc("/user",
		authMiddleware(
			helpers.MakeHTTPHandleFunc(handlers.GetUser, s.store, []string{"GET"}),
		),
	)*/

	router.HandleFunc("/transaction", authMiddleware(helpers.MakeHTTPHandleFunc(handlers.CreateTransaction, s.store, []string{"POST"})))
	router.HandleFunc("/transactions", authMiddleware(helpers.MakeHTTPHandleFunc(handlers.GetTransactions, s.store, []string{"GET"})))

	router.HandleFunc("/stocks", helpers.MakeHTTPHandleFunc(handlers.GetAllStocks, s.store, []string{"GET"}))

	router.HandleFunc("/update-profile", authMiddleware(helpers.MakeHTTPHandleFunc(handlers.UpdateUserProfile, s.store, []string{"GET"})))
	router.HandleFunc("/get-recommendations", authMiddleware(helpers.MakeHTTPHandleFunc(handlers.GetRecommendations, s.store, []string{"GET"})))
	router.HandleFunc("/advice", authMiddleware(helpers.MakeHTTPHandleFunc(handlers.GetAdvice, s.store, []string{"GET"})))
	router.HandleFunc("/asset-recommendation/{symbol}", authMiddleware(helpers.MakeHTTPHandleFunc(handlers.GetAssetRecommendation, s.store, []string{"GET"})))
	// Debugging routes
	//router.HandleFunc("/users", helpers.MakeHTTPHandleFunc(handlers.GetAllUsers, s.store, []string{"GET"}))
	//router.HandleFunc("/all-transactions", helpers.MakeHTTPHandleFunc(handlers.GetAllTransactions, s.store, []string{"GET"}))
	//router.HandleFunc("/products", helpers.MakeHTTPHandleFunc(handlers.GetProducts, s.store, []string{"GET"}))
	s.router = router
}
