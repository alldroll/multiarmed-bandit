package multiarmed_bandit

import (
	"log"
	"testing"
)

func TestSimple(t *testing.T) {
	persistent := NewStorage()

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

	snapshot := NewSnapshotStorage(persistent)

	log.Printf("%v", snapshot.FindAll())

	snapshot.IncrementVariant("test2", 0, 1, 0)
	snapshot.IncrementVariant("test2", 0, 1, 0)
	snapshot.IncrementVariant("test2", 0, 0, 1)

	log.Printf("%+v", snapshot)

	err := snapshot.Update()
	if err != nil {
		t.Errorf("Unxpected error %v", err)
	}

	log.Printf("%+v", persistent.Find("test2").GetVariants()[0])
	log.Printf("%+v", snapshot.Find("test2").GetVariants()[0])

	snapshot.IncrementVariant("test2", 0, 1000000, 300000)
	log.Printf("%+v", snapshot)
	err = snapshot.Update()
	log.Printf("%+v", snapshot.Find("test2").GetVariants()[0])
}
