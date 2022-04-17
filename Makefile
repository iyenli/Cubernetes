build_path = build
cubectl_path = cmd/cubectl

cubectl: ${cubectl_path}/config.go ${cubectl_path}/cubectl.go
	go build -o ${build_path}/cubectl ${cubectl_path}/config.go ${cubectl_path}/cubectl.go