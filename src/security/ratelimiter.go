package security

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/shaileshhb/equisplit/src/db"
)

func (a *Authentication) RateLimiter(c *fiber.Ctx) error {

	value, err := a.rdb.Get(db.Ctx, c.IP()).Result()
	if err != nil {
		fmt.Println("-------------- err in get")
		return err
	}
	fmt.Println("inside rate limiter middleware", value)
	return nil
}
