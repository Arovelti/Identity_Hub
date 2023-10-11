package repository

import (
	"errors"
	"sync"

	models "github.com/Arovelti/identityhub/profile_service/models"
	"github.com/google/uuid"
)

type Repository interface {
	// Create and edit profile
	Create(profile *models.Profile) error
	Update(id uuid.UUID, name string, profile *models.Profile) error
	Delete(id uuid.UUID) error

	// Get profile or profiles
	GetByID(id uuid.UUID) (*models.Profile, error)
	GetByUsername(username string) (*models.Profile, error)
	List() []*models.Profile

	// Generate Test Profiles
	GenerateTestProfiles()
}

type InMemoryRepository struct {
	mutex sync.RWMutex
	users map[models.UniqueInfo]*models.Profile
}

func New() Repository {
	return &InMemoryRepository{
		users: make(map[models.UniqueInfo]*models.Profile),
	}
}

// Create - creates new user
func (r *InMemoryRepository) Create(profile *models.Profile) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	uniqueInfo := models.UniqueInfo{
		ID:   profile.ID,
		Name: profile.Name,
	}

	if _, exist := r.users[uniqueInfo]; exist {
		return errors.New("user with the same ID or username already exists")
	}

	profile.ID = uuid.New()
	r.users[uniqueInfo] = profile
	return nil
}

// GetByID - returns user by ID
func (r *InMemoryRepository) GetByID(id uuid.UUID) (*models.Profile, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, v := range r.users {
		if v.ID == id {
			return v, nil
		}
	}

	return nil, errors.New("user not found")
}

// List - returns all profiles
func (r *InMemoryRepository) List() []*models.Profile {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	profiles := make([]*models.Profile, 0, len(r.users))
	for _, profile := range r.users {
		profiles = append(profiles, profile)
	}

	return profiles
}

// GetByUsername - returns pforile by username
func (r *InMemoryRepository) GetByUsername(username string) (*models.Profile, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Name == username {
			return user, nil
		}
	}

	return nil, errors.New("user not found")
}

// Update - update user profile
func (r *InMemoryRepository) Update(id uuid.UUID, name string, profile *models.Profile) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	uniqueInfo := models.UniqueInfo{
		ID:   id,
		Name: name,
	}

	if _, exists := r.users[uniqueInfo]; !exists {
		return errors.New("user not found")
	}

	r.users[uniqueInfo] = profile
	return nil
}

// Delete - delete user by ID
func (r *InMemoryRepository) Delete(id uuid.UUID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, user := range r.users {
		if user.ID == id {
			delete(r.users, models.UniqueInfo{ID: id, Name: user.Name})
			return nil
		}
	}

	return nil
}

// GenerateTestProfiles
func (r *InMemoryRepository) GenerateTestProfiles() {
	r.users[models.UniqueInfo{
		ID:   uuid.MustParse("7f2de087-c871-409a-b93b-b4049bf46aef"),
		Name: "Generated_Test_User",
	}] = &models.Profile{
		ID:       uuid.MustParse("7f2de087-c871-409a-b93b-b4049bf46aef"),
		Name:     "Generated_Test_User",
		Email:    "Generated_Test_Usser_Email",
		Password: "Generated_Test_User_Password",
		Admin:    true,
	}

	r.users[models.UniqueInfo{
		ID:   uuid.MustParse("7e34649d-765e-4418-a7a8-35434a003188"),
		Name: "Test_Name_1",
	}] = &models.Profile{
		Name:     "Test_Name_1",
		ID:       uuid.MustParse("7e34649d-765e-4418-a7a8-35434a003188"),
		Email:    "Test_Email_1",
		Password: "Test_Password_1",
		Admin:    true,
	}

	r.users[models.UniqueInfo{
		ID:   uuid.MustParse("f77d0f18-6851-11ee-8c99-0242ac120002"),
		Name: "Test_Name_2",
	}] = &models.Profile{
		ID:       uuid.MustParse("f77d0f18-6851-11ee-8c99-0242ac120002"),
		Name:     "Test_Name_2",
		Email:    "Test_Email_2",
		Password: "Test_Password_2",
		Admin:    false,
	}
}
