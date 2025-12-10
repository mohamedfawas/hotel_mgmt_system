package hotels

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mohamedfawas/hotel_mgmt_system/pkg/apiresponse"
	"github.com/mohamedfawas/hotel_mgmt_system/pkg/apperror"
	"github.com/mohamedfawas/hotel_mgmt_system/pkg/constants"
	"github.com/mohamedfawas/hotel_mgmt_system/pkg/pagination"
)

type Handler struct {
	svc HotelService
}

func NewHandler(svc HotelService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateHotel(c *gin.Context) {
	var req CreateHotelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, apperror.ErrInvalidHotelCreationRequest, nil)
		return
	}

	if req.Name == "" || req.Address == "" {
		apiresponse.Error(c, apperror.ErrEmptyHotelCreationFields, nil)
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

	newHotel := &Hotel{
		BusinessID: businessIDUUID,
		Name:       req.Name,
		Address:    req.Address,
	}

	createdHotel, err := h.svc.CreateHotel(c.Request.Context(), newHotel)
	if err != nil {
		log.Printf("failed to create hotel: %w", err)
		apiresponse.Error(c, err, nil)
		return
	}

	createdHotelResponse := CreateHotelResponse{
		HotelCode:    createdHotel.HotelCode,
		HotelName:    createdHotel.Name,
		HotelAddress: createdHotel.Address,
	}

	apiresponse.Created(c, "Hotel created successfully", createdHotelResponse)
}


func (h *Handler) ListHotels(c *gin.Context) {
	// get business id (tenant) from context
	businessIDRaw, ok := c.Get(constants.ContextBusinessID)
	if !ok {
		apiresponse.Error(c, apperror.ErrBusinessIDMissing, nil)
		return
	}
	businessID, err := uuid.Parse(businessIDRaw.(string))
	if err != nil {
		apiresponse.Error(c, fmt.Errorf("failed to parse business ID: %v", err), nil)
		return
	}

	// parse pagination parameters
	p := pagination.FromRequest(c)

	// call service
	hotelsWithRooms, err := h.svc.ListHotelsWithRooms(c.Request.Context(), businessID, p.Limit, p.Offset)
	if err != nil {
		log.Printf("failed to list hotels: %v", err)
		apiresponse.Error(c, err, nil)
		return
	}

	apiresponse.Success(c, "hotels fetched", hotelsWithRooms)
}