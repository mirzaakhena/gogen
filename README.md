# Gogen (Clean Architecture Code Generator)
Helping generate your boiler plate and code structure based on clean architecure.

## Introduction
Have you ever wondered how to apply the clean architecture properly and how to manage the layout of your go project's folder structure? This tools will help you to make it.

This generator have basic structure like this
```
application/
controller/
gateway/
infrastructure/
usecase/
main.go
```

The complete folder layout will be like this

```
application/
  registry/
    reg1.go
    reg2.go
  application.go
controller/
  version1/
    Usecase1.go
    Usecase2.go
  version2/
    Usecase1.go
    Usecase2.go
  interceptor.go    
gateway/
  hardcode/
    Usecase1.go
    Usecase2.go 
  inmemory/
    Usecase1.go
    Usecase2.go
  persistence/
    Usecase1.go
    Usecase2.go
  younameit/
    Usecase1.go
    Usecase2.go
infrastructure/
  httpserver/
    gracefully_shutdown.go
    http_handler.go		
usecase/
  usecase1/
    port/
      inport.go
      outport.go
    interactor.go
  usecase2/
    port/
      inport.go
      outport.go
    interactor.go
main.go
```

## Clean Architecture Concept
The main purpose of this architecture is separation between infrastructure part and logic part.
It must decoupling from any framework or library. 

## How to use the gogen?
You need to start by creating an usecase, then you continue create the gateway, then create the controller, and the last step is bind those three (usecase + gateway + controller) in registry part. That's it. 

To create the usecase, you need to understand the concept of usecase according to Uncle Bob's Clean Architecture article. Usecase has 3 main part.
* Input Port (Inport)
* Interactor
* Output Port (Outport)

```
Controller is using an Inport.
Inport is implemented by Interactor
Interactor is using an Outport
Outport is implemented by Gateway
```

*Inport* is an interface that has only one method (named `Execute`) that will be called by *Controller*. The method in interface define all the required (request and response) parameter to run the specific usecase. *Inport* will implemented by *Interactor*. Request and response struct is allowed to share to *Outport* under the same usecase, but must not shared to other *Inport* or *Outport* usecase.

*Interactor* is the place that you can define your logic flow. When a usecase logic need a data, it will ask the *Outport* to provide it. *Interactor* also use *Outport* to send data, store data or do some action to other service. *Interactor* only have one *Outport* field. We must not adding new *Outport* field to *Interactor* to keep a simplicity and consistency. 

*Outport* is a data and action provider for *Interactor*. *Outport* never know how it is implemented. The implementor (in this case a *Gateway*) will decide how to provide a data or do an action. This *Outport* is very exclusive for specific usecase (in this case *Interactor*) and must not shared to other usecase. By having exclusive *Outport* it will isolate the testing for usecase. 

By organize this usecase in a such structure, we can easily change the *Controller* or the *Gateway* in very flexible way. without worry to change the logic part. This is how the logic and infrastructure separation is working.

How it is different with common three layer architecture (Controller -> Service -> Repository) pattern?
The main different is 
* *Service* allowed to have many Repository. But in Clean Architecture, *Interactor* only have one *Outport*.
* *Service* have many method grouped by the domain. In Clean Architecture, we focus per usecase. One usecase for One Class to achieve *Single Responsibility Principle*.
* In *Repository* you often see CRUD pattern. Every developer can added new method if they think they need it. In reality this *Repository* is shared to different *Service* that may not use that method. In *Outport* you will strictly to adding method that guarantee used. Even adding new method or updating existing method will not interfere another usecase. 


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
usecase/error.go
```

`usecase/createorder/port/inport.go`
```
package port

import (
  "context"
)

// CreateOrderInport ...
type CreateOrderInport interface {
  Execute(ctx context.Context, req CreateOrderRequest) (*CreateOrderResponse, error)
}

// CreateOrderRequest ...
type CreateOrderRequest struct { 
}

// CreateOrderResponse ...
type CreateOrderResponse struct { 
}
```

`usecase/createorder/port/outport.go`
```
package port  

import (
  "context"
) 

// CreateOrderOutport ...
type CreateOrderOutport interface { 
  CreateOrder(ctx context.Context, req CreateOrderRequest) (*CreateOrderResponse, error) 
}
```

`usecase/createorder/interactor.go`
```
package createorder

import (
  "context"

  "your/go/path/project/usecase/createorder/port"
)

//go:generate mockery --dir port/ --name CreateOrderOutport -output mocks/

