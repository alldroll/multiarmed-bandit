package multiarmed_bandit

//
type Algorithm interface {
	//
	Suggest(name string) (uint32, error)
	//
	Show(name string, choice uint32) error
	//
	Reward(name string, choice uint32) error
}
