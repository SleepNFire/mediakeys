package pkg

type AdvertData struct {
	Id          string
	Title       string
	Description *string
	Link        string
	PrintNumber uint64
}

func (advert *AdvertData) Validation() error {
	return nil
}
