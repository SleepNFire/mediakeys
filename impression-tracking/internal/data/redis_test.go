package data_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestAdvert struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Link  string `json:"link"`
}

func TestAccessorRedis(t *testing.T) {
	if os.Getenv("FUNCTIONNEL_TEST") != "1" {
		return
	}
	t.Run("connection", Test_Connection)
	t.Run("push_and_pull_cache", Test_PushAndPull)
}

func Test_Connection(t *testing.T) {
	status := Redis.Client.Ping(context.Background())
	assert.NoError(t, status.Err())
}

func Test_PushAndPull(t *testing.T) {
	testAdvert, _ := Redis.Find("1")
	//assert.NoError(t, err)
	assert.Equal(t, uint64(0), testAdvert)

	advert := &TestAdvert{
		Id:    "1",
		Title: "some_title",
		Link:  "some_link",
	}

	jsonstr, err := json.Marshal(advert)
	assert.NoError(t, err)

	redisErr := Redis.Client.Set(context.Background(), "advertising:1", jsonstr, time.Second)
	assert.NoError(t, redisErr.Err())

	testAdvert, err = Redis.Find("1")
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), testAdvert)

	Redis.Inc("1")

	testAdvert, err = Redis.Find("1")
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), testAdvert)

	Redis.Inc("1")

	testAdvert, err = Redis.Find("1")
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), testAdvert)

	time.Sleep(3 * time.Second)

	testAdvert, err = Redis.Find("1")
	assert.ErrorContains(t, err, "the id does not exist")
	assert.Equal(t, uint64(0), testAdvert)

	exist, _ := Redis.Client.Exists(context.Background(), "advertising:1").Result()
	assert.Equal(t, int64(0), exist)
}
