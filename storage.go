package multiarmed_bandit

// Storage represents layer for manipulating entities
type Storage interface {
	// Find returns Experiment by given name
	Find(name string) Experiment
	// FindAll returns content
	FindAll() []Experiment
	// Save saves experiment in storage
	Save(experiment Experiment) error
	// IncrementVariant increments shows and rewards for given experiment's variant.
	IncrementVariant(name string, id uint32, showsDelta, rewardsDelta uint32) error
	// Increment increments shows and rewards for given experiment's variants.
	Increment(name string, showsDelta, rewardsDelta []uint32) error
}
