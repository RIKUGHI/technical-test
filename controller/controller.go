package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/rikughi/technical-test/db"
	"github.com/rikughi/technical-test/dto"
	"github.com/rikughi/technical-test/helper"
	"github.com/rikughi/technical-test/model"
)

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) SetupA(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		c.Read(w, r)
	} else if r.Method == http.MethodPost {
		c.Create(w, r)
	} else {
		errorResponse := map[string]string{
			"error": "Method not allowed",
		}

		helper.WriteToResponseBody(w, errorResponse, http.StatusOK)
	}
}

func (c *Controller) SetupB(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		c.GetById(w, r)
	} else if r.Method == http.MethodPut {
		c.Update(w, r)
	} else if r.Method == http.MethodDelete {
		c.Delete(w, r)
	} else {
		errorResponse := map[string]string{
			"error": "Method not allowed",
		}

		helper.WriteToResponseBody(w, errorResponse, http.StatusOK)
	}
}

func (c *Controller) Read(w http.ResponseWriter, r *http.Request) {
	modelQ := r.URL.Query().Get("model")
	techQ := r.URL.Query().Get("tech")

	var filteredCollection []*model.Data
	for _, item := range db.Collection.Data {
		if (modelQ == "" || item.Model == modelQ) && (techQ == "" || db.Contains(item.Tech, techQ)) {
			filteredCollection = append(filteredCollection, item)
		}
	}

	helper.WriteToResponseBody(w, filteredCollection, http.StatusOK)
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	dto := new(dto.DataCreateRequest)

	errorResponse := map[string]any{
		"error": "",
	}

	err := helper.ReadFromRequestBody(r, dto)
	if err != nil {
		errorResponse["error"] = fmt.Sprintf("Invalid JSON: %+v", err)
		helper.WriteToResponseBody(w, errorResponse, http.StatusOK)
		return
	}

	if dto.Code == "" || dto.Name == "" || dto.Model == "" || dto.Status == "" || len(dto.Tech) == 0 {
		errorResponse["error"] = "Missing required fields: Code, Name, Model, Tech, and Status are required."
		helper.WriteToResponseBody(w, errorResponse, http.StatusBadRequest) // HTTP 400 Bad Request
		return
	}

	for _, item := range db.Collection.Data {
		if item.Code == dto.Code {
			errorResponse["error"] = fmt.Sprintf("Code %s already exists.", dto.Code)
			helper.WriteToResponseBody(w, errorResponse, http.StatusBadRequest)
			return
		}
	}

	newItem := &model.Data{
		Code:        dto.Code,
		Name:        dto.Name,
		Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua ",
		Model:       dto.Model,
		Tech:        dto.Tech,
		Status:      dto.Status,
	}

	db.Collection.Data = append(db.Collection.Data, newItem)

	helper.WriteToResponseBody(w, newItem, http.StatusOK)
}

func (c *Controller) GetById(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/data/")

	errorResponse := map[string]any{
		"error": "",
	}

	data := db.FirstOrNil(code)

	if data != nil {
		helper.WriteToResponseBody(w, data, http.StatusOK)
		return
	}

	errorResponse["error"] = fmt.Sprintf("Code %s is not exists.", code)
	helper.WriteToResponseBody(w, errorResponse, http.StatusBadRequest)
}

func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/data/")
	dto := new(dto.DataCreateRequest)

	errorResponse := map[string]any{
		"error": "",
	}

	err := helper.ReadFromRequestBody(r, dto)
	if err != nil {
		errorResponse["error"] = fmt.Sprintf("Invalid JSON: %+v", err)
		helper.WriteToResponseBody(w, errorResponse, http.StatusOK)
		return
	}

	data := db.FirstOrNil(code)

	if data == nil {
		errorResponse["error"] = fmt.Sprintf("Code %s is not exists.", code)
		helper.WriteToResponseBody(w, errorResponse, http.StatusBadRequest)
		return
	}

	data.Code = dto.Code
	data.Name = dto.Name
	data.Description = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua "
	data.Model = dto.Model
	data.Tech = dto.Tech
	data.Status = dto.Status

	if err := db.UpdateByCode(code, data); err != nil {
		errorResponse["error"] = err.Error()
		helper.WriteToResponseBody(w, errorResponse, http.StatusBadRequest)
		return
	}

	helper.WriteToResponseBody(w, map[string]any{
		"data": "Update berhasil",
	}, http.StatusOK)
}

func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/data/")

	errorResponse := map[string]any{
		"error": "",
	}

	data := db.FirstOrNil(code)

	if data == nil {
		errorResponse["error"] = fmt.Sprintf("Code %s is not exists.", code)
		helper.WriteToResponseBody(w, errorResponse, http.StatusBadRequest)
		return
	}

	if err := db.DeleteByCode(code); err != nil {
		errorResponse["error"] = err.Error()
		helper.WriteToResponseBody(w, errorResponse, http.StatusBadRequest)
		return
	}

	helper.WriteToResponseBody(w, map[string]any{
		"data": "Delete berhasil",
	}, http.StatusOK)
}
