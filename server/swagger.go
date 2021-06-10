package server

import (
	"encoding/json"
	"reflect"
	"strings"
)

const (
	SwaggerDescriptionTagName   = "swagger-description"
	SwaggerParametersInBody     = "body"
	SwaggerParametersInFormData = "formData"
	SwaggerParametersInPath     = "path"
)

func NewSwagger() *Swagger {
	return &Swagger{
		Swagger: "2.0",
		Info: SwaggerInfo{
			Title:   "",
			Version: "",
		},
		tagName: "json",
		Schemes: []string{
			"http",
			"https",
		},
		Consumes: []string{
			"application/json",
			"multipart/form-data",
			"application/x-www-form-urlencoded",
		},
		Produces: []string{
			"application/json",
			"application/x-www-form-urlencoded",
		},
		Paths:       make(map[string]SwaggerPath, 0),
		Definitions: make(map[string]SwaggerDefinition),
	}
}

type Swagger struct {
	tagName     string
	Swagger     string                       `json:"swagger"`
	Info        SwaggerInfo                  `json:"info"`
	Host        []string                     `json:"host"`
	BasePath    string                       `json:"basePath"`
	Schemes     []string                     `json:"schemes"`
	Consumes    []string                     `json:"consumes"`
	Produces    []string                     `json:"produces"`
	Paths       map[string]SwaggerPath       `json:"paths"`
	Definitions map[string]SwaggerDefinition `json:"definitions"`
}

func (s *Swagger) SetInfoTitle(v string) {
	s.Info.Title = v
}

func (s *Swagger) ToJSON() string {
	b, _ := json.Marshal(s)
	return string(b)
}

func (s *Swagger) isStructArray(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Kind() == reflect.Struct
}

func (s *Swagger) scanStruct(inType reflect.Type) {
	for inType.Kind() == reflect.Ptr {
		inType = inType.Elem()
	}
	num := inType.NumField()
	definition := SwaggerDefinition{
		Type:        "object",
		Properties:  make(map[string]SwaggerDefinitionProperty),
		Description: "",
	}

	for i := 0; i < num; i++ {
		f := inType.Field(i)

		fieldName := f.Name
		if f.Tag.Get(s.tagName) != "" {
			fieldName = strings.Split(f.Tag.Get(s.tagName), ",")[0]
		}
		description := f.Tag.Get(SwaggerDescriptionTagName)

		fieldType := f.Type
		for fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		if fieldType.Kind() == reflect.Slice {
			v := reflect.New(fieldType.Elem())
			vType := v.Type()
			if vType.Kind() == reflect.Ptr {
				vType = vType.Elem()
			}

			if s.isStructArray(vType) {
				definition.Properties[fieldName] = SwaggerDefinitionProperty{
					Type: "array",
					Items: map[string]string{
						"$ref": "#/definitions/" + s.getName(vType),
					},
					Description: description,
				}
				s.scanStruct(vType)
			} else {
				definition.Properties[fieldName] = SwaggerDefinitionProperty{
					Type: "array",
					Items: map[string]string{
						"type": vType.Name(),
					},
					Description: description,
				}
			}
		} else if fieldType.Kind() == reflect.Struct {
			s.scanStruct(fieldType)
			definition.Properties[fieldName] = SwaggerDefinitionProperty{
				Ref:         "#/definitions/" + s.getName(fieldType),
				Description: description,
			}
		} else {
			typeName, typeFormat := s.formatType(fieldType)
			definition.Properties[fieldName] = SwaggerDefinitionProperty{
				Type:        typeName,
				Format:      typeFormat,
				Description: description,
			}
		}
	}
	s.Definitions[ s.getName(inType)] = definition
}

func (s *Swagger) formatType(t reflect.Type) (string, string) {
	typeName := t.String()
	typeFormat := ""
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Int:
		typeName = "integer"
		typeFormat = t.String()
	case reflect.Float32, reflect.Float64:
		typeName = "number"
		typeFormat = "float"
	case reflect.Bool:
		typeName = "boolean"
		typeFormat = ""
	}
	return typeName, typeFormat
}

func (s *Swagger) getName(t reflect.Type) string {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return strings.Replace(t.PkgPath(), "/", ".", -1) + "." + t.Name()
}

func (s *Swagger) buildParameters(in string, inType reflect.Type) []SwaggerParameter {
	parameters := make([]SwaggerParameter, 0)

	for inType.Kind() == reflect.Ptr {
		inType = inType.Elem()
	}
	if in == SwaggerParametersInBody {
		param := SwaggerParameter{
			Name: in,
			IN:   in,
		}
		param.Schema = SwaggerParameterSchema{
			Ref: "#/definitions/" + s.getName(inType),
		}
		parameters = append(parameters, param)
		return parameters
	}

	// @TODO for form data
	num := inType.NumField()
	for i := 0; i < num; i++ {
		f := inType.Field(i)
		name := f.Tag.Get(s.tagName)
		if name == "" {
			continue
		}
		name = strings.Split(name, ",")[0]
		param := SwaggerParameter{
			Name: name,
			IN:   in,
		}
		parameters = append(parameters, param)
	}
	return parameters
}

