package transformer

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/rantav/go-archetype/inputs"
)

func Read(transformationsFile string) (*Transformations, error) {
	yamlFile, err := ioutil.ReadFile(transformationsFile)
	if err != nil {
		return nil, err
	}
	var spec transformationsSpec
	err = yaml.Unmarshal(yamlFile, &spec)
	if err != nil {
		return nil, err
	}
	return FromSpec(spec)
}

func FromSpec(spec transformationsSpec) (*Transformations, error) {
	return &Transformations{
		ignore:       spec.Ignore,
		transformers: transformersFromSpec(spec.Transformations),
		prompters:    inputs.FromSpec(spec.Inputs),
		userInputs:   make(map[string]inputs.PromptResponse),
	}, nil
}

func transformersFromSpec(transformationSpecs []transformationSpec) []Transformer {
	var transformers []Transformer
	for _, t := range transformationSpecs {
		transformers = append(transformers, newTransformer(t))
	}
	return transformers
}

func newTransformer(spec transformationSpec) Transformer {
	return newTextReplacer(spec) // TODO: Add types here
}