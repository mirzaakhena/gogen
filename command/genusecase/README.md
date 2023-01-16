# Gogen Usecase

Call `gogen usecase RunProductCreate` will

* Create file Inport, Interactor and Outport
 
  ```
  └── domain_yourdomainname
      └── usecase
          └── runordercreate
              ├── README.md
              ├── inport.go
              ├── interactor.go
              └── outport.go
  ```

* There are specialization for Inport. 

  if using prefix `GetAll` then InportRequest and InportResponse will have
  ```go
  type InportRequest struct {
    Page int
    Size int
  }
  
  type InportResponse struct {
    Count int64
    Items []any
  } 
  ```
  otherwise it will be empty
  ```go
  type InportRequest struct {

  }
  
  type InportResponse struct {

  } 
  ```
  



