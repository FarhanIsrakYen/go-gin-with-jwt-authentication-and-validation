package routes

import (
	"fmt"
	"go-gin-with-jwt-authentication-and-validation/controllers"
	"io/ioutil"
	"log"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type Route struct {
    Path       string `yaml:"path"`
    Method     string `yaml:"method"`
    Controller string `yaml:"controller"`
}

type RoutesConfig struct {
    Routes []Route `yaml:"routes"`
}

func SetupRoutes(r *gin.Engine) error {
    routesConfig, err := loadRoutes("routes/routes.yaml")
    if err != nil {
        return err
    }

    for _, route := range routesConfig.Routes {
		fmt.Println(route.Path)
        controllerFunc := getControllerFunction(route.Controller)
        if controllerFunc == nil {
            log.Printf("Controller function %s not found", route.Controller)
            continue
        }

        switch route.Method {
        case "GET":
            r.GET(route.Path, controllerFunc)
        case "POST":
            r.POST(route.Path, controllerFunc)
        case "PUT":
            r.PUT(route.Path, controllerFunc)
        case "DELETE":
            r.DELETE(route.Path, controllerFunc)
        default:
            log.Printf("Unsupported method %s for path %s", route.Method, route.Path)
        }
    }

    return nil
}

func loadRoutes(filePath string) (*RoutesConfig, error) {
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, err
    }

    var routesConfig RoutesConfig
    err = yaml.Unmarshal(data, &routesConfig)
    if err != nil {
        return nil, err
    }

    return &routesConfig, nil
}

func getControllerFunction(controller string) gin.HandlerFunc {
    parts := strings.Split(controller, ".")
    if len(parts) != 2 {
        log.Printf("Invalid controller format: %s", controller)
        return nil
    }

    packageName := parts[0]
    methodName := parts[1]

    var ctrlInstance reflect.Value
    switch packageName {
    case "AuthController":
        ctrlInstance = reflect.ValueOf(controllers.AuthController{})
    case "UserController":
        ctrlInstance = reflect.ValueOf(controllers.UserController{})
    default:
        log.Printf("Unknown controller package: %s", packageName)
        return nil
    }

    method := ctrlInstance.MethodByName(methodName)
    if !method.IsValid() {
        log.Printf("Method %s not found in controller %s", methodName, packageName)
        return nil
    }

    return func(c *gin.Context) {
        method.Call([]reflect.Value{reflect.ValueOf(c)})
    }
}
