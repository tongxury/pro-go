package invitecode

import (
	"fmt"
	"testing"
)

func TestDecode(t *testing.T) {
	g, e := NewGenerator(6)
	if e != nil {
		t.Fatal(e)
	}

	//fmt.Println(g.Encode(5459275635))
	fmt.Println(g.Decode("kCq6so"))

}

func TestGenerator_Encode(t *testing.T) {
	g, e := NewGenerator(6)
	if e != nil {
		t.Fatal(e)
	}

	test := func(id uint64) bool {
		code, e := g.Encode(id)
		if e != nil {
			t.Error(id, e)
			return false
		}
		nid := g.Decode(code)
		//t.Logf("解析邀请码结果：code:%s id:%d 是否相等:%t", code, nid, id == nid)
		//t.Log("=========================")
		if nid != id {
			t.Error(id, nid)
			return false
		}
		return true
	}

	var _min, _max uint64 = 0, 20
	for id := _min; id <= _max; id++ {
		if !test(id) {
			return
		}
	}
}
