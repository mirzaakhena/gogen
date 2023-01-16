# Gogen Repository

Call `gogen repository SaveOrder Order RunOrderCreate` will

* Create an entity Order (if not exist)
  ```
  └── domain_yourdomainname
      └── model
          ├── entity
          │   └── order.go
          └── vo
              └── order_id.go
  ```

* Create repository SaveOrderRepo (if not exist)
  ```
  └── domain_yourdomainname
      └── model
          └── repository
              └── repository.go
  ```

* Inject code into Outport
  ```
  type Outport interface {
    repository.SaveOrderRepo
  }
  ```
  
* Inject code into Interactor. It will replace the `//!` flag
  ```
  func (r *runOrderCreateInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {
  
    res := &InportResponse{}
    
    // code your usecase definition here ...
    
    orderObj, err := entity.NewOrder(entity.OrderCreateRequest{})
    if err != nil {
      return nil, err
    }
    
    err = r.outport.SaveOrder(ctx, orderObj)
    if err != nil {
      return nil, err
    }
    
    //!
    
    return res, nil
  }
  ```


  
