# Gogen

Helping generate your boiler plate and code structure based on clean architecure.

## Download it
```
$ go get github.com/mirzaakhena/gogen
```
Install it into your local system
```
$ go install $GOPATH/src/github.com/mirzaakhena/gogen/
```

## Step by step to working with gogen

## 1. Create your basic usecase structure

So you want to create your first usecase. Let say the usecase name is  a `CreateOrder`. We always create our usecase name with `PascalCase`. This is a `gogen convention`. Now let's try our gogen code generator to create this usecase for us.
```
$ gogen usecase CreateOrder
```

After you run this command, you will have those files generated for you

```
usecase/createorder/port/inport.go
usecase/createorder/port/outport.go
usecase/createorder/interactor.go
```

Usecase name will be used as a package name under usecase folder by lowercasing the usecase name.

`port/inport.go` is an interface with one method that will implement by your usecase. The standart method name is a `Execute`.

`port/outport.go` is an interface which has many method that will be used by your usecase. It must not shared to another usecase.

`interactor.go` is the core implementation of the usecase (handle your bussiness application). It implement the method from inport and call the method from outport.

You also will found that both inport and outport share the same Request and Response struct from the inport.
The Request struct name is `CreateOrderRequest` and Response struct name is `CreateOrderResponse`.
You may create different Request Response struct for outport by using this command. We will start over the usecase creation by removing it

```
$ rm -rf usecase/
$ gogen usecase CreateOrder Save Publish
```

Now you will find outport has 2 new methods: Save and Publish. Each of the method has it own Request Response struct. SaveRequest, SaveResponse, PublishRequest, PublishResponse


## 3. Add new outport method

Calling this command will add new method on your `outport`'s interface
```
$ gogen outports CreateOrder Validate CheckLastOrder
```
Open the outport file you will found that there are 3 new methods defined.

If you want those new method is shown in interactor, you may delete the existing interactor file and run again the usecase command
```
$ rm usecase/interactor.go
$ gogen usecase CreateOrder
```
Open your new interactor file. Now you see you have the template of outport method called in your interactor's file. 

## 4. Create your usecase test file

For test mock struct, we use mockery. So you need to install it first. See the guideline how to install it in https://github.com/vektra/mockery

If you already have the mockery, call this command will create the test file for the respective usecase
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

or just simply delete the mock/ folder and call the `gogen test CreateOrder` again.

## 5. Create datasource for your usecase

Datasource is the struct to implement your outport interface. You need to set a name for your datasource. In this example we will set name : Production
```
$ gogen datasource Production CreateOrder
```
This command will generate
```
datasource/production/CreateOrder.go
```

You can give any datasource name you like. Maybe you want to experiment with just simple database with SQLite for testing purpose, you can create the "experimental" datasource version. 
```
$ gogen datasource Experimental CreateOrder
```

Or you want to have hardcode version, the you can run this command
```
$ gogen datasource Hardcode CreateOrder
```

## 6. Create controller for your usecase

In gogen, we define controller as input gateway from outside world. It can be rest api, grpc, consumer for event handling, or anything. For now i only implement the gin framework version for restapi.

Call this command for create a controller. Restapi is your controller name. You can name it whatever you want.
```
$ gogen controller Restapi CreateOrder
```

It will generate
```
controller/restapi/CreateOrder.go
controller/interceptor.go
```

You also will get the global interceptor for all of controller


## 7. Glue your controller, usecase, and datasource together

After generate the usecase, datasource and controller, we need to bind them all by calling this command.
```
$ gogen registry Default Restapi Production CreateOrder
```
Default is the registry name. You cann name it wahtever you want. After you call the command, you will get some of those file generated
```
application/registry/Default.go
application/application.go
application/gracefully_shutdown.go
application/http_handler.go
```

Then open file `application/registry/Default.go` then you will find this
```
package registry

import (
	"your/apps/path/application"
	"your/apps/path/controller/restapi"
	"your/apps/path/datasource/production"
	"your/apps/path/usecase/createorder"
)

type defaultRegistry struct {
	application.GinHTTPFramework
}

func NewDefaultRegistry() application.RegistryVersion {

	app := defaultRegistry{ //
		GinHTTPFramework: application.NewGinHTTPFramework(":8080"),
	}

	return &app

}

// RegisterUsecase is implementation of RegistryVersion.RegisterUsecase()
func (r *defaultRegistry) RegisterUsecase() {
	r.createOrderHandler()
	//code_injection function call
}

func (r *defaultRegistry) createOrderHandler() {
	outport := production.NewCreateOrderDatasource()
	inport := createorder.NewCreateOrderUsecase(outport)
	r.Router.POST("/createorder", restapi.CreateOrder(inport))
}

//code_injection function declaration
```
gogen registry will inject some code in `registry.go`. Basically you can also write it by yourself, 
but you need to notice that we have two comment line. 
```
//code_injection function call

and

//code_injection function declaration
```
gogen will look at that comment to do the code injection.
if you remove that comment line, gogen registry will no longer work anymore.


Another improvement is:
- create outport request response in separate file. We can create the filename like:
	outport-Save.go

- listing the existing usecases by just call `gogen usecase`

- define the technology in the begining when call `gogen init`

Any other interesting idea?

feature:
  gracefully shutdown
  interceptor
  
  config
  log, rotate log
  database
  message_broker
  http_client

  error

  extractor

  token
  bcrypt password
  
  controller_query
    paging
    sorting
    filtering

  router
  page
    input_dialog, tabel_list
  auth
  login, register, forgot pass
  

  http://blog.opus.ch/2019/01/ddd-concepts-and-patterns-service-and-repository/

Domain Service
Application Service
Infrastructure Service


Tactical Pattern
  Entity
  ValueObject
  Service
  Repository
  Factory
  Event
  Aggregate

Strategic design
  BoundedContext