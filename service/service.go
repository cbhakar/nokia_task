package service

import (
	"database/sql"
	"errors"
	"github.com/segmentio/ksuid"
	"log"
	"nokia_task/model"
	"nokia_task/redis"
	"nokia_task/store"
)

func GetUID() string {
	return ksuid.New().String()
}

func AddUser(req model.User) (resp model.User, err error) {

	req.Uid = GetUID()

	count, err := store.GetUserCountByName(req.Name, req.MobileNo)
	if err != nil && err != sql.ErrNoRows {
		err = errors.New("error getting user details")
		return
	} else {
		err = nil
	}
	if count > 0 {
		err = errors.New("user with same name already exists")
		return
	}

	r, err := store.AddNewUser(store.User{
		Uid:      req.Uid,
		Name:     req.Name,
		Age:      req.Age,
		Address:  req.Address,
		MobileNo: req.MobileNo,
	})
	if err != nil {
		err = errors.New("error inserting new user record")
		return
	}
	err = redis.SetUserDataToRedis(r)

	resp.Uid = req.Uid
	resp.Name = req.Name
	resp.Age = req.Age
	resp.Address = req.Address
	resp.MobileNo = req.MobileNo

	return
}

func GetUser(limit int, offset int) (r model.LimitedUsers, err error) {

	users, count, err := redis.GetUserDataFromRedis(offset, (limit+offset)-1)
	if err != nil || len(users) == 0 {
		users, err = store.GetUsersWithPagination(limit, offset)
		if len(users) < 1 {
			err = errors.New("no user found")
			return
		}
		if err != nil {
			err = errors.New("unable to fetch user details")
			return
		}
		count, err = store.GetAllUsersCount()
		if err != nil {
			err = errors.New("unable to fetch user details")
			return
		}

	}
	if len(users) > 0 {
		resp := make([]model.User, len(users))
		for idx, r := range users{
			resp[idx].Uid = r.Uid
			resp[idx].Name = r.Name
			resp[idx].Age = r.Age
			resp[idx].Address = r.Address
			resp[idx].MobileNo = r.MobileNo
		}
		r.Users = resp
		r.Count = count
	}else {
		err = errors.New("no data found")
	}

	return
}

func ReloadDataToRedis() (err error) {
	log.Println("flushing redis")
	err = redis.DeleteAllUserFromRedis()
	if err != nil {
		err = errors.New("error flushing cache data")
		return
	}
	log.Println("restoring user details to redis")
	allUsers, err := store.GetAllUsers()
	if err == sql.ErrNoRows {
		err = errors.New("no user found ")
		return
	}
	if err != nil {
		err = errors.New("unable to fetch user details")
		return
	}

	for _, user := range allUsers {
		go redis.SetUserDataToRedis(user)
	}
	return
}

