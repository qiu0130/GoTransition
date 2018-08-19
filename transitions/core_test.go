package transitions

import (
	"testing"
)

func TestInfo(t *testing.T) {
	test := "hello world %s, %s"

	Info(test, "xx", "yyy")
	Warning(test, "xx", "yy")
	Error(test, "xx", "yy")

	test1 := "%s, %s 哈哈"
	Info(test1, "我", "xx")
	Warning(test1, "你", "xx")
	Error(test1, "他", "xxx")

}
