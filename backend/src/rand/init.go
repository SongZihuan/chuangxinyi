package rand

import (
	errors "github.com/wuntsong-org/wterrors"
	"math/rand"
	"time"
)

var GlobalRander *rand.Rand

func InitRander() errors.WTError {
	GlobalRander = rand.New(rand.NewSource(time.Now().UnixNano()))
	return nil
}
