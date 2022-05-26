package action_brain

import (
	"Cubernetes/pkg/actionbrain/controller"
	"Cubernetes/pkg/actionbrain/informer"
	"log"
	"sync"
)

type ActionBrain interface {
	Run()
}

type actionBrainManager struct {
	actionController controller.ActionController
	actorInformer    informer.ActorInformer
	actionInformer   informer.ActionInformer

	wg sync.WaitGroup
}

func NewActionBrain() (ActionBrain, error) {
	wg := sync.WaitGroup{}

	actorInformer, _ := informer.NewActorInformer()
	actionInformer, _ := informer.NewActionInformer()

	actionController, err := controller.NewActionController(
		actorInformer, actionInformer, wg)
	if err != nil {
		log.Printf("fail to create ActionController: %v\n", err)
		return nil, err
	}

	return &actionBrainManager{
		actionController: actionController,
		actorInformer:    actorInformer,
		actionInformer:   actionInformer,
		wg:               wg,
	}, nil
}

func (abm *actionBrainManager) Run() {
	go abm.actionController.Run()

	abm.wg.Wait()

	go abm.actionInformer.ListAndWatchActionsWithRetry()
	abm.actorInformer.ListAndWatchActorsWithRetry()

	log.Fatalln("Unreachable here")
}
