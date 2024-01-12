package app

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartApplication(t *testing.T) {
	if os.Getenv("FUNCTIONNEL_TEST") != "1" {
		return
	}
	app := Init()

	err := app.Start(context.Background())
	assert.NoError(t, err)

	app.Done()
	err = app.Stop(context.Background())
	assert.NoError(t, err)
}
