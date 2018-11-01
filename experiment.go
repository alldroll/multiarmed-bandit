package multiarmed_bandit

//
type Experiment interface {
	//
	GetName() string
	//
	GetGamma() float32
	//
	GetVariants() []Variant
	//
	//Equals(other Experiment) bool
}

//
type Variant interface {
	//
	GetWeigth() float32
	//
	GetReward() uint32
	//
	GetShows() uint32
}

//
type experiment struct {
	name     string
	gamma    float32
	variants []Variant
}

//
func NewExperiment(name string, gamma float32, variants []Variant) *experiment {
	return &experiment{
		name:     name,
		gamma:    gamma,
		variants: variants,
	}
}

//
func (e *experiment) GetName() string {
	return e.name
}

//
func (e *experiment) GetGamma() float32 {
	return e.gamma
}

//
func (e *experiment) GetVariants() []Variant {
	return e.variants
}

//
func (e *experiment) Equals(other Experiment) {

}

//
type variant struct {
	weigth  float32
	rewards uint32
	shows   uint32
}

//
func NewVariant(weigth float32, rewards, shows uint32) *variant {
	return &variant{
		weigth:  weigth,
		rewards: rewards,
		shows:   shows,
	}
}

//
func (v *variant) GetWeigth() float32 {
	return v.weigth
}

//
func (v *variant) GetReward() uint32 {
	return v.rewards
}

//
func (v *variant) GetShows() uint32 {
	return v.shows
}
