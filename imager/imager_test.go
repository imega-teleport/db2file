package imager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetImageInfo(t *testing.T) {
	info, err := GetImageInfo("teleport.png")
	assert.NoError(t, err)

	assert.Equal(t, 200, info.Height)
	assert.Equal(t, 200, info.Width)
	assert.Equal(t, "teleport.png", info.Name)
	assert.Equal(t, "image/png", info.Mime)
}
