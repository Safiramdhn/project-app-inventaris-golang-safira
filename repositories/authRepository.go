package repositories

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/Safiramdhn/project-app-inventaris-golang-safira/models"
)

type AuthRepository interface {
	Register(userDTO *models.UserDTO) (*models.User, error)
	Login(loginRequest *models.LoginRequest) (*models.User, error)
	CreateSession(sessionInput *models.Session) (*models.Session, error)
	ValidateSession(sessionToken string) (*models.Session, error)
	InvalidateSession(sessionToken string) error
	RefreshSession(oldSessionToken string) (*models.Session, error)
}

type authRepository struct {
	DB *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{DB: db}
}

// CreateSession implements AuthRepository.
func (a *authRepository) CreateSession(sessionInput *models.Session) (*models.Session, error) {
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

	sqlStatement := `INSERT INTO sessions (user_id, session_token, expires_at) VALUES ($1, $2, $3) RETURNING session_token`
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
func (a *authRepository) InvalidateSession(sessionToken string) error {
	tx, err := a.DB.Begin()
	if err != nil {
		log.Printf("Error starting invalidate transaction: %v\n", err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	updateStatement := `UPDATE sessions SET is_active = false WHERE session_token = $1`
	result, err := tx.Exec(updateStatement, sessionToken)
	if err != nil {
		log.Printf("Error updating session: %v\n", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching rows affected: %v\n", err)
		return err
	}
	if rowsAffected == 0 {
		log.Println("Session not found")
		return errors.New("session not found")
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Error committing validate session transaction: %v\n", err)
		return err
	}
	return nil
}

// Login implements AuthRepository.
func (a *authRepository) Login(loginRequest *models.LoginRequest) (*models.User, error) {
	sqlStatement := `SELECT id, username, email, password_hash FROM users WHERE username=$1`
	var user models.User
	err := a.DB.QueryRow(sqlStatement, loginRequest.Username).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
	if err == sql.ErrNoRows {
		log.Println("User not found")
		return nil, errors.New("user not found")
	} else if err != nil {
		log.Printf("Error querying user: %v\n", err)
		return nil, err
	}
	return &user, nil
}

// RefreshSession implements AuthRepository.
func (a *authRepository) RefreshSession(oldSessionToken string) (*models.Session, error) {
	panic("unimplemented")
}

// Register implements AuthRepository.
func (a *authRepository) Register(userDTO *models.UserDTO) (*models.User, error) {
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
func (a *authRepository) ValidateSession(sessionToken string) (*models.Session, error) {
	// Prepare SQL statement
	sqlStatement := `SELECT session_token, user_id, expires_at FROM sessions WHERE session_token = $1`

	var session models.Session

	// Execute the query
	err := a.DB.QueryRow(sqlStatement, sessionToken).Scan(&session.SessionToken, &session.UserID, &session.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid session token")
		}
		log.Printf("Error querying session: %v\n", err)
		return nil, err
	}

	// Check if the session has expired
	if session.ExpiresAt.Before(time.Now()) {
		// Invalidate the expired session
		if err := a.InvalidateSession(session.SessionToken); err != nil {
			log.Printf("Error invalidating expired session: %v\n", err)
			return nil, err
		}
		return nil, errors.New("session expired")
	}

	// Return the valid session
	return &session, nil
}
