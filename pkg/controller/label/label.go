package label

import (
	"strings"

	"k8s.io/klog/v2"
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

func LabelResources(labeler Labeler, resources []Resource) {
	if labeler != nil && len(resources) > 0 {
		for _, resource := range resources {
			labeler.Label(resource)
		}
	} else {
		klog.Warningf("label resource labeler is %v and resources lens %d", labeler, len(resources))
	}
}
