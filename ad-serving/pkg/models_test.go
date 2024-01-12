package pkg

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdvertData_Validation(t *testing.T) {
	descrition := "some_decription"

	tests := []struct {
		name        string
		advert      AdvertData
		expectedErr error
	}{
		{
			name: "good advert",
			advert: AdvertData{
				Id:          "some_id",
				Title:       "some_title",
				Description: &descrition,
				Link:        "some_link",
			},
			expectedErr: nil,
		},
		{
			name: "no Id",
			advert: AdvertData{
				Title:       "some_title",
				Description: &descrition,
				Link:        "some_link",
			},
			expectedErr: fmt.Errorf("there is no Id"),
		},
		{
			name: "no title",
			advert: AdvertData{
				Id:          "some_id",
				Description: &descrition,
				Link:        "some_link",
			},
			expectedErr: fmt.Errorf("there is no Title"),
		},
		{
			name: "no link",
			advert: AdvertData{
				Id:          "some_id",
				Title:       "some_title",
				Description: &descrition,
			},
			expectedErr: fmt.Errorf("there is no Link"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.advert.Validation()
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
