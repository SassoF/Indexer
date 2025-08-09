package session

import (
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
)

var SessionManager *scs.SessionManager

func Init() {

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	pool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			opts := []redis.DialOption{
				redis.DialConnectTimeout(5 * time.Second),
			}

			if redisPassword != "" {
				opts = append(opts, redis.DialPassword(redisPassword))
			}

			return redis.Dial("tcp", redisHost+":"+redisPort, opts...)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	SessionManager = scs.New()
	SessionManager.Store = redisstore.New(pool)
	SessionManager.Lifetime = 24 * time.Hour
	SessionManager.Cookie.Secure = true
	SessionManager.Cookie.HttpOnly = true
	SessionManager.Cookie.SameSite = http.SameSiteStrictMode
	SessionManager.Cookie.Persist = true
}
