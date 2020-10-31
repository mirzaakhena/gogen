package gogen

type Generator interface {
	Generate(args ...string) error
}
