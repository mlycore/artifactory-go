package artifactory

import (
	"fmt"
	"encoding/json"
	"net/http"
)

type storage struct {
	r *artifactoryRequest
	c Client
	data interface{}
}

func (s *storage)GetImageFolderInfo() *storage {
	s.setGetImageFolderInfo()
	return s
}

func (s *storage)setGetImageFolderInfo()  {
	s.r.method = http.MethodGet
	s.r.prefix = API_STORAGE
	s.r.kind = KIND_GET_IMAGEFOLDERINFO
}

func (s *storage)getImageFolderInfo(req *artifactoryRequest, getdata interface{}) (int, string ,error) {
	path := req.protocol + req.host + req.prefix + req.repository + "/" + req.image + "/" + req.tag
	code, data, err := s.c.Get(path)
	return code, data, err
}

func (s *storage)GetImageFileList() *storage {
	s.setGetImageFileList()
	return s
}

func (s *storage)setGetImageFileList() {
	s.r.method = http.MethodGet
	s.r.prefix = API_STORAGE
	s.r.params = READ_FILELIST
	s.r.kind = KIND_GET_IMAGEFILELIST
}

func (s *storage)getImageFileList(req *artifactoryRequest, getdata interface{}) (int, string, error) {
	path := req.protocol + req.host + req.prefix + req.repository + "/" + req.image + "/" + req.tag + "/" + req.params
	code, data, err := s.c.Get(path)
	return code, data, err
}

const (
	API_STORAGE = "/artifactory/api/storage/"
	READ_FILELIST = "?list"
)

const (
	KIND_GET_IMAGEFILELIST int = iota
	KIND_GET_IMAGEFOLDERINFO
)

func (s *storage)GetObject() (interface{}, error) {
	var err error
	var data string
	var obj interface{}

	switch s.r.kind{
	case KIND_GET_IMAGEFOLDERINFO:
		_, data, err = s.getImageFolderInfo(s.r, s.data)
		obj = &FolderInfo{}
	case KIND_GET_IMAGEFILELIST:
		_, data, err = s.getImageFileList(s.r, s.data)
		obj = &FileList{}
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

	s.data = obj
	return obj, nil
}
