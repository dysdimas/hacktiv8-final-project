package controllers

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/mervick/aes-everywhere/go/aes256"
	"github.com/sirupsen/logrus"
	"hacktiv8-golang-assignment-final/api"
	"hacktiv8-golang-assignment-final/utils"
	"net/http"
	"strings"
)

type AuthController struct {
	store  *sessions.CookieStore
	logger *logrus.Logger
}

func NewAuthController(store *sessions.CookieStore, logger *logrus.Logger) *AuthController {
	return &AuthController{
		store:  store,
		logger: logger,
	}
}

func (c *AuthController) CheckSession(w http.ResponseWriter, r *http.Request) {
	session, _ := c.store.Get(r, utils.SessionKey)

	// validate session is no data
	if len(session.Values) == 0 {
		api.Err401Unauthorized(w, "no session")
	} else {
		api.Ok(w, map[string]string{
			"username": fmt.Sprintf("%v", session.Values["username"]),
			"name":     fmt.Sprintf("%v", session.Values["name"]),
		}, nil)
	}
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimSpace(r.PostFormValue("username"))
	password := strings.TrimSpace(r.PostFormValue("password"))
	name := strings.TrimSpace(r.PostFormValue("name"))

	if username == "" || password == "" || name == "" {
		api.Err400BR(w, "please fill form")
		return
	} else {

		// Check Duplicate Data
		var IsUserExist = false
		if Users != nil && len(Users) > 0 {
			for _, user := range Users {
				if strings.ToLower(user.Email) == strings.ToLower(username) {
					IsUserExist = true
					break
				}
			}
		}

		if IsUserExist {
			api.Err400BR(w, fmt.Sprintf("user [%s] already exist", username))
			return
		} else {
			// Register Process
			new_data := UserDto{
				Email:    username,
				Password: aes256.Encrypt(password, utils.SessionKey),
				Name:     name,
			}
			Users = append(Users, new_data)
			api.Ok(w, nil, nil)
		}
	}
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimSpace(r.PostFormValue("username"))
	password := strings.TrimSpace(r.PostFormValue("password"))

	if username == "" || password == "" {
		api.Err400BR(w, "please fill form")
		return
	} else {
		// Check User Data
		var userData = &UserDto{}
		userData = nil
		if Users != nil && len(Users) > 0 {
			for _, user := range Users {
				if strings.ToLower(user.Email) == strings.ToLower(username) && aes256.Decrypt(user.Password, utils.SessionKey) == password {
					userData = &user
					break
				}
			}
		}

		if userData == nil {
			api.Err400BR(w, "invalid user or password")
			return
		} else {

			// Login Success
			session, _ := c.store.Get(r, utils.SessionKey)

			session.Values["username"] = userData.Email
			session.Values["password"] = userData.Password
			session.Values["name"] = userData.Name

			// Store session
			err := session.Save(r, w)
			if err != nil {
				api.Err500ISE(w, err.Error())
				return
			} else {
				api.Ok(w, nil, nil)
			}

		}
	}
}

func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	// Process get session
	session, _ := c.store.Get(r, utils.SessionKey)
	// Process to expired session
	session.Options.MaxAge = -1
	err := session.Save(r, w)
	if err != nil {
		api.Err500ISE(w, err.Error())
		return
	}
	api.Ok(w, nil, nil)
}
