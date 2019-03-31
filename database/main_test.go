package database

import (
	"testing"

	"github.com/devenney/form3/common"
)

func TestMain(m *testing.M) {
	common.InitConfig()

	m.Run()
}
