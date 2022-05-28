package actor_runtime

import (
	"Cubernetes/pkg/apiserver/objfile"
	"Cubernetes/pkg/cubelet/actorruntime/options"
	"Cubernetes/pkg/object"
	"log"
	"os"
	"path"
)

type ScriptManager interface {
	EnsureScriptExist(actor *object.Actor) error
	GetScriptDirPath(actor *object.Actor) string
}

func NewScriptManager() ScriptManager {
	regDir := options.ScriptRegistryPath
	if _, err := os.Stat(regDir); err != nil && os.IsNotExist(err) {
		os.Mkdir(regDir, 0666)
	}
	return &scriptManager{scriptRegistryPath: regDir}
}

type scriptManager struct {
	scriptRegistryPath string
}

func (sm *scriptManager) EnsureScriptExist(actor *object.Actor) error {
	dir := path.Join(sm.scriptRegistryPath, actor.Spec.ActionName)
	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		os.Mkdir(dir, 0666)
	} else if err != nil {
		return err
	}

	script := path.Join(dir, actor.Spec.ScriptFile)
	if _, err := os.Stat(script); err != nil && os.IsNotExist(err) {
		// Pull script file from apiserver
		if err = objfile.GetActionFile(actor.Spec.ActionName, script); err != nil {
			log.Printf("fail to pull file %s from APIServer: %v", script, err)
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

func (sm *scriptManager) GetScriptDirPath(actor *object.Actor) string {
	return path.Join(sm.scriptRegistryPath, actor.Spec.ActionName)
}
