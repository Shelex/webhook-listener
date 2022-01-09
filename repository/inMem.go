package repository

import (
	"errors"
	"time"

	"github.com/Shelex/webhook-listener/entities"
)

type inMem struct {
	hooks map[string][]entities.Hook
}

func NewInMem() (Storage, error) {
	DB = &inMem{
		hooks: map[string][]entities.Hook{},
	}
	return DB, nil
}

func (i inMem) Add(hook entities.Hook) error {
	hook.Created_at = time.Now().UTC().Unix()
	i.hooks[hook.Channel] = append(i.hooks[hook.Channel], hook)
	return nil
}

func (i inMem) Get(channel string, pagination Pagination) ([]entities.Hook, int64, error) {
	if _, ok := i.hooks[channel]; !ok {
		return nil, 0, errors.New("hooks for channel not found")
	}

	start, end := Paginate(pagination, int64(len(i.hooks[channel])))
	page := i.hooks[channel][start:end]

	return page, int64(len(i.hooks[channel])), nil
}

func (i inMem) Delete(channel string) error {
	_, ok := i.hooks[channel]
	if !ok {
		return errors.New("channel not found")
	}
	delete(i.hooks, channel)
	return nil
}

func (i inMem) ClearExpired() error {
	expired := GetExpiryDate()

	for channel, hooks := range i.hooks {
		deleted := 0
		for index := range hooks {
			j := index - deleted
			if hooks[j].Created_at <= expired {
				i.hooks[channel] = hooks[:j+copy(hooks[j:], hooks[j+1:])]
				deleted += 1
			}
		}
		if len(i.hooks[channel]) == 0 {
			delete(i.hooks, channel)
		}
	}
	return nil
}
