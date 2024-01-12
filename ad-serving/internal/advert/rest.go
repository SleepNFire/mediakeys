package advert

import (
	"net/http"

	"github.com/SleepNFire/mediakeys/ad-serving/pkg"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const (
	BasePath = "api/v1/advert"
)

type CacheAccessor interface {
	Store(advert *pkg.AdvertData) error
	Find(id string) (*pkg.AdvertData, error)
}

type AdvertEndpoint struct {
	Redis CacheAccessor
}

func NewAdvertEndpoint(redis CacheAccessor) (*AdvertEndpoint, error) {
	return &AdvertEndpoint{
		Redis: redis,
	}, nil
}

func (advert *AdvertEndpoint) RegisterEndpoints(router *gin.Engine) {
	router.POST(BasePath, advert.SaveAdvert)
	router.GET(BasePath+"/:id", advert.GetAdvert)
}

func (adEndpoint *AdvertEndpoint) SaveAdvert(c *gin.Context) {
	ctx := c.Request.Context()

	var advertData pkg.AdvertData

	if err := c.ShouldBindJSON(&advertData); err != nil {
		log.Ctx(ctx).Info().Err(err).Msg("object unknow")
		c.JSON(http.StatusBadRequest, pkg.RestMessage{Message: pkg.ErrObjectUnknown.Error()})
		return
	}

	if err := advertData.Validation(); err != nil {
		log.Ctx(ctx).Info().Err(err).Msg("object not valid")
		c.JSON(http.StatusBadRequest, pkg.RestMessage{Message: err.Error()})
		return
	}

	err := adEndpoint.Redis.Store(&advertData)
	if err != nil {
		switch {
		case err == pkg.ErrAdvertAlreadyExist:
			log.Ctx(ctx).Info().Err(err).Msg("The advert already existe")
			c.JSON(http.StatusBadRequest, pkg.RestMessage{Message: err.Error()})
		default:
			log.Ctx(ctx).Error().Err(err).Msg("object unknow")
			c.JSON(http.StatusInternalServerError, pkg.RestMessage{Message: pkg.ErrInternalError.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, advertData)
}

func (adEndpoint *AdvertEndpoint) GetAdvert(c *gin.Context) {
	ctx := c.Request.Context()

	advertId := c.Param("id")

	advertData, err := adEndpoint.Redis.Find(advertId)
	if err != nil {
		switch {
		case err == pkg.ErrNotFound:
			log.Ctx(ctx).Info().Err(err).Msg("The advert already existe")
			c.JSON(http.StatusNotFound, pkg.RestMessage{Message: pkg.ErrNotFound.Error()})
		default:
			log.Ctx(ctx).Error().Err(err).Msg("unexpected error")
			c.JSON(http.StatusInternalServerError, pkg.RestMessage{Message: pkg.ErrInternalError.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, advertData)
}
