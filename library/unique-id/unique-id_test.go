package unique_id

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateUniqueId(t *testing.T) {
	ids, err := GenerateUniqueId(UniqueBizTypeMatch, 1)
	fmt.Println(ids)
	assert2 := assert.New(t)
	assert2.True(err == nil && len(ids) == 1)
}

func TestGetUniqueIdInfo(t *testing.T) {
	info, err := GetUniqueIdInfo(74647063872875648)
	assert2 := assert.New(t)
	assert2.True(err == nil && info != nil && info.BizType == UniqueBizTypeMatch)
}
