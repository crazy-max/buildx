package lib

import "github.com/zclconf/go-cty/cty/function"

var Functions = map[string]function.Function{
	"gitTag": GitTagFunc,
}
