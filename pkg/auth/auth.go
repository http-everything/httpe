package auth

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"net/http"

	"http-everything/httpe/pkg/rules"
)

// IsRequestAuthenticated checks if the request is authenticated based on the provided list of users, hashing algorithm,
// and the request's basic authentication credentials.
// If the list of users is empty, the request is considered authenticated.
// If the request is missing basic authentication credentials, the request is considered unauthenticated.
// For each user in the list, the request password is hashed using the specified algorithm and compared with the user's stored password.
// If a matching user and password is found, the request is considered authenticated.
// Returns a boolean indicating if the request is authenticated and any error that occurred during the authentication process.
func IsRequestAuthenticated(users []rules.User, hashing string, r *http.Request) (ok bool, err error) {
	if len(users) == 0 {
		// The rule doesn't require the request to be authenticated
		return true, nil
	}
	u, p, ok := r.BasicAuth()
	if !ok {
		// The request is missing basic authentication credentials
		return false, nil
	}
	for _, user := range users {
		password, err := hashPassword(p, hashing)
		if err != nil {
			return false, err
		}
		if user.Username == u && user.Password == password {
			return true, nil
		}
	}
	return false, nil
}

// hashPassword hashes the given password using the specified algorithm.
// If the algorithm is an empty string, the password is returned as is.
// If the algorithm is "sha256", the password is hashed using the SHA-256 algorithm.
// If the algorithm is "sha512", the password is hashed using the SHA-512 algorithm.
// Returns the hashed password as a hexadecimal string and any error that occurred during hashing.
// If an unknown hashing algorithm is specified, returns an error.
func hashPassword(password string, algorithm string) (hash string, err error) {
	switch algorithm {
	case "":
		return password, nil
	case "sha256":
		h := sha256.Sum256([]byte(password))
		return hex.EncodeToString(h[:]), nil
	case "sha512":
		h := sha512.Sum512([]byte(password))
		return hex.EncodeToString(h[:]), nil
	default:
		return hash, errors.New("unknown hashing algorithm")
	}
}
