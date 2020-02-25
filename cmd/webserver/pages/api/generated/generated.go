// Package generated provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package generated

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi"
	"net/http"
	"strings"
)

// AppSchema defines model for app-schema.
type AppSchema struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// ErrorSchema defines model for error-schema.
type ErrorSchema struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// PaginationSchema defines model for pagination-schema.
type PaginationSchema struct {
	Limit        int64 `json:"limit"`
	Offset       int64 `json:"offset"`
	PagesTotal   int   `json:"pagesTotal"`
	Pagescurrent int   `json:"pagescurrent"`
	RowsFiltered int64 `json:"rowsFiltered"`
	RowsTotal    int64 `json:"rowsTotal"`
}

// PlayerSchema defines model for player-schema.
type PlayerSchema struct {
	Id   *int    `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

// SucccessSchema defines model for succcess-schema.
type SucccessSchema struct {
	Message string `json:"message"`
}

// App defines model for app.
type App AppSchema

// Apps defines model for apps.
type Apps struct {
	Apps       []AppSchema      `json:"apps"`
	Pagination PaginationSchema `json:"pagination"`
}

// Error defines model for error.
type Error ErrorSchema

// PaginationResponse defines model for pagination-response.
type PaginationResponse PaginationSchema

// Player defines model for player.
type Player PlayerSchema

// Players defines model for players.
type Players struct {
	Pagination PaginationSchema `json:"pagination"`
	Players    []PlayerSchema   `json:"players"`
}

// GetAppsParams defines parameters for GetApps.
type GetAppsParams struct {
	Ids  *[]interface{} `json:"ids,omitempty"`
	Tags *[]interface{} `json:"tags,omitempty"`

	// Offset
	Offset *int `json:"offset,omitempty"`

	// Limit
	Limit *int `json:"limit,omitempty"`
}

type ServerInterface interface {
	// List apps (GET /apps)
	GetApps(w http.ResponseWriter, r *http.Request)
	// Retrieve app (GET /apps/{id})
	GetAppsId(w http.ResponseWriter, r *http.Request)
	// Update a player (POST /players/{id})
	PostPlayersId(w http.ResponseWriter, r *http.Request)
}

// ParamsForGetApps operation parameters from context
func ParamsForGetApps(ctx context.Context) *GetAppsParams {
	return ctx.Value("GetAppsParams").(*GetAppsParams)
}

// GetApps operation middleware
func GetAppsCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var err error

		ctx = context.WithValue(ctx, "key-query.Scopes", []string{""})

		ctx = context.WithValue(ctx, "key-cookie.Scopes", []string{""})

		ctx = context.WithValue(ctx, "key-header.Scopes", []string{""})

		// Parameter object where we will unmarshal all parameters from the context
		var params GetAppsParams

		// ------------- Optional query parameter "ids" -------------
		if paramValue := r.URL.Query().Get("ids"); paramValue != "" {

		}

		err = runtime.BindQueryParameter("form", true, false, "ids", r.URL.Query(), &params.Ids)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid format for parameter ids: %s", err), http.StatusBadRequest)
			return
		}

		// ------------- Optional query parameter "tags" -------------
		if paramValue := r.URL.Query().Get("tags"); paramValue != "" {

		}

		err = runtime.BindQueryParameter("form", true, false, "tags", r.URL.Query(), &params.Tags)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid format for parameter tags: %s", err), http.StatusBadRequest)
			return
		}

		// ------------- Optional query parameter "offset" -------------
		if paramValue := r.URL.Query().Get("offset"); paramValue != "" {

		}

		err = runtime.BindQueryParameter("form", true, false, "offset", r.URL.Query(), &params.Offset)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid format for parameter offset: %s", err), http.StatusBadRequest)
			return
		}

		// ------------- Optional query parameter "limit" -------------
		if paramValue := r.URL.Query().Get("limit"); paramValue != "" {

		}

		err = runtime.BindQueryParameter("form", true, false, "limit", r.URL.Query(), &params.Limit)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid format for parameter limit: %s", err), http.StatusBadRequest)
			return
		}

		ctx = context.WithValue(ctx, "GetAppsParams", &params)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetAppsId operation middleware
func GetAppsIdCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var err error

		// ------------- Path parameter "id" -------------
		var id int32

		err = runtime.BindStyledParameter("simple", false, "id", chi.URLParam(r, "id"), &id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid format for parameter id: %s", err), http.StatusBadRequest)
			return
		}

		ctx = context.WithValue(ctx, "id", id)

		ctx = context.WithValue(ctx, "key-header.Scopes", []string{""})

		ctx = context.WithValue(ctx, "key-query.Scopes", []string{""})

		ctx = context.WithValue(ctx, "key-cookie.Scopes", []string{""})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// PostPlayersId operation middleware
func PostPlayersIdCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var err error

		// ------------- Path parameter "id" -------------
		var id int64

		err = runtime.BindStyledParameter("simple", false, "id", chi.URLParam(r, "id"), &id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid format for parameter id: %s", err), http.StatusBadRequest)
			return
		}

		ctx = context.WithValue(ctx, "id", id)

		ctx = context.WithValue(ctx, "key-cookie.Scopes", []string{""})

		ctx = context.WithValue(ctx, "key-header.Scopes", []string{""})

		ctx = context.WithValue(ctx, "key-query.Scopes", []string{""})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerFromMux(si, chi.NewRouter())
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r *chi.Mux) http.Handler {
	r.Group(func(r chi.Router) {
		r.Use(GetAppsCtx)
		r.Get("/apps", si.GetApps)
	})
	r.Group(func(r chi.Router) {
		r.Use(GetAppsIdCtx)
		r.Get("/apps/{id}", si.GetAppsId)
	})
	r.Group(func(r chi.Router) {
		r.Use(PostPlayersIdCtx)
		r.Post("/players/{id}", si.PostPlayersId)
	})

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8xXTW/jNhD9K8K0RyVyskWx0KkpihZpF2jQtCfDh6k0lrkrfixJZWME/u/FkPq0ZcdB",
	"jcVeYtMczpt5fDOcvEChpdGKlHeQv4AlZ7RyFBZoDH8UWnlSvv2lFgV6oVX20WnFv7liQxL52/eW1pDD",
	"d9ngM4u7LkNjrlrL3W6XQkmusMKwJ8jhTiUMtksZwb0JdOrog3A+0Wv25tLki/CbxGAlVDgNKRirDVkv",
	"+gTDp/Ak3VsySMFvDUEOaC1ueT1CecXPYDkQkoKlz42wVEK+hEnEIcbVDGfjVDkCslbbi91X8Hb6xiLg",
	"JPerTj8Xi2OOrYNgHrCiRKi1tjLSxkHVuKXL8RHdnSIkaRF77LfpeKrM/yWnSQBnaXsvu315n9JnB3VK",
	"op0NW7SQLSFXxwgQJf9twxDKUxWpVShptOO8Fao6iFCU0Jquuso4ilTokuaxJDmH1RlwwcVgv5rWxDHg",
	"WkgRhBF1G7F//AHSmVD0eu3oXGODFbm/tcd6Pq+wXzTWtsI8tLD6i/tV1J5CgmeB8pEe81X7PQLb9NKW",
	"k7G3vWAm2e2lsuqVfzFZsWCboijIuaM+z5bJoA92S0Vjhd8+stfo6BNtrwqtP4ngS3AFtcsuQDYZ2EQj",
	"/qDw+vDJDWEZG1442S7POfm5IbvtD8bV6XPcIbBykC9fOrv2GerXfc3zrXBv7rohFkF0JFHUkMNHIQmr",
	"mn6q+IfrQssB+/c6bEEKjWXbjffG5VlWoaTy32utaqEo65xyUMLXfPDRE8rkl5+Tu4d7SOGJrIsd6SYU",
	"kyGFRkAO764X14sgI78JV5B1I0EVq43vOVTxfQk5/Eb+jvf5gEVJPvTY5TxxomTDQTTjmgiyS0Hi831s",
	"zzeLxWHTnfcbiH+D4zm/0zb9Z1d9c3B9aY5HrjU2tYd8kYIUSshGhu+HRX44nsXyngPqSn8Gh5OQ+ByB",
	"AlU97M0M7CqdzrG3i8VXmmOnM1loHlIiF1fcwqietnb62S6NusteRLl7TXz35RH5sYbH6oNx8/G2oWOa",
	"eXcL3yyh3T8GEyr/Im8FPVHYmmWzbT49oUa7GUYftPMP0fLyrIZnT+LzB1KV30B++5V53X+2Zsh9bILF",
	"Hrv/mBI9JdjNswPB01Gve8ACV+Ona7maPkjdun1mlivO25F96ogemnueZbUusN5o5/P3i/c3Gffp3Wr3",
	"XwAAAP//eN3svqgOAAA=",
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file.
func GetSwagger() (*openapi3.Swagger, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error loading Swagger: %s", err)
	}
	return swagger, nil
}