package label

type Resource interface {
	GetLabel() string
	Label(label string)
}

type Labeler interface {
	LabelResource(res Resource)
	IsLabeled(res Resource) bool
}
