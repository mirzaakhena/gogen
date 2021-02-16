# Gogen (Clean Architecture Code Generator)
Helping generate your boiler plate and code structure based on clean architecure.

## Introduction
Have you ever wondered how to apply the clean architecture properly and how to manage the layout of your go project's folder structure? This tools will help you to make it.


## Sample Apps 
https://github.com/mirzaakhena/oms 
(The sample is not finished yet)


## Structure
This generator have basic structure like this
```
application/
controller/
domain/
gateway/
infrastructure/
shared/
usecase/
config.toml
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
You always start from creating an usecase, then you continue create the gateway, then create the controller, and the last step is bind those three (usecase + gateway + controller) in registry part. That's it. 

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

By organize this usecase in a such structure, we can easily change the *Controller* or the *Gateway* in very flexible way without worry to change the logic part. This is how the logic and infrastructure separation is working.

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

So you will create your first usecase. Let say the usecase name is  a `CreateOrder`. We will always create our usecase name with `PascalCase`. Now let's try our gogen code generator to create this usecase for us.
```
$ gogen usecase CreateOrder
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
$ gogen test CreateOrder
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

You also will get the global interceptor for all of controller.

## 6. Glue your controller, usecase, and gateway together

After generate the usecase, gateway and controller, we need to bind them all by calling this command.
```
$ gogen registry Default Restapi Production CreateOrder
```
Default is the registry name. You can name it whatever you want. After calling the command, some of those file generated will generated for you
