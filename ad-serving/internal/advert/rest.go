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
	Exist(id string) bool
}

type ImpressionAccessor interface {
	GetNumber(id string) (string, error)
	Inc(id string) (string, error)
}

type AdvertEndpoint struct {
	Redis   CacheAccessor
	ImpGrpc ImpressionAccessor
}

type AdvertImp struct {
	ImpressionUrl    string `json:"impression_url"`
	ImpressionNumber string `json:"impression_number"`
}

func NewAdvertEndpoint(redis CacheAccessor, impGrpc ImpressionAccessor) (*AdvertEndpoint, error) {
	return &AdvertEndpoint{
		Redis:   redis,
		ImpGrpc: impGrpc,
	}, nil
}

func (advert *AdvertEndpoint) RegisterEndpoints(router *gin.Engine) {
	router.POST(BasePath, advert.SaveAdvert)
	router.GET(BasePath+"/:id", advert.GetAdvert)
	router.GET(BasePath+"/:id/impression", advert.GetImpression)
	router.POST(BasePath+"/serve", advert.ServeAd)
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
			log.Ctx(ctx).Info().Err(err).Msg("The advert doesn't existe")
			c.JSON(http.StatusNotFound, pkg.RestMessage{Message: pkg.ErrNotFound.Error()})
		default:
			log.Ctx(ctx).Error().Err(err).Msg("unexpected error")
			c.JSON(http.StatusInternalServerError, pkg.RestMessage{Message: pkg.ErrInternalError.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, advertData)
}

func (adEndpoint *AdvertEndpoint) GetImpression(c *gin.Context) {
	ctx := c.Request.Context()

	advertId := c.Param("id")

	advertData, err := adEndpoint.Redis.Find(advertId)
	if err != nil {
		switch {
		case err == pkg.ErrNotFound:
			log.Ctx(ctx).Info().Err(err).Msg("The advert doesn't existe")
			c.JSON(http.StatusNotFound, pkg.RestMessage{Message: pkg.ErrNotFound.Error()})
		default:
			log.Ctx(ctx).Error().Err(err).Msg("unexpected error")
			c.JSON(http.StatusInternalServerError, pkg.RestMessage{Message: pkg.ErrInternalError.Error()})
		}
		return
	}

	value, err := adEndpoint.ImpGrpc.GetNumber(advertId)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("error during retrieving impression number")
		value = err.Error()
	}
	c.JSON(http.StatusOK, AdvertImp{
		ImpressionUrl:    advertData.Link,
		ImpressionNumber: value,
	})
}

func (adEndpoint *AdvertEndpoint) ServeAd(c *gin.Context) {
	var advertData pkg.AdvertId

	if err := c.ShouldBindJSON(&advertData); err != nil {
		log.Info().Err(err).Msg("object unknow")
		c.JSON(http.StatusBadRequest, pkg.RestMessage{Message: pkg.ErrObjectUnknown.Error()})
		return
	}

	if advertData.Id == "" {
		log.Info().Msg("object not valid")
		c.JSON(http.StatusBadRequest, pkg.RestMessage{Message: "no advert id"})
		return
	}

	if !adEndpoint.Redis.Exist(advertData.Id) {
		log.Info().Msg("The advert doesn't existe")
		c.JSON(http.StatusNotFound, pkg.RestMessage{Message: pkg.ErrNotFound.Error()})
		return
	}

	response, err := adEndpoint.ImpGrpc.Inc(advertData.Id)
	if err != nil {
		log.Error().Err(err).Str("value", response).Msg("unexpected error")
		c.JSON(http.StatusInternalServerError, pkg.RestMessage{Message: pkg.ErrInternalError.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
