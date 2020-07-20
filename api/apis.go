package api

import (
	"encoding/json"
	"errors"
	"github.com/matryer/respond"
	"net/http"
	"nokia_task/model"
	"nokia_task/service"
)

func AddUserHandler(w http.ResponseWriter, r *http.Request) {

	req := model.User{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = errors.New("unable to get payload")
	}
	defer r.Body.Close()

	if err := req.Validate(); err != nil {
		err := map[string]interface{}{"message": err.Error()}
		w.Header().Set("Content-type", "applciation/json")
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}

	resp, err := service.AddUser(req)
	if err != nil {
		if err.Error() == "user with same name already exists" {
			err := map[string]interface{}{"message": err.Error()}
			w.Header().Set("Content-type", "applciation/json")
			respond.With(w, r, http.StatusConflict, err)
		} else {
			err := map[string]interface{}{"message": err.Error()}
			w.Header().Set("Content-type", "applciation/json")
			respond.With(w, r, http.StatusBadRequest, err)
		}
	} else {
		w.Header().Set("Content-type", "applciation/json")
		respond.With(w, r, http.StatusCreated, resp)
		return
	}
}

func ReloadHandler(w http.ResponseWriter, r *http.Request) {
	err := service.ReloadDataToRedis()
	if err != nil {
		err := map[string]string{"message": err.Error()}
		w.Header().Set("Content-type", "applciation/json")
		respond.With(w, r, http.StatusBadRequest, err)
		return
	} else {
		w.Header().Set("Content-type", "applciation/json")
		respond.With(w, r, http.StatusCreated, map[string]string{"status": "reload successful"})
		return
	}
}

func GetUserWithPaginationHandler(w http.ResponseWriter, r *http.Request) {

	pagination := model.Pagination{}
	err := pagination.BindFrom(r)
	if err != nil {
		err := map[string]string{"message": err.Error()}
		w.Header().Set("Content-type", "applciation/json")
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}
	resp, err := service.GetUser(pagination.PageSize, pagination.GetOffset())
	if err != nil {
		err := map[string]string{"message": err.Error()}
		w.Header().Set("Content-type", "applciation/json")
		respond.With(w, r, http.StatusNotFound, err)
	} else {
		w.Header().Set("Content-type", "applciation/json")
		respond.With(w, r, http.StatusOK, resp)
	}
}
