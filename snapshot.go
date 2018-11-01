package multiarmed_bandit

import (
	"fmt"
	"sync/atomic"
)

//
type tableSnapshot = map[string]Experiment

//
type snapshotStorage struct {
	persistentStorage Storage
	snapshotRef       atomic.Value
	deltaMap          deltaMap
}

//
func NewSnapshotStorage(persistentStorage Storage) *snapshotStorage {
	storage := &snapshotStorage{
		persistentStorage: persistentStorage,
		snapshotRef:       atomic.Value{},
		deltaMap:          deltaMap{},
	}

	storage.takeSnapshot()
	return storage
}

//
func (s *snapshotStorage) Find(name string) Experiment {
	table := s.snapshotRef.Load().(tableSnapshot)
	return table[name]
}

//
func (s *snapshotStorage) FindAll() []Experiment {
	table := s.snapshotRef.Load().(tableSnapshot)
	list := make([]Experiment, 0, len(table))

	for _, entity := range table {
		list = append(list, entity)
	}

	return list
}

//
func (s *snapshotStorage) Save(experiment Experiment) error {
	return s.persistentStorage.Save(experiment)
}

//
func (s *snapshotStorage) IncrementVariant(name string, id uint32, showsDelta, rewardsDelta uint32) error {
	return s.deltaMap.increment(name, id, showsDelta, rewardsDelta)
}

//
func (s *snapshotStorage) Increment(name string, showsDelta, rewardsDelta []uint32) error {
	if len(showsDelta) != len(rewardsDelta) {
		return fmt.Errorf("len(showsDelta) should be equal to len(rewardsDelta)")
	}

	for i := 0; i < len(showsDelta); i++ {
		if err := s.IncrementVariant(name, uint32(i), showsDelta[i], rewardsDelta[i]); err != nil {
			return err
		}
	}

	return nil
}

// Update implements Updatable interface
func (s *snapshotStorage) Update() error {
	if err := s.storeDeltas(); err != nil {
		return err
	}

	// take snapshot from persistent storage
	s.takeSnapshot()
	return nil
}

// storeDeltas saves deltas to persitent storage
func (s *snapshotStorage) storeDeltas() error {
	return s.deltaMap.store(s.persistentStorage)
}

// takeSnapshot from persistent storage
func (s *snapshotStorage) takeSnapshot() {
	snapshot := make(tableSnapshot)

	for _, experiment := range s.persistentStorage.FindAll() {
		snapshot[experiment.GetName()] = experiment
	}

	// we have to actualize deltas with given snapshot
	s.deltaMap.actualize(snapshot)

	// store new instance of snapshot
	s.snapshotRef.Store(snapshot)
}
