package multiarmed_bandit

type Algorithm interface {
	//
	Choose(name string) (uint32, error)
	//
	Reward(name string, choice uint32) error
}
