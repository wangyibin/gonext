
```go
func main() {
  // Echo instance
	e := echo.New()
  // do some init

  // create router for apidoc
  router := apidoc.NewRouter(e)

  // first clone https://github.com/swagger-api/swagger-ui.git
  // copy dist folder to your project dir /public/swagger
  // change doc url in public/swagger/index.html
  // publish swagger docs
	e.Static("/swagger", "public/swagger")

  controllers.bindPets(router)

  // bind the error handler for validation error
  e.SetHTTPErrorHandler(func(err error, c *echo.Context) {
		e.Logger().Error(err)

		if c.Response().Committed() {
			return
		}

		httpStatus := http.StatusInternalServerError
		code := "000"
		var msg interface{} = http.StatusText(httpStatus)

		if validateError, ok := err.(validator.ValidationErrors); ok {
			httpStatus = http.StatusBadRequest
			code = "900"
			msg = validateError
		} else if he, ok := err.(*echo.HTTPError); ok {
			httpStatus = he.Code()
			msg = he.Error()
		} else {
			msg = err.Error()
		}
		c.JSON(httpStatus, map[string]interface{}{
			"code":  code,
			"error": msg,
		})
	})
}
```

```go
package controllers

import (
	"fmt"
	"go-apidoc"

	"github.com/labstack/echo"
)

type (
	getPetsRequest struct {
		Name  string
		Age   int
		Color *int // use pointer for field which is not required
	}

	addPetRequest struct {
		Body *pet
	}

	getPetRequest struct {
		ID string
	}
	// we use https://github.com/go-playground/validator for struct validation
	// the validate method will be called automatically
	pet struct {
		Name string `desc:"名称" validate:"required,gte=3"`
	}

	accessToken struct {
		Token string
	}

	petAddSuccess struct {
		Name  string
		Token string
	}
)

// c *echo.Context， is optional
func getPets(req *getPetsRequest, c *echo.Context) ([]pet, error) {
	if req.Color == nil {
		fmt.Println("req.Color is Nil")
	} else {
		fmt.Printf("req.Color is %d\n", *req.Color)
	}
	return []pet{pet{Name: "p1"}, pet{Name: "p2"}}, nil
}

func getPet(req *getPetRequest) (*pet, error) {
	return &pet{Name: "p1"}, nil
}

// do authenticate here, and pass the access token to next function
func auth(c *echo.Context) *accessToken {
	return &accessToken{Token: "sss"}
}

// get access token by previous function
func addPet(req *addPetRequest, token *accessToken) *petAddSuccess {
	return &petAddSuccess{Name: "hello " + req.Body.Name, Token: token.Token}
}

func bindPets(r *apidoc.Router) {
	// create group
	// first params: is the tag, this is required. And we use it as swagger tag
	// sencond param: is the description for the group
	// the third param: is the path
	g := r.Group("pet", "pet desc", "/api/pets")
	// you can add two string parameter in the parameter list
	// the first one will be recognized as summary for the handler
	// the second one will be the description
	g.Get("/", getPets, "Get pets list", "write some implementation notes here")
	// call the handler one by one, if one of the handler return a error
	// the following handler will not be execute
	g.Post("/", auth, addPet)
	g.Get("/:id", getPet, "Get pet by ID", "bala bala bala....")
}

```

change url in `public/swagger/index.html`
```javascript
// line 35 to line 39
if (url && url.length > 1) {
  url = decodeURIComponent(url[1]);
} else {
  url = "/api-docs";//"http://petstore.swagger.io/v2/swagger.json";
}
```
