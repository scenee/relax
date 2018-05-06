package relfile

import (
	"fmt"
	"testing"
)

func TestUpdatedMap(t *testing.T) {
	old := map[string]interface{}{
		"Foo": "Hello",
		"Bar": 1,
		"Baz": map[string]interface{}{
			"Foo": "World",
			"Bar": 2,
		},
	}

	new := map[string]interface{}{
		"Fiz": true,
		"Foo": "Hello1",
		"Baz": map[string]interface{}{
			"Foo": "World1",
			"Bar": 2,
		},
	}

	fmt.Println(updatedMap(old, new))
}
