package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
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
	mateID := c.QueryParam("mate_id")

	db := h.db

	switch mateID {
	case "":
		db = db.Preload("Shares")
	default:
		db = db.Preload("Shares", func(d *gorm.DB) *gorm.DB {
			return d.Where("mate_id = ?", mateID)
		})
	}

	var res []model.Bill

	err := db.Where("owner_id = ?", ownerID).Order("created_at ASC").Find(&res).Error
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, &utils.HTTPResponse{
			Message: err.Error(),
		})
	}

	var ids []string
	for _, b := range res {
		ids = append(ids, b.ID)
	}

	var totalTaken []model.TotalShareQty
	err = db.Raw("SELECT bill_id, SUM(qty) as qty FROM shares WHERE bill_id IN (?) GROUP BY bill_id;", ids).Scan(&totalTaken).Error
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, &utils.HTTPResponse{
			Message: err.Error(),
		})
	}

	for i := 0; i < len(res); i++ {
		for _, t := range totalTaken {
			if t.BillID == res[i].ID {
				res[i].Taken = t.Qty
			}
		}
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

	req.ExpiredAt = time.Now().Add(7 * 24 * time.Hour)

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

func (h *Handler) CreateShare(c echo.Context) error {
	var req model.Share

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

	reqType := c.QueryParam("type")

	tx := h.db.Begin()

	var bill model.Bill
	err := tx.Where("id = ?", req.BillID).First(&bill).Error
	if err != nil {
		tx.Rollback()
		return utils.Response(c, http.StatusInternalServerError, &utils.HTTPResponse{
			Message: err.Error(),
		})
	}

	var share model.Share
	err = tx.Where("bill_id = ? AND mate_id = ?", bill.ID, req.MateID).First(&share).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return utils.Response(c, http.StatusInternalServerError, &utils.HTTPResponse{
			Message: err.Error(),
		})
	}

	if share.ID != "" {
		req.ID = share.ID
		req.Price += share.Price

		switch reqType {
		case "SUB":
			req.Qty = share.Qty - req.Qty
		default:
			req.Qty += share.Qty
		}
	}

	if req.ID == "" {
		req.ID = uuid.New().String()
	}

	if req.Qty == 0 {
		err := h.db.Where("id = ?", req.ID).Delete(&model.Share{}).Error
		if err != nil {
			return utils.Response(c, http.StatusInternalServerError, &utils.HTTPResponse{
				Message: err.Error(),
			})
		}

		tx.Commit()
		return utils.Response(c, http.StatusCreated, &utils.HTTPResponse{
			Data: req.ID,
		})
	}

	err = tx.Save(&req).Error
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

func (h *Handler) FindAllWithShare(c echo.Context) error {
	ownerID := c.QueryParam("owner_id")

	var bills []model.Bill
	err := h.db.Where("owner_id = ?", ownerID).Order("created_at ASC").Find(&bills).Error
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, &utils.HTTPResponse{
			Message: err.Error(),
		})
	}

	res := []model.AssignedBillResponse{}
	for _, bill := range bills {
		assigned := model.NewAssignedBill(bill)
		res = append(res, assigned)
	}

	return utils.Response(c, http.StatusOK, &utils.HTTPResponse{
		Data: res,
	})
}

func (h *Handler) SplitBill(c echo.Context) error {
	ownerID := c.Param("owner_id")

	if ownerID == "" {
		return utils.Response(c, http.StatusBadRequest, &utils.HTTPResponse{
			Message: errors.New("invalid owner id").Error(),
		})
	}

	var owner model.User
	err := h.db.Where("id = ?", ownerID).Preload("Mates.User").Find(&owner).Error
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, &utils.HTTPResponse{
			Message: err.Error(),
		})
	}

	mateIDs := []string{}
	for _, m := range owner.Mates {
		mateIDs = append(mateIDs, m.UserID)
	}

	var mates []model.Mate
	err = h.db.Table("users").Where("id IN (?)", mateIDs).Preload("Shares.Bill").Find(&mates).Error
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, &utils.HTTPResponse{
			Message: err.Error(),
		})
	}

	split := h.transformToOutputFormat(owner, mates)
	entity := split.ToEntity()

	err = h.db.Save(&entity).Error
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, &utils.HTTPResponse{
			Message: err.Error(),
		})
	}

	return utils.Response(c, http.StatusCreated, &utils.HTTPResponse{
		Data: split,
	})
}

func (h *Handler) transformToOutputFormat(owner model.User, mates []model.Mate) model.Split {
	var split model.Split

	split.ID = uuid.New().String()
	split.OwnerID = owner.ID
	split.CreatedAt = time.Now()

	for _, mate := range mates {
		var splitMate model.SplitMate

		splitMate.ID = mate.ID
		splitMate.Name = mate.Name

		for _, share := range mate.Shares {
			var item model.SplitItem

			item.ID = share.BillID
			item.Name = share.Bill.Name
			item.Qty = share.Qty
			item.Price = share.Price
			item.Total = share.Qty * share.Price

			splitMate.GrandTotal += item.Total
			splitMate.SplitItems = append(splitMate.SplitItems, item)
		}

		split.SplitMates = append(split.SplitMates, splitMate)
	}

	return split
}

func (h *Handler) FindSplit(c echo.Context) error {
	id := c.Param("id")

	var split model.SplitEntity
	err := h.db.Where("id = ?", id).Preload("OwnerDetail.UserBanks").First(&split).Error
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, &utils.HTTPResponse{
			Message: err.Error(),
		})
	}

	res := model.Split{
		ID:          split.ID,
		OwnerID:     split.OwnerID,
		OwnerDetail: split.OwnerDetail,
		CreatedAt:   split.CreatedAt,
	}

	err = json.Unmarshal([]byte(split.SplitMates), &res.SplitMates)
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, &utils.HTTPResponse{
			Message: err.Error(),
		})
	}

	return utils.Response(c, http.StatusOK, &utils.HTTPResponse{
		Data: res,
	})
}
