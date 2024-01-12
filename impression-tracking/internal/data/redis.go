package data

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/SleepNFire/mediakeys/impression-tracking/config"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

const (
	prefix      = "print"
	advertising = "advertising"
)

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
		fmt.Println("error during connecting to Redis")
		return nil, errors.New("error during connecting to Redis")
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

func CreateAdKey(id string) string {
	return fmt.Sprintf("%s:%s", advertising, id)
}

func (redisAccessor *RedisAccessor) AdKeyExist(id string) bool {
	if exists, _ := redisAccessor.Client.Exists(context.Background(), CreateAdKey(id)).Result(); exists != 0 {
		return true
	}

	key := CreateKey(id)
	if exists, _ := redisAccessor.Client.Exists(context.Background(), key).Result(); exists != 0 {
		redisAccessor.Client.Del(context.Background(), key)
	}
	return false
}

func (redisAccessor *RedisAccessor) Find(id string) (uint64, error) {

	if !redisAccessor.AdKeyExist(id) {
		return 0, fmt.Errorf("the advert does not existe")
	}

	key := CreateKey(id)

	advertValue, err := redisAccessor.Client.Get(context.Background(), key).Result()
	if err != nil {
		return 0, err
	}

	advertUint64, err := strconv.ParseUint(advertValue, 10, 64)
	if err != nil {
		return 0, err
	}

	return advertUint64, nil
}

func (redisAccessor *RedisAccessor) Increment(id string) error {

	if !redisAccessor.AdKeyExist(id) {
		return fmt.Errorf("the advert does not existe")
	}

	key := CreateKey(id)

	exists, err := redisAccessor.Client.Exists(context.Background(), key).Result()
	if err != nil {
		return err
	}

	if exists == 0 {
		err := redisAccessor.Client.Set(context.Background(), key, 0, 0).Err()
		if err != nil {
			return err
		}
		return nil
	}

	_, err = redisAccessor.Client.Incr(context.Background(), key).Result()
	if err != nil {
		return err
	}

	return nil
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
