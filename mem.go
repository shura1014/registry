package registry

import (
	"net"
	"time"
)

type Data struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Address struct {
	Address string `json:"address"`
}

func (addr *Address) Addr() string {
	return addr.Address
}

func (addr *Address) IsValid() bool {
	return true
}

type MemRegister struct {
	registerInfo *Node
}

func NewMemRegister() *MemRegister {
	registerInfo := NewNode("/")
	return &MemRegister{
		registerInfo: registerInfo,
	}
}

func (register *MemRegister) registryData(data Data) {
	node := register.registerInfo.Get(data.Key)
	if node != nil && node.IsEnd && node.Data != nil {
		*node.Data = append(*node.Data, Address{
			Address: data.Value,
		})
	} else {
		addrs := []Address{{
			Address: data.Value,
		}}
		register.registerInfo.Put(data.Key, &addrs)
	}

	checkNode := register.registerInfo.Get(data.Key)

	register.notify(checkNode.RouterName, checkNode.Data)
	Info("register %s %s", checkNode.RouterName, checkNode.Data)
}

func (register *MemRegister) deleteData(data Data) {
	register.registerInfo.delete(data.Key)
}

func (register *MemRegister) Get(data Data) *[]Address {
	return register.registerInfo.Get(data.Key).Data
}

// 应用下线检查
func (register *MemRegister) checkAddress() {
	for {
		Iterator(register.registerInfo, func(key string, Data *[]Address) {
			var addr []Address
			var isNotify bool
			for _, address := range *Data {
				conn, err := net.DialTimeout("tcp", address.Address, 1000*time.Millisecond)
				if conn != nil || err == nil {
					addr = append(addr, address)
					_ = conn.Close()
				} else {
					isNotify = true
				}
			}
			if isNotify {
				register.notify(key, Data)
			}
			*Data = addr
		})
		time.Sleep(2 * time.Second)
	}
}
