package events

import "github.com/Papiermond/eventbus"

type SceneChangeEvent struct {
	Scene string
	Data  []any
}

func (e SceneChangeEvent) GetType() eventbus.EventType {
	return SCENE_CHANGE_EVENT
}
