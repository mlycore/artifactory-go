package artifactory

import (
	"fmt"
	"encoding/json"
	"strings"
	"net/http"
)

type properties struct {
	r          *artifactoryRequest
	c          Client
	data       interface{}
}

func (p *properties)queryImagesByProperties(req *artifactoryRequest, tagsdata interface{})(int, string, error)  {
	path := req.protocol + req.host + req.prefix + req.params
	tags := tagsdata.(map[string]string)
	for k, v := range tags {
		path += k + "=" + v + "&"
	}
	path = strings.TrimSuffix(path, "&")
	code, data, err := p.c.Get(path)
	return code, data, err
}

func (p *properties)queryPropertiesByImage(req *artifactoryRequest, tagsdata interface{})(int, string, error)  {
	path := req.protocol + req.host + req.prefix + req.repository + "/" + req.image + "/" + req.tag + req.params
	code, data, err := p.c.Get(path)
	return code, data, err
}

func (p *properties)addProperty(req *artifactoryRequest, tagsdata interface{}) (int, string, error) {
	path := req.protocol + req.host + req.prefix + req.repository + "/" + req.image + "/" + req.tag + req.params
	tags := tagsdata.(map[string]string)
	for k, v := range tags{
		path += k + "=" + v + ";"
	}
	code, data, err := p.c.Put(path, nil)
	return code, data, err
}

func (p *properties)deleteProperty(req *artifactoryRequest, tagsdata interface{}) (int, string, error) {
	path := req.protocol + req.host + req.prefix + req.repository + "/" + req.image + "/" + req.tag + req.params
	tags := tagsdata.([]string)
	for _, k := range tags {
		path += k + ","
	}
	code, data, err := p.c.Delete(path)
	return code, data, err
}

func (p *properties) Add(tags map[string]string) *properties {
	for k, v := range tags {
		p.data.(map[string]string)[k] = v
	}
	p.setAdd()
	return p
}

func (p *properties) Update(tags map[string]string) *properties {
	for k, v := range tags {
		p.data.(map[string]string)[k] = v
	}
	p.setUpdate()
	return p
}

func (p *properties) Delete(tags []string) *properties {
	p.data = tags
	p.setDelete()
	return p
}

const (
	WRITE_PROPERTY = "?properties="

	API_SEARCH = "/artifactory/api/search/"
	READ_PROP = "prop?"
	READ_PROPERTY = "?properties"

)

const (
	KIND_ADD_PROPERTY int = iota
	KIND_UPDATE_PROPERTY
	KIND_DELETE_PROPERTY
	KIND_QUERY_PROPERTYBYIMAGE
	KIND_QUERY_IMAGEBYPROPERTY
)

func (p *properties) Set() (interface{}, error) {
	var err error
	var data string

	switch p.r.kind{
	case KIND_ADD_PROPERTY:
		_, data, err = p.addProperty(p.r, p.data)
	case KIND_UPDATE_PROPERTY:
		_, data, err = p.addProperty(p.r, p.data)
	case KIND_DELETE_PROPERTY:
		_, data, err = p.deleteProperty(p.r, p.data)
	default:
		return nil, fmt.Errorf("unknown method")
	}

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (p *properties)GetObject() (interface{}, error) {
	var err error
	var data string
	var obj interface{}

	switch p.r.kind{
	case KIND_QUERY_PROPERTYBYIMAGE:
		_, data, err = p.queryPropertiesByImage(p.r, p.data)
		obj = &QueryPropertiesByImage{}
	case KIND_QUERY_IMAGEBYPROPERTY:
		_, data, err = p.queryImagesByProperties(p.r, p.data)
		obj = &QueryImagesByProperties{}
	default:
		return nil, fmt.Errorf("property set error")
	}

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(data), obj)
	if err != nil {
		return nil, err
	}

	p.data = obj
	return obj, nil
}

func (p *properties)QueryImages(tags map[string]string) *properties {
	for k, v := range tags {
		p.data.(map[string]string)[k] = v
	}
	p.setQueryImagesByProperties()
	return p
}

func (p *properties)QueryProperties() *properties {
	p.setQueryPropertiesByImage()
	return p
}


func (p *properties)setAdd()   {
	p.r.method = http.MethodPut
	p.r.prefix = API_STORAGE
	p.r.params = WRITE_PROPERTY
	p.r.kind = KIND_ADD_PROPERTY
}

func (p *properties)setUpdate()  {
	p.r.method = http.MethodPut
	p.r.prefix = API_STORAGE
	p.r.params = WRITE_PROPERTY
	p.r.kind = KIND_UPDATE_PROPERTY
}

func (p *properties)setDelete() {
	p.r.method = http.MethodDelete
	p.r.prefix = API_STORAGE
	p.r.params = WRITE_PROPERTY
	p.r.kind = KIND_DELETE_PROPERTY
}

func (p *properties)setQueryImagesByProperties()   {
	p.r.method = http.MethodGet
	p.r.prefix = API_SEARCH
	p.r.params = READ_PROP
	p.r.kind = KIND_QUERY_IMAGEBYPROPERTY
}

func (p *properties)setQueryPropertiesByImage()   {
	p.r.method = http.MethodGet
	p.r.prefix = API_STORAGE
	p.r.params = READ_PROPERTY
	p.r.kind = KIND_QUERY_PROPERTYBYIMAGE
}
