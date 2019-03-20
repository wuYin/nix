package velar

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

type AddService interface {
	Add(n, diff int) int
	AddAll(base int, diffs []int) int
}

type AddServiceProxy struct{}

func (s AddServiceProxy) Add(n, diff int) int {
	return n + diff
}

func (s AddServiceProxy) AddAll(base int, diffs []int) int {
	for _, d := range diffs {
		base += d
	}
	return base
}

func TestCallAdd(t *testing.T) {
	addService := Service{
		Uri:       "add-service",
		Proxy:     AddServiceProxy{},
		Interface: reflect.TypeOf((*AddService)(nil)).Elem(),
	}

	manager := NewServiceManager([]Service{addService})
	p := CallReq{
		ServiceUri: "add-service",
		Method:     "Add",
		Args:       []interface{}{1, 1},
		Timeout:    2 * time.Second,
	}
	resp := manager.Call(p)
	if resp.Ec != 0 {
		t.Fatalf("call failed, ec: %d, em: %s", resp.Ec, resp.Em)
	}

	fmt.Println(p.Method, p.Args, resp.Res) // Add [1 1] 2  // bingo 调用成功
}

func TestCallAddAll(t *testing.T) {
	addService := Service{
		Uri:       "add-service",
		Proxy:     &AddServiceProxy{},
		Interface: reflect.TypeOf((*AddService)(nil)).Elem(),
	}
	manager := NewServiceManager([]Service{addService})
	p := CallReq{
		ServiceUri: "add-service",
		Method:     "AddAll",
		Args:       []interface{}{1, []int{10, 20}},
		Timeout:    2 * time.Second,
	}

	resp := manager.Call(p)
	fmt.Printf("%+v\n", resp)
}
