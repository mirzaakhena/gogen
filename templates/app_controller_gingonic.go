package templates

import _ "embed"

// Registry
var (
	//go:embed controller/gingonic/app-registry._go
	ApplicationRegistryGinFile string
)

var (
	//go:embed controller/gingonic/response._go
	ControllerResponseFile string

	//go:embed controller/gingonic/inport._go
	ControllerInportFile string

	//go:embed controller/gingonic/router-struct._go
	ControllerGinFile string

	//go:embed controller/gingonic/router-register._go
	ControllerBindRouterGinFile string

	//go:embed controller/gingonic/interceptor._go
	ControllerInterceptorGinFile string

	//go:embed controller/gingonic/handler-func._go
	ControllerFuncGinFile string
)

// Server
var (
	//go:embed controller/gingonic/handler-server._go
	ServerFile string
)
