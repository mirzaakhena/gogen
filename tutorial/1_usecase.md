# Usecase

Usecase follow the usecase definition in Uncle bob Architecture. 
Basically Usecase is a simple interface with one method to do only one more thing (to support the single responsibility principle). 
The basic contract form of usecase is in `shared/usecase/usecase.go`

```go
type Inport[REQUEST, RESPONSE any] interface {
	Execute(ctx context.Context, req REQUEST) (*RESPONSE, error)
}
```

When you create an usecase, let say the usecase name is RunOrderSubmit with this command

```shell
$ gogen usecas RunOrderSubmit
```

Then you will have this three file `inport.go`, `interactor.go` and `outport.go` under `domain_yourdomain/usecase/runordersubmit` folder

### inport.go

Inport define the InportRequest and InportResponse. 
InportRequest define all the required argument to run the usecase and 
InportResponse define all the result value after run the usecase. 

```go
type Inport usecase.Inport[context.Context, InportRequest, InportResponse]

type InportRequest struct {
}

type InportResponse struct {
}
```

### interactor.go

Interactor have a `struct` and the `struct method` that implement the Inport interface.
In the struct we have one (and only one) outport field with type interface

```go
type runOrderSubmitInteractor struct {
	outport Outport
}
```

Interactor basically does not hold any state when running. It is recommended not to add other field in the struct except the outport field.

We also have the 'Constructor' function that create an interactor object. 

```go
func NewUsecase(outputPort Outport) Inport {
	return &runOrderSubmitInteractor{
		outport: outputPort,
	}
}
```

And the struct method. In this method, developer can write the needed logic or algorithm 

```go
func (r *runOrderSubmitInteractor) Execute(ctx workflow.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}

	// code your usecase definition here ...
	//!

	return res, nil
}
```

### outport.go

The last part is Outport. Outport is an interface which has a several required method that will be called by an interactor method. By default this interface is empty.

```go
type Outport interface {
}
```

