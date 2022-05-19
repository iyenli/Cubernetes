package localstorage

import (
	cubeconfig "Cubernetes/config"
	"Cubernetes/pkg/object"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

type Metadata struct {
	Node     object.Node
	MasterIP string
}

// SaveMeta save node metadata to log meta
func SaveMeta(meta Metadata) error {
	_, err := os.Stat(cubeconfig.MetaDir)
	if err != nil {
		err = os.MkdirAll(cubeconfig.MetaDir, 0666)
		if err != nil {
			log.Println("[FATAL]: fail to make metadata dir, err: ", err)
			return err
		}
	}

	file, err := os.OpenFile(cubeconfig.MetaFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Println("[FATAL] fail to open metadata file, err: ", err)
		return err
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Println("[WARNING] fail to close metadata file, err: ", err)
		}
	}(file)

	buf, err := json.Marshal(meta)
	if err != nil {
		log.Println("[FATAL] fail to marshal metadata, err: ", err)
		return err
	}

	write, err := file.Write(buf)
	if err != nil || write != len(buf) {
		log.Println("[FATAL] fail to write metadata file, err: ", err)
		return err
	}
	return nil
}

func LoadMeta() (Metadata, error) {
	f, err := os.OpenFile(cubeconfig.MetaFile, os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		log.Println("[FATAL] fail to open metadata file, err: ", err)
		return Metadata{}, err
	}

	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			log.Println("[WARNING] fail to close metadata file, err: ", err)
		}
	}(f)

	content, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println("[FATAL] fail to read metadata file, err: ", err)
		return Metadata{}, err
	}

	if len(content) == 0 {
		log.Println("[FATAL] metadata file is empty")
		return Metadata{}, errors.New("empty file")
	}

	var meta Metadata
	err = json.Unmarshal(content, &meta)
	if err != nil {
		log.Println("[FATAL] fail to parse metadata, err: ", err)
	}

	return meta, err
}

// TryLoadMeta load metadata without warnings
func TryLoadMeta() (Metadata, error) {
	f, err := os.OpenFile(cubeconfig.MetaFile, os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		return Metadata{}, err
	}

	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			log.Println("[WARNING] fail to close metadata file, err: ", err)
		}
	}(f)

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return Metadata{}, err
	}

	if len(content) == 0 {
		return Metadata{}, errors.New("empty file")
	}

	var meta Metadata
	err = json.Unmarshal(content, &meta)
	return meta, err
}

func ClearMeta() error {
	log.Println("Clearing local metadata...")
	err := os.RemoveAll(cubeconfig.MetaDir)
	if err != nil {
		log.Println("[FATAL] fail to clear metadata, err: ", err)
	}
	return err
}
