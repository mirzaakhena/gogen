# Gogen

Helping generate your boiler structure and boiler plate code


## Step by step to working with gogen

### 1. Create basic project structure
```
gogen init .
```
it will create the basic project structure
```
.application_schema/usecases/
appliation/runner.go
appliation/setup.go
appliation/wiring.go
controllers/
datasources/mocks
entities/
usecases/<usecase_name>/inport/
usecases/<usecase_name>/interactor/
usecases/<usecase_name>/outport/
repositories/
services/
utils/
config.toml
main.go
README.md
```

### 2. Create your basic usecase structure
```
gogen usecase CreateOrder
```
When you run this command in the first time, you will have those files generated for you
- .application_schema/usecases/CreateOrder.yml
- inport/CreateOrder.go
- interactor/CreateOrder.go
- interactor/CreateOrder_test.go
- outport/CreateOrder.go
- datasources/mocks/CreateOrder.go

You can start run test your code by running the interactor/CreateOrder_test.go test file

### 3. Define your usecase input port and output port
Open file inport/CreateOrder.go you will find empty request and empty response struct

Open file outport/CreateOrder.go you also will find empty request and empty response struct

Open file .application_schema/usecases/CreateOrder.yml 

update the inport (input port) and outport (output port) 

run the command again
```
gogen usecase CreateOrder
```

Then you will get inport/CreateOrder.go, outport/CreateOrder.go and datasources/mocks/CreateOrder.go is updated

But this command will not change the interactor/CreateOrder.go and interactor/CreateOrder_test.go to protect your current logic adn test code. So you need to update it manually.


### 4. Create datasource
```
gogen datasource production CreateOrder
```
production is sample datasource. It can be anything you like. It will create one file under one folder
- datasources/production/CreateOrder.go



### 5. Create controller
```
gogen controller restapi gin CreateOrder
```
- controllers/restapi/CreateOrder.go