package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	models "github.com/Arovelti/identityhub/profile_service/models"
	"github.com/Arovelti/identityhub/repository"
	"github.com/google/uuid"
)

type Service struct {
	Repo      repository.Repository
	AuthToken string
}

func New(repo *repository.Repository) Service {
	return Service{
		Repo:      *repo,
		AuthToken: "",
	}
}

func (s *Service) CreateProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var p models.Profile
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode JSON: %v", err), http.StatusBadRequest)
		return
	}

	if !p.Admin {
		http.Error(w, "Access denied", http.StatusForbidden)
	}

	err = s.Repo.Create(&p)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create profile: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Service) GetProfilByIDeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid profile ID", http.StatusBadRequest)
		return
	}

	profile, err := s.Repo.GetByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get profile: %v", err), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(profile)
}

func (s *Service) ListProfilesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	profiles := s.Repo.List()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(profiles)
}

func (s *Service) GetProfileByUsernameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	profile, err := s.Repo.GetByUsername(username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get profile by username: %v", err), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(profile)
}
