[![Go Report Card](https://goreportcard.com/badge/github.com/mirzaakhena/gogen)](https://goreportcard.com/report/github.com/mirzaakhena/gogen)

# Gogen (Clean Architecture Code Generator)
Provide code structure based on clean architecure


## The problem gogen wants to solve
CLEAN ARCHITECTURE (CA) is a concept of "composing and organizing folder and code files in the very maintainable ways" which has the benefit of separating code logic and infrastructure very neatly so that it is easy to test, easy to mock and easy to switch between technologies with very few changes.

This concept is agnostic against the technology. Means it does not depend on specific programming language. 

If we are googling looking for how the CA structure, then there will be many code examples that implement CA in many programming language. But in general it is just a project template. If we want to apply it, we must imitate it.

I'm trying to go one step further. Instead of copying the project template, why not create a code generator that applies this CA concept. These tools will help bootstrap the code so that developers don't have to struggle with imitating previous projects that aren't necessarily proven.

The process of drafting the concept and making it is also not easy and requires time and coding experience. I had to do some research by reading and studying dozens of articles on the internet about Clean Architecture, Clean Code, Solid Design, Domain Driven Design.

I try to empathize with programmers. Try to feel what they think when they want to write code. For example, "what should I create first? where should I put the controller? what is the proper name for this file?" I try to pour all these feelings and thoughts into this tool. So that it can guide programmers in coding activities.

Some principles I apply are
1. These tools should not be "know-it-all" tools. The programmer should still be the master of design. Because I don't want these tools to drive logic programmers instead. This tool only helps to guide to write standard code templates with clear names and conventions. The rest we still give the programmer space to work.
2. This tool has several alternatives to choose the technology. So if the programmer has better technology or is more familiar, the programmer can easily replace it.
3. I apply the Scream Architecture concept in it so that the generated code can speak for itself to the developers about what their role is and what they are doing (helping the learning process).

Some benefits that can be obtained if you apply this tool are:
1. These tools can become standard in a team. I love innovation and improvisation. However, if innovation and improvisation do not have a clear concept, it is feared that it will mislead the development process and complicate the process of changing or adding requirements in the future.
2. Because it has become a standard, this tool help the communication process between developers QA, project manager, and product owner, 
3. Help the handover process and knowledge transfer with new programmers because it is easy to learn and imitated.
4. Speed up the code review process and minimize code conflicts during code merges.
5. The code generated results in a readable, simple structure with few directories and a minimum depth that has been calculated very carefully.
6. Facilitate the creation of story cards. a standard structure will help shape the mindset of project managers when making stories. For example every member of developer team can have task per usecase. 

However, this is just a tools. The important things to remember is we must follow the basic principles of clean architecture itself. You may copy and paste the existing code if you think it is easier. But remember you have to be careful anytime you do that. 

## Video Tutorial how to use it
https://www.youtube.com/playlist?list=PLWBGlxJNCxvoONzJPKfLZxFQ0CDxkbECa

## Sample Apps
- https://github.com/mirzaakhena/userprofile
- https://github.com/mirzaakhena/danarisan
- https://github.com/mirzaakhena/kraicklist
- https://github.com/mirzaakhena/mywallet

> For those video and sample apps structure has become obsoleted due to the gogen evolution.

## Documentation (still under development)
https://mirzaakhena.github.io

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
  ├── domain_payment                
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
  │   ├── driver          
  │   ├── gateway   
  │   ├── infrastructure          
  │   │   ├── cache
  │   │   ├── config
  │   │   ├── database
  │   │   ├── logger
  │   │   ├── messaging
  │   │   ├── remoting
  │   │   ├── server
  │   │   └── token 
  │   ├── model
  │   ├── usecase                                       
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

![gogen architecture](https://github.com/mirzaakhena/gogen/blob/master/gogen-architecture-1.png)



## Clean Architecture Concept
The main purpose of this architecture is :
* Separation concern between INFRASTRUCTURE part and LOGIC Part
* Independent of Framework. Free to swap to any framework
* Independent of UI. Free to swap UI. For ex: from Web UI to Console UI
* Independent of Database. Free to swap to any database (data storage)
* Independent of any external agency. Business rule doesn’t know anything about outside world
* Testable. The business rules can be tested without the UI, Database, Web Server, or any other external element


## How to use the gogen?
You always start from creating an usecase, then you continue to create the gateway, then create the controller, and the last step is bind those three (usecase + gateway + controller) in application part. That's it.

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

*Inport* is an interface that has only one method (named `Execute`) that will be called by *Controller*. The method in an interface define all the required (request and response) parameter to run the specific usecase. *Inport* will implemented by *Interactor*. 

*Interactor* is the place that you can define your business process flow by involving entity, valueobject, service and repository. When an usecase logic need a data, it will ask the *Outport* to provide it. *Interactor* also use *Outport* to send data, store data or do some action to other service. *Interactor* only have one *Outport* field. We must not adding new *Outport* field to *Interactor* to keep a simplicity and consistency.

*Outport* is a data and action provider for *Interactor*. *Outport* never know how it is implemented. The implementor (in this case a *Gateway*) will decide how to provide a data or do an action. This *Outport* is very exclusive for specific usecase (in this case *Interactor*) and must not shared to other usecase. By having exclusive *Outport* it will isolate the testing for usecase.

By organize this usecase in a such structure, we can easily change the *Controller*, or the *Gateway* in very flexible way without worry to change the logic part. This is how the logic and infrastructure separation is working.

## Comparison with three layer architecture (Controller -> Service -> Repository) pattern
* *Controller* is the same controller for both architecture.
* *Service* is similar like *Interactor* with additional strict rule. *Service* allowed to have many repositories. In gogen Clean Architecture, *Interactor* only have one *Outport*.
* *Service* have many method grouped by the domain. In Clean Architecture, we focus per usecase. One usecase for One Class to achieve *Single Responsibility Principle*.
* In *Repository* you often see the CRUD pattern. Every developer can added new method if they think they need it. In reality this *Repository* is shared to different *Service* that may not use that method. In *Outport* you will strictly to adding method that guarantee used. Even adding new method or updating existing method will not interfere another usecase.
* *Repository* is an *Outport* with *Gateway* as its implementation.

## Gogen Convention
* Usecase is first class citizen. 
* We always start the code from the usecase.
* Usecase doesn't care about what the technology will be used
* Usecase basically manage interaction between entity, value object, repository and service
* As interface, inport has one and only one method which handle one usecase
* Interactor is a manager for entity and outport
* Interactor must not decide to have any technology. All technology is provided by Outport.
* Outport is a servant for an interactor. It will provide any data from external source
* As interface, Outport at least have one method and can have multiple method
* All method in Outport is guarantee used by interactor
* Inport and Outport must not shared to other usecase. They are exclusive for specific usecase
* Interactor have one and only one outport. Mutiple outport is prohibited
* Entity is mandatory to make a validation for any input. Controller optionaly handle the validation
* Interactor is the one who made a decision mostly based on Entity consideration.
* Entity must not produce unpredictible value like current time, or UUID value. Unpredictible value is provided by Inport or Outport (external Entity)
* Since a log is technology, log only found in Controller and Gateway not in Interactor or Entity
* To avoid a log polution, Log only printing the coming request, leaving response or error response
* Error code can be produced by anyone and will printed in log
* If somehow Gateway produce an error it may log the error, forward back the error,
  forward back the error with new error message or all the possibility
* Error code at least have messaged and code (imitate the http protocol response code)
* Error enum must accessed by the developer and Error code can read by end user
* Interactor and Entity is prioritized to be tested first rather than Controller and Gateway
* Controller name can be an actor name who is using the system

## Why you (will) need gogen?
- we want to separate logic code and infrastructure code
- save time because we think less for naming (one of "hardest" think in programming), conventions and structure
- Increase readability, scream architecture
- built for lazy developer
- consistent structure
- gogen is zero dependency. Your code will not have dependency to gogen at all
- gogen is not engine it just a "well written code" so there is no performance issue
- gogen is good templating tools, because deleting is easier than creating right
- gogen support multiple application in one repo
- gogen already implement trace id in every usecase Request
- support lazy documentation. Your interactor is telling everything about how it's work. no need to work twice only to write/update doc
- suitable for new project and revamp existing project per service
- allow you to do code modification for experimental purpose without changing the current implementation
- there is no automagically in gogen. 

## Install Gogen
Install it into your local system
```
$ go install -v github.com/mirzaakhena/gogen@v0.0.2
```

## Step by step to working with gogen

## Create a domain

You need to create the domain first. Let say you want to create order domain
```
$ gogen domain order
```

Then you will see some file and folder created for you

## Create your basic usecase structure

So you will create your first usecase. Let say the usecase name is  a `RunOrderCreate`. We will always create our usecase name with `PascalCase`. Now let's try our gogen code generator to create this usecase for us.
```
$ gogen usecase RunOrderCreate
```

But wait, why the name is very awkward?

In Gogen we have some usecase name convention

```
Run<SomeEntityName><Action>
On<SomeEntityName><Action>
GetAll<SomeEntityName><Action>
GetOne<SomeEntityName><Action>
```

Run prefix used for command usecase that is actively called by external service, something like restapi. 
GetAll prefix used for usecase that return a list. 
GetOne (or just Get) prefix used for usecase that return single object.
On prefix used for command usecase that is passively called by internal service, something like messagebroker event.

Usecase name will be used as a package name under usecase folder by lowercasing the usecase name.

- `domain_order/usecase/runordercreate/inport.go` is an interface with one method that will implement by your usecase. The standart method name is a `Execute`.
- `domain_order/usecase/runordercreate/outport.go` is an interface which has many methods that will be used by your usecase. It must not shared to another usecase.
- `domain_order/usecase/runordercreate/interactor.go` is the core implementation of the usecase (handle your bussiness application). It implements the method from inport and call the method from outport.

## Create your usecase test file
```
$ gogen test normal RunOrderCreate
```
normal is the test name and RunOrderCreate is the usecase name.
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
- Inject code template into `domain_order/usecase/runordercreate/interactor.go`. Injected code will be appear if `//!` is found in interactor's file.

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
   gogen repository SaveOrder Order CreateOrder
     'SaveOrder'   is a repository func name
     'Order'       is an entity name
     'CreateOrder' is an usecase name

   # Create a repository without inject the template code into usecase
   gogen repository SaveOrder Order
     'SaveOrder' is a repository func name
     'Order'     is an entity name
     
```

It will show sample command to remind you on how to use it. In this case, repository has 2 type of command.  

You can create another domain let say you want to create domain payment

```
$ gogen domain payment
```

Then you want to create a usecase under domain payment
```
$ gogen usecase RunPaymentCreate
```

But, hei the usecase is created under the order domain! How to create it under payment domain?

open the folder `.gogen` you will see file `domain`. You need to update the file from this

```
-order
payment
```

into this

```
order
-payment
```

Now you are working under payment domain. 

Re-run the command again
```
$ gogen usecase RunPaymentCreate
```

Now you got the usecase boilerplate code created under payment domain. (don't forget to delete your previous RunPaymentCreate usecase under the order domain manually)

---

