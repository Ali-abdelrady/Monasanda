package osioutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetRepoRootDir(t *testing.T) {
	pvForward, pvError := GetRepoRootDir()

	assert.Nil(t,pvError)
	assert.Equal(t,"d:/firnasSoft/pkgs/firnaspkgs",pvForward)
}
