package utils

import (
	"Cubernetes/pkg/cubeproxy/utils/options"
	"github.com/otiai10/copy"
	"log"
	"net"
	"os"
)

func CreateNginxConfig(host string, paths []string, serviceIPs []string, ports []string) error {
	if len(paths) != len(serviceIPs) || len(serviceIPs) != len(ports) {
		log.Fatalln("[Fatal]: different length of service and path")
	}
	if len(paths) == 0 {
		log.Println("[INFO]: No paths in dns")
		return nil
	}

	elements := make([]string, len(paths))
	for idx, path := range paths {
		dst := "http://" + net.JoinHostPort(serviceIPs[idx], ports[idx]) + "/"
		elements[idx] = CreateLocation(path, true, dst, false)
	}

	server := CreateServer(host, elements)
	config := server.String()

	log.Println("[INFO]: New DNS object created, nginx config file is:")
	log.Println(config)
	log.Printf("[INFO]: The file will be stored at '/etc/cubernetes/cubeproxy/nginx/%v/site-enabled/'\n", host)

	err := PrepareNginxFile(host, &config)
	if err != nil {
		log.Println("[Error]: save nginx config failed")
		return err
	}
	return nil
}

func PrepareNginxFile(hostname string, config *string) error {
	configFile := options.NginxFile + hostname
	if configFile[len(configFile)-1] != '/' {
		configFile = configFile + "/"
	}

	if _, err := os.Stat(configFile); err == nil {
		// remove stale config
		err := os.RemoveAll(configFile)
		if err != nil {
			log.Println("[Error]: delete exist config file failed")
			return err
		}
	}

	err := os.MkdirAll(configFile, 0666)
	if err != nil {
		log.Println("[Error]: Create directory failed")
		return err
	}

	err = copy.Copy(options.CopySrc, configFile)
	if err != nil {
		log.Println("[Error]: copy config folder failed")
		return err
	}

	// delete exist default file
	toDelete := configFile + options.SiteEnabled + options.DefaultFile
	err = DeleteIfExist(toDelete)
	if err != nil {
		return err
	}

	toDelete = configFile + options.SiteAvailable + options.DefaultFile
	err = DeleteIfExist(toDelete)
	if err != nil {
		return err
	}

	f, err := os.Create(configFile + options.DefaultFile)
	if err != nil {
		log.Println("[Error]: replace config file failed")
		return err
	}

	// write into
	writeString, err := f.WriteString(options.DefaultConfContent)
	if err != nil {
		log.Println("[Error]: write default config file failed")
		return err
	}

	if writeString != len(*config) {
		log.Println("[Warn]: not write whole file")
	}

	configFile = configFile + options.SiteEnabled
	err = os.MkdirAll(configFile, 0666)
	if err != nil {
		log.Println("[Error]: create config file failed")
		return err
	}

	configFile = configFile + hostname + options.Suffix
	log.Println("[INFO]: Finally, config is written to", configFile)
	// rm if exist
	err = DeleteIfExist(configFile)
	if err != nil {
		return err
	}

	f, err = os.Create(configFile)
	if err != nil {
		log.Println("[Error]: create config file failed")
		return err
	}

	// write into
	writeString, err = f.WriteString(*config)
	if err != nil {
		log.Println("[Error]: write config file failed")
		return err
	}

	if writeString != len(*config) {
		log.Println("[Warn]: not write whole file")
	}
	return nil
}

func DeleteIfExist(file string) error {
	if _, err := os.Stat(file); err == nil {
		err := os.Remove(file)
		if err != nil {
			log.Println("[Error]: delete exist config file failed")
			return err
		}
	}
	return nil
}
