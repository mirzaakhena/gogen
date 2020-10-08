# Gogen

Helping generate your boiler structure and boiler plate code based on clean architecure


## Step by step to working with gogen

## 1. Create basic project structure
```
gogen init .
```
this command will create the basic project structure in the current directory
```
application/application.go
application/registry.go
application/runner.go
application/schema.go
controller/
datasource/
model/
usecase/
util/
config.toml
main.go
README.md
```

## 2. Create your basic usecase structure

Let say you have a usecase named CreateOrder (It is better to have PascalCase for usecase name). This usecase is for (of course) to create an order. We can easily recognize it as a "command" usecase. Let's use our gogen code generator to create it for us.
```
gogen usecase command CreateOrder
```

When you run this command, you will have those files generated for you

```
usecase/createorder/port/inport.go
usecase/createorder/port/outport.go
usecase/createorder/interactor.go
```

## 3. Create your usecase test file

For test mock, we will need mockery. So you need to install it first
```
gogen test CreateOrder
```
This command will add new files
```
usecase/createorder/mocks/CreateOrderOutport.go
usecase/createorder/interactor_test.go
```

If you want to update your mock file you can use
```
$ cd /usecase/createorder
$ go generate
```

or just simply delete the mock/ folder and call the `gogen test CreteOrder` again.

## 4. Create datasource for your usecase
```
gogen datasource production CreateOrder
```
This will generate
```
datasource/production/CreateOrder.go
```

## 5. Create controller for your usecase
```
gogen controller restapi.gin CreateOrder
```
This will generate

```
controller/restapi/CreateOrder.go
```

## 6. Glue your usecase, datasource, and controller together
open file `application/registry.go`

```
package application

import (
	"your/golang/path/controller/restapi"
	"your/golang/path/datasource/production"
	"your/golang/path/usecase/createorder"
)

func (a *Application) RegisterUsecase() {
	createOrder(a)
}

func createorder(a *Application) {
	outport := production.NewCreateOrderDatasource()
	inport := createorder.NewCreateOrderUsecase(outport)
	a.Router.POST("/createorder", restapi.CreateOrder(inport))
}
```
For now this require manual effort to type manually (coding). Well actually we can do code injection. But this can be another feature to add in the future