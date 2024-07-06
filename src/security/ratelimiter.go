package security

// TokenBucketRateLimiter is a rate limiter implementation using token bucket algorithm.
// Here the bucket get refilled with X number of API_QUOTA. If Bucket becomes empty 429 error is thrown.
// func (a *Authentication) TokenBucketRateLimiter(c *fiber.Ctx) error {

// 	ip := c.IP()
// 	value, err := a.rdb.Get(db.Ctx, ip).Result()
// 	if err != nil && err != redis.Nil {
// 		a.log.Error().Err(err).Msg("")
// 		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
// 			"error": err.Error(),
// 		})
// 	}

// 	// this indicates that no entry exist for the specified IP in cache.
// 	// if err == redis.Nil {
// 	// 	fmt.Println("setting new value for specified ip")
// 	// 	err := a.rdb.Set(db.Ctx, ip, os.Getenv("API_QUOTA"), 60*time.Second).Err()
// 	// 	if err != nil {
// 	// 		a.log.Error().Err(err).Msg("")
// 	// 		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
// 	// 			"error": err.Error(),
// 	// 		})
// 	// 	}
// 	// 	value = os.Getenv("API_QUOTA")
// 	// }

// 	var valueInt int

// 	if len(value) > 0 {
// 		valueInt, err = strconv.Atoi(value)
// 		if err != nil {
// 			a.log.Error().Err(err).Msg("")
// 			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
// 				"error": err.Error(),
// 			})
// 		}
// 	}

// 	if valueInt <= 0 {
// 		// limit, err := a.rdb.TTL(db.Ctx, ip).Result()
// 		// if err != nil {
// 		// 	a.log.Error().Err(err).Msg("")
// 		// 	return c.Status(http.StatusBadRequest).JSON(fiber.Map{
// 		// 		"error": err.Error(),
// 		// 	})
// 		// }

// 		// fmt.Println("==============limit after calling ttl", limit, limit/time.Nanosecond/time.Minute)
// 		a.log.Error().Err(errors.New("rate limit exceeded")).Msg("")
// 		return c.Status(http.StatusTooManyRequests).JSON(fiber.Map{
// 			"error": "rate limit exceeded",
// 		})
// 	}

// 	err = a.rdb.Set(db.Ctx, ip, valueInt-1, 60*time.Second).Err()
// 	if err != nil {
// 		a.log.Error().Err(err).Msg("")
// 		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
// 			"error": err.Error(),
// 		})
// 	}

// 	return c.Next()
// }
