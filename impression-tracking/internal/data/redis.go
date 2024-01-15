package data

import (
	"fmt"
	"net/http"

	"github.com/SleepNFire/mediakeys/impression-tracking/config"
	"github.com/SleepNFire/mediakeys/impression-tracking/pkg"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

const prefix = "printing"
const advertising = "advertising"

type RedisAccessor struct {
	Client *redis.Client
	Config config.Redis
}

func NewRedisAccessor(globalConf *config.Config) (*RedisAccessor, error) {
	redisAccessor := &RedisAccessor{
		Config: globalConf.Redis,
	}
	client, err := redisAccessor.connectToRedis()
	if err != nil {
		log.Error().Interface("config", redisAccessor.Config).Err(err).Msg("there is an error during the connection on redis")
		return nil, pkg.ErrRedisUnaccessible
	}

	redisAccessor.Client = client
	return redisAccessor, nil
}

func (redisAccessor *RedisAccessor) connectToRedis() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisAccessor.Config.Host, redisAccessor.Config.Port),
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func CreateKeys(id string) (key string, advertKey string) {
	return fmt.Sprintf("%s:%s", prefix, id), fmt.Sprintf("%s:%s", advertising, id)
}

func (redisAccessor *RedisAccessor) Inc(id string) error {
	key, advertKey := CreateKeys(id)
	exists, err := redisAccessor.Client.Exists(context.Background(), advertKey).Result()
	if err != nil {
		log.Error().Err(err).Msg("unexpected error")
		return err
	} else if exists == 0 {
		log.Info().Err(err).Msg("Impossible to inc a advert that does exist")
		redisAccessor.Client.Del(context.Background(), key).Result()
		return pkg.ErrNotFound
	}

	_, err = redisAccessor.Client.Exists(context.Background(), key).Result()
	if err != nil {
		log.Error().Err(err).Msg("unexpected error")
		return nil
	}

	err = redisAccessor.Client.Incr(context.Background(), key).Err()
	if err != nil {
		log.Error().Err(err).Msg("error saving the key with data")
		return err
	}

	return nil
}
func (redisAccessor *RedisAccessor) Find(id string) (uint64, error) {
	key, advertKey := CreateKeys(id)
	exists, err := redisAccessor.Client.Exists(context.Background(), advertKey).Result()
	if err != nil {
		log.Error().Err(err).Msg("unexpected error")
		return 0, err
	} else if exists == 0 {
		log.Info().Err(err).Msg("Impossible to inc a advert that does exist")
		redisAccessor.Client.Del(context.Background(), key).Result()
		return 0, pkg.ErrNotFound
	}

	value, err := redisAccessor.Client.Get(context.Background(), key).Uint64()
	if err != redis.Nil && err != nil {
		log.Error().Err(err).Msg("failed to retrieve value from Redis")
		return 0, err
	}

	return value, nil
}

func (redisAccessor *RedisAccessor) RegisterEndpoints(router *gin.Engine) {
	router.GET("/redis/health", redisAccessor.Ping)
}

func (redisAccessor *RedisAccessor) Ping(c *gin.Context) {
	_, err := redisAccessor.Client.Ping(context.Background()).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "internal error: connection to Redis")
		return
	}
	c.JSON(http.StatusOK, "Redis is good")
}
