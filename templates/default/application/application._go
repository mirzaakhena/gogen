package application

type RegistryContract interface {
	RegisterUsecase()
	RunApplication()
}

func Run(rv RegistryContract) {
	if rv != nil {
		rv.RegisterUsecase()
		rv.RunApplication()
	}
}
