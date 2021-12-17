package entity

type Component interface {
	Update()
}

type Entity struct {
	Components []Component
}
