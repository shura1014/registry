package registry

import (
	"encoding/json"
	"github.com/shura1014/common"
	"github.com/shura1014/common/utils/stringutil"
	"net"
	"net/http"
	"strings"
	"time"
)

var (
	registerPath = "/mem/register"
	discoverPath = "/mem/discover"
	monitorPath  = "/mem/monitor"
	probePath    = "/mem/probe"
)

// Start 启动注册中心
func (register *MemRegister) Start(address ...string) {
	addr := ":9999"
	if len(address) > 0 {
		addr = address[0]
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	go register.checkAddress()
	_ = http.Serve(l, register)
}

// 延时渲染
var monitors map[string][]*common.DeferredResultWriter

func init() {
	monitors = make(map[string][]*common.DeferredResultWriter)
}

// 监听
func (register *MemRegister) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	switch path {
	case probePath:
		_, _ = w.Write([]byte{'O', 'K'})
	case registerPath:
		register.handlerRegister(w, req)
	case discoverPath:
		register.handlerDiscover(w, req)
	case monitorPath:
		register.handlerMonitor(w, req)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

// /mem/register
// {"appName":"upc","key":"upc_get","value":"127.0.0.1"}
// curl  -POST -d '{"appName":"upc","key":"upc_get","value":"127.0.0.1"}' http://127.0.0.1:9999/mem/register
func (register *MemRegister) handlerRegister(w http.ResponseWriter, req *http.Request) {
	body := &Data{}
	flag := GetBody(w, req, body)
	if !flag {
		return
	}
	register.registryData(*body)
	return
}

func (register *MemRegister) handlerDiscover(w http.ResponseWriter, req *http.Request) {
	body := &Data{}
	flag := GetBody(w, req, body)
	if !flag {
		return
	}
	ips := register.Get(*body)
	marshal, err := json.Marshal(ips)
	assertError(err, w, func() bool {
		_, err = w.Write(marshal)
		return assertError(err, w)
	})
}

// GetBody 获取请求体，有报错返回false
func GetBody(w http.ResponseWriter, req *http.Request, data any) bool {
	body := req.Body
	if req == nil || body == nil {
		return false
	}
	decoder := json.NewDecoder(body)
	err := decoder.Decode(data)
	return assertError(err, w)
}

// 长轮训 30秒没有响应结果返回304，否则返回正常结果
func (register *MemRegister) handlerMonitor(w http.ResponseWriter, req *http.Request) {
	var appNames []string
	flag := GetBody(w, req, &appNames)
	if !flag {
		return
	}
	common.NewDeferredResultWriter(30*time.Second, map[string]any{"code": 304}, w, func(deferred *common.DeferredResultWriter) {
		for _, appName := range appNames {
			if !strings.HasPrefix(appName, defaultSplit) {
				appName = defaultSplit + appName
			}
			monitors[appName] = append(monitors[appName], deferred)
		}
	})
}

// 有地址下线
func (register *MemRegister) notify(appName string, data *[]Address) {
	writer := monitors[appName]
	for _, w := range writer {
		_ = w.SetResult(map[string]any{"code": 200, "appName": appName, "data": data})
	}
}

func assertError(err error, w http.ResponseWriter, fun ...func() bool) bool {
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(stringutil.StringToBytes(err.Error()))
		return false
	}
	if len(fun) > 0 {
		return fun[0]()
	}
	return true
}
