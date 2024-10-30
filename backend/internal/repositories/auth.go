package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/legendary-acp/chimecast/internal/models"
	"github.com/legendary-acp/chimecast/internal/utils"
	"github.com/mattn/go-sqlite3"
)

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (a *AuthRepository) RegisterUser(user models.User) error {
	// Use transaction to handle user creation
	tx, err := a.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Attempt to insert user
	_, err = tx.Exec("INSERT INTO user (userName, id, email, name, hashedPassword, createdAt) VALUES (?, ?, ?, ?, ?, ?)", user.UserName, user.ID, user.Email, user.Name, user.HashedPassword, user.CreatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			log.Println("User already exists:", user)
			return utils.ErrUserAlreadyExists
		}
		return fmt.Errorf("failed to register user: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Println("User registered:", user)
	return nil
}

func (a *AuthRepository) Login(userName string) (*models.User, error) {
	var user models.User

	// Prepare and execute the SQL statement
	stmt, err := a.db.Prepare("SELECT userName, email, name, hashedPassword, createdAt FROM user WHERE username = ?")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query: %v", err)
	}
	defer stmt.Close()

	// Execute the query
	err = stmt.QueryRow(userName).Scan(&user.UserName, &user.Email, &user.Name, &user.HashedPassword, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to query user: %v", err)
	}

	return &user, nil
}

func isUniqueViolation(err error) bool {
	if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.Code == sqlite3.ErrConstraint {
		return true
	}
	return false
}
