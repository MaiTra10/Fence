package internal

type Trader struct {
	Name     string
	ImageURL string
	Tasks    []Task
}

type Task struct {
	Name             string
	WikiURL          string
	Objectives       []string
	Rewards          []string
	PrereqTasks      [][]string
	OtherChoices     [][]string
	RequiredForKappa bool
}
