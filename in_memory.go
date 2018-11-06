package multiarmed_bandit

import (
	"fmt"
)

type memoryStorage struct {
	data map[string]Experiment
}

func NewInMemoryStorage() *memoryStorage {
	return &memoryStorage{
		data: map[string]Experiment{},
	}
}

func (m *memoryStorage) Find(name string) (Experiment, error) {
	return m.data[name], nil
}

func (m *memoryStorage) FindAll() ([]Experiment, error) {
	data := []Experiment{}

	for _, experiment := range m.data {
		data = append(data, experiment)
	}

	return data, nil
}

func (m *memoryStorage) Save(experiment Experiment) error {
	m.data[experiment.GetName()] = experiment
	return nil
}

func (m *memoryStorage) IncrementVariant(name string, id uint32, showsDelta, rewardDelta uint32) error {
	experiment, _ := m.Find(name)
	if experiment == nil {
		return fmt.Errorf("Experiment %s is not found", name)
	}

	variants := experiment.GetVariants()
	variantId := int(id)
	if variantId >= len(variants) {
		return fmt.Errorf("There is not variant %d in experiment %s is not found", id, name)
	}

	variant := variants[variantId]
	variants[variantId] = NewVariant(
		variant.GetWeigth(),
		variant.GetReward()+rewardDelta,
		variant.GetShows()+showsDelta,
	)

	return nil
}

func (m *memoryStorage) Increment(name string, showsDelta, rewardsDelta []uint32) error {
	if len(showsDelta) != len(rewardsDelta) {
		return fmt.Errorf("len(showsDelta) should be equal to len(rewardsDelta)")
	}

	for i := 0; i < len(showsDelta); i++ {
		if err := m.IncrementVariant(name, uint32(i), showsDelta[i], rewardsDelta[i]); err != nil {
			return err
		}
	}

	return nil
}
