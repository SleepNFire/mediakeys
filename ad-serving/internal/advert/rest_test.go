package advert

import (
	"fmt"
	"net/http/httptest"
	"testing"

	advert_mock "github.com/SleepNFire/mediakeys/ad-serving/internal/mocks/advert"
	"github.com/SleepNFire/mediakeys/ad-serving/pkg"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAdvertEndpoint_SaveAdvert(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRedis := advert_mock.NewMockCacheAccessor(ctrl)

	ad, err := NewAdvertEndpoint(mockRedis)
	assert.NoError(t, err)

	endpoint := gin.New()
	ad.RegisterEndpoints(endpoint)

	srv := httptest.NewServer(endpoint)
	defer srv.Close()

	tests := []struct {
		name             string
		body             interface{}
		expectedCode     int
		expectedResponse pkg.AdvertData
		expectedErr      error
	}{
		{
			name: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var response pkg.AdvertData
			var errApi string

			client := resty.New()

			resp, err := client.R().
				SetBody(tt.body).
				SetResult(&response).
				SetError(&errApi).
				Post(fmt.Sprintf("%s/api/v1/advert", srv.URL))

			assert.NoError(t, err)

			switch {
			case resp.StatusCode() == tt.expectedCode:
				assert.Equal(t, tt.expectedResponse, response)
				assert.Equal(t, tt.expectedErr, errApi)
			default:
				t.Log(tt, response, errApi)
				assert.Fail(t, "we don't manage this error")
			}
		})
	}
}
