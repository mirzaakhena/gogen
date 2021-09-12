package application

type RegistryContract interface {
	SetupController()
	RunApplication()
}

func Run(rv RegistryContract) {
	if rv != nil {
		rv.SetupController()
		rv.RunApplication()
	}
}
