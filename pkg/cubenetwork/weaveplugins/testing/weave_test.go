package testing

import (
	"Cubernetes/pkg/cubenetwork/weaveplugins"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"log"
	"net"
	osexec "os/exec"
	"strings"
	"testing"
)

// TestAddNode JCloud env
func TestAddNode(t *testing.T) {
	err := PrepareTest()
	assert.NoError(t, err)

	host1 := weaveplugins.Host{IP: net.ParseIP("192.168.1.9")}
	host2 := weaveplugins.Host{IP: net.ParseIP("192.168.1.5")}

	err = weaveplugins.AddNode(host1, host2)
	assert.NoError(t, err)
}

// TestInitNode WSL env
func TestInitNode(t *testing.T) {
	err := PrepareTest()
	assert.NoError(t, err)

	err = weaveplugins.InitWeave()
	assert.NoError(t, err)
}

func TestInstallWeave(t *testing.T) {
	err := weaveplugins.InstallWeave()
	assert.NoError(t, err)
}

func TestWeaveStatus(t *testing.T) {
	output, err := weaveplugins.CheckPeers()
	assert.NoError(t, err)

	t.Log(string(output))
}

func TestWeaveStop(t *testing.T) {
	err := weaveplugins.CloseNetwork()
	assert.NoError(t, err)
}

func TestAddPod(t *testing.T) {
	id := RunContainer()
	network, err := weaveplugins.AddPodToNetwork(id)
	assert.NoError(t, err)

	t.Logf("Container ID: %v, IP: %v", id, network.String())
	err = weaveplugins.DeletePodFromNetwork(id)
	assert.NoError(t, err)
}

func TestParseIP(t *testing.T) {
	ip := net.ParseIP("10.40.0.0")
	assert.NotNil(t, ip)
}

func TestDNSEntry(t *testing.T) {
	id := RunContainer()
	newString := uuid.NewString()[:5]

	err := weaveplugins.AddDNSEntry(newString, id)
	assert.NoError(t, err)

	err = weaveplugins.DeleteDNSEntry(id)
	assert.NoError(t, err)
}

func TestSetWeaveEnv(t *testing.T) {
	err := weaveplugins.SetWeaveEnv()
	assert.NoError(t, err)
}

func TestGetIP(t *testing.T) {
	ip, err := weaveplugins.GetPodIPByID("0e4c6bcbb8ad")
	assert.NoError(t, err)
	assert.Equal(t, true, net.IP.Equal(ip, net.ParseIP("10.32.0.2")))

}

const (
	weaveName = "weave"
	launch    = "launch"
	sudo      = "sudo"
)

// delete containers and close weave
func PrepareTest() error {
	path, err := osexec.LookPath(weaveName)
	if err != nil {
		log.Println("Weave Not found.")
		return err
	}

	cmd := osexec.Command(path, "stop")
	err = cmd.Run()
	if err != nil {
		log.Panicf("Weave stop error: %s\n", err)
		return err
	}

	cmd = osexec.Command("docker", "ps", "-aq")
	byteOutput, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("rm docker error: %s\n", err)
		return err
	}

	// Warning: If you have docker running, be careful of this test
	if len(string(byteOutput)) != 0 {
		output := strings.ReplaceAll(string(byteOutput), "\n", " ")

		cmd = osexec.Command("docker", "stop", output)
		err = cmd.Run()

		cmd = osexec.Command("docker", "rm", output)
		err = cmd.Run()

		if err != nil {
			return err
		}
	}
	return nil
}

func RunContainer() string {
	cmd := osexec.Command("docker", "run", "-d", "-ti", "weaveworks/ubuntu")
	byteOutput, _ := cmd.CombinedOutput()

	return string(byteOutput)
}
