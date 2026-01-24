package internal

var eftWikiBaseURL = "https://escapefromtarkov.fandom.com"

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
	PrereqTasks      []RelatedTask
	OtherChoices     []RelatedTask
	RequiredForKappa bool
}

type RelatedTask struct {
	Name    string
	WikiURL string
}
