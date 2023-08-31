package registry

import (
	"github.com/shura1014/common"
	"net"
	"testing"
)

// curl http://127.0.0.1:9999/mem/probe
func TestRegister(t *testing.T) {
	register := NewMemRegister()
	register.Start()
}

// curl  -POST -d '{"key":"upc","value":"127.0.0.1:8888"}' http://127.0.0.1:9999/mem/register
func TestServiceRegistry(t *testing.T) {
	_, _ = net.Listen("tcp", ":8888")
	common.Wait()
}

func TestServiceRegistry2(t *testing.T) {
	_, _ = net.Listen("tcp", ":8889")
	common.Wait()
}
