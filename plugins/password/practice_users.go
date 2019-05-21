package password

import (
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	password  string = "password12345!"
	apiKey    string = "1234567890key-admin"
	apiSecret string = "1234567890secret-admin"
)

// OAuth type of password
type OAuth struct {
	BCryptCost int
}

// User type of password
type User struct {
	BCryptCost int
}

// Retriever sets Practice Password
type Retriever interface {
	generatePasswordHash() ([]byte, error)
}

// Retrieve the password for a given type
func Retrieve(r Retriever) (string, error) {
	logContext := "[Plugins Practice Password Users]"
	passwordHash, err := r.generatePasswordHash()
	if err != nil {
		return "", errors.Wrap(err, logContext)
	}
	//log.Info("Password: ", string(passwordHash))
	return string(passwordHash), nil
}

// generatePasswordHash create the password hash for Admin and Internal Admin
func (o *OAuth) generatePasswordHash() ([]byte, error) {
	o.BCryptCost = 4
	p := []byte(password)
	passwordHash, err := bcrypt.GenerateFromPassword(p, o.BCryptCost)
	if err != nil {
		return nil, fmt.Errorf("Issue Generating Password Hash: %v", err)
	}

	return passwordHash, nil
}

// generatePasswordHash create the password hash for users
func (u *User) generatePasswordHash() error {
	u.BCryptCost = 10

	return nil
}
