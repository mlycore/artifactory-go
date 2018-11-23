package artifactory

import (
	"net/http"
	"fmt"
	"encoding/json"
	"strconv"
)

type buildInfo struct{
	r *artifactoryRequest
	c Client
	data interface{}
}

const (
	API_BUILD = "/artifactory/api/build/"
)

const (
	KIND_GET_BUILDINFO int = iota
)

func (b *buildInfo)GetBuildInfo(buildNum int32) *buildInfo  {
	b.data = buildNum
	b.setGetBuildInfo()
	return b
}

func (b *buildInfo)setGetBuildInfo() {
	b.r.method = http.MethodGet
	b.r.prefix = API_BUILD
	b.r.kind = KIND_GET_BUILDINFO
}

func (b *buildInfo)getBuildInfo(req *artifactoryRequest, getdata interface{}) (int, string, error) {
	buildname := req.image
	buildnum := strconv.Itoa(int(getdata.(int32)))
	path := req.protocol + req.host + req.prefix + buildname + "/" + buildnum
	fmt.Printf("%s\n", path)
	code, data, err := b.c.Get(path)
	fmt.Printf("%s\n", data)
	return code, data, err
}

func (b *buildInfo)GetObject()(interface{}, error)  {
	var err error
	var data string
	var obj interface{}

	switch b.r.kind{
	case KIND_GET_BUILDINFO:
		_, data, err = b.getBuildInfo(b.r, b.data)
		obj = &BuildInfo{}
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

	b.data = obj
	return obj, nil

}