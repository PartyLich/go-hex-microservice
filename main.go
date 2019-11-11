package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/PartyLich/hex-microservice/api"
	mongoRepo "github.com/PartyLich/hex-microservice/repository/mongodb"
	redisRepo "github.com/PartyLich/hex-microservice/repository/redis"
	"github.com/PartyLich/hex-microservice/shortUrl"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func httpPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	return fmt.Sprintf(":%s", port)
}

func chooseRepo() shortUrl.RedirectRepository {
	fmt.Println(os.Getenv("URL_DB"))
	switch os.Getenv("URL_DB") {
	case "redis":
		redisURL := os.Getenv("REDIS_URL")
		repo, err := redisRepo.NewRedisRepository(redisURL)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	case "mongo":
		fallthrough
	default:
		mongoURL := os.Getenv("MONGO_URL")
		mongodb := os.Getenv("MONGO_DB")
		mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
		repo, err := mongoRepo.NewMongoRepository(mongoURL, mongodb, mongoTimeout)
		if err != nil {
			log.Fatal(err)
		}

		return repo
	}
}

func main() {
	repo := chooseRepo()
	service := shortUrl.NewRedirectService(repo)
	handler := api.NewHandler(service)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Routes
	router.Get("/{code}", handler.Get)
	router.Post("/", handler.Post)

	errs := make(chan error, 2)

	// Start http server
	go func() {
		port := httpPort()
		fmt.Printf("Listening on port %s", port)
		errs <- http.ListenAndServe(port, router)
	}()

	// listen for exit notification
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	// drain error channel
	fmt.Printf("Terminated %s", <-errs)
}
