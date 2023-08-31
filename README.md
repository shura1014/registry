# registry

简单的内存注册中心

# 快速使用

## 启动

默认启动端口 9999

```go
register := NewMemRegister()
register.Start()
```

curl http://127.0.0.1:9999/mem/probe
OK


## 服务注册

注册一个服务 ： 建议使用应用级注册

例如注册upc应用地址 127.0.0.1:8888

```go
// curl  -POST -d '{"key":"upc","value":"127.0.0.1:8888"}' http://127.0.0.1:9999/mem/register
func TestServiceRegistry(t *testing.T) {
	// 模拟监听一个端口
	_, _ = net.Listen("tcp", ":8888")
	common.Wait()
}
```

> 启动服务

> 访问

curl  -POST -d '{"key":"upc","value":"127.0.0.1:8888"}' http://127.0.0.1:9999/mem/register

> 查看注册中心日志

[registry] 2023-08-31 11:05:27 /Users/wendell/GolandProjects/shura/registry/mem.go:52 INFO register /upc &[{127.0.0.1:8888}]

## 服务发现

curl  -POST -d '{"key":"upc"}' http://127.0.0.1:9999/mem/discover

[{"address":"127.0.0.1:8888"}]

## 监听地址变化

客户端通过与注册中心建立长轮询监听服务地址，30秒一次，如果地址有变动立刻返回，否则30秒后返回304

> 例如

curl  -POST -d '["upc"]' http://127.0.0.1:9999/mem/monitor         
{"code":304}

> 有新的服务注册

curl  -POST -d '["upc"]' http://127.0.0.1:9999/mem/monitor

{"appName":"/upc","code":200,"data":[{"address":"127.0.0.1:8888"},{"address":"127.0.0.1:8889"}]}

[registry] 2023-08-31 11:16:32 /Users/wendell/GolandProjects/shura/registry/mem.go:52 INFO register /upc &[{127.0.0.1:8888} {127.0.0.1:8889}] 


> 有服务下线

手动关闭一个实例

curl  -POST -d '["upc"]' http://127.0.0.1:9999/mem/monitor

{"appName":"/upc","code":200,"data":[{"address":"127.0.0.1:8888"}]}

另一个也关掉

curl  -POST -d '["upc"]' http://127.0.0.1:9999/mem/monitor

{"appName":"/upc","code":200,"data":null}
