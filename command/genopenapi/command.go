package genopenapi

import (
	"encoding/json"
	"fmt"
	"github.com/mirzaakhena/gogen/utils"
	"io/ioutil"
)

// ObjTemplate ...
type ObjTemplate struct {
	GomodPath     string
	DefaultDomain string
}

type OASSchema struct {
	Ref        string   `json:"$ref,omitempty"`
	Type       string   `json:"type,omitempty"`
	Properties []any    `json:"properties,omitempty"`
	Required   bool     `json:"required,omitempty"`
	Default    string   `json:"default,omitempty"`
	Enum       []string `json:"enum"`
}

type OASParameter struct {
	Name        string                `json:"name"`
	In          string                `json:"in"`
	Description string                `json:"description"`
	Required    bool                  `json:"required"`
	Schema      OASSchema             `json:"schema"`
	Example     any                   `json:"example,omitempty"`
	Examples    map[string]OASExample `json:"examples,omitempty"`
}

type OASExample struct {
	Description   string `json:"description"`
	Summary       string `json:"summary"`
	Value         any    `json:"value"`
	ExternalValue string `json:"externalValue"`
}

type OASMediaType struct {
	Schema   OASSchema             `json:"schema,omitempty"`
	Example  any                   `json:"example,omitempty"`
	Examples map[string]OASExample `json:"examples,omitempty"`
}

type OASRequestBody struct {
	Description string                  `json:"description,omitempty"`
	Content     map[string]OASMediaType `json:"content"`
	Required    bool                    `json:"required,omitempty"`
}

type OASResponse struct {
	Description string                  `json:"description,omitempty"`
	Content     map[string]OASMediaType `json:"content"`
}

type OASOperationObject struct {
	Tags        []string               `json:"tags"`
	Summary     string                 `json:"summary,omitempty"`
	OperationID string                 `json:"operationId"`
	Parameters  []OASParameter         `json:"parameters"`
	RequestBody OASRequestBody         `json:"requestBody"`
	Responses   map[string]OASResponse `json:"responses"`
}

type OASPathItem struct {
	Get         OASOperationObject `json:"get,omitempty"`
	Post        OASOperationObject `json:"post,omitempty"`
	Put         OASOperationObject `json:"put,omitempty"`
	Delete      OASOperationObject `json:"delete,omitempty"`
	Summary     string             `json:"summary,omitempty"`
	Description string             `json:"description,omitempty"`
	Ref         string             `json:"$ref"`
	Servers     OASServer          `json:"servers,omitempty"`
	Parameters  OASParameter       `json:"parameters,omitempty"`
}

type OASExternalDocumentation struct {
	Description string `json:"description"`
	URL         string `json:"url"`
}

type OASServerVariable struct {
	Description string   `json:"description,omitempty"`
	Default     string   `json:"default"`
	Enum        []string `json:"enum,omitempty"`
}

type OASServer struct {
	Description string                       `json:"description"`
	Url         string                       `json:"url"`
	Variables   map[string]OASServerVariable `json:"variables"`
}

type OASLicense struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type OASContact struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Url   string `json:"url"`
}

type OASInfo struct {
	Title          string     `json:"title"`
	Version        string     `json:"version"`
	Description    string     `json:"description"`
	TermsOfService string     `json:"termsOfService"`
	Contact        OASContact `json:"contact"`
	License        OASLicense `json:"license"`
}

type OASSecurityScheme struct {
}

type OASComponent struct {
	SecuritySchemes map[string]OASSecurityScheme `json:"securitySchemes"`
	Parameters      map[string]OASParameter      `json:"parameters,omitempty"`
	Schemes         map[string]OASSchema         `json:"schemes"`
}

type OASRoot struct {
	OpenAPI      string                   `json:"openapi"`
	Info         OASInfo                  `json:"info"`
	Servers      []OASServer              `json:"servers"`
	ExternalDocs OASExternalDocumentation `json:"externalDocs"`
	Paths        map[string]OASPathItem   `json:"paths"`
	Components   map[string]OASComponent  `json:"components"`
}

func Run(inputs ...string) error {

	domainName := utils.GetDefaultDomain()

	data := OASRoot{
		OpenAPI: "3.0.3",
		Info: OASInfo{
			Title:          "Application Restful API",
			Version:        "1",
			Description:    "Application Restful API",
			TermsOfService: "https://gogen.com/tnc",
			Contact: OASContact{
				Name:  "Mirza Akhena",
				Email: "mirza.akhena@gmail.com",
				Url:   "mirzaakhena.com",
			},
			License: OASLicense{
				Name: "APACHE 2.0",
				Url:  "https://www.apache.org/licenses/LICENSE-2.0",
			},
		},
		Servers: []OASServer{
			{
				Description: fmt.Sprintf("%s", domainName),
				Url:         "{protocol}://{environment}:{port}/api/{version}",
				Variables: map[string]OASServerVariable{
					"protocol": {
						Description: "Protocol",
						Default:     "http",
						Enum:        []string{"http", "https"},
					},
					"environment": {
						Description: "Environment",
						Default:     "localhost",
						Enum:        []string{"localhost", "dev", "qa", "prod"},
					},
					"port": {
						Description: "Port",
						Default:     "8080",
						Enum:        []string{"8080", "8081", "8082", "80", "443"},
					},
					"version": {
						Description: "Version",
						Default:     "v1",
						Enum:        nil,
					},
				},
			},
		},
		ExternalDocs: OASExternalDocumentation{
			Description: fmt.Sprintf("Documentation for %s", domainName),
			URL:         fmt.Sprintf("https://bitbucket.org/%s", domainName),
		},
		Paths:      map[string]OASPathItem{},
		Components: map[string]OASComponent{},
	}
	data.Paths["/products"] = OASPathItem{
		Get: OASOperationObject{
			Tags:        []string{"product"},
			Summary:     "Query All Products",
			OperationID: "GetAllProduct",
			Parameters: []OASParameter{
				{
					Name:        "id",
					In:          "query",
					Description: "product id",
					Required:    false,
					Schema: OASSchema{
						Type: "string",
					},
					Example: "PRD1234001",
				},
			},
			Responses: map[string]OASResponse{
				"200": {
					Description: "",
					Content: map[string]OASMediaType{
						"application/json": {
							Schema: OASSchema{
								Type: "string",
							},
						},
					},
				},
			},
		},
		Post: OASOperationObject{
			Tags:        []string{"product"},
			Summary:     "Insert New Product",
			OperationID: "RunProductCreate",
			Parameters:  nil,
			RequestBody: OASRequestBody{
				Description: "",
				Content:     nil,
				Required:    false,
			},
			Responses: nil,
		},
		Put:         OASOperationObject{},
		Delete:      OASOperationObject{},
		Summary:     "",
		Description: "",
		Ref:         "",
		Servers:     OASServer{},
		Parameters:  OASParameter{},
	}

	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("domain_%s/openapi.json", domainName), file, 0644)
	if err != nil {
		return err
	}

	return nil
}
