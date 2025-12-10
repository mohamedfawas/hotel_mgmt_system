package apperror

import (
	"errors"
	"net/http"

	"github.com/mohamedfawas/hotel_mgmt_system/pkg/constants"
)

var (
	ErrBusinessIDMissing = &AppError{
		Err:            errors.New("business id is missing"),
		Code:           constants.BadRequestError,
		HTTPStatusCode: http.StatusBadRequest,
		PublicMsg:      "Business ID is required",
	}
	ErrBusinessNotFound = &AppError{
		Err:            errors.New("business not found"),
		Code:           constants.UnauthorizedError,
		HTTPStatusCode: http.StatusNotFound,
		PublicMsg:      "Given business id not found",
	}
	ErrInvalidHotelCreationRequest = &AppError{
		Err:            errors.New("invalid hotel creation request"),
		Code:           constants.BadRequestError,
		HTTPStatusCode: http.StatusBadRequest,
		PublicMsg:      "Give valid 'name' and 'address' fields",
	}
	ErrEmptyHotelCreationFields = &AppError{
		Err:            errors.New("empty hotel creation fields"),
		Code:           constants.BadRequestError,
		HTTPStatusCode: http.StatusBadRequest,
		PublicMsg:      "Give valid 'name' and 'address' fields",
	}
	ErrInvalidRoomCreationRequest = &AppError{
		Err:            errors.New("invalid room creation request"),
		Code:           constants.BadRequestError,
		HTTPStatusCode: http.StatusBadRequest,
		PublicMsg:      "Give valid 'room_number' and 'room_type' fields and valid 'hotel_id' ",
	}
	ErrRequestedHotelNotFound = &AppError{
		Err:            errors.New("invalid hotel code"),
		Code:           constants.NotFoundError,
		HTTPStatusCode: http.StatusNotFound,
		PublicMsg:      "The requested hotel id doesn't exist, please provide a valid hotel id.",
	}
	ErrTenantNotAuthorized = &AppError{
		Err:            errors.New("uauthorized to access the resource"),
		Code:           constants.UnauthorizedError,
		HTTPStatusCode: http.StatusUnauthorized,
		PublicMsg:      "You are not authorized to access this resource",
	}
)
