package util

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Parser helps in parsing the data from the URL params.
type Parser struct {
	ctx *fiber.Ctx
	// Params gin.Params
	// Form   url.Values
}

// NewParser will call request.ParseForm() and create a new instance of parser.
func NewParser(c *fiber.Ctx) *Parser {
	return &Parser{
		ctx: c,
		// Params: c.,
		// Form: ctx.Request.Form,
	}
}

// GetParameter will get parameter from the given paramName in URL params.
func (p *Parser) GetParameter(paramName string) string {
	return p.ctx.Params(paramName, "")
}

// ParseLimitAndOffset will parse limit and offset from query params.
func (p *Parser) ParseLimitAndOffset() (limit, offset int) {
	limitparam := p.ctx.Query("limit", "10")
	offsetparam := p.ctx.Query("offset", "0")

	var err error
	limit = 30
	if len(limitparam) > 0 {
		limit, err = strconv.Atoi(limitparam)
		if err != nil {
			return
		}
	}
	if len(offsetparam) > 0 {
		offset, err = strconv.Atoi(offsetparam)
		if err != nil {
			return
		}
	}
	return
}
