package transformers

import (
	"fmt"
	"log"
	"strings"
)

type mutateFunc func(interface{}) (interface{}, error)

func mutateField(
	m map[string]interface{},
	pathToField []string,
	createIfNotPresent bool,
	fns ...mutateFunc) error {
	if len(pathToField) == 0 {
		return nil
	}

	_, found := m[pathToField[0]]
	if !found {
		if !createIfNotPresent {
			return nil
		}
		m[pathToField[0]] = map[string]interface{}{}
	}

	if len(pathToField) == 1 {
		var err error
		for _, fn := range fns {
			m[pathToField[0]], err = fn(m[pathToField[0]])
			if err != nil {
				return err
			}
		}
		return nil
	}

	v := m[pathToField[0]]
	newPathToField := pathToField[1:]
	switch typedV := v.(type) {
	case nil:
		log.Printf(
			"nil value at `%s` ignored in mutation attempt",
			strings.Join(pathToField, "."))
		return nil
	case map[string]interface{}:
		return mutateField(typedV, newPathToField, createIfNotPresent, fns...)
	case []interface{}:
		for i := range typedV {
			item := typedV[i]
			typedItem, ok := item.(map[string]interface{})
			if !ok {
				return fmt.Errorf("%#v is expected to be %T", item, typedItem)
			}
			err := mutateField(typedItem, newPathToField, createIfNotPresent, fns...)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("%#v is not expected to be a primitive type", typedV)
	}
}
