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

`application.go` is the place you can instantiate the technology you want to use in your system. For example you want to setup gin, echo, message broker, database connection or anything.

`registry.go` is where you can bind your usecase, datasource and controller also you define the router here.

`runner.go` is the file for running the apps. This already have gracefully shutdown feature.


`schema.go` is the file for define your model to migrate into database.


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

`inport.go` is just an interface with only one method that will be a incoming gateway for your usecase. The method name is a gogen convention. You must not change the method name.

`outport.go` is an interface with many method that will be required only by your usecase. It must not shared to another usecase.

`interactor.go` is the core implementation of the usecase. It implement the method from inport and call the method from outport.

gogen only gives you the basic "template code". If you want to add/change your request response in your inport.go or outport.go, and want to have the "new method and the fields" appear in your interactor.go, then after you add the new field in the request or response struct (inport or outport), you can just simply delete the interactor.go, then call the `gogen usecase CreateOrder` again. It Will generate the new "helper code" for you.

Another variant of this usecase command is
```
gogen usecase query ShowOrder
```

For now we only have `command` and `query` usecase. You can try it to see the different.

## 3. Create your custom outport method

After you run the `gogen usecase` you need to delete the `port/outport.go` and `interactor.go` file. Then call this command
```
gogen outports CreateOrder CheckOrderID SaveOrder PublishOrder
```
Open the outport file you will found that there are 3 new methods defined.
After that run again the `gogen usecase command CreateOrder`

## 4. Create your usecase test file

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

## 5. Create datasource for your usecase

Datasource is the implemetor of your outport. You need to define where is your datasource. You can give any datasource name you like.
For the example you just want to hardcode any data for your apps. You can simply create the "hardcode" datasource version. Maybe you want to experiment with just simple database with SQLite for testing purpose, you can create the "testing" datasource version. For now, we will try to generate code for "production" datasource version.

```
gogen datasource production CreateOrder
```
This will generate
```
datasource/production/CreateOrder.go
```
You can start define what the data we must provide to run our usecase here.

## 6. Create controller for your usecase

In gogen, i define controller as any technology that will receive input from outside world. It can be rest api, grpc, consumer for event handling, or anything.

For now i only implement the gin framework version.

```
gogen controller restapi.gin CreateOrder
```
This will generate

```
controller/restapi/CreateOrder.go
```

## 7. Glue your usecase, datasource, and controller together

After generate the usecase, datasource and controller, we need to bind them all by calling this command:
```
gogen registry restapi production CreateOrder
```
Then open file `application/registry.go` then you will find this

```
package application

import (
	"your/golang/path/controller/restapi"
	"your/golang/path/datasource/production"
	"your/golang/path/usecase/createorder"
)

func (a *Application) RegisterUsecase() {
	createOrder(a)
  //code_injection function call
}

func createorder(a *Application) {
	outport := production.NewCreateOrderDatasource()
	inport := createorder.NewCreateOrderUsecase(outport)
	a.Router.POST("/createorder", restapi.CreateOrder(inport))
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

## 8. Create your model
This will simply create Order struct. That's it.
```
gogen model Order
```


Any other interesting idea?