package services

import (
	"errors"
	"log"
	"time"

	"github.com/Safiramdhn/project-app-inventaris-golang-safira/models"
	"github.com/Safiramdhn/project-app-inventaris-golang-safira/repositories"
	"github.com/Safiramdhn/project-app-inventaris-golang-safira/utils"
	"github.com/Safiramdhn/project-app-inventaris-golang-safira/validations"
)

type AuthService struct {
	AuthRepo repositories.AuthRepository
}

func NewAuthService(repo repositories.AuthRepository) *AuthService {
	return &AuthService{AuthRepo: repo}
}

func (as *AuthService) RegisterUser(userDTO *models.UserDTO) (*models.User, error) {
	// check email and username not empty
	if userDTO.Email == "" || userDTO.Username == "" {
		log.Println("email and username is required")
		return nil, errors.New("email and username is required")
	}
	// check password not empty then hash it
	if userDTO.Password == "" {
		log.Println("password is required")
		return nil, errors.New("password is required")
	}
	hashedPassword, err := validations.HashPassword(userDTO.Password)
	if err != nil {
		log.Println("error hashing password")
		return nil, err
	}
	userDTO.Password = hashedPassword

	return as.AuthRepo.Register(userDTO)
}

func (as *AuthService) LoginUser(loginRequest *models.LoginRequest) (*models.Session, error) {
	user, err := as.AuthRepo.Login(loginRequest)
	if err != nil {
		return nil, err
	}

	passwordValidation := validations.CheckPassword(user.PasswordHash, loginRequest.Password)
	if !passwordValidation {
		return nil, errors.New("invalid password")
	}
	sessionInput := models.Session{}
	sessionInput.UserID = user.ID
	sessionInput.SessionToken = utils.GenerateToken()
	sessionInput.ExpiresAt = time.Now().Add(time.Hour * 6)

	session, err := as.AuthRepo.CreateSession(&sessionInput)
	if err != nil {
		log.Printf("error creating session: %v\n", err.Error())
		return nil, err
	}
	return session, nil
}

func (as *AuthService) GetSession(sessionToken string) (*models.Session, error) {
	return as.AuthRepo.ValidateSession(sessionToken)
}
