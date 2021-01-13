# Gogen (Clean Architecture Code Generator)
Helping generate your boiler plate and code structure based on clean architecure.

## Introduction
Have you wondering how to apply the clean architecture properly and how to manage your go project folder structure layout? This tools will help you to generate it.

This generator have basic structure like this
```
application/
controller/
gateway/
usecase/
main.go
```

The complete folder layout will be like this

```
application/
  registry
    reg1.go
    reg2.go
    reg3.go
  application.go
  gracefully_shutdown.go
  http_handler.go
controller/
  version1/
    Usecase1.go
    Usecase2.go
    Usecase3.go
  version2/
    Usecase1.go
    Usecase2.go
    Usecase3.go
  interceptor.go    
gateway/
  hardcode/
    Usecase1.go
    Usecase2.go
    Usecase3.go    
  inmemory/
    Usecase1.go
    Usecase2.go
    Usecase3.go  
  persistence/
    Usecase1.go
    Usecase2.go
    Usecase3.go  
  younameit/
    Usecase1.go
    Usecase2.go
    Usecase3.go  
usecase/
  usecase1/
    port
      inport.go
      outport.go
    interactor.go
  usecase2/
    port
      inport.go
      outport.go
    interactor.go  
  usecase3/    
    port
      inport.go
      outport.go
    interactor.go  
main.go
```

## Clean Architecture Concept
The main goal of this architecture is separation between infrastructure part and logic part

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

## 5. Create gateway for your usecase

Gateway is the struct to implement your outport interface. You need to set a name for your gateway. In this example we will set name : Production
```
$ gogen gateway Production CreateOrder
```
This command will generate
```
gateway/production/CreateOrder.go
```

You can give any gateway name you like. Maybe you want to experiment with just simple database with SQLite for testing purpose, you can create the "experimental" gateway version. 
```
$ gogen gateway Experimental CreateOrder
```

Or you want to have hardcode version, the you can run this command
```
$ gogen gateway Hardcode CreateOrder
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


## 7. Glue your controller, usecase, and gateway together

After generate the usecase, gateway and controller, we need to bind them all by calling this command.
```
$ gogen registry Default Restapi Production CreateOrder
```
Default is the registry name. You can name it whatever you want. After calling the command, some of those file generated will generated for you
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
	"your/apps/path/gateway/production"
	"your/apps/path/usecase/createorder"
)

type defaultRegistry struct {
	application.HTTPHandler
}

func NewDefaultRegistry() application.RegistryVersion {

	app := defaultRegistry{ //
		HTTPHandler: application.NewHTTPHandler(":8080"),
	}

	return &app

}

// RegisterUsecase is implementation of RegistryVersion.RegisterUsecase()
func (r *defaultRegistry) RegisterUsecase() {
	r.createOrderHandler()
}

func (r *defaultRegistry) createOrderHandler() {
  outport := production.NewCreateOrderGateway()
  inport := createorder.NewCreateOrderUsecase(outport)
  r.ServerMux.HandleFunc("/createorder", controller.Authorized(restapi.CreateOrderHandler(inport)))
}
```

## 8. Create your own template

You can customize your own template by call this command in your local path
```
$ gogen init
```
You will have '.gogen/templates/default' folder which have all the template needed base on default gogen cleann architecture. You may update or change the template localy and run the previous command using your template.

another feature will coming
```
  gracefully shutdown
  interceptor
  config
  log
    file
    rotate
    session id
  swagger
  crud
  database
  message_broker
  http_client
  error collection
  token
  bcrypt password
  user login auth
  controller_query
    paging
    sorting
    filtering
  router
  page
    input_dialog, tabel_list
  auth
  login, register, forgot pass
  Domain Driven Design
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
```