package advert

import (
	"fmt"
	"net/http"
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
		prepareMock      func(mr *advert_mock.MockCacheAccessor)
		expectedCode     int
		expectedResponse pkg.AdvertData
		expectedErr      pkg.RestMessage
	}{
		{
			name:         "Object unknown",
			body:         nil,
			prepareMock:  func(mr *advert_mock.MockCacheAccessor) {},
			expectedCode: 400,
			expectedErr:  pkg.RestMessage{Message: pkg.ErrObjectUnknown.Error()},
		},
		{
			name: "Object not valid",
			body: pkg.AdvertData{
				Id: "some_id",
			},
			prepareMock:  func(mr *advert_mock.MockCacheAccessor) {},
			expectedCode: 400,
			expectedErr:  pkg.RestMessage{Message: "there is no Title"},
		},
		{
			name: "Already exist",
			body: pkg.AdvertData{
				Id:    "some_id",
				Title: "some_title",
				Link:  "some_link",
			},
			prepareMock: func(mr *advert_mock.MockCacheAccessor) {
				mr.EXPECT().Store(&pkg.AdvertData{
					Id:    "some_id",
					Title: "some_title",
					Link:  "some_link",
				}).Return(pkg.ErrAdvertAlreadyExist)
			},
			expectedCode: 400,
			expectedErr:  pkg.RestMessage{Message: pkg.ErrAdvertAlreadyExist.Error()},
		},
		{
			name: "internal error",
			body: pkg.AdvertData{
				Id:    "some_id",
				Title: "some_title",
				Link:  "some_link",
			},
			prepareMock: func(mr *advert_mock.MockCacheAccessor) {
				mr.EXPECT().Store(&pkg.AdvertData{
					Id:    "some_id",
					Title: "some_title",
					Link:  "some_link",
				}).Return(pkg.ErrInternalError)
			},
			expectedCode: 500,
			expectedErr:  pkg.RestMessage{Message: pkg.ErrInternalError.Error()},
		},
		{
			name: "success",
			body: pkg.AdvertData{
				Id:    "some_id",
				Title: "some_title",
				Link:  "some_link",
			},
			prepareMock: func(mr *advert_mock.MockCacheAccessor) {
				mr.EXPECT().Store(&pkg.AdvertData{
					Id:    "some_id",
					Title: "some_title",
					Link:  "some_link",
				}).Return(nil)
			},
			expectedCode: 200,
			expectedResponse: pkg.AdvertData{
				Id:    "some_id",
				Title: "some_title",
				Link:  "some_link",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var response pkg.AdvertData
			var errApi pkg.RestMessage

			tt.prepareMock(mockRedis)

			client := resty.New()

			resp, err := client.R().
				SetBody(tt.body).
				SetResult(&response).
				SetError(&errApi).
				Post(fmt.Sprintf("%s/api/v1/advert", srv.URL))

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode())

			switch resp.StatusCode() {
			case http.StatusOK:
				assert.Equal(t, tt.expectedResponse, response)
			default:
				assert.Equal(t, tt.expectedErr, errApi)
			}
		})
	}
}

func TestAdvertEndpoint_GetAdvert(t *testing.T) {
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
		param            string
		prepareMock      func(mr *advert_mock.MockCacheAccessor)
		expectedCode     int
		expectedResponse pkg.AdvertData
		expectedErr      pkg.RestMessage
	}{
		{
			name:  "Error not found",
			param: "some_param",
			prepareMock: func(mr *advert_mock.MockCacheAccessor) {
				mr.EXPECT().Find("some_param").Return(nil, pkg.ErrNotFound)
			},
			expectedCode: 404,
			expectedErr:  pkg.RestMessage{Message: pkg.ErrNotFound.Error()},
		},
		{
			name:  "internal error",
			param: "some_param",
			prepareMock: func(mr *advert_mock.MockCacheAccessor) {
				mr.EXPECT().Find("some_param").Return(nil, pkg.ErrInternalError)
			},
			expectedCode: 500,
			expectedErr:  pkg.RestMessage{Message: pkg.ErrInternalError.Error()},
		},
		{
			name:  "success",
			param: "some_param",
			prepareMock: func(mr *advert_mock.MockCacheAccessor) {
				mr.EXPECT().Find("some_param").Return(&pkg.AdvertData{
					Id:    "some_id",
					Title: "some_title",
					Link:  "some_link",
				}, nil)
			},
			expectedCode: 200,
			expectedResponse: pkg.AdvertData{
				Id:    "some_id",
				Title: "some_title",
				Link:  "some_link",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var response pkg.AdvertData
			var errApi pkg.RestMessage

			tt.prepareMock(mockRedis)

			client := resty.New()

			resp, err := client.R().
				SetPathParam("id", tt.param).
				SetResult(&response).
				SetError(&errApi).
				Get(fmt.Sprintf("%s/api/v1/advert/{id}", srv.URL))

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode())

			switch resp.StatusCode() {
			case http.StatusOK:
				assert.Equal(t, tt.expectedResponse, response)
			default:
				assert.Equal(t, tt.expectedErr, errApi)
			}
		})
	}
}
