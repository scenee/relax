package relfile

import (
	"gopkg.in/yaml.v2"
	"testing"
)

func TestMergeMap(t *testing.T) {
	old := map[string]interface{}{
		"Foo": "Hello",
		"Bar": 1,
		"Baz": map[string]interface{}{
			"Foo": "World",
			"Bar": 2,
		},
	}

	new := map[string]interface{}{
		"Foo": "Hello1",
		"Fiz": true,
		"Baz": map[string]interface{}{
			"Bar": "ok",
		},
	}

	merged := mergeMap(old, new)

	if v, ok := merged["Foo"]; !ok || v != "Hello1" {
		t.Error("Invalid merged value", merged)
	}

	if v, ok := merged["Bar"]; !ok || v != 1 {
		t.Error("Invalid merged value", merged)
	}

	if v, ok := merged["Fiz"]; !ok || v != true {
		t.Error("Invalid merged value", merged)
	}

	baz := merged["Baz"].(map[string]interface{})

	if v, ok := baz["Foo"]; !ok || v != "World" {
		t.Error("Invalid merged value", merged)
	}
	if v, ok := baz["Bar"]; !ok || v != "ok" {
		t.Error("Invalid merged value", merged)
	}
}

type Sample struct {
	Protocols []interface{} `yaml:"protocols"`
}

func TestYAMLArrayMerge(t *testing.T) {
	var (
		err error
		s   Sample
	)

	yamlData := []byte(`
# https://github.com/yaml/yaml/issues/35
base: &default
    - foo
    - bar

protocols:
    - <: *default
    - baz
`)

	err = yaml.Unmarshal(yamlData, &s)
	if err != nil {
		t.Error(err)
	}
	t.Log(s)
}
