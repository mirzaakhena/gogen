[![Go Report Card](https://goreportcard.com/badge/github.com/mirzaakhena/gogen)](https://goreportcard.com/report/github.com/mirzaakhena/gogen)

# Gogen Framework
A **Code Generator** and **Code Structure** provider which follow the **Clean Architecture** and **Domain Driven Design** concept

## The Problem Gogen Want To Solve
If we are googling looking for how the Clean Architecture code structure, 
then there will be many code examples in google, github and stackoverflow 
which claimed apply the Clean Architecture concept.

But in general, it is all just a project template. If you want to apply it, you must clone, imitate it or can just follow the folder and code structure. You are still doing anything manually like create a folder, create a file, naming a struct, an interface, basically everything.  

Gogen trying to go one step further. Instead of copying the project template, 
Gogen help you to generate a code that applies this Clean Architecture and Domain Driven Design concept. 

These tools not only bootstrap the code once in the beginning, 
you can use it progressively to generate other component like 
entity, value object, repository, use case, controller
so that developers don't have to struggle with imitating the template projects that aren't necessarily proven.

It helps you to focus on the core of logic and business process rather than manually create a file, name it, create a struct, interface or any other file.

## Architecture
![gogen architecture](https://github.com/mirzaakhena/gogen/blob/master/gogen-architecture-1.png)

## Structure
This generator has basic structure like this

```
  ├── .gogen      
  ├── application
  │   ├── app_one.go
  │   ├── app_two.go  
  │   └── app_three.go    
  ├── domain_order                    
  │   ├── controller
  │   │   └── restapi
  │   │       ├── controller.go
  │   │       ├── handler_getallorder.go
  │   │       ├── handler_getoneorder.go    
  │   │       ├── handler_onwebhook.go
  │   │       ├── handler_runordersubmit.go            
  │   │       ├── interceptor.go
  │   │       └── router.go                       
  │   ├── gateway
  │   │   └── prod
  │   │       └── gateway.go             
  │   ├── model
  │   │   ├── entity
  │   │   │   ├── order.go
  │   │   │   └── order_history.go
  │   │   ├── enum
  │   │   │   └── order_status.go   
  │   │   ├── errorenum
  │   │   │   └── error_enum.go          
  │   │   ├── service
  │   │   │   └── order_service.go
  │   │   ├── repository
  │   │   │   └── order_repo.go    
  │   │   └── vo
  │   │       └── order_id.go            
  │   ├── usecase
  │   │   ├── getallorder
  │   │   │   ├── inport.go
  │   │   │   ├── interactor.go
  │   │   │   ├── outport.go   
  │   │   │   └── README.md                 
  │   │   ├── getoneorder  
  │   │   ├── onpaymentsuccess
  │   │   ├── onpaymentfail
  │   │   └── runordersubmit   
  │   └── README.md                     
  ├── domain_product                
  │   ├── controller          
  │   ├── gateway
  │   ├── model             
  │   ├── usecase
  │   └── README.md     
  ├── domain_shipment                   
  │   ├── controller          
  │   ├── gateway   
  │   ├── model     
  │   ├── usecase              
  │   └── README.md    
  ├── shared
  │   ├── config
  │   ├── gogen            
  │   ├── infrastructure          
  │   │   ├── cache
  │   │   ├── database
  │   │   ├── logger
  │   │   ├── messaging
  │   │   ├── remoting
  │   │   ├── server
  │   │   └── token 
  │   ├── model                             
  │   └── util
  ├── .gitignore    
  ├── config.json    
  ├── config.sample.json  
  ├── docker-compose.yml
  ├── Dockerfile                  
  ├── go.mod 
  ├── main.go                        
  └── README.md
```

## Gogen Features
- Simple And Flexible Config (`shared/config/data.go`)
- Interchangeable log with trace_id (`shared/infrastructure/logger`)
- Adjustable response code (`shared/infrastructure/model/payload/response.go`)
- Error message catalog with error code (`domain_yourdomain/model/errorenum`)
- Customizable code template (`.gogen/templates`)
- Infrastructure Agnostic (you can choose your own framework and library)

## What is Clean Architecture ?
Clean Architecture is an architectural pattern for designing software systems that aims to achieve separation of concerns and maintainability by keeping the codebase independent of any particular UI framework, database, or external service. It is based on the principles of "SOLID" and "DDD" (Domain-Driven Design).

The purpose of Clean Architecture is to create a software system that is easy to understand, maintain, and change over time. It accomplishes this by structuring the system into layers that enforce a clear separation of concerns, where each layer has a specific responsibility and dependencies only flow inwards.

The core idea of Clean Architecture is to establish a clear separation of concerns between the business logic, the application logic, and the infrastructure details. The business logic, which encapsulates the essential rules and behaviors of the system, should be independent of any implementation details, such as UI, database, or third-party services. This allows the business logic to be tested and developed independently from the infrastructure.

## Video Tutorial how to use it
https://www.youtube.com/playlist?list=PLWBGlxJNCxvoONzJPKfLZxFQ0CDxkbECa

## Sample New Code
- https://github.com/mirzaakhena/theitem Code sample for CRUD implementation
- https://github.com/mirzaakhena/gogen_grpc_graphql Code sample for gRPC and GraphQL Implementation
- https://github.com/mirzaakhena/gogen_pubsub Code sample for Redis PubSub and Kafka Implementation

## Sample Old Code 
those code are obsolete due to gogen evolution. But the concept is still relevant
- https://github.com/mirzaakhena/userprofile
- https://github.com/mirzaakhena/danarisan
- https://github.com/mirzaakhena/kraicklist
- https://github.com/mirzaakhena/mywallet

## Install Gogen
Install it into your local system
```
$ go install -v github.com/mirzaakhena/gogen@v0.0.22
```
After calling that command, gogen will be installed in `/Users/<username>/go/bin`. 

Try to call 
```shell
$ gogen
```

It must shown something like this
```shell
Try one of this command to learn how to use it
  gogen domain
  gogen entity
  gogen valueobject
  gogen valuestring
  gogen enum
  gogen usecase
  gogen repository
  gogen service
  gogen test
  gogen gateway
  gogen controller
  gogen error
  gogen application
  gogen crud
  gogen webapp
  gogen web
```

### Troubleshooting


If you find something like
```shell
Command 'gogen' not found
```

You need to make sure that You are already adding `/Users/<username>/go/bin` into your path.
If you are in mac/linux, add this like into you .bashrc or .bash_profile or .
```text
PATH=$PATH:/usr/local/go/bin
PATH=$PATH:$HOME/go/bin

export PATH
```

## Step by step to working with gogen

Before we start, make sure you already have `go.mod` file.
If you don't have it, create the new one with this command
(for example your module name is : "your/awesome/project")
```shell
$ go mod init your/awesome/project
```

There is minimal 7 step to work with gogen
1. create the domain
2. create the entity, value object, enum or error
3. create a use case
4. create a repository or service
5. create a gateway as outport implementation
6. create a controller which call the inport
7. create a application which bind controller, use case and gateway


## Create a domain

You need to create the domain first. Let say you want to create order domain
```
$ gogen domain order
```

Then you will see some file and folder created for you

## Create your basic use case structure

So you will create your first use case. Let say the use case name is  a `RunOrderCreate`. We will always create our use case name with `PascalCase`. Now let's try our gogen code generator to create this use case for us.
```
$ gogen usecase RunOrderCreate
```

In Gogen we have some use case name convention

```
Run<SomeEntityName><Action>
GetAll<SomeEntityName><Action>
GetOne<SomeEntityName><Action>
```

- `Run` prefix used for command use case that is actively called by external service, something like restapi, or grpc.
- `GetAll` prefix used for use case that return a list. 
- `GetOne` (or just Get) prefix used for use case that return single object.

Use case name will be used as a package name under use case folder by lower-casing the use case name.

- `domain_order/usecase/runordercreate/inport.go` is an interface with one method that will implement by your usecase. The standart method name is a `Execute`.
- `domain_order/usecase/runordercreate/outport.go` is an interface which has many methods that will be used by your usecase. It must not shared to another usecase.
- `domain_order/usecase/runordercreate/interactor.go` is the core implementation of the usecase (handle your bussiness application). It implements the method from inport and call the method from outport.

## Create your use case test file
```
$ gogen test normal RunOrderCreate
```
normal is the test name and RunOrderCreate is the use case name.
This command will help you
- create a test file `domain_order/usecase/runordercreate/testcase_normal_test.go`

## Create a repository 
```
$ gogen repository SaveOrder Order RunOrderCreate
```
This command will help you 
- create a Repository named `SaveOrderRepo` under `domain_order/usecase/model/repository/repository.go` (if it not exists yet)
- create an Entity with name `Order` under `domain_order/usecase/model/entity/order.go` (if it not exists yet)
- inject `repository.SaveOrderRepo` into `domain_order/usecase/runordercreate/outport.go` 
- and inject code template into `domain_order/usecase/runordercreate/interactor.go`. Injected code will be appear if `//!` is found in interactor's file.

Usecase name in Create a Repository command is optional so you can call it too without injecting it to the usecase
```
$ gogen repository SaveOrder Order
```

## Create a service
```
$ gogen service PublishMessage RunOrderCreate
```
This command will help you to
- create a Service named `PublishMessageService` under `domain_order/usecase/model/service/service.go`
- inject `service.PublishMessageService` into `domain_order/usecase/runordercreate/outport.go`
- Inject code template into `domain_order/usecase/runordercreate/interactor.go`. Injected code will be appeared if `//!` is found in interactor's file.

## Create a gateway for your usecase

Gateway is the struct to implement your outport interface. You need to set a name for your gateway. 
In this example we will set name : prod
```
$ gogen gateway prod
```
This command will read the Outport of `runordercreate` usecase and implement all the method needed in `domain_order/gateway/prod/gateway.go`

## Create a controller for your usecase

In gogen, we define a controller as trigger of the usecase. 
It can be rest api, grpc, consumer for event handling, or any other source input. 
By default, it only uses gin/gonic web framework. 
Call this command for create a controller. 
`restapi` is your controller name. Controller name can be grouped by client who use the API
```
$ gogen controller restapi
```

You need to download some dependency after this step by calling
```
$ go mod tidy

go: finding module for package github.com/matoous/go-nanoid
go: finding module for package github.com/gin-contrib/cors
go: finding module for package github.com/gin-gonic/gin
go: found github.com/gin-gonic/gin in github.com/gin-gonic/gin v1.7.7
go: found github.com/gin-contrib/cors in github.com/gin-contrib/cors v1.3.1
go: found github.com/matoous/go-nanoid in github.com/matoous/go-nanoid v1.5.0

```

## Glue your controller, usecase, and gateway together

After generate the usecase, gateway and controller, we need to bind them all by calling this command.
```
$ gogen application appone
```
appone is the application name. After calling the command, some of those file generated will generate for you in `application/app_one.go`

Now you can run the application by opening the terminal then type

```
$ go run main.go

You may try 'go run main.go <app_name>' :
 - appone
```
You will see that message, tell you to run the complete command. type once again the complete command like this

```
$ go run main.go appone

[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /ping                     --> gogendemo/shared/infrastructure/server.NewGinHTTPHandler.func1 (3 handlers)
[GIN-debug] POST   /runordersubmit           --> gogendemo/domain_order/controller/restapi.(*Controller).runOrderSubmitHandler.func1 (5 handlers)
{"appName":"appone","appInstID":"MIRZ","start":"2022-05-27 22:19:40","severity":"INFO","message":"0000000000000000 server is running at :8080","location":"server.(*GracefullyShutdown).RunWithGracefullyShutdown:40","time":"2022-05-27 22:19:40"}
```

Congratulation! your application is running

You can also adding another component into your code.

## Create entity
entity is a mutable object that has an identifier. This command will create new entity struct under `domain_order/usecase/model/entity/` folder
```
$ gogen entity OrderHistory
```

## Create valueobject
valueobject is an immutable object that has no identifier. This command will create new valueobject under `domain_order/usecase/model/vo/` folder
```
$ gogen valueobject FullName FirstName LastName 
```

## Create valuestring
valuestring is a valueobject simple string type. This command will create new valuestring struct under `domain_order/usecase/model/vo/` folder
```
$ gogen valuestring OrderID
```

## Create enum
enum is a single immutable value. This command will create new enum struct under `domain_order/usecase/model/enum/` folder
```
$ gogen enum PaymentMethod DANA Gopay Ovo LinkAja
```

## Create error enum
error enum is a shared error collection. This command will add new error enum line in `domain_order/usecase/model/errorenum/error_enum.go` file
```
$ gogen error SomethingGoesWrongError
```

## Too many command and I cannot remember it!

Phew!! I knew you cannot remember those all command. Don't worry if you forget the command, then you can just call 'gogen'
```
$ gogen

Try one of this command to learn how to use it
  gogen application
  gogen web
  gogen valueobject
  gogen repository
  gogen service
  gogen gateway
  gogen error
  gogen crud
  gogen webapp
  gogen domain
  gogen usecase
  gogen entity
  gogen valuestring
  gogen enum
  gogen controller
  gogen test
  gogen openapi
```

It will help you to remember what command do you want to use. Let say you forgot how to use the `gogen repository` command

Just call like this
```
$ gogen repository

   # Create a repository and inject the template code into interactor file with '//!' flag
   gogen repository SaveOrder Order RunOrderCreate
     'SaveOrder'      is a repository method name
     'Order'          is an entity name
     'RunOrderCreate' is an usecase name

   # Create a repository without inject the template code into usecase
   gogen repository SaveOrder Order
     'SaveOrder' is a repository func name
     'Order'     is an entity name
     
```

It will show sample command to remind you on how to use it. In this case, repository has 2 type of command.  

## Working in Another domain 

You can create another domain let say you want to create domain payment

```
$ gogen domain payment
```

Then you want to create a use case under domain payment
```
$ gogen usecase RunPaymentCreate
```
After you run it, you will find the use case still created under domain order.

How to work under payment domain?

open the folder `.gogen` you will see file `gogenrc.json`. Just update the domain field, from `order`

```
{
  "domain": "order",
  "controller": "gin",
  "gateway": "simple",
  "crud": "gin"
}
```

into `payment`

```
{
  "domain": "payment",
  "controller": "gin",
  "gateway": "simple",
  "crud": "gin"
}
```

That's it, now you are working under payment domain. 

Re-run the command again
```
$ gogen usecase RunPaymentCreate
```

Now you got the usecase boilerplate code created under payment domain. (don't forget to delete your previous RunPaymentCreate use case under the order domain manually)

---

