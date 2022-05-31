build_path = build
static_path = /etc/cubernetes/static
targets = cubectl apiserver cubelet cuberoot cubeproxy manager scheduler gpuserver gateway brain

all: ${targets} gpuexamples static

.PHONY:clean
clean:
	rm -f $(addprefix ${build_path}/, ${targets})
	rm -f ${build_path}/*.tar.gz
	rm -rf ${static_path}

cubectl: cmd/cubectl/cubectl.go
	go build -o ${build_path}/cubectl cmd/cubectl/cubectl.go

apiserver: cmd/apiserver/apiserver.go
	go build -o ${build_path}/apiserver cmd/apiserver/apiserver.go

cubelet: cmd/cubelet/cubelet.go
	go build -o ${build_path}/cubelet cmd/cubelet/cubelet.go

cuberoot: cmd/cuberoot/cuberoot.go
	go build -o ${build_path}/cuberoot cmd/cuberoot/cuberoot.go

cubeproxy: cmd/cubeproxy/cubeproxy.go
	go build -o ${build_path}/cubeproxy cmd/cubeproxy/cubeproxy.go

manager: cmd/controller_manager/manager.go
	go build -o ${build_path}/manager cmd/controller_manager/manager.go

scheduler: cmd/scheduler/scheduler.go
	go build -o ${build_path}/scheduler cmd/scheduler/scheduler.go

gpuserver: cmd/gpujobserver/gpujobserver.go
	go build -o ${build_path}/gpuserver cmd/gpujobserver/gpujobserver.go

gateway: cmd/gateway/gateway.go
	go build -o ${build_path}/gateway cmd/gateway/gateway.go

brain: cmd/action_brain/brain.go
	go build -o ${build_path}/brain cmd/action_brain/brain.go

gpu_path = example/gpujob
gpu_files = cublashello matmult matadd
gpuexamples: $(addprefix ${gpu_path}/, ${gpu_files})
	$(foreach file, ${gpu_files}, tar zcvf ${build_path}/${file}.tar.gz -C ${gpu_path} ${file};)

static: ./static/*
	cp -r ./static ${static_path}