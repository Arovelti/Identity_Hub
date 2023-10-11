package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Arovelti/identityhub/playground/http/service"
	"github.com/Arovelti/identityhub/logger"
	models "github.com/Arovelti/identityhub/profile_service/models"
	"github.com/Arovelti/identityhub/repository"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

func init() {
	rand.New(rand.NewSource(rand.Int63()))
}

func main() {
	// Logger
	l := logger.InitLogger()

	dir, err := os.Getwd()
	if err != nil {
		l.Error("current directory can be reached")
	}

	l.Info("slog info message", slog.String("path", dir))
	l.Debug("slog debug message", slog.String("path", dir))

	// With context
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Without context
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		takeSig := <-sigChan
		cancel()
		l.Info("Shutting down gracefully", slog.String("signal", takeSig.String()))
	}()

	// Repository
	//repo := repository.NewInMemoryRepository()
	repo := TestRepositoryMethods(l)

	// Envs
	host := os.Getenv("LOCAL_HOST")
	port:= os.Getenv("LOCAL_HTTP_PORT")
	localhost := fmt.Sprintf("%s:%s", host, port)

	srv := service.New(&repo)

	// Handler
	mux := http.NewServeMux()
	mux.Handle("/login", srv.AuthMiddleware(http.HandlerFunc(securedHandler)))

	// common endpoints
	mux.HandleFunc("/profiles", srv.ListProfilesHandler)

	// admin endpoints
	mux.HandleFunc("/profile", srv.CreateProfileHandler)
	mux.HandleFunc("/profile/id/{id}", srv.GetProfilByIDeHandler)
	mux.HandleFunc("/profile/name/{username}", srv.GetProfileByUsernameHandler)
	// delete user

	// Server
	server := http.Server{
		Addr:           localhost,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	l.Info("Start serving on", slog.String(host, port))

	// Graseful shutdown

	group, groupCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return server.ListenAndServe()
	})
	group.Go(func() error {
		<-groupCtx.Done()
		return server.Shutdown(context.Background())
	})

	if err := group.Wait(); err != nil {
		l.Info("exit error group", slog.String("error", err.Error()))
	}
}

func securedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "welcome to the identity hub: this is a secure endpoint"}`))
}

func TestRepositoryMethods(l *slog.Logger) repository.Repository {
	var err error

	repo := repository.New()
	if err := repo.Create(&testProfile); err != nil {
		l.Error("create", slog.String("error", err.Error()))
	}

	repo.GenerateTestProfiles()

	if err := repo.Create(&testProfile2); err != nil {
		l.Error("create", slog.String("error", err.Error()))
	}
	p := repo.List()
	for _, v := range p {
		fmt.Printf("%+v\n", v)
	}

	p2, err2 := repo.GetByUsername("Test_Name")
	if err2 != nil {
		l.Error(err2.Error())
	}
	fmt.Println(p2)

	p3, err3 := repo.GetByID(uuid.MustParse("7f2de087-c871-409a-b93b-b4049bf46aef"))
	if err3 != nil {
		fmt.Println(err3)
	}
	fmt.Println(p3)

	p3.Email = "Updated_email"
	err = repo.Update(p3.ID, p3.Name, p3)
	if err != nil {
		fmt.Println(err)
	}

	err4 := repo.Delete(uuid.MustParse("7f2de087-c871-409a-b93b-b4049bf46aef"))
	if err4 != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Deleted - OK")
	}

	p = repo.List()
	for _, v := range p {
		fmt.Printf("%+v\n", v)
	}

	return repo
}

var testProfile = models.Profile{
	Name:     "Test_Name",
	Email:    "Test_Email",
	Password: "Test_Password",
	Admin:    true,
}

var testProfile2 = models.Profile{
	Name:     "Test_Name_2",
	Email:    "Test_Email_2",
	Password: "Test_Password_2",
	Admin:    false,
}

// r := httprouter.New()
// r.GET("/", HomeHandler)

// r.GET("/posts", PostsIndexHandler)
// r.POST("/posts", PostsCreateHandler)