func (s *Swagger) buildResponse(inType reflect.Type) SwaggerResponse {
	rsp := SwaggerResponse{
		Description: "",
		Schema: SwaggerSchema{
			Type:       "object",
			Properties: make(map[string]SwaggerDefinitionProperty),
		},
	}
	for inType.Kind() == reflect.Ptr {
		inType = inType.Elem()
	}
	if inType.Kind() != reflect.Struct {
		return rsp
	}

	num := inType.NumField()
	for i := 0; i < num; i++ {
		f := inType.Field(i)
		name := f.Tag.Get(s.tagName)
		if name == "" {
			continue
		}
		description := f.Tag.Get(SwaggerDescriptionTagName)
		name = strings.Split(name, ",")[0]

		fieldType := f.Type
		for fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		if fieldType.Kind() == reflect.Slice {
			v := reflect.New(fieldType.Elem())
			vType := v.Type()
			if vType.Kind() == reflect.Ptr {
				vType = vType.Elem()
			}
			if s.isStructArray(vType) {
				rsp.Schema.Properties[name] = SwaggerDefinitionProperty{
					Type: "array",
					Items: map[string]string{
						"$ref": "#/definitions/" + s.getName(vType),
					},
					Description: description,
				}
			} else {
				typeName, typeFormat := s.formatType(vType)
				rsp.Schema.Properties[name] = SwaggerDefinitionProperty{
					Type: "array",
					Items: map[string]string{
						"type":   typeName,
						"format": typeFormat,
					},
					Description: description,
				}
			}
		} else if fieldType.Kind() == reflect.Struct {
			rsp.Schema.Properties[name] = SwaggerDefinitionProperty{
				Ref:         "#/definitions/" + s.getName(fieldType),
				Description: description,
			}
		} else {
			typeName, typeFormat := s.formatType(fieldType)
			rsp.Schema.Properties[name] = SwaggerDefinitionProperty{
				Type:        typeName,
				Description: f.Tag.Get(SwaggerDescriptionTagName),
				Format:      typeFormat,
			}
		}
	}

	return rsp
}

func (s *Swagger) AddHandler(method string, path string, req reflect.Type, rsp reflect.Type, opts *RouterHandlerOptions) error {
	s.scanStruct(req)
	s.scanStruct(rsp)
	parameters := s.buildParameters(opts.ParametersIn, req)
	method = strings.ToLower(method)
	if _, ok := s.Paths[path]; ok {
		s.Paths[path][method] = SwaggerPathEndpoint{
			Summary: opts.Summary,
			Tags:    opts.Tags,
			Responses: map[string]SwaggerResponse{
				"200": SwaggerResponse{
					Description: "",
					Schema: SwaggerSchema{
						Type:       "object",
						Properties: make(map[string]SwaggerDefinitionProperty),
					},
				},
			},
			Parameters: parameters,
		}
	} else {
		s.Paths[path] = SwaggerPath{
			method: SwaggerPathEndpoint{
				Summary: opts.Summary,
				Tags:    opts.Tags,
				Responses: map[string]SwaggerResponse{
					"200": s.buildResponse(rsp),
				},
				Parameters: parameters,
			},
		}
	}
	if opts.ResponseWrapper != nil {
		wrapperType := reflect.TypeOf(opts.ResponseWrapper)
		s.scanStruct(wrapperType)
		wrapperRsp := s.buildResponse(wrapperType)

		description := ""
		field, found := s.getFieldByTagNameName(wrapperType, opts.ResponseWrapperDataNodeName)
		if found {
			description = field.Tag.Get(SwaggerDescriptionTagName)
		}
		wrapperRsp.Schema.Properties[opts.ResponseWrapperDataNodeName] = SwaggerDefinitionProperty{
			Ref:         "#/definitions/" + s.getName(rsp),
			Description: description,
		}
		s.Paths[path][method].Responses["200"] = wrapperRsp
	}
	return nil
}

func (s *Swagger) getFieldByTagNameName(t reflect.Type, tagName string) (reflect.StructField, bool) {
	num := t.NumField()
	for i := 0; i < num; i++ {
		f := t.Field(i)
		if f.Tag.Get(s.tagName) == tagName {
			return f, true
		}
	}
	return reflect.StructField{}, false
}

type SwaggerInfoContact struct {
	Email string `json:"email"`
}

type SwaggerInfo struct {
	Title       string             `json:"title"`
	Version     string             `json:"version"`
	Description string             `json:"description"`
	Contact     SwaggerInfoContact `json:"contact"`
}

type SwaggerSchema struct {
	Type       string                               `json:"type"`
	Properties map[string]SwaggerDefinitionProperty `json:"properties"`
}

type SwaggerResponse struct {
	Description string        `json:"description"`
	Schema      SwaggerSchema `json:"schema"`
}

type SwaggerParameterSchema struct {
	Ref string `json:"$ref"`
}

type SwaggerParameter struct {
	Name     string                 `json:"name"`
	IN       string                 `json:"in"`
	Required bool                   `json:"required"`
	Schema   SwaggerParameterSchema `json:"schema"`
}
type SwaggerPathEndpoint struct {
	Summary    string                     `json:"summary"`
	Responses  map[string]SwaggerResponse `json:"responses"`
	Parameters []SwaggerParameter         `json:"parameters"`
	Tags       []string                   `json:"tags"`
}
type SwaggerPath map[string]SwaggerPathEndpoint

type SwaggerDefinitionProperty struct {
	Type        string            `json:"type,omitempty"`
	Res         map[string]string `json:"res,omitempty"`
	Format      string            `json:"format,omitempty"` // 当为number时,可以指定format为float等
	Ref         string            `json:"$ref,omitempty"`   // 当类型为引用时使用
	Items       map[string]string `json:"items,omitempty"`  // 当类型为数组时使用
	Description string            `json:"description"`      // 属性说明
}

type SwaggerDefinition struct {
	Type        string                               `json:"type"`
	Properties  map[string]SwaggerDefinitionProperty `json:"properties"`
	Description string                               `json:"description"`
}
