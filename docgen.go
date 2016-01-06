package apidoc

import (
	"net/http"
	"reflect"

	"github.com/labstack/echo"
)

func propertiesOfEntity(bodyType reflect.Type) map[string]interface{} {
	properties := make(map[string]interface{})
	for i := 0; i < bodyType.NumField(); i++ {
		typeField := bodyType.Field(i)

		properties[typeField.Name] = map[string]interface{}{
			"type": typeField.Type.String(),
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
		})
	}

}
