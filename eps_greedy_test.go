package multiarmed_bandit

import "testing"

func TestFlow(t *testing.T) {
	storage, err := NewSnapshotStorage(NewInMemoryStorage())
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	storage.Save(
		NewExperiment(
			"test1",
			0.1,
			[]Variant{
				NewVariant(1.0, 0, 0),
				NewVariant(1.0, 0, 0),
			},
		),
	)

	storage.Update()

	algo := NewEpsilonGreedy(0.1, storage)

	for i := 0; i < 100; i++ {
		choice, err := algo.Choose("test1")
		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}

		if i%4 == 0 {
			storage.Update()
		}

		if choice == 0 && i%7 <= 1 {
			algo.Reward("test1", 0)
		}

		if choice == 1 && i%7 <= 0 {
			algo.Reward("test1", 1)
		}
	}

	storage.Update()

	choice, err := algo.Choose("test1")

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if choice != 0 {
		t.Errorf("Expected 0, got %v", choice)
	}
}
