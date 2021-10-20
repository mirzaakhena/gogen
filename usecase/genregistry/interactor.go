package genregistry

import (
  "context"
  "fmt"
  "github.com/mirzaakhena/gogen/domain/entity"
  "github.com/mirzaakhena/gogen/domain/service"
  "github.com/mirzaakhena/gogen/infrastructure/util"
)

//go:generate mockery --name Outport -output mocks/

type genRegistryInteractor struct {
  outport Outport
}

// NewUsecase is constructor for create default implementation of usecase GenRegistry
func NewUsecase(outputPort Outport) Inport {
  return &genRegistryInteractor{
    outport: outputPort,
  }
}

// Execute the usecase GenRegistry
func (r *genRegistryInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

  res := &InportResponse{}

  var objGateway *entity.ObjGateway
  var objController *entity.ObjController

  packagePath := r.outport.GetPackagePath(ctx)

  err := service.CreateEverythingExactly("default/", "infrastructure/server", map[string]string{}, struct {
    PackagePath string
  }{PackagePath: packagePath})
  if err != nil {
    return nil, err
  }

  driverName := "gin"

  tem := r.outport.GetServerFileTemplate(ctx, driverName)
  out := fmt.Sprintf("infrastructure/server/http_server_%s.go", driverName)
  _, err = r.outport.WriteFileIfNotExist(ctx, tem, out, struct{}{})
  if err != nil {
    return nil, err
  }

  // if controller name is not given, then we will do auto controller discovery strategy
  if req.ControllerName == "" {

    // look up the controller by foldername
    objControllers, err := entity.FindAllObjController()
    if err != nil {
      return nil, err
    }

    // if there is more than one controller
    if len(objControllers) > 1 {
      names := make([]string, 0)

      // collect all the controller name
      for _, g := range objControllers {
        names = append(names, g.ControllerName.String())
      }

      // return error
      return nil, fmt.Errorf("select one of this controller %v", names)
    }

    // currently, we are expecting only one gateway
    objController = objControllers[0]

  } else {

    var err error

    // when controller name is given
    objController, err = entity.FindControllerByName(req.ControllerName)
    if err != nil {
      return nil, err
    }

    // in case the controller name is not found
    if objController == nil {
      return nil, fmt.Errorf("no controller with name %s found", req.ControllerName)
    }

  }

  // if gateway name is not given, then we will do auto gateway discovery strategy
  if req.GatewayName == "" {

    // look up the gateway by foldername
    objGateways, err := entity.FindAllObjGateway()
    if err != nil {
      return nil, err
    }

    // if there is more than one gateway
    if len(objGateways) > 1 {
      names := make([]string, 0)

      // collect all the gateway name
      for _, g := range objGateways {
        names = append(names, g.GatewayName.String())
      }

      // return error
      return nil, fmt.Errorf("select one of this gateways %v", names)
    }

    // currently, we are expecting only one gateway
    objGateway = objGateways[0]

  } else {

    var err error

    // when gateway name is given
    objGateway, err = entity.FindGatewayByName(req.GatewayName)
    if err != nil {
      return nil, err
    }

    // in case the gateway name is not found
    if objGateway == nil {
      return nil, fmt.Errorf("no gateway with name %s found", req.GatewayName)
    }
  }

  usecaseNames, err := objController.FindAllUsecaseInportNameFromController()
  if err != nil {
    return nil, err
  }

  objRegistry, err := entity.NewObjRegistry(entity.ObjGatewayRequest{
    RegistryName:  req.RegistryName,
    ObjController: objController,
    ObjGateway:    objGateway,
    UsecaseNames:  usecaseNames,
  })
  if err != nil {
    return nil, err
  }

  if !util.IsFileExist(objRegistry.GetRegistryFileName()) {

    //err = service.CreateEverythingExactly("default/", "application", map[string]string{
    //  "registryname": objRegistry.RegistryName.LowerCase(),
    //}, objRegistry.GetData(packagePath))
    //if err != nil {
    //  return nil, err
    //}

    tem := r.outport.GetApplicationFileTemplate(ctx, driverName)
    out := fmt.Sprintf("application/registry/%s.go", objRegistry.RegistryName.LowerCase())
    _, err = r.outport.WriteFileIfNotExist(ctx, tem, out, objRegistry.GetData(packagePath))
    if err != nil {
      return nil, err
    }

    // reformat registry.go
    err = r.outport.Reformat(ctx, out, nil)
    if err != nil {
      return nil, err
    }

  } else {

    //err = objRegistry.InjectUsecaseInportField()
    //if err != nil {
    //  return nil, err
    //}

  }

   templateFile := r.outport.GetMainFileTemplate(ctx)
   _, err = r.outport.WriteFileIfNotExist(ctx, templateFile, "main.go", objRegistry.GetData(packagePath))
   if err != nil {
     return nil, err
   }

  return res, nil
}
