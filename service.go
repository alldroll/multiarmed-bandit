package multiarmed_bandit

import (
	"database/sql"
)

type Service interface {
	//
	GetStorage() Storage
	//
	GetAlgorithm() Algorithm
	//
	Update() error
}

type service struct {
	storage   Storage
	updatable Updatable
	algo      Algorithm
}

func NewService(db *sql.DB) (*service, error) {
	storage := NewSqlStorage(db)
	updatable, err := NewSnapshotStorage(storage)

	if err != nil {
		return nil, err
	}

	algo := NewEpsilonGreedy(0.1, updatable)

	return &service{
		storage:   storage,
		updatable: updatable,
		algo:      algo,
	}, nil
}

//
func (s *service) GetStorage() Storage {
	return s.storage
}

//
func (s *service) GetAlgorithm() Algorithm {
	return s.algo
}

//
func (s *service) Update() error {
	return s.updatable.Update()
}
