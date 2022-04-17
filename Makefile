build_path = build

cubectl: pkg/objconfig/objconfig.go cmd/cubectl/cubectl.go
	go build -o ${build_path}/cubectl cmd/cubectl/cubectl.go