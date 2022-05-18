package testing

import (
	"Cubernetes/pkg/object"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"testing"
)

// deprecated function
//func TestDNSCheck(t *testing.T) {
//	str := "example.com/"
//	ans, err := utils.CheckDNSHostName(str)
//	assert.NoError(t, err)
//	assert.Equal(t, "example.com", ans)
//}

// Warning: not universal
func TestMarshall(t *testing.T) {
	var DNS object.Dns
	file, err := ioutil.ReadFile("/home/lee/Cubernetes/Cubernetes/example/yaml/dns/dns.yaml")
	assert.NoError(t, err)

	err = yaml.Unmarshal(file, &DNS)
	assert.NoError(t, err)

	file, err = json.Marshal(DNS)
	assert.NoError(t, err)
	log.Println(string(file))
}
