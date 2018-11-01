package multiarmed_bandit

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// deltas (give me more description)
type deltas = []uint64

// deltaMap (give me more description)
type deltaMap struct {
	sync.Map
}

// pack packs 2 uint32 in uint64
func pack(a, b uint32) uint64 {
	return (uint64(a) << 32) | uint64(b)
}

// unpack explode uint64 into 2 uint32
func unpack(v uint64) (uint32, uint32) {
	return uint32(v >> 32), uint32(v & 0xFFFFFF)
}

// increment adds showsDelta, rewardDelta to given experiment's variant
func (d *deltaMap) increment(name string, id uint32, showsDelta, rewardDelta uint32) error {
	v, ok := d.Load(name)
	if !ok {
		return fmt.Errorf("Given experiment %s is not found", name)
	}

	deltas := v.(deltas)
	variant := int(id)

	if variant >= len(deltas) {
		return fmt.Errorf("Given variant %d for experiment %s is not found", id, name)
	}

	atomic.AddUint64(&deltas[variant], pack(showsDelta, rewardDelta))
	return nil
}

// store persists deltas to given storage
func (d *deltaMap) store(storage Storage) error {
	var (
		err                    error
		pack                   uint64
		shows, rewards         uint32
		showsList, rewardsList []uint32 = []uint32{}, []uint32{}
	)

	d.Range(func(name, v interface{}) bool {
		showsList = showsList[:0]
		rewardsList = rewardsList[:0]
		variantsDeltas := v.(deltas)

		for variant, _ := range variantsDeltas {
			pack = atomic.SwapUint64(&variantsDeltas[variant], 0)
			shows, rewards = unpack(pack)
			showsList = append(showsList, shows)
			rewardsList = append(rewardsList, rewards)
		}

		err = storage.Increment(name.(string), showsList, rewardsList)
		return err == nil
	})

	return err
}

// actualize deltas structure with given snapshot
func (d *deltaMap) actualize(snapshot tableSnapshot) {
	var (
		variantsLen int
		v           interface{}
		ok          bool
	)

	for name, experiment := range snapshot {
		variantsLen = len(experiment.GetVariants())
		v, ok = d.Load(name)

		// if there is not the experiment or the variantsLen is not the same
		// store new value
		if !ok || len(v.(deltas)) != variantsLen {
			d.Store(
				name,
				make(deltas, variantsLen),
			)
		}
	}

	// clean unused deltas
	d.Range(func(name, v interface{}) bool {
		if _, ok := snapshot[name.(string)]; !ok {
			d.Delete(name)
		}

		return true
	})
}
