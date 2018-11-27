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
		err            error
		s              Sample
		yamlData       []byte
		protos         []interface{}
		expectedProtos []string
		failed         bool
	)

	yamlData = []byte(`
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
	protos = cleanupMapValue(s.Protocols).([]interface{})
	expectedProtos = []string{"foo", "bar", "baz"}

	failed = false
	for i, s := range expectedProtos {
		_s, ok := protos[i].(string)
		if ok && s == _s {
			continue
		}
		failed = true
	}
	if failed {
		t.Error("Does not match", protos, expectedProtos)
	}

	yamlData = []byte(`
base: &default
    - foo
    - bar

protocols:
    - baz
    - <: *default
`)

	err = yaml.Unmarshal(yamlData, &s)
	if err != nil {
		t.Error(err)
	}
	protos = cleanupMapValue(s.Protocols).([]interface{})
	expectedProtos = []string{"baz", "foo", "bar"}

	failed = false
	for i, s := range expectedProtos {
		_s, ok := protos[i].(string)
		if ok && s == _s {
			continue
		}
		failed = true
	}
	if failed {
		t.Error("Does not match", protos, expectedProtos)
	}

	yamlData = []byte(`
base: &default
    - foo
    - bar

protocols:
    - <: *default
    - baz
    - <: *default
`)

	err = yaml.Unmarshal(yamlData, &s)
	if err != nil {
		t.Error(err)
	}
	protos = cleanupMapValue(s.Protocols).([]interface{})
	expectedProtos = []string{"foo", "bar", "baz", "foo", "bar"}

	failed = false
	for i, s := range expectedProtos {
		_s, ok := protos[i].(string)
		if ok && s == _s {
			continue
		}
		failed = true
	}
	if failed {
		t.Error("Does not match", protos, expectedProtos)
	}
}
