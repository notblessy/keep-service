package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/notblessy/model"
	"github.com/notblessy/utils"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) FindAllBill(c echo.Context) error {
	ownerID := c.QueryParam("owner_id")

	var res []model.Bill

	err := h.db.Where("owner_id = ?", ownerID).Order("created_at ASC").Find(&res).Error
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, &utils.HTTPResponse{
			Message: err.Error(),
		})
	}

	return utils.Response(c, http.StatusOK, &utils.HTTPResponse{
		Data: res,
	})
}

func (h *Handler) CreateBill(c echo.Context) error {
	var req model.Bill

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

func (h *Handler) DeleteBill(c echo.Context) error {
	id := c.Param("id")

	err := h.db.Where("id = ?", id).Delete(&model.Bill{}).Error
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, &utils.HTTPResponse{
			Message: err.Error(),
		})
	}

	return utils.Response(c, http.StatusCreated, &utils.HTTPResponse{
		Data: id,
	})
}