// NewCreateOrderUsecase ...
func NewCreateOrderUsecase(outputPort port.CreateOrderOutport) port.CreateOrderInport {
  return &createOrderInteractor{
    gateway: outputPort,
  }
}

type createOrderInteractor struct {
  gateway port.CreateOrderOutport
}

// Execute ...
func (_r *createOrderInteractor) Execute(ctx context.Context, req port.CreateOrderRequest) (*port.CreateOrderResponse, error) { 

  var res port.CreateOrderResponse  
	
  {
    resOutport, err := _r.gateway.CreateOrder(ctx, port.CreateOrderRequest { // 
    })

    if err != nil {
      return nil, err
    }
		
    _ = resOutport 
  } 

  return &res, nil
}
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
$ gogen usecase CreateOrder Check Save
```

`usecase/createorder/port/outport.go`
```
package port  

import (
	"context"
) 

// CreateOrderOutport ...
type CreateOrderOutport interface { 
	Check(ctx context.Context, req CheckRequest) (*CheckResponse, error) 
	Save(ctx context.Context, req SaveRequest) (*SaveResponse, error) 
}
 
// CheckRequest ...
type CheckRequest struct { 
} 

// CheckResponse ...
type CheckResponse struct { 
} 

// SaveRequest ...
type SaveRequest struct { 
} 

// SaveResponse ...
type SaveResponse struct { 
} 
```

`usecase/createorder/interactor.go`
```
package createorder

import (
	"context"

	"your/go/path/project/usecase/createorder/port"
)

//go:generate mockery --dir port/ --name CreateOrderOutport -output mocks/

// NewCreateOrderUsecase ...
func NewCreateOrderUsecase(outputPort port.CreateOrderOutport) port.CreateOrderInport {
	return &createOrderInteractor{
		gateway: outputPort,
	}
}

type createOrderInteractor struct {
	gateway port.CreateOrderOutport
}

// Execute ...
func (_r *createOrderInteractor) Execute(ctx context.Context, req port.CreateOrderRequest) (*port.CreateOrderResponse, error) { 

	var res port.CreateOrderResponse  
	
	{
		resOutport, err := _r.gateway.Check(ctx, port.CheckRequest { // 
		})

		if err != nil {
			return nil, err
		}
		
		_ = resOutport 
	} 
	
	{
		resOutport, err := _r.gateway.Save(ctx, port.SaveRequest { // 
		})

		if err != nil {
			return nil, err
		}
		
		_ = resOutport 
	} 

	return &res, nil
}
```

Now you will find outport has 2 new methods: Check and Save. Each of the method has it own Request Response struct in *Outport*. 


## 2. Add new outport method

Calling this command will add new method on your *Outport*'s interface
```
$ gogen outports CreateOrder ValidateLastOrder Publish
```
Open the *Outport* file you will found that there are 2 new methods defined. Now you have 4 methods : Check, Save, ValidateLastOrder and Publish


## 3. Create your usecase test file

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

`usecase/createorder/interactor_test.go`
```
package createorder

import (
	"context"
	"testing"

	"your/go/path/project/usecase/createorder/mocks"
	"your/go/path/project/usecase/createorder/port"
	"github.com/stretchr/testify/assert"
)

func Test_CreateOrder_Normal(t *testing.T) {

	ctx := context.Background()

	outputPort := mocks.CreateOrderOutport{} 
	{
		call := outputPort.On("Check", ctx, port.CheckRequest{ // 
		})
		call.Return(&port.CheckResponse{ // 
		}, nil)
	} 
	{
		call := outputPort.On("Save", ctx, port.SaveRequest{ // 
		})
		call.Return(&port.SaveResponse{ // 
		}, nil)
	} 
	{
		call := outputPort.On("Publish", ctx, port.PublishRequest{ // 
		})
		call.Return(&port.PublishResponse{ // 
		}, nil)
	} 

	res, err := NewCreateOrderUsecase(&outputPort).Execute(ctx, port.CreateOrderRequest{ // 
	})

	assert.Nil(t, err)

	assert.Equal(t, &port.CreateOrderResponse{ // 
	}, res)

}


```

If you want to update your mock file you can use
```
$ cd /usecase/createorder
$ go generate
```

or just simply delete the mock/ folder and call the `gogen test CreateOrder` again.

## 4. Create gateway for your usecase

Gateway is the struct to implement your outport interface. You need to set a name for your gateway. In this example we will set name : Production
```
$ gogen gateway Production CreateOrder
```
This command will generate
```
gateway/production/CreateOrder.go
```

`gateway/production/CreateOrder.go`
```
package production

