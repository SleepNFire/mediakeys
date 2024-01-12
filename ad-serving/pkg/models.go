package pkg

import "fmt"

type AdvertData struct {
	Id          string  `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Link        string  `json:"link"`
}

func (advert *AdvertData) Validation() error {
	if advert.Id == "" {
		return fmt.Errorf("there is no Id")
	}
	if advert.Title == "" {
		return fmt.Errorf("there is no Title")
	}
	if advert.Link == "" {
		return fmt.Errorf("there is no Link")
	}
	return nil
}
