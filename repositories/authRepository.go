package repositories

import (
	"database/sql"
	"log"

	"github.com/Safiramdhn/project-app-inventaris-golang-safira/models"
)

type AuthRepository interface {
	Register(userDTO *models.UserDTO) (*models.User, error)
	Login(loginRequest *models.LoginRequest) (*models.User, error)
	CreateSession(sessionInput *models.Session) (*models.Session, error)
	ValidateSession(sessionToken string) (*models.User, error)
	InvalidateSession(sessionToken string) error
	RefreshSession(oldSessionToken string) (*models.Session, error)
}

type authRepositoryImpl struct {
	DB *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepositoryImpl{DB: db}
}

// CreateSession implements AuthRepository.
func (a *authRepositoryImpl) CreateSession(sessionInput *models.Session) (*models.Session, error) {
	tx, err := a.DB.Begin()
	if err != nil {
		log.Printf("Error starting register transaction: %v\n", err)
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	sqlStatement := `INSERT INTO sessions VALUES (user_id, session_token, expired_at) VALUES ($1, $2, $3) RETURNING session_token`
	err = tx.QueryRow(sqlStatement, sessionInput.UserID, sessionInput.SessionToken, sessionInput.ExpiresAt).Scan(&sessionInput.SessionToken)
	if err != nil {
		log.Printf("Error inserting session: %v\n", err)
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing register transaction: %v\n", err)
		return nil, err
	}
	return sessionInput, nil
}

// InvalidateSession implements AuthRepository.
func (a *authRepositoryImpl) InvalidateSession(sessionToken string) error {
	panic("unimplemented")
}

// Login implements AuthRepository.
func (a *authRepositoryImpl) Login(loginRequest *models.LoginRequest) (*models.User, error) {
	sqlStatement := `SELECT id, username, email, password_hash FROM users WHERE username=$1`
	var user models.User
	err := a.DB.QueryRow(sqlStatement, loginRequest.Username, loginRequest.Password).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Printf("Error querying user: %v\n", err)
		return nil, err
	}
	return &user, nil
}

// RefreshSession implements AuthRepository.
func (a *authRepositoryImpl) RefreshSession(oldSessionToken string) (*models.Session, error) {
	panic("unimplemented")
}

// Register implements AuthRepository.
func (a *authRepositoryImpl) Register(userDTO *models.UserDTO) (*models.User, error) {
	tx, err := a.DB.Begin()
	if err != nil {
		log.Printf("Error starting register transaction: %v\n", err)
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	user := models.User{}
	sqlStatement := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id, email`
	err = tx.QueryRow(sqlStatement, userDTO.Username, userDTO.Email, userDTO.Password).Scan(&user.ID, &user.Email)
	if err != nil {
		log.Printf("Error inserting user: %v\n", err)
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing register transaction: %v\n", err)
		return nil, err
	}
	return &user, nil
}

// ValidateSession implements AuthRepository.
func (a *authRepositoryImpl) ValidateSession(sessionToken string) (*models.User, error) {
	panic("unimplemented")
}
