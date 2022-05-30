package actor_runtime

import (
	"Cubernetes/pkg/apiserver/objfile"
	"Cubernetes/pkg/cubelet/actorruntime/options"
	"Cubernetes/pkg/object"
	"log"
	"os"
	"path"
	"sync"
)

type ScriptManager interface {
	EnsureScriptExist(actor *object.Actor) error
	GetScriptDirPath(actor *object.Actor) string
}

const (
	metaFileName = "SCRIPT_META"
)

func NewScriptManager() ScriptManager {
	regDir := options.ScriptRegistryPath
	if _, err := os.Stat(regDir); err != nil && os.IsNotExist(err) {
		os.Mkdir(regDir, 0666)
	}
	return &scriptManager{scriptRegistryPath: regDir}
}

type scriptManager struct {
	scriptRegistryPath string
	lock               sync.Mutex
}

func (sm *scriptManager) EnsureScriptExist(actor *object.Actor) error {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	dir := path.Join(sm.scriptRegistryPath, actor.Spec.ActionName)
	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		os.Mkdir(dir, 0666)
	} else if err != nil {
		return err
	}

	script := path.Join(dir, "action.py")
	if _, err := os.Stat(script); err != nil && os.IsNotExist(err) {
		// script not exist, pull file from apiserver
		if err = objfile.GetActionFile(actor.Spec.ActionName, script); err != nil {
			log.Printf("fail to pull file %s from APIServer: %v", script, err)
			return err
		}

		meta := path.Join(dir, metaFileName)
		if err = os.WriteFile(meta, []byte(actor.Spec.ScriptUID), 0666); err != nil {
			log.Printf("fail to write script meta of %s: %v", actor.Spec.ActionName, err)
			return err
		}

		return nil
	} else if err != nil {
		return err
	}

	// script exists, check if UID changed
	meta := path.Join(dir, metaFileName)
	if uidBytes, err := os.ReadFile(meta); err == nil {
		if string(uidBytes) != actor.Spec.ScriptUID {
			// script updated, get script from apiserver
			if err = objfile.GetActionFile(actor.Spec.ActionName, script); err != nil {
				log.Printf("fail to pull file %s from APIServer: %v", script, err)
				return err
			}

			if err = os.WriteFile(meta, []byte(actor.Spec.ScriptUID), 0666); err != nil {
				log.Printf("fail to write script meta of %s: %v", actor.Spec.ActionName, err)
				return err
			}
		}
	} else if os.IsNotExist(err) {
		log.Printf("script of %s exists but meta not! why???????????????\n", actor.Spec.ActionName)
		return nil
	} else {
		return err
	}

	return nil
}

func (sm *scriptManager) GetScriptDirPath(actor *object.Actor) string {
	return path.Join(sm.scriptRegistryPath, actor.Spec.ActionName)
}
