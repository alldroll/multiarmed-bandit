package multiarmed_bandit

import (
	"fmt"
	"gonum.org/v1/gonum/stat/distuv"
)

func NewThompsonSampling(prior [2]float64, storage Storage) *thompsonSampling {
	return &thompsonSampling{
		prior:   prior,
		storage: storage,
	}
}

type thompsonSampling struct {
	storage Storage
	prior   [2]float64
}

//
func (t *thompsonSampling) Suggest(name string) (uint32, error) {
	experiment, err := t.storage.Find(name)

	if err != nil {
		return 0, err
	}

	if experiment == nil {
		return 0, fmt.Errorf("Experiment %v is not exists", name)
	}

	bestValue := float64(0.0)
	choice := uint32(0)
	weights := t.computeWeights(experiment)

	for i, w := range weights {
		if w > bestValue {
			bestValue = w
			choice = uint32(i)
		}
	}

	return choice, nil
}

//
func (t *thompsonSampling) Show(name string, choice uint32) error {
	return t.storage.IncrementVariant(name, choice, 1, 0)
}

//
func (t *thompsonSampling) Reward(name string, choice uint32) error {
	return t.storage.IncrementVariant(name, choice, 0, 1)
}

//
func (t *thompsonSampling) computeWeights(experiment Experiment) []float64 {
	weights := []float64{}

	for _, variant := range experiment.GetVariants() {
		dist := distuv.Beta{
			Alpha: t.prior[0] + float64(variant.GetReward()),
			Beta:  t.prior[1] + float64(variant.GetShows()-variant.GetReward()),
		}

		weights = append(weights, dist.Rand())
	}

	return weights
}
