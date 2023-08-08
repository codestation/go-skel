// Package oapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.13.4 DO NOT EDIT.
package oapi

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /auth/login)
	Login(ctx echo.Context) error
	// Create a new delay job request
	// (POST /background/delay)
	ProcessBackground(ctx echo.Context) error
	// Check if the app is started
	// (GET /health/live)
	LiveCheck(ctx echo.Context, params LiveCheckParams) error
	// Check if the app is ready to accept connections
	// (GET /health/ready)
	ReadyCheck(ctx echo.Context, params ReadyCheckParams) error
	// Retrieve a list of profiles
	// (GET /profiles)
	ListProfiles(ctx echo.Context, params ListProfilesParams) error
	// Create a new profile
	// (POST /profiles)
	SaveProfile(ctx echo.Context) error
	// Delete a profile by ID
	// (DELETE /profiles/{id})
	RemoveProfile(ctx echo.Context, id ProfileId) error
	// Get a profile by ID
	// (GET /profiles/{id})
	GetProfile(ctx echo.Context, id ProfileId) error
	// Update a profile by ID
	// (PATCH /profiles/{id})
	UpdateProfile(ctx echo.Context, id ProfileId) error

	// (GET /queues/{name}/tasks/{id})
	GetTask(ctx echo.Context, name QueueName, id TaskId) error

	// (GET /queues/{name}/tasks/{id}/response)
	GetTaskResponse(ctx echo.Context, name QueueName, id TaskId) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// Login converts echo context to params.
func (w *ServerInterfaceWrapper) Login(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.Login(ctx)
	return err
}

// ProcessBackground converts echo context to params.
func (w *ServerInterfaceWrapper) ProcessBackground(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ProcessBackground(ctx)
	return err
}

// LiveCheck converts echo context to params.
func (w *ServerInterfaceWrapper) LiveCheck(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params LiveCheckParams

	// ------------- Optional query parameter "verbose" -------------

	err = runtime.BindQueryParameter("form", true, false, "verbose", ctx.QueryParams(), &params.Verbose)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter verbose: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.LiveCheck(ctx, params)
	return err
}

// ReadyCheck converts echo context to params.
func (w *ServerInterfaceWrapper) ReadyCheck(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params ReadyCheckParams

	// ------------- Optional query parameter "verbose" -------------

	err = runtime.BindQueryParameter("form", true, false, "verbose", ctx.QueryParams(), &params.Verbose)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter verbose: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ReadyCheck(ctx, params)
	return err
}

// ListProfiles converts echo context to params.
func (w *ServerInterfaceWrapper) ListProfiles(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params ListProfilesParams

	paramsMap := map[string]bool{
		"before":   true,
		"after":    true,
		"page":     true,
		"q":        true,
		"limit":    true,
		"includes": true,
		"filters":  true,
		"fields":   true,
		"sort":     true,
	}

	// ------------- Optional query parameter "before" -------------

	err = runtime.BindQueryParameter("form", true, false, "before", ctx.QueryParams(), &params.Before)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter before: %s", err))
	}

	// ------------- Optional query parameter "after" -------------

	err = runtime.BindQueryParameter("form", true, false, "after", ctx.QueryParams(), &params.After)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter after: %s", err))
	}

	// ------------- Optional query parameter "page" -------------

	err = runtime.BindQueryParameter("form", true, false, "page", ctx.QueryParams(), &params.Page)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter page: %s", err))
	}

	// ------------- Optional query parameter "q" -------------

	err = runtime.BindQueryParameter("form", true, false, "q", ctx.QueryParams(), &params.Q)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter q: %s", err))
	}

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// ------------- Optional query parameter "includes" -------------

	err = runtime.BindQueryParameter("form", true, false, "includes", ctx.QueryParams(), &params.Includes)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter includes: %s", err))
	}

	// ------------- Optional query parameter "filters" -------------

	params.Filters = &Filters{}
	for key, values := range ctx.QueryParams() {
		if !paramsMap[key] {
			(*params.Filters)[key] = values[0]
		}
	}

	// ------------- Optional query parameter "fields" -------------

	err = runtime.BindQueryParameter("form", true, false, "fields", ctx.QueryParams(), &params.Fields)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter fields: %s", err))
	}

	// ------------- Optional query parameter "sort" -------------

	err = runtime.BindQueryParameter("form", true, false, "sort", ctx.QueryParams(), &params.Sort)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter sort: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListProfiles(ctx, params)
	return err
}

// SaveProfile converts echo context to params.
func (w *ServerInterfaceWrapper) SaveProfile(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.SaveProfile(ctx)
	return err
}

// RemoveProfile converts echo context to params.
func (w *ServerInterfaceWrapper) RemoveProfile(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id ProfileId

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.RemoveProfile(ctx, id)
	return err
}

// GetProfile converts echo context to params.
func (w *ServerInterfaceWrapper) GetProfile(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id ProfileId

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetProfile(ctx, id)
	return err
}

// UpdateProfile converts echo context to params.
func (w *ServerInterfaceWrapper) UpdateProfile(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id ProfileId

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.UpdateProfile(ctx, id)
	return err
}

// GetTask converts echo context to params.
func (w *ServerInterfaceWrapper) GetTask(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "name" -------------
	var name QueueName

	err = runtime.BindStyledParameterWithLocation("simple", false, "name", runtime.ParamLocationPath, ctx.Param("name"), &name)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter name: %s", err))
	}

	// ------------- Path parameter "id" -------------
	var id TaskId

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetTask(ctx, name, id)
	return err
}

// GetTaskResponse converts echo context to params.
func (w *ServerInterfaceWrapper) GetTaskResponse(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "name" -------------
	var name QueueName

	err = runtime.BindStyledParameterWithLocation("simple", false, "name", runtime.ParamLocationPath, ctx.Param("name"), &name)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter name: %s", err))
	}

	// ------------- Path parameter "id" -------------
	var id TaskId

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetTaskResponse(ctx, name, id)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/auth/login", wrapper.Login)
	router.POST(baseURL+"/background/delay", wrapper.ProcessBackground)
	router.GET(baseURL+"/health/live", wrapper.LiveCheck)
	router.GET(baseURL+"/health/ready", wrapper.ReadyCheck)
	router.GET(baseURL+"/profiles", wrapper.ListProfiles)
	router.POST(baseURL+"/profiles", wrapper.SaveProfile)
	router.DELETE(baseURL+"/profiles/:id", wrapper.RemoveProfile)
	router.GET(baseURL+"/profiles/:id", wrapper.GetProfile)
	router.PATCH(baseURL+"/profiles/:id", wrapper.UpdateProfile)
	router.GET(baseURL+"/queues/:name/tasks/:id", wrapper.GetTask)
	router.GET(baseURL+"/queues/:name/tasks/:id/response", wrapper.GetTaskResponse)

}
