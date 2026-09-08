package languages

import "fmt"

type Definition struct {
	ID          string
	DisplayName string
	FileName    string
	Command     []string
}

type Registry struct {
	definitions map[string]Definition
}

func NewRegistry(pythonCommand string) *Registry {
	if pythonCommand == "" {
		pythonCommand = "python"
	}
	return &Registry{definitions: map[string]Definition{
		"python": {
			ID:          "python",
			DisplayName: "Python",
			FileName:    "main.py",
			Command:     []string{pythonCommand},
		},
	}}
}

func (r *Registry) Get(id string) (Definition, error) {
	definition, ok := r.definitions[id]
	if !ok {
		return Definition{}, fmt.Errorf("linguagem não suportada: %s", id)
	}
	return definition, nil
}
