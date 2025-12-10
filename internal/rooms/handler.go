package rooms

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mohamedfawas/hotel_mgmt_system/pkg/apiresponse"
	"github.com/mohamedfawas/hotel_mgmt_system/pkg/apperror"
	"github.com/mohamedfawas/hotel_mgmt_system/pkg/constants"
)

type Handler struct {
	svc RoomService
}

func NewHandler(svc RoomService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateRoom(c *gin.Context) {
	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, apperror.ErrInvalidRoomCreationRequest, nil)
		return
	}

	hotelCodeStr := c.Param("hotel_id")
	hotelCode, err := strconv.ParseInt(hotelCodeStr, 10, 64)
	if err != nil {
		apiresponse.Error(c, apperror.ErrInvalidRoomCreationRequest, nil)
		return
	}

	if req.RoomNumber <= 0 || req.RoomType == "" {
		apiresponse.Error(c, apperror.ErrInvalidRoomCreationRequest, nil)
		return
	}

	businessID, ok := c.Get(constants.ContextBusinessID)
	if !ok {
		apiresponse.Error(c, apperror.ErrBusinessIDMissing, nil)
		return
	}

	businessIDUUID, err := uuid.Parse(businessID.(string))
	if err != nil {
		apiresponse.Error(c, fmt.Errorf("failed to parse business ID: %v", err), nil)
		return
	}

	newRoom := &Room{
		BusinessID: businessIDUUID,
		RoomNumber: req.RoomNumber,
		RoomType:   req.RoomType,
	}

	createdRoom, err := h.svc.CreateRoom(c.Request.Context(), hotelCode, newRoom)
	if err != nil {
		if apperror.ShouldLogError(err) {
			log.Printf("failed to create room: %w", err)
		}
		apiresponse.Error(c, err, nil)
		return
	}

	createdRoomResponse := CreateRoomResponse{
		RoomNumber: createdRoom.RoomNumber,
		RoomType:   createdRoom.RoomType,
	}

	apiresponse.Created(c, "Room created successfully", createdRoomResponse)

}
