package lxCrypt_test

import (
	"testing"
	"github.com/litixsoft/lx-golib/crypt"
	"github.com/stretchr/testify/assert"
)

func TestCrypt_GeneratePassword(t *testing.T) {
	c := lxCrypt.Crypt{}
	cryptPwd, err := c.GeneratePassword("plain-pwd")
	assert.NoError(t, err)

	err = c.ComparePassword(cryptPwd, "plain-pwd")
	assert.NoError(t, err)
}