import (
	"context"
	"log"

	"your/go/path/project/usecase/createorder/port"
)

type createOrder struct {
}

// NewCreateOrderGateway ...
func NewCreateOrderGateway() port.CreateOrderOutport {
	return &createOrder{}
}

// Check ...
func (_r *createOrder) Check(ctx context.Context, req port.CheckRequest) (*port.CheckResponse, error) {
	log.Printf("Gateway Check Request  %v", req) 

	var res port.CheckResponse 
	
	log.Printf("Gateway Check Response %v", res)
	return &res, nil
} 

// Save ...
func (_r *createOrder) Save(ctx context.Context, req port.SaveRequest) (*port.SaveResponse, error) {
	log.Printf("Gateway Save Request  %v", req) 

	var res port.SaveResponse 
	
	log.Printf("Gateway Save Response %v", res)
	return &res, nil
} 

// Publish ...
func (_r *createOrder) Publish(ctx context.Context, req port.PublishRequest) (*port.PublishResponse, error) {
	log.Printf("Gateway Publish Request  %v", req) 

	var res port.PublishResponse 
	
	log.Printf("Gateway Publish Response %v", res)
	return &res, nil
} 
```

You can give any gateway name you like. Maybe you want to experiment with just simple database with SQLite for testing purpose, you can create the "experimental" gateway version. 
```
$ gogen gateway Experimental CreateOrder
```

Or you want to have hardcode version, the you can run this command
```
$ gogen gateway Hardcode CreateOrder
```

## 5. Create controller for your usecase

In gogen, we define controller as trigger of the usecase. It can be rest api, grpc, consumer for event handling, or any other source input. By default it only use net/http restapi. 

Call this command for create a controller. Restapi is your controller name. You can name it whatever you want.
```
$ gogen controller restapi CreateOrder
```

It will generate
```
controller/restapi/CreateOrder.go
controller/interceptor.go
```

`controller/interceptor.go`
```
package restapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"your/go/path/project/usecase/createorder/port"
)

// CreateOrder ...
func CreateOrderHandler(inputPort port.CreateOrderInport) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		jsonReq, _ := ioutil.ReadAll(r.Body)

		log.Printf("Controller CreateOrderHandler Request  %v", string(jsonReq))

		var req port.CreateOrderRequest

		if err := json.Unmarshal(jsonReq, &req); err != nil {
			log.Printf("Controller CreateOrderHandler Response %v", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, err := inputPort.Execute(context.Background(), req)
		if err != nil {
			log.Printf("Controller CreateOrderHandler Response %v", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		jsonRes, _ := json.Marshal(res)
		fmt.Fprint(w, string(jsonRes))

		log.Printf("Controller CreateOrderHandler Response %v", string(jsonRes))

	}
}
```

You also will get the global interceptor for all of controller.


## 6. Glue your controller, usecase, and gateway together

After generate the usecase, gateway and controller, we need to bind them all by calling this command.
```
$ gogen registry Default Restapi Production CreateOrder
```
Default is the registry name. You can name it whatever you want. After calling the command, some of those file generated will generated for you
```
application/registry/Default.go
application/application.go
infrastructure/httpserver/gracefully_shutdown.go
infrastructure/httpserver/http_handler.go
```


`application/registry/Default.go`
```
package registry

import (
	"your/go/path/project/application"
	"your/go/path/project/infrastructure/httpserver"
	"your/go/path/project/controller"	
	"your/go/path/project/controller/restapi"
	"your/go/path/project/gateway/production"
	"your/go/path/project/usecase/createorder"
)

type defaultRegistry struct {
	httpserver.HTTPHandler
}

func NewDefaultRegistry() application.RegistryContract {

	app := defaultRegistry{ //
		HTTPHandler: httpserver.NewHTTPHandler(":8080"),
	}

	return &app

}

// RegisterUsecase is implementation of RegistryContract.RegisterUsecase()
func (r *defaultRegistry) RegisterUsecase() {
	r.createOrderHandler()
}

func (r *defaultRegistry) createOrderHandler() {
  outport := production.NewCreateOrderGateway()
  inport := createorder.NewCreateOrderUsecase(outport)
  r.ServerMux.HandleFunc("/createorder", controller.Authorized(restapi.CreateOrderHandler(inport)))
}

```

## 7. Create your own template

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