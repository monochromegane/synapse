package main

import (
	"strings"

	"github.com/monochromegane/synapse"
)

type associator struct {
}

func (a associator) Associate(ctx synapse.Context, profile synapse.Profile) (synapse.Association, error) {
	y := strings.Split(profile["birth"], "-")[0]
	categoryID := string([]rune(y)[3])
	association := synapse.Association{"category_id": categoryID}
	return association, nil
}

func (a associator) Initialize() error {
	return nil
}

func (a associator) Finalize() error {
	return nil
}

func NewAssociator() synapse.Associator {
	return associator{}
}
