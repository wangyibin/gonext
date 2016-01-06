package apidoc

import (
	"net/http"
	"reflect"

	"github.com/labstack/echo"
)

func propertiesOfEntity(bodyType reflect.Type) map[string]interface{} {
	properties := make(map[string]interface{})
	for i := 0; i < bodyType.NumField(); i++ {
		field := bodyType.Field(i)
		propertyName := lowCamelStr(field.Name)
		fieldType := field.Type.String()
		required := true
		if field.Type.Kind() == reflect.Ptr {
			required = false
			fieldType = field.Type.Elem().String()
		}

		properties[propertyName] = map[string]interface{}{
			"type":     fieldType,
			"required": required,
		}
	}
	return properties
}
func getdoc(router *Router) echo.HandlerFunc {

	return func(c *echo.Context) error {
		paths := make(map[string]interface{})
		definitions := make(map[string]interface{})
		for _, resource := range router.resources {
			paths[resource.path] = resource.toJSON()
			if resource.requestBody != nil {
				definitions[resource.requestBody.Name()] = map[string]interface{}{
					"type":       "object",
					"properties": propertiesOfEntity(resource.requestBody),
				}
			}
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"basePath": "/api",
			"host":     "localhost:3000",
			"swagger":  "2.0",
			"info": map[string]interface{}{
				"title":          "Swagger Sample App",
				"description":    "This is a sample server Petstore server.",
				"termsOfService": "http://swagger.io/terms/",
				"contact": map[string]string{
					"name":  "API Support",
					"url":   "http://www.swagger.io/support",
					"email": "support@swagger.io",
				},
				"license": map[string]string{
					"name": "Apache 2.0",
					"url":  "http://www.apache.org/licenses/LICENSE-2.0.html",
				},
				"version": "1.0.1",
			},
			"paths":       paths,
			"definitions": definitions,
			"tags": []map[string]string{
				map[string]string{
					"name":        "pet",
					"description": "Everything about your Pets",
				},
			},
		})
	}

}
