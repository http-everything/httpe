package auth_test

import (
	"errors"
	"http-everything/httpe/pkg/auth"
	"http-everything/httpe/pkg/rules"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsRequestAuthenticated(t *testing.T) {
	users := []rules.User{
		{Username: "user", Password: "password"},
		{Username: "user2", Password: "5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8"},
		{Username: "user3", Password: "b109f3bbbc244eb82441917ed06d618b9008dd09b3befd1b5e07394c706a8bb" +
			"980b1d7785e5976ec049b46df5f1326af5a2ea6d103fd07c95385ffab0cacbc86"},
	}
	cases := []struct {
		name     string
		users    []rules.User
		hashing  string
		username string
		password string
		wantOk   bool
		wantErr  error
	}{
		{
			name:     "No Users",
			users:    nil,
			hashing:  "",
			username: "user",
			password: "password",
			wantOk:   true,
			wantErr:  nil,
		},
		{
			name:     "No Auth",
			users:    users,
			hashing:  "",
			username: "",
			password: "",
			wantOk:   false,
			wantErr:  nil,
		},
		{
			name:     "Incorrect User",
			users:    users,
			hashing:  "",
			username: "wrongUser",
			password: "password",
			wantOk:   false,
			wantErr:  nil,
		},
		{
			name:     "Incorrect Password",
			users:    users,
			hashing:  "",
			username: "user",
			password: "wrongPassword",
			wantOk:   false,
			wantErr:  nil,
		},
		{
			name:     "Correct User and Clear-Text-Password",
			users:    users,
			hashing:  "",
			username: "user",
			password: "password",
			wantOk:   true,
			wantErr:  nil,
		},
		{
			name:     "Correct User and SHA256-Password",
			users:    users,
			hashing:  "sha256",
			username: "user2",
			password: "password",
			wantOk:   true,
			wantErr:  nil,
		},
		{
			name:     "Correct User and SHA512-Password",
			users:    users,
			hashing:  "sha512",
			username: "user3",
			password: "password",
			wantOk:   true,
			wantErr:  nil,
		},
		{
			name:     "Wrong Algorithm",
			users:    users,
			hashing:  "wrong-algorithm",
			username: "user",
			password: "password",
			wantOk:   false,
			wantErr:  errors.New("unknown hashing algorithm"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup request with Basic Auth
			r, _ := http.NewRequest("GET", "/", nil)
			if tc.username != "" && tc.password != "" {
				r.SetBasicAuth(tc.username, tc.password)
			}

			ok, err := auth.IsRequestAuthenticated(tc.users, tc.hashing, r)

			assert.Equal(t, tc.wantOk, ok)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
