package lxTestHelper

import (
	"testing"
	"github.com/golang/mock/gomock"
)

func GetMockController(t *testing.T) *gomock.Controller {
	return gomock.NewController(t)

}
