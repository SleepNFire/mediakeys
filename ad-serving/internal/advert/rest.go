package advert

import (
	"net/http"

	"github.com/SleepNFire/mediakeys/ad-serving/pkg"
	"github.com/gin-gonic/gin"
)

type CacheAccessor interface {
	Store(advert *pkg.AdvertData) error
	Find(id string) (*pkg.AdvertData, error)
}

type AdvertEndpoint struct {
	Redis CacheAccessor
}

func Init(redis CacheAccessor) *AdvertEndpoint {
	return &AdvertEndpoint{
		Redis: redis,
	}
}

func (advert *AdvertEndpoint) RegisterEndpoints(router *gin.Engine) {

}

func (advert *AdvertEndpoint) SaveAdvert(c *gin.Context) {
	var advertData pkg.AdvertData

	if err := c.ShouldBindJSON(&advertData); err != nil {
		c.JSON(http.StatusBadRequest, "object unkown")
		return
	}

	if 
}
