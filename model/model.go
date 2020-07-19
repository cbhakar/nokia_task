package model

import (
	"errors"
	"net/http"
	"strconv"
)

var (
	DefaultPageNum   = 1
	DefaultPageSize  = 10
	MaxPageSize      = 100

	ErrPageSizeNotAllowed  = errors.New("page size invalid")
	ErrInvalidPagination   = errors.New("invalid pagination data")
	ErrInvalidPgQueryValue = errors.New("invalid pagination values in query")
	)

type User struct {
	Id       uint64 `json:"id,omitempty" db:"id"`
	Uid      string `json:"uid,omitempty" db:"uid"`
	Name     string `json:"name" db:"name"`
	Age      int    `json:"age" db:"age"`
	Address  string `json:"address" db:"address"`
	MobileNo string `json:"mobile_no" db:"mobile"`
}
type LimitedUsers struct {
	Users   []User	`json:"users"`
	Count    int64 `json:"count"`
}

type Pagination struct {
	PageNum  int `query:"page"`
	PageSize int `query:"size"`
}

func (gr *User) Validate() (err error) {

	if len(gr.Name) < 1 {
		return errors.New("the 'name' field is required")
	}

	if gr.Age < 18 && gr.Age > 65 {
		return errors.New("age must be in range 18-65")
	}
	if len(gr.Address) < 1 {
		return errors.New("the 'address' field is required")
	}
	if len(gr.MobileNo) != 10 {
		return errors.New("mobile number must contain 10 digits")
	}
	return
}

func (p *Pagination) BindFrom(r *http.Request) (err error) {
	pageQP := r.URL.Query()["page"]
	sizeQP := r.URL.Query()["size"]
	var pg, sz int64
	if len(pageQP) != 0 {
		pg, err = strconv.ParseInt(pageQP[0], 10, 32)
		if err != nil {
			return ErrInvalidPgQueryValue
		}
	}

	if len(sizeQP) != 0 {
		sz, err = strconv.ParseInt(sizeQP[0], 10, 32)
		if err != nil {
			return ErrInvalidPgQueryValue
		}
	}

	p.PageNum = int(pg)
	p.PageSize = int(sz)

	err = p.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (p *Pagination) Validate() error {
	if p.PageNum < 1 {
		p.PageNum = DefaultPageNum
	}

	if p.PageSize < 1 {
		p.PageSize = DefaultPageSize
	}

	if p.PageSize > MaxPageSize {
		return ErrPageSizeNotAllowed
	}
	return nil
}

func (p Pagination) GetOffset() (offset int) {
	p.mustValid()
	offset = ((int(p.PageNum) - 1) * int(p.PageSize))
	return
}

func (p Pagination) mustValid() {
	if p.PageNum < 1 || p.PageSize < 1 || p.PageSize > MaxPageSize {
		panic( "incorrect pagination")
	}
}