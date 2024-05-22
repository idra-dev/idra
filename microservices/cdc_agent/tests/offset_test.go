package tests

import (
	"microservices/libraries/etcd"
	"testing"
)

func TestOffset(t *testing.T){
	manager := etcd.IntOffset{}
	manager.SetOffsetId("aaa", 12)
	offset := manager.GetOffsetId("aaa")
	if offset != 12 {
		t.Errorf("got %q, wanted %q", offset, 12)
	}
}
