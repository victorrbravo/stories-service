package graph

type Link struct {
	To    string `json:"to"`
	Query string `json:"query"`
}

type Task struct {
	Name      string   `json:"name"`
	TokenLink string   `json:"token_link"`
	Pattern   string   `json:"pattern"`
	Links     []Link   `json:"links"`
	Needs     []string `json:"needs"`
}

type Workflow struct {
	Name     string   `json:"name"`
	Variable string   `json:"variable"`
	Token    string   `json:"token"`
	Needs    []string `json:"needs"`
	Tasks    []Task   `json:"tasks"`
}
