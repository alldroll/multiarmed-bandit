package multiarmed_bandit

// Storage represents layer for manipulating entities
type Storage interface {
	// Find returns Experiment by given name
	Find(name string) (Experiment, error)
	// FindAll returns content
	FindAll() ([]Experiment, error)
	// Save saves experiment in storage
	Save(experiment Experiment) error
	// Delete deletes experiment by given name
	Delete(name string) error
	// IncrementVariant increments shows and rewards for given experiment's variant.
	IncrementVariant(name string, id uint32, showsDelta, rewardsDelta uint32) error
	// Increment increments shows and rewards for given experiment's variants.
	Increment(name string, showsDelta, rewardsDelta []uint32) error
}
