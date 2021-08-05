# Gogen (Clean Architecture Code Generator)
Provide code structure based on clean architecure and domain driven design

## Introduction
Gogen is a tools for generate code structure and boiler plate code. There are many commands for generate an usecase, repository, service, entity and etc.


## Sample Apps
https://github.com/mirzaakhena/userprofile

## Structure
This generator has basic structure like this
```
application/apperror
application/registry

controller/

domain/entity
domain/repository
domain/service
domain/vo

gateway/

infrastructure/log

usecase/

main.go
```

## Clean Architecture Concept
The main purpose of this architecture is :
* Separation concern between INFRASTRUCTURE part and LOGIC Part
* Independent of Framework. Free to swap to any framework
* Independent of UI. Free to swap UI. For ex: from Web UI to Console UI
* Independent of Database. Free to swap to any database (data storage)
* Independent of any external agency. Business rule doesn’t know anything about outside world
* Testable. The business rules can be tested without the UI, Database, Web Server, or any other external element


## How to use the gogen?
You always start from creating an usecase, then you continue to create the gateway, then create the controller, and the last step is bind those three (usecase + gateway + controller) in registry part. That's it.

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

*Inport* is an interface that has only one method (named `Execute`) that will be called by *Controller*. The method in a interface define all the required (request and response) parameter to run the specific usecase. *Inport* will implemented by *Interactor*. Request and response struct is allowed to share to *Outport* under the same usecase, but must not shared to other *Inport* or *Outport* usecase.

*Interactor* is the place that you can define your logic flow. When an usecase logic need a data, it will ask the *Outport* to provide it. *Interactor* also use *Outport* to send data, store data or do some action to other service. *Interactor* only have one *Outport* field. We must not adding new *Outport* field to *Interactor* to keep a simplicity and consistency.

*Outport* is a data and action provider for *Interactor*. *Outport* never know how it is implemented. The implementor (in this case a *Gateway*) will decide how to provide a data or do an action. This *Outport* is very exclusive for specific usecase (in this case *Interactor*) and must not shared to other usecase. By having exclusive *Outport* it will isolate the testing for usecase.

By organize this usecase in a such structure, we can easily change the *Controller*, or the *Gateway* in very flexible way without worry to change the logic part. This is how the logic and infrastructure separation is working.

## Comparison with three layer architecture (Controller -> Service -> Repository) pattern
* *Controller* is the same controller for both architecture.
* *Service* is similar like *Interactor* with additional strict rule. *Service* allowed to have many repositories. In Clean Architecture, *Interactor* only have one *Outport*.
* *Service* have many method grouped by the domain. In Clean Architecture, we focus per usecase. One usecase for One Class to achieve *Single Responsibility Principle*.
* In *Repository* you often see the CRUD pattern. Every developer can added new method if they think they need it. In reality this *Repository* is shared to different *Service* that may not use that method. In *Outport* you will strictly to adding method that guarantee used. Even adding new method or updating existing method will not interfere another usecase.
* *Repository* is an *Outport* with *Gateway* as its implementation.

## Gogen Convention
* As interface, inport has one and only one method which handle one usecase
* Interactor is a manager for entity and outport
* Interactor must not decide to have any technology. All technology is provided by Outport.
* Outport is a servant for an interactor. It will provide any data from external source
* As interface, Outport at least have one method and can have multiple method
* All method in Outport is guarantee used by interactor
* Inport and Outport must not shared to other usecase. They are exclusive for specific usecase
* Request or Response DTO from Inport may shared to Outport in the same Usecase (Interactor)
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
* Registry name can be an application name. You can utilize it if you want to develop microservice with mono repo


## Download it
```
$ go get github.com/mirzaakhena/gogen
```
Install it into your local system (make sure you are in gogen directory)
```
$ go install
```

## Step by step to working with gogen

## Create your basic usecase structure

So you will create your first usecase. Let say the usecase name is  a `CreateOrder`. We will always create our usecase name with `PascalCase`. Now let's try our gogen code generator to create this usecase for us.
```
$ gogen usecase CreateOrder
```

Usecase name will be used as a package name under usecase folder by lowercasing the usecase name.

`usecase/createorder/inport.go` is an interface with one method that will implement by your usecase. The standart method name is a `Execute`.

`usecase/createorder/outport.go` is an interface which has many methods that will be used by your usecase. It must not shared to another usecase.

`usecase/createorder/interactor.go` is the core implementation of the usecase (handle your bussiness application). It implements the method from inport and call the method from outport.

## Create your usecase test file
```
$ gogen test normal CreateOrder
```
normal is the test name and CreateOrder is the usecase name

## Create a repository 
```
$ gogen repository SaveOrder Order CreateOrder
```

## Create a service
```
$ gogen service PublishMessage CreateOrder
```

## Create a gateway for your usecase

Gateway is the struct to implement your outport interface. You need to set a name for your gateway. 
In this example we will set name : inmemory
```
$ gogen gateway inmemory CreateOrder
```

## Create a controller for your usecase

In gogen, we define a controller as trigger of the usecase. It can be rest api, grpc, consumer for event handling, or any other source input. By default, it only uses net/http restapi. 
Call this command for create a controller. Restapi is your controller name. You can name it whatever you want.
```
$ gogen controller restapi CreateOrder
```

You also will get the global interceptor for all controllers.

## Glue your controller, usecase, and gateway together

After generate the usecase, gateway and controller, we need to bind them all by calling this command.
```
$ gogen registry appone restapi inmemory CreateOrder
```
appone is the registry name. You can name it whatever you want. After calling the command, some of those file generated will generate for you in `application/registry`

## Create entity
entity is a mutable object that has an identifier. This command will create new entity struct under `domain/entity/` folder
```
$ gogen entity Order
```

## Create valueobject
valueobject is an immutable object that has no identifier. This command will create new valueobject under `domain/vo/` folder
```
$ gogen valueobject FullName FirstName LastName 
```

## Create valuestring
valuestring is a valueobject simple string type. This command will create new valuestring struct under `domain/vo/` folder
```
$ gogen valuestring OrderID
```

## Create enum
enum is a single immutable value. This command will create new enum struct under `domain/vo/` folder
```
$ gogen enum PaymentMethod DANA Gopay Ovo LinkAja
```

## Create error enum
error enum is a shared error collection. This command will add new error enum line in `application/apperror/error_enum.go` file
```
$ gogen error SomethingGoesWrongError
```