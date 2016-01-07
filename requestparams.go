package apidoc

import (
	"fmt"
	"reflect"
	"strings"
)

// RequestParam struct
type RequestParam struct {
	PathParams  []Param
	QueryParams []Param
	RequestBody reflect.Type
}

// Param struct
type Param struct {
	Name        string
	Type        reflect.Type
	Description string
	Required    bool
}

func (p *Param) String() string {
	return fmt.Sprintf("name[%s] type:%s, required:%t\n", p.Name, p.Type, p.Required)
}

// ToSwaggerJSON func
func (p *Param) ToSwaggerJSON(position string) map[string]interface{} {
	typ, format := GoTypeToSwaggerType(p.Type)
	return map[string]interface{}{
		"name":        p.Name,
		"in":          position,
		"format":      format,
		"required":    p.Required,
		"type":        typ,
		"description": p.Description,
	}
}

// BuildRequestParam func
func BuildRequestParam(path string, inType reflect.Type) *RequestParam {
	requestType := inType
	if requestType.Kind() == reflect.Ptr {
		requestType = requestType.Elem()
	}
	pnames := PathNames(path)
	var pathParams []Param
	var queryParams []Param
	var requestBody reflect.Type
	fmt.Printf("%s\n\trequest type %v\n", path, requestType)
	for i := 0; i < requestType.NumField(); i++ {
		typeField := requestType.Field(i)

		if strings.ToUpper(typeField.Name) != "BODY" {
			param := Param{Name: typeField.Name, Type: typeField.Type, Required: typeField.Type.Kind() != reflect.Ptr}
			if containsIgnoreCase(pnames, typeField.Name) {
				pathParams = append(pathParams, param)
				fmt.Printf("\tPath Params %s", param.String())
			} else {
				queryParams = append(queryParams, param)
				fmt.Printf("\tQuery Params %s", param.String())
			}
		} else {
			if typeField.Type.Kind() == reflect.Ptr {
				requestBody = typeField.Type.Elem()
			} else {
				requestBody = typeField.Type
			}
			fmt.Printf("\tRequestBody[%s]\n", requestBody)
		}
	}
	return &RequestParam{PathParams: pathParams, QueryParams: queryParams, RequestBody: requestBody}
}

// ToSwaggerJSON func
func (req *RequestParam) ToSwaggerJSON() []map[string]interface{} {
	var parameters []map[string]interface{}
	for _, pathParam := range req.PathParams {
		parameters = append(parameters, pathParam.ToSwaggerJSON("path"))
	}
	for _, queryParam := range req.QueryParams {
		parameters = append(parameters, queryParam.ToSwaggerJSON("query"))
	}
	if req.RequestBody != nil {
		parameters = append(parameters, map[string]interface{}{
			"in":       "body",
			"name":     "body",
			"required": true,
			"schema": map[string]string{
				"$ref": "#/definitions/" + req.RequestBody.Name(),
			},
		})
	}
	return parameters
}
