package modes

import (
	"github.com/go-xorm/xorm"
	"go_util/lib"
)

func Db(index int) *xorm.Engine {
	return lib.UseHand(index)
}
