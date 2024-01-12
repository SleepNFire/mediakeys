package data_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	// testAdvert, err := Redis.Find("1")
	// assert.ErrorContains(t, err, "the id does not exist")
	// assert.Nil(t, testAdvert)

	// advert := &pkg.AdvertData{
	// 	Id:    "1",
	// 	Title: "some_title",
	// 	Link:  "some_link",
	// }

	// err = Redis.Store(advert)
	// assert.NoError(t, err)

	// testAdvert, err = Redis.Find("1")
	// assert.NoError(t, err)
	// assert.Equal(t, advert, testAdvert)

	// time.Sleep(2 * time.Second)

	// testAdvert, err = Redis.Find("1")
	// assert.ErrorContains(t, err, "the id does not exist")
	// assert.Nil(t, testAdvert)
}
