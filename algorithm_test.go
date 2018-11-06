package multiarmed_bandit

import (
	"math/rand"
	"testing"
)

func TestFlow(t *testing.T) {
	storage, err := NewSnapshotStorage(NewInMemoryStorage())
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	algos := []Algorithm{
		NewEpsilonGreedy(0.1, storage),
		NewThompsonSampling([2]float64{1, 1}, storage),
	}

	experiment := NewExperiment(
		"test1",
		0.1,
		[]Variant{
			NewVariant(1.0, 0, 0),
			NewVariant(1.0, 0, 0),
		},
	)

	probs := []float32{0.2, 0.25}

	for _, algo := range algos {
		storage.Save(experiment)
		storage.Update()

		for i := 0; i < 1000; i++ {
			choice, err := algo.Suggest("test1")
			if err != nil {
				t.Errorf("Unexpected error %v", err)
			}

			if err := algo.Show("test1", choice); err != nil {
				t.Errorf("Unexpected error %v", err)
			}

			if i%10 == 0 {
				storage.Update()
			}

			r := rand.Float32()

			if r <= probs[choice] {
				algo.Reward("test1", choice)
			}
		}

		storage.Update()
		choice, err := algo.Suggest("test1")

		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}

		if choice != 1 {
			t.Errorf("Expected 1, got %v", choice)
		}
	}
}
