package testing

import (
	"Cubernetes/pkg/cubenetwork/servicenetwork/utils"
	"log"
	"net"
	"testing"
)

func TestUtils(t *testing.T) {
	ip := net.ParseIP("192.168.255.254")
	for i := 0; i < 10; i++ {
		log.Println("IP1: ", ip.String())
		k := utils.Ip2int(ip)
		ip = utils.Int2ip(k + 1)
	}
}
