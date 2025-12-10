package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Params struct {
	Limit  int
	Offset int
	Page   int
	Per    int
}

const (
	DefaultPage    = 1
	DefaultPerPage = 20
	MaxPerPage     = 100
)

func FromRequest(c *gin.Context) Params {
	page := DefaultPage
	per := DefaultPerPage
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			if v > MaxPerPage {
				per = MaxPerPage
			} else {
				per = v
			}
		}
	}
	offset := (page - 1) * per
	return Params{
		Limit:  per,
		Offset: offset,
		Page:   page,
		Per:    per,
	}
}
