package handler

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/notblessy/model"
	"github.com/notblessy/utils"
	"gorm.io/gorm"
)

func (h *Handler) CreateUser(c echo.Context) error {
	var req model.User

	if err := c.Bind(&req); err != nil {
		return utils.Response(c, http.StatusBadRequest, &utils.HTTPResponse{
			Message: fmt.Sprintf("error bind request: %s", err.Error()),
		})
	}

	if err := c.Validate(&req); err != nil {
		return utils.Response(c, http.StatusBadRequest, &utils.HTTPResponse{
			Message: fmt.Sprintf("error validate: %s", err.Error()),
		})
	}

	err := h.db.Save(&req).Error
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, &utils.HTTPResponse{
			Message: err.Error(),
		})
	}

	return utils.Response(c, http.StatusCreated, &utils.HTTPResponse{
		Data: req.ID,
	})
}

func (h *Handler) CreateMates(c echo.Context) error {
	var req model.User

	if err := c.Bind(&req); err != nil {
		return utils.Response(c, http.StatusBadRequest, &utils.HTTPResponse{
			Message: fmt.Sprintf("error bind request: %s", err.Error()),
		})
	}

	if err := c.Validate(&req); err != nil {
		return utils.Response(c, http.StatusBadRequest, &utils.HTTPResponse{
			Message: fmt.Sprintf("error validate: %s", err.Error()),
		})
	}

	ownerID := c.QueryParam("owner_id")

	tx := h.db.Begin()

	err := tx.Create(&req).Error
	if err != nil {
		tx.Rollback()
		return utils.Response(c, http.StatusInternalServerError, &utils.HTTPResponse{
			Message: err.Error(),
		})
	}

	mate := model.UserMate{
		ID:      uuid.New().String(),
		OwnerID: ownerID,
		UserID:  req.ID,
	}

	err = tx.Create(&mate).Error
	if err != nil {
		tx.Rollback()
		return utils.Response(c, http.StatusInternalServerError, &utils.HTTPResponse{
			Message: err.Error(),
		})
	}

	tx.Commit()
	return utils.Response(c, http.StatusCreated, &utils.HTTPResponse{
		Data: req.ID,
	})
}

func (h *Handler) FindOneUser(c echo.Context) error {
	id := c.Param("id")

	var res model.User
	err := h.db.Where("id = ?", id).Preload("Mates.User").Find(&res).Error
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, &utils.HTTPResponse{
			Message: err.Error(),
		})
	}

	return utils.Response(c, http.StatusOK, &utils.HTTPResponse{
		Data: res,
	})
}

func (h *Handler) DeleteMate(c echo.Context) error {
	id := c.Param("mate_id")

	var mate model.UserMate
	err := h.db.Where("id = ?", id).First(&mate).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return utils.Response(c, http.StatusInternalServerError, &utils.HTTPResponse{
			Message: err.Error(),
		})
	}

	tx := h.db.Begin()

	err = tx.Where("id = ?", mate.ID).Delete(&model.UserMate{}).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return utils.Response(c, http.StatusInternalServerError, &utils.HTTPResponse{
			Message: err.Error(),
		})
	}

	err = tx.Where("id = ?", mate.UserID).Delete(&model.User{}).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return utils.Response(c, http.StatusInternalServerError, &utils.HTTPResponse{
			Message: err.Error(),
		})
	}

	tx.Commit()
	return utils.Response(c, http.StatusCreated, &utils.HTTPResponse{
		Data: id,
	})
}
