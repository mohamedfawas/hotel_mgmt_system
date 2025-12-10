package hotels

type CreateHotelRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type CreateHotelResponse struct {
	HotelCode    int64  `json:"hotel_code"`
	HotelName    string `json:"hotel_name"`
	HotelAddress string `json:"hotel_address"`
}
