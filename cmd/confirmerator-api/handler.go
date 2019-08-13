package main

import (
	"bytes"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/bahelit/confirmerator/api/account"
	"github.com/bahelit/confirmerator/api/device"
	"github.com/bahelit/confirmerator/api/user"
)

// ErrResponse renderer type for handling errors.
//
// In the best case scenario, errors package helps reveal information on the error, setting it on Err,
// and in the Render() method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

var (
	ErrNotFound       = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}
	ErrBadRequest     = &ErrResponse{HTTPStatusCode: 406, StatusText: "Bad request."}
	ErrNotImplemented = &ErrResponse{HTTPStatusCode: 501, StatusText: "Not implemented."}
)

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		render.JSON(w, r, ErrBadRequest)
	}

	userID, err := user.UpdateUserAccount(client, buf)
	if err != nil {
		render.JSON(w, r, ErrBadRequest)
	} else {
		response := make(map[string]string)
		response["message"] = "Success"
		if len(userID) > 0 {
			response["id"] = userID
		}
		render.JSON(w, r, response)
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	selectedUser, err := user.GetUserAccount(client, userID)
	if err != nil {
		render.JSON(w, r, ErrNotFound)
	} else {
		render.JSON(w, r, selectedUser) // A chi router helper for serializing and returning json
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, ErrNotImplemented)
}

//func GetAllUsers(w http.ResponseWriter, r *http.Request) {
//	render.JSON(w, r, accounts)
//}

func UpdateAccount(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		render.JSON(w, r, ErrBadRequest)
	}

	accountID, err := account.UpdateAccount(client, buf)
	if err != nil {
		render.JSON(w, r, ErrBadRequest)
	} else {
		response := make(map[string]string)
		response["message"] = "Success"
		if len(accountID) > 0 {
			response["id"] = accountID
		}
		render.JSON(w, r, response)
	}
}

func GetAccount(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	selectedAccount, err := account.GetAccountsForUser(client, userID)
	if err != nil {
		render.JSON(w, r, ErrNotFound)
	} else {
		render.JSON(w, r, selectedAccount)
	}
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]string)
	response["message"] = "Deleted TODO successfully"
	render.JSON(w, r, response) // Return some demo response
}

func UpdateDevice(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		render.JSON(w, r, ErrBadRequest)
	}

	deviceID, err := device.UpdateDevice(client, buf)
	if err != nil {
		render.JSON(w, r, ErrBadRequest)
	} else {
		response := make(map[string]string)
		response["message"] = "Success"
		if len(deviceID) > 0 {
			response["id"] = deviceID
		}
		render.JSON(w, r, response)
	}
}

func GetDevice(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	selectedAccount, err := device.GetDevices(client, userID)
	if err != nil {
		render.JSON(w, r, ErrNotFound)
	} else {
		render.JSON(w, r, selectedAccount)
	}
}

func DeleteDevice(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "id")

	err := device.Delete(client, deviceID)
	if err != nil {
		render.JSON(w, r, ErrBadRequest)
	} else {
		response := make(map[string]string)
		response["message"] = "Success"
		render.JSON(w, r, response)
	}
}
