package multiarmed_bandit

import "testing"

func TestSimple(t *testing.T) {
	persistent := NewInMemoryStorage()

	persistent.Save(
		NewExperiment(
			"test1",
			0.1,
			[]Variant{
				NewVariant(
					1.0,
					0,
					0,
				),
				NewVariant(
					1.0,
					0,
					0,
				),
			},
		),
	)

	persistent.Save(
		NewExperiment(
			"test2",
			0.1,
			[]Variant{
				NewVariant(
					1.0,
					0,
					0,
				),
				NewVariant(
					1.0,
					0,
					0,
				),
				NewVariant(
					1.0,
					0,
					0,
				),
			},
		),
	)

	snapshot, err := NewSnapshotStorage(persistent)
	if err != nil {
		t.Errorf("Unxpected error %v", err)
	}

	snapshot.IncrementVariant("test2", 0, 1, 0)
	snapshot.IncrementVariant("test2", 0, 1, 0)
	snapshot.IncrementVariant("test2", 0, 0, 1)

	err = snapshot.Update()
	if err != nil {
		t.Errorf("Unxpected error %v", err)
	}

	snapshot.IncrementVariant("test2", 0, 1000000, 300000)

	err = snapshot.Update()
	if err != nil {
		t.Errorf("Unxpected error %v", err)
	}

	// TODO check it
}
