

Usecase
    create inport
    create outport
    create interactor

Entity
    create an entity

Method
    read entity
    create an entity's method

ValueObject
    create a valueobject

ValueString
    create a valuestring

Enum
    create an enum
    create error

Error
    create error

Repository
    creating repo interface
    creating entity
    inject code into outport
    inject code into interactor
    create error

Service
    create service interface
    inject code to outport
    inject code to interactor

Controller
    read inport
    create controller

Gateway
    read all method interface in outport
    read all import
    if it is extend external interface then read this external interface
    if gateway never exist, create gateway with implementation
    if gateway already created before, then just add unimplemented method only

Registry
    read controller
    read usecase
    read gateway
    bind all together

Test
    read inport
    read outport
    read interactor

Init

Config

Template


semua append code harus lewat ast

https://betterprogramming.pub/rpc-in-golang-19661033942
https://medium.com/rungo/building-rpc-remote-procedure-call-network-in-go-5bfebe90f7e9



reconstruct method

