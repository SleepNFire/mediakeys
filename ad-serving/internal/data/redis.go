package data

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SleepNFire/mediakeys/ad-serving/config"
	"github.com/SleepNFire/mediakeys/ad-serving/pkg"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

const prefix = "advertising"

type RedisAccessor struct {
	Client *redis.Client
	Config config.Redis
}

func NewRedisAccessor(globalConf *config.Config) (*RedisAccessor, error) {
	log.Error().Msg("Redis started")
	redisAccessor := &RedisAccessor{
		Config: globalConf.Redis,
	}
	client, err := redisAccessor.connectToRedis()
	if err != nil {
		log.Error().Err(err).Msg("there is anerror during the connection on redis")
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

func CreateKey(id string) string {
	return fmt.Sprintf("%s:%s", prefix, id)
}

func (redisAccessor *RedisAccessor) Store(advert *pkg.AdvertData) error {
	key := CreateKey(advert.Id)
	exists, _ := redisAccessor.Client.Exists(context.Background(), key).Result()
	if exists != 0 {
		return pkg.ErrAdvertAlreadyExist
	}

	advertJSON, err := json.Marshal(advert)
	if err != nil {
		log.Error().Err(err).Msg("error marshal the json")
		return err
	}

	err = redisAccessor.Client.Set(context.Background(), key, string(advertJSON), redisAccessor.Config.Expiration).Err()
	if err != nil {
		log.Error().Err(err).Msg("error saving the key with data")
		return err
	}

	return nil
}

func (redisAccessor *RedisAccessor) Find(id string) (*pkg.AdvertData, error) {
	key := CreateKey(id)
	advertJSON, err := redisAccessor.Client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil, pkg.ErrNotFound
	} else if err != nil {
		log.Error().Err(err).Msg("error during retriving the key in redis")
		return nil, err
	}

	var advert pkg.AdvertData
	err = json.Unmarshal([]byte(advertJSON), &advert)
	if err != nil {
		log.Error().Err(err).Msg("error unmarshal the json")
		return nil, err
	}

	return &advert, nil
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
