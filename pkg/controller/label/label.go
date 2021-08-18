package label

import (
	"strings"
)

type Resource interface {
	GetLabel() string
	AddLabel(label string)
}

type Labeler interface {
	GetLabelName() string
	Label(res Resource)
	IsLabeled(res Resource) bool
}

type BaseLabeler struct {
	LabelName string
}

func (l BaseLabeler) GetLabelName() string {
	return l.LabelName
}

func (l BaseLabeler) Label(res Resource) {
	res.AddLabel(l.LabelName)
}

func (l BaseLabeler) IsLabeled(res Resource) bool {
	if res != nil {
		return strings.Compare(l.LabelName, res.GetLabel()) == 0
	}
	return false
}

func (l BaseLabeler) Equals(labeler Labeler) bool {
	if labeler == nil {
		return false
	}
	return strings.Compare(l.LabelName, labeler.GetLabelName()) == 0
}
