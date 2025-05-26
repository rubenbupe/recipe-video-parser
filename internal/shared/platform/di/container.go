package di

import (
	"os"

	di "github.com/sarulabs/di/v2"
)

var env = os.Getenv("ENV")

var envDefs = map[string][]di.Def{
	"dev": devDefs,
}

type DiContainer struct {
	Container *di.Container
}

var container *DiContainer

func buildContainer() (*di.Container, error) {
	builder, err := di.NewEnhancedBuilder()
	if err != nil {
		return nil, err
	}

	for _, def := range defaultDefs {
		if err := builder.Add(&def); err != nil {
			return nil, err
		}
	}

	if envDefs, ok := envDefs[env]; ok {
		for _, def := range envDefs {
			if err := builder.Add(&def); err != nil {
				return nil, err
			}
		}
	}

	container, err := builder.Build()
	if err != nil {
		return nil, err
	}

	return &container, nil
}

func Instance() *DiContainer {
	if container == nil {
		con, err := buildContainer()
		if err != nil {
			print("Error building DI container")
			panic(err)
		}
		container = &DiContainer{Container: con}
	}

	return container
}

func (container DiContainer) GetByTag(tag string) []di.Def {
	defs := container.Container.Definitions()
	var result []di.Def
	for _, def := range defs {
		if def.Tags != nil {
			for _, t := range def.Tags {
				if t.Name == tag {
					result = append(result, def)
				}
			}
		}
	}
	return result
}
