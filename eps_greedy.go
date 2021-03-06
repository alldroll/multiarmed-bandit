package multiarmed_bandit

import (
	"fmt"
	"math/rand"
)

//
func NewEpsilonGreedy(epsilon float32, storage Storage) *epsilonGreedy {
	return &epsilonGreedy{
		epsilon: epsilon,
		storage: storage,
	}
}

type epsilonGreedy struct {
	storage Storage
	epsilon float32
}

//
func (e *epsilonGreedy) Suggest(name string) (uint32, error) {
	experiment, err := e.storage.Find(name)

	if err != nil {
		return 0, err
	}

	if experiment == nil {
		return 0, fmt.Errorf("Experiment %v is not exists", name)
	}

	variants := experiment.GetVariants()
	lenVariants := len(variants)

	if rand.Float32() < e.epsilon {
		return uint32(rand.Intn(lenVariants)), nil
	}

	bestValue := float32(0.0)
	best := make([]int, 0)

	for i, variant := range variants {
		ctr := float32(variant.GetReward()+1) / float32(variant.GetShows()+1)

		if ctr > 1 {
			ctr = 1.0
		}

		if ctr > bestValue {
			best = best[:0]
			best = append(best, i)
			bestValue = ctr
		} else if ctr == bestValue {
			best = append(best, i)
		}
	}

	bestIdx := rand.Intn(len(best))
	return uint32(best[bestIdx]), nil
}

//
func (e *epsilonGreedy) Show(name string, choice uint32) error {
	return e.storage.IncrementVariant(name, choice, 1, 0)
}

//
func (e *epsilonGreedy) Reward(name string, choice uint32) error {
	return e.storage.IncrementVariant(name, choice, 0, 1)
}
