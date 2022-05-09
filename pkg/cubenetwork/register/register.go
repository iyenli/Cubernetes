package register

import (
	cubeconfig "Cubernetes/config"
	"Cubernetes/pkg/object"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net"
	"os"
)

func HTTPRegister(req object.NodeRegisterRequest) (object.NodeRegisterResponse, error) {
	log.Println(req.UUID)
	log.Println(req.IP.String())
	return object.NodeRegisterResponse{
		UUID: uuid.NewString(),
	}, nil
}

func RegistryMaster(args []string) {
	if len(args) == 3 { // not same machine with api server
		cubeconfig.APIServerIp = args[2]
	}

	// TODO: Register
	uuid, err := loadUUID()
	IP := net.ParseIP(args[2])
	if IP == nil {
		log.Panicln("Parse IP error: Master IP is illegal")
	}

	resp, err := HTTPRegister(object.NodeRegisterRequest{
		IP:   IP,
		UUID: uuid,
	})

	check(err)
	err = saveUUID(resp.UUID)
	check(err)
}

// Save uuid to log meta
func saveUUID(uuid string) error {
	meta, err := os.OpenFile(MetaFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	check(err)
	defer func(meta *os.File) {
		err := meta.Close()
		if err != nil {
			log.Panicln(err)
		}
	}(meta)

	write, err := meta.Write([]byte(uuid))
	check(err)
	if write != len(uuid) {
		log.Panicln("Write meta failed")
	}
	return nil
}

func loadUUID() (string, error) {
	_, err := os.Stat(MetaDir)
	if err != nil {
		err = os.Mkdir(MetaDir, 0666)
		check(err)
	}

	f, err := os.OpenFile(MetaFile, os.O_CREATE|os.O_RDONLY, 0666)
	check(err)

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Panicln(err)
		}
	}(f)

	content, err := ioutil.ReadAll(f)
	check(err)

	return string(content), nil
}

func check(e error) {
	if e != nil {
		log.Panicln(e.Error())
	}
}
