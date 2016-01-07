package apidoc

import (
	"fmt"
	"reflect"
	"strings"
)

// SwaggerDefinitions cahce
var SwaggerDefinitions = make(map[string]interface{})

// SwaggerPaths cache
var SwaggerPaths = make(map[string]interface{})

// SwaggerTags cache
var SwaggerTags = make(map[string]string)

// SwaggerPath struct
type SwaggerPath struct {
	Path string
	JSON map[string]interface{}
}

// SwaggerPathDefine struct
type SwaggerPathDefine struct {
	Tag         string
	Method      string
	Summary     string
	Description string
	Path        string
	Handler1    interface{}
	Handler2    interface{}
	Handler3    interface{}
}

// MountSwaggerPath func
func MountSwaggerPath(pathDefine *SwaggerPathDefine) error {
	fmt.Printf("%s  %s\n", pathDefine.Method, pathDefine.Path)
	newPath, err := BuildSwaggerPath(pathDefine)
	if err != nil {
		return err
	}
	if exist, ok := SwaggerPaths[newPath.Path]; !ok {
		SwaggerPaths[newPath.Path] = newPath.JSON
	} else {
		for k, v := range newPath.JSON {
			exist.(map[string]interface{})[k] = v
		}
	}
	return nil
}

// BuildSwaggerPath func
func BuildSwaggerPath(pathDefine *SwaggerPathDefine) (*SwaggerPath, error) {
	resultPath := pathDefine.Path
	for _, pname := range PathNames(pathDefine.Path) {
		resultPath = strings.Replace(resultPath, ":"+pname, "{"+pname+"}", -1)
	}

	inType, outType, err := validateChain(pathDefine.Handler1, pathDefine.Handler2, pathDefine.Handler3)

	if err != nil {
		return nil, err
	}
	operationID := getOperationID(inType, pathDefine.Handler1, pathDefine.Handler2, pathDefine.Handler3)

	successResponse := map[string]interface{}{
		"description": "successful operation",
	}
	if outType != nil {
		successResponse = map[string]interface{}{
			"description": "successful operation",
			"schema":      SwaggerEntitySchemaRef(outType),
		}
	}
	json := map[string]interface{}{
		strings.ToLower(pathDefine.Method): map[string]interface{}{
			"tags":        []string{pathDefine.Tag},
			"summary":     pathDefine.Summary,
			"description": pathDefine.Description,
			"produces":    []string{"application/json"},
			"consumes":    []string{"application/json"},
			"operationId": operationID,
			"parameters":  BuildRequestParam(pathDefine.Path, inType).ToSwaggerJSON(),
			"responses": map[string]interface{}{
				"200": successResponse,
				"500": map[string]interface{}{
					"description": "Interal Server Error",
				},
			},
		},
	}

	if inType != nil {
		actualInType := inType
		if inType.Kind() == reflect.Ptr {
			actualInType = inType.Elem()
		}
		for i := 0; i < actualInType.NumField(); i++ {
			typeField := actualInType.Field(i)
			if strings.ToUpper(typeField.Name) == "BODY" {
				MountSwaggerDefinition(typeField.Type)
			}
		}
	}
	if outType != nil {
		MountSwaggerDefinition(outType)
	}
	return &SwaggerPath{Path: resultPath, JSON: json}, nil
}

func propertiesOfEntity(bodyType reflect.Type) map[string]interface{} {
	properties := make(map[string]interface{})
	var requiredFields []string
	for i := 0; i < bodyType.NumField(); i++ {
		field := bodyType.Field(i)
		propertyName := lowCamelStr(field.Name)
		fieldType := field.Type
		if field.Type.Kind() == reflect.Ptr {
			fieldType = field.Type.Elem()
		} else {
			requiredFields = append(requiredFields, propertyName)
		}
		typ, format := GoTypeToSwaggerType(fieldType)

		description := field.Tag.Get("desc")

		switch typ {
		case "array":
			prefix := "type"
			if format[0] == '#' {
				prefix = "$ref"
			}
			properties[propertyName] = map[string]interface{}{
				"type":        "array",
				"description": description,
				"items": map[string]interface{}{
					prefix: format,
				},
			}
		case "object":
			properties[propertyName] = map[string]interface{}{
				"description": description,
				"$ref":        format,
			}
		default:
			properties[propertyName] = map[string]interface{}{
				"description": description,
				"type":        typ,
				"format":      format,
			}
		}
	}
	return map[string]interface{}{
		"type":       "object",
		"required":   requiredFields,
		"properties": properties,
	}
}

// MountSwaggerDefinition func
func MountSwaggerDefinition(typ reflect.Type) {
	entityType := typ
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}
	if entityType.Kind() == reflect.Map {
		return
	}

	if entityType.Kind() == reflect.Array || entityType.Kind() == reflect.Slice {
		MountSwaggerDefinition(entityType.Elem())
		return
	}
	if _, ok := SwaggerDefinitions[entityType.Name()]; !ok && entityType.Kind() == reflect.Struct {
		SwaggerDefinitions[entityType.Name()] = propertiesOfEntity(entityType)
		for i := 0; i < entityType.NumField(); i++ {
			field := entityType.Field(i)
			if field.Type.Kind() == reflect.Struct {
				MountSwaggerDefinition(field.Type)
			}
		}
	}
}

// SwaggerEntitySchemaRef used in parameter object and response object
func SwaggerEntitySchemaRef(inType reflect.Type) map[string]interface{} {
	entityType := inType
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}
	typ, format := GoTypeToSwaggerType(entityType)

	switch typ {
	case "array":
		prefix := "type"
		if format[0] == '#' {
			prefix = "$ref"
		}
		return map[string]interface{}{
			"type": "array",
			"items": map[string]interface{}{
				prefix: format,
			},
		}
	case "object":
		return map[string]interface{}{
			"$ref": format,
		}
	default:
		return map[string]interface{}{
			"type":   typ,
			"format": format,
		}
	}
}

// GoTypeToSwaggerType func
// http://swagger.io/specification/#parameterObject
func GoTypeToSwaggerType(typ reflect.Type) (string, string) {
	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "integer", "int32"
	case reflect.Int64, reflect.Uint64:
		return "integer", "int64"
	case reflect.String:
		return "string", "string"
	case reflect.Float32:
		return "number", "float"
	case reflect.Float64:
		return "number", "double"
	case reflect.Bool:
		return "boolean", "boolean"
	case reflect.Array, reflect.Slice:
		t, f := GoTypeToSwaggerType(typ.Elem())
		format := t
		if t == "object" {
			format = f
		}
		return "array", format
	case reflect.Struct:
		return "object", "#/definitions/" + typ.Name()
	default:
		return "string", "string"
	}
}
