package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/totorialman/avito-backend-autumn-2025/internal/config"
	userHandler "github.com/totorialman/avito-backend-autumn-2025/internal/handler/user"
	userRepository "github.com/totorialman/avito-backend-autumn-2025/internal/repository/user"
	userService "github.com/totorialman/avito-backend-autumn-2025/internal/service/user"

	teamHandler "github.com/totorialman/avito-backend-autumn-2025/internal/handler/team"
	teamRepository "github.com/totorialman/avito-backend-autumn-2025/internal/repository/team"
	teamService "github.com/totorialman/avito-backend-autumn-2025/internal/service/team"

	prHandler "github.com/totorialman/avito-backend-autumn-2025/internal/handler/pr"
	prRepository "github.com/totorialman/avito-backend-autumn-2025/internal/repository/pr"
	prService "github.com/totorialman/avito-backend-autumn-2025/internal/service/pr"

	statsHandler "github.com/totorialman/avito-backend-autumn-2025/internal/handler/stats"
	statsRepository "github.com/totorialman/avito-backend-autumn-2025/internal/repository/stats"
	statsService "github.com/totorialman/avito-backend-autumn-2025/internal/service/stats"
)

func main() {
	log.SetOutput(os.Stdout)

	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("Warning: %v", err)
	}

	port, err := getPort()
	if err != nil {
		log.Fatalf("Failed to get port: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	dbPool := config.MustInitDB(ctx)
	defer dbPool.Close()

	userRepository := userRepository.NewUserRepository(dbPool)
	teamRepository := teamRepository.NewTeamRepository(dbPool)
	prRepository := prRepository.NewPrRepository(dbPool)

	userService := userService.NewUserService(userRepository, prRepository)
	teamService := teamService.NewTeamService(teamRepository, userRepository)
	prService := prService.NewPrService(prRepository, userRepository)

	userHandler := userHandler.NewUserHandler(userService)
	teamHandler := teamHandler.NewTeamHandler(teamService)
	prHandler := prHandler.NewPrHandler(prService)

	statsRepository := statsRepository.NewStatsRepository(dbPool)
	statsService := statsService.NewStatsService(statsRepository)
	statsHandler := statsHandler.NewStatsHandler(statsService)

	router := setupRouter(userHandler, teamHandler, prHandler, statsHandler)

	if err := runServer(ctx, router, port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

const defaultPort = "8080"

func getPort() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		log.Printf("Using default port: %s", defaultPort)
		port = defaultPort
	}

	_, err := strconv.Atoi(port)
	if err != nil {
		return "", fmt.Errorf("invalid port (not an integer): %s", port)
	}
	return port, nil
}

func setupRouter(userHandler *userHandler.UserHandler, teamHandler *teamHandler.TeamHandler, prHandler *prHandler.PrHandler, statsHandler *statsHandler.StatsHandler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/pullRequest/create", prHandler.CreatePR).Methods("POST")
	r.HandleFunc("/pullRequest/merge", prHandler.MergePR).Methods("POST")
	r.HandleFunc("/pullRequest/reassign", prHandler.Reassign).Methods("POST")

	r.HandleFunc("/team/add", teamHandler.AddTeam).Methods("POST")
	r.HandleFunc("/team/get", teamHandler.GetTeam).Methods("GET")

	r.HandleFunc("/users/setIsActive", userHandler.SetActive).Methods("POST")
	r.HandleFunc("/users/getReview", userHandler.GetUserReviews).Methods("GET")

	r.HandleFunc("/stats/reviewer-assignments", statsHandler.GetReviewerAssignmentStats).Methods("GET")

	return r
}

func runServer(ctx context.Context, handler http.Handler, port string) error {
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Starting server on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server stopped with error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down service-courier")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	log.Println("Server exited gracefully")
	return nil
}
