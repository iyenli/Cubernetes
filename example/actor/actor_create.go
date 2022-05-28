package main

import (
	actor_runtime "Cubernetes/pkg/cubelet/actorruntime"
	"Cubernetes/pkg/object"
	"fmt"
	"time"
)

func main() {
	actor := object.Actor{
		TypeMeta: object.TypeMeta{
			APIVersion: "v1",
			Kind:       "Actor",
		},
		ObjectMeta: object.ObjectMeta{
			Name: "add_fuck",
			UID:  "fake_uid",
			Labels: map[string]string{
				"cubernetes.action.uid": "1145141919810"},
		},
		Spec: object.ActorSpec{
			ActionName:    "add",
			ScriptFile:    "action.py",
			InvokeActions: []string{},
		},
		Status: &object.ActorStatus{
			Phase: object.ActorCreated,
		},
	}
	
	runtime, err := actor_runtime.NewActorRuntime()
	if err != nil {
		panic(err)
	}

	if err = runtime.CreateActor(&actor); err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 10)

	if phase, err := runtime.InspectActor("fake_uid"); err != nil {
		panic(err)
	} else {
		fmt.Printf("phase is %v\n", phase)
	}
}