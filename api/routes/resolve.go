package routes

import (
	"github.com/akhil/fiber-url-shortener/database"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)


func ResolveURL(c *fiber.Ctx) error {
	
	url := c.Params("url")
	// query the db to find the original URL, if a match is found
	// increment the redirect counter and redirect to the original URL
	// else return error message
	r := database.CreateClient(0)
	defer r.Close()

	value, err := r.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "short not found on database",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot connect to DB",
		})
	}
	
	rInr := database.CreateClient(1)
	defer rInr.Close()
	_ = rInr.Incr(database.Ctx, "counter")
	// redirect to original URL
	return c.Redirect(value, 301)
}