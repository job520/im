package utils

import "testing"

func TestSet(t *testing.T) {
	s := NewSet(1, 2, 3)
	s.Add(4)
	s.Add(5)
	if !s.Has(5) {
		t.Error("set 设置失败")
	}
	t.Log("测试通过")
}
