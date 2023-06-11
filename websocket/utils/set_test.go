package utils

import (
	"fmt"
	"testing"
)

func TestSet(t *testing.T) {
	s := NewSet("a", "b", "c")
	s.Add("d")
	s.Add("e")
	if !s.Has("e") {
		fmt.Println("测试未通过")
		return
	}
	fmt.Println("测试通过")
}
