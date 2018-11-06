package multiarmed_bandit

import (
	"fmt"
	"sync/atomic"
)

type Updatable interface {
	//
	Update() error
}

//
type tableSnapshot = map[string]Experiment

//
type snapshotStorage struct {
	persistentStorage Storage
	snapshotRef       atomic.Value
	deltaMap          deltaMap
}

//
func NewSnapshotStorage(persistentStorage Storage) (*snapshotStorage, error) {
	storage := &snapshotStorage{
		persistentStorage: persistentStorage,
		snapshotRef:       atomic.Value{},
		deltaMap:          deltaMap{},
	}

	// TODO should return err
	err := storage.takeSnapshot()
	if err != nil {
		return nil, err
	}

	return storage, nil
}

//
func (s *snapshotStorage) Find(name string) (Experiment, error) {
	table := s.snapshotRef.Load().(tableSnapshot)
	return table[name], nil
}

//
func (s *snapshotStorage) FindAll() ([]Experiment, error) {
	table := s.snapshotRef.Load().(tableSnapshot)
	list := make([]Experiment, 0, len(table))

	for _, entity := range table {
		list = append(list, entity)
	}

	return list, nil
}

//
func (s *snapshotStorage) Save(experiment Experiment) error {
	// TODO need to call Update?
	return s.persistentStorage.Save(experiment)
}

//
func (s *snapshotStorage) Delete(name string) error {
	// TODO need to call Update?
	return s.persistentStorage.Delete(name)
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
	return s.takeSnapshot()
}

// storeDeltas saves deltas to persitent storage
func (s *snapshotStorage) storeDeltas() error {
	return s.deltaMap.store(s.persistentStorage)
}

// takeSnapshot from persistent storage
func (s *snapshotStorage) takeSnapshot() error {
	snapshot := make(tableSnapshot)
	experiments, err := s.persistentStorage.FindAll()

	if err != nil {
		return err
	}

	for _, experiment := range experiments {
		snapshot[experiment.GetName()] = experiment
	}

	// we have to actualize deltas with given snapshot
	s.deltaMap.actualize(snapshot)

	// store new instance of snapshot
	s.snapshotRef.Store(snapshot)

	return nil
}
