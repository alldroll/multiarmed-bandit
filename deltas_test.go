package multiarmed_bandit

import (
	"reflect"
	"testing"
)

func TestIncrement(t *testing.T) {
	m := deltaMap{}

	cases := []struct {
		name             string
		v, shows, reward uint32
	}{
		{
			"test1",
			1, 1, 0,
		},
		{
			"test2",
			1, 1, 0,
		},
		{
			"test1",
			1, 10, 0,
		},
		{
			"test2",
			1, 4, 3,
		},
		{
			"test1",
			2, 4, 3,
		},
	}

	m.Store("test1", make(deltas, 3))
	m.Store("test2", make(deltas, 2))

	for _, c := range cases {
		err := m.increment(c.name, c.v, c.shows, c.reward)
		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}
	}

	v, _ := m.Load("test1")
	expected := deltas{0, pack(11, 0), pack(4, 3)}

	if !reflect.DeepEqual(v.(deltas), expected) {
		t.Errorf("Test fail, expected %v, got %v", expected, v)
	}

	v, _ = m.Load("test2")
	expected = deltas{0, pack(5, 3)}

	if !reflect.DeepEqual(v.(deltas), expected) {
		t.Errorf("Test fail, expected %v, got %v", expected, v)
	}
}

func TestIncrementOnNonexistentData(t *testing.T) {
	m := deltaMap{}

	if err := m.increment("test1", 1, 10, 1); err == nil {
		t.Errorf("Expected error, got nil")
	}

	m.Store("test1", make(deltas, 2))

	if err := m.increment("test1", 2, 10, 1); err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestActualize(t *testing.T) {
	experiments := []Experiment{
		NewExperiment(
			"test1",
			0.1,
			[]Variant{
				NewVariant(1.0, 10, 0),
				NewVariant(1.0, 10, 5),
			},
		),
		NewExperiment(
			"test2",
			0.2,
			[]Variant{
				NewVariant(1.0, 0, 0),
				NewVariant(1.0, 0, 0),
			},
		),
	}

	snapshot := make(tableSnapshot)
	for _, experiment := range experiments {
		snapshot[experiment.GetName()] = experiment
	}

	// actualize called first time
	m := deltaMap{}
	m.actualize(snapshot)

	m.increment("test1", 0, 10, 1)
	m.actualize(snapshot)

	// deltas structure should not be changed
	v, _ := m.Load("test1")
	expected := deltas{pack(10, 1), 0}

	if !reflect.DeepEqual(v.(deltas), expected) {
		t.Errorf("Test fail, expected %v, got %v", expected, v)
	}

	// deltas structure should be reseted, because we have changed test1's variants
	snapshot["test1"] = NewExperiment(
		"test2",
		0.2,
		[]Variant{
			NewVariant(1.0, 0, 0),
			NewVariant(1.0, 0, 0),
			NewVariant(1.0, 0, 0),
		},
	)
	m.actualize(snapshot)

	v, _ = m.Load("test1")
	expected = deltas{0, 0, 0}

	if !reflect.DeepEqual(v.(deltas), expected) {
		t.Errorf("Test fail, expected %v, got %v", expected, v)
	}
}

func TestStore(t *testing.T) {
	experiments := []Experiment{
		NewExperiment(
			"test1",
			0.1,
			[]Variant{
				NewVariant(1.0, 10, 0),
				NewVariant(1.0, 10, 5),
			},
		),
		NewExperiment(
			"test2",
			0.2,
			[]Variant{
				NewVariant(1.0, 0, 0),
				NewVariant(1.0, 0, 0),
			},
		),
	}

	m := deltaMap{}
	storage := NewStorage()
	snapshot := make(tableSnapshot)

	for _, experiment := range experiments {
		snapshot[experiment.GetName()] = experiment
		storage.Save(experiment)
	}

	m.actualize(snapshot)
	m.increment("test1", 0, 10, 3)
	m.increment("test1", 1, 31, 10)
	m.increment("test2", 0, 1000, 3)

	if err := m.store(storage); err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	/*

		expected1 := NewExperiment(
			"test1",
			0.1,
			[]Variant{
				NewVariant(1.0, 20, 3),
				NewVariant(1.0, 41, 15),
			},
		)
		experiment1 := storage.Find("test1")

		if !reflect.DeepEqual(expected1, experiment1) {
			t.Errorf("Test fail, expected %v, got %v", expected1, experiment1)
		}*/
}

func BenchmarkIncrement(b *testing.B) {
	m := deltaMap{}
	m.Store("test1", make(deltas, 3))

	for i := 1; i < b.N; i++ {
		m.increment("test1", uint32(i%3), 1, uint32(i%2))
	}
}
