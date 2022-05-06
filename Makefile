build_path = build

all: cubectl apiserver cubelet

clean:
	rm ${build_path}/cubectl $(build_path)/apiserver $(build_path)/cubelet

cubectl: cmd/cubectl/cubectl.go
	@go build -o ${build_path}/cubectl cmd/cubectl/cubectl.go

apiserver: cmd/apiserver/apiserver.go
	@go build -o $(build_path)/apiserver cmd/apiserver/apiserver.go

cubelet: cmd/cubelet/cubelet.go
	@go build -o $(build_path)/cubelet cmd/cubelet/cubelet.go

cuberoot: cmd/cuberoot/cuberoot.go
	@go build -o $(build_path)/cuberoot cmd/cuberoot/cuberoot.go

cubeproxy: cmd/cubeproxy/cubeproxy.go
	@go build -o $(build_path)/cubeproxy cmd/cubeproxy/cubeproxy.go
