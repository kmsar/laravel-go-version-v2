package tests

import (
	"fmt"
	"github.com/laravel-go-version/v2/pkg/Illuminate/Support/Utils"

	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNoPanic(t *testing.T) {
	var err = Utils.NoPanic(func() {
		panic("报错")
	})
	fmt.Println(err)
	assert.Error(t, err)
}
