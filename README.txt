

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

Gateway
    read all method interface in outport
    read all import
    if it is extend external interface then read this external interface
    if gateway never exist, create gateway with implementation
    if gateway already created before, then just add unimplemented method only

Controller
    read inport
    create controller

Registry
    read controller
    read usecase
    read gateway
    bind all together

Test
    read inport
    read outport
    read interactor

Mapper

Init

Config

Template


https://betterprogramming.pub/rpc-in-golang-19661033942
https://medium.com/rungo/building-rpc-remote-procedure-call-network-in-go-5bfebe90f7e9

kalo implementasi gateway ada di struct lain? kita harus bisa handle dan tidak perlu di extend lagi

Controller dikategorikan berdasarkan actor yang mengakses usecase
    user api
    backoffice api
    dev api
    webhook api

<ControllerName>Controller
    RegisterHandler()
    handle<UsecaseName>(inport)


Read Gateway
    Kita cari outport interface dari usecase yang di inginkan
    baca seluruh field interface nya
    ada 4 kemungkinan dari field tsb
    1. berupa extension interface lain yang berada di package yang berbeda
    2. berupa extension interface lain yang berada di dalam package dan file yang sama : BELUM DIHANDLE
    3. berupa extension interface lain yang berada di dalam package namun file yang berbeda : BELUM DIHANDLE
    4. berupa direct method
        jika parameter punya type yang berada di packgae yg sama : BELUM DIHANDLE

    tujuannya adalah mengkoleksi semua method yang dibutuhkan




