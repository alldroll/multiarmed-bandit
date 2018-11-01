package multiarmed_bandit

import (
	"fmt"
	"math/rand"
)

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
func (e *epsilonGreedy) Choose(name string) (uint32, error) {
	experiment := e.storage.Find(name)
	if experiment == nil {
		return 0, fmt.Errorf("Experiment %v is not exists", name)
	}

	variants := experiment.GetVariants()
	lenVariants := len(variants)

	if rand.Float32() < e.epsilon {
		choice := uint32(rand.Intn(lenVariants))

		if err := e.storage.IncrementVariant(name, choice, 1, 0); err != nil {
			return 0, err
		}

		return choice, nil
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
	choice := uint32(best[bestIdx])

	if err := e.storage.IncrementVariant(name, choice, 1, 0); err != nil {
		return 0, err
	}

	return choice, nil
}

//
func (e *epsilonGreedy) Reward(name string, choice uint32) error {
	return e.storage.IncrementVariant(name, choice, 0, 1)
}
