package lxCrypt

import (
	"github.com/litixsoft/lx-golib/helper"
)

// ICrypt,
// interface for mapping bcrypt
type ICrypt interface {
	GeneratePassword(plainPwd string) (string, error)
	ComparePassword(hashedPwd, plainPwd string) error
}

// Crypt,
// type for bcrypt mapper
type Crypt struct{}

// GeneratePassword,
// mapper for create new encrypt password from plain password
func (c *Crypt) GeneratePassword(plainPwd string) (string, error) {
	return lxHelper.GenerateFromPassword(plainPwd)
}

// ComparePassword,
// mapper for compare encrypt password with plain password
func (c *Crypt) ComparePassword(cryptPwd, plainPwd string) error {
	return lxHelper.CompareHashAndPassword(cryptPwd, plainPwd)
}