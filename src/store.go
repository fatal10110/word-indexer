package main

import (
	"log"
	"sync"
)

type Store interface {
	GetStat(word string) (int, bool)
	UpdateStats(stats IndexResults) error
	Clear() error
}

func NewMapStore() Store {
	return &mapStore{store: make(IndexResults)}
}

type mapStore struct {
	sync.RWMutex
	store IndexResults
}


// GetStat locks the map to prevent writing, then reads the stat for specific word
func (ms *mapStore) GetStat(word string) (int, bool) {
	ms.RLock()
	defer ms.RUnlock()

	return ms.store.Get(word)
}

// UpdateStats locks the map to prevent read and write, then updates the map data
func (ms *mapStore) UpdateStats(stats IndexResults) error {
	ms.Lock()
	defer ms.Unlock()

	for word, stat := range stats {
		ms.store.Add(word, stat)
	}

	log.Println("store updated")

	return nil
}

func (ms *mapStore) Clear() error {
	ms.store = make(IndexResults)

	return nil
}