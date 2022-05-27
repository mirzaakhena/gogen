package command

type Runner interface {
	Run(inputs ...string) error
}
