package spec_validate

import (
	"context"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/legacy"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Validator struct {
	spec map[string]*openapi3.T
}

func Load(service string, data []byte) (*Validator, error) {
	loader := &openapi3.Loader{Context: context.Background(), IsExternalRefsAllowed: true}
	spec, err := loader.LoadFromData(data)
	if err != nil {
		return nil, err
	}
	return &Validator{spec: map[string]*openapi3.T{service: spec}}, nil
}

func (v *Validator) Validate(req *http.Request, resp *http.Response) (err error) {
	var route *routers.Route
	var pathParams map[string]string
	for _, doc := range v.spec {
		var router routers.Router
		router, err = legacy.NewRouter(doc)
		if err != nil {
			return
		}
		route, pathParams, err = router.FindRoute(req)
		if err != nil {
			continue
		}
		break
	}

	if route == nil {
		fmt.Printf("can't find any route for request %s\n", req.URL.Path)
		collectValidateRecord("route.Path", strings.ToUpper("route.Method"), "miss route", -1)
		return
	}

	inputRequest := &openapi3filter.RequestValidationInput{
		Request:     req,
		PathParams:  pathParams,
		QueryParams: req.URL.Query(),
		Route:       route,
	}
	err = openapi3filter.ValidateRequest(context.Background(), inputRequest)
	if err != nil {
		collectValidateRecord(route.Path, strings.ToUpper(route.Method), fmt.Sprintf("validate request failed: %s", err.Error()), resp.StatusCode)
		return
	}

	if resp == nil {
		err = errors.New("response for validating should not be nil")
		collectValidateRecord(route.Path, strings.ToUpper(route.Method), "miss flow response", resp.StatusCode)
		return
	}

	inputResponse := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: inputRequest,
		Status:                 resp.StatusCode,
		Header:                 resp.Header,
		Body:                   resp.Body,
	}
	err = openapi3filter.ValidateResponse(context.Background(), inputResponse)
	if err != nil {
		collectValidateRecord(route.Path, strings.ToUpper(route.Method), fmt.Sprintf("validate response failed: %s", err.Error()), resp.StatusCode)
		return
	}
	collectValidateRecord(route.Path, strings.ToUpper(route.Method), "pass", resp.StatusCode)
	return
}

var ValidatePool = make(map[string]map[string]map[string]string)

func collectValidateRecord(path, method, cover string, statusCode int) {
	if ValidatePool[path] == nil {
		ValidatePool[path] = make(map[string]map[string]string)
	}
	if ValidatePool[path][method] == nil {
		ValidatePool[path][method] = make(map[string]string)
	}
	ValidatePool[path][method][strconv.Itoa(statusCode)] = cover
}

func (v *Validator) Report() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"path", "method", "status", "cover"})
	fmt.Println("Contract Report")
	pass := true
	for _, doc := range v.spec {
		for path, item := range doc.Paths {
			for method, opt := range item.Operations() {
				for statusCode, _ := range opt.Responses {
					if _, ok := ValidatePool[path]; ok {
						if _, ok := ValidatePool[path][method]; ok {
							if cover, ok := ValidatePool[path][method][statusCode]; ok {
								if cover != "pass" {
									pass = false
								}
								table.Append([]string{path, method, statusCode, cover})
							} else {
								pass = false
								table.Append([]string{path, method, statusCode, "miss status"})
							}
						} else {
							pass = false
							table.Append([]string{path, method, statusCode, "miss method"})
						}
					} else {
						pass = false
						table.Append([]string{path, method, statusCode, "miss path"})
					}
				}
			}
		}
	}
	table.Render()
	if !pass {
		os.Exit(2)
	}
}
