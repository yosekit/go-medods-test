package database

import (
	"errors"

	db "goauth/internal/database/postgres"

	"golang.org/x/crypto/bcrypt"
)


func GetRefreshToken(userID string) ([]byte, error) {
	query := "SELECT hashed_token FROM refresh_tokens WHERE user_id=$1"

	row := db.Client.QueryRow(query, userID)

	var token []byte
	if err := row.Scan(&token); err != nil {
		return nil, err
	}

	return token, nil
}

func SaveRefreshToken(userID string, clientIP string, refreshToken []byte) error {
	query := "INSERT INTO refresh_tokens (user_id, client_ip, hashed_token) VALUES ($1, $2, $3) RETURNING id"

	stmt, err := db.Client.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	hashedRefreshToken, _ := bcrypt.GenerateFromPassword(refreshToken, bcrypt.DefaultCost)

	_, err = stmt.Exec(userID, clientIP, hashedRefreshToken)
	if err != nil {
		return err
	}

	return nil
}

func UpdateRefreshToken(userID string, clientIP string, refreshToken []byte) error {
	query := `
		UPDATE refresh_tokens
		SET hashed_token = $1, client_ip = $2
		WHERE user_id = $3
	`

	stmt, err := db.Client.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	hashedRefreshToken := bcryptHash(refreshToken)

	_, err = stmt.Exec(hashedRefreshToken, clientIP, userID)
	if err != nil {
		return err
	}

	return nil
}

func ValidateRefreshToken(userID string, refreshToken []byte) error {
	storedRefreshToken, err := GetRefreshToken(userID)
	if err != nil {
		return errors.New("Refresh token not found")
	}

	err = bcryptCompare(storedRefreshToken, refreshToken)
	if err != nil {
		return errors.New("Invalid refresh token")
	}

	return nil
}


func bcryptHash(token []byte) []byte {
	hashedToken, _ := bcrypt.GenerateFromPassword(token, bcrypt.DefaultCost)
	return hashedToken
}

func bcryptCompare(hashedToken []byte, token []byte) error {
	return bcrypt.CompareHashAndPassword(hashedToken, token)
}