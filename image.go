package artifactory
import(
	"net/http"
)


type imageList struct {
	r *artifactoryRequest
	c Client
	List	[]string `json:"list"`
}

const (
	API_DOCKER = "/artifactory/api/docker/"
	API_VERSION = "v2"
	IMAGE_LIST = "/_catalog"
	TAG_LIST = "/tags/list"
)


func (l *imageList)setGetImageList()  {
	l.r.method = http.MethodGet
	l.r.prefix = API_DOCKER
	l.r.params = API_VERSION + IMAGE_LIST
}

func (l *imageList)GetImageList() *imageList {
	l.setGetImageList()
	return l
}

func (l *imageList)GetObject() (interface{}, error) {
	var err error
	var data interface{}

	_, data, err = l.getImageList(l.r)
	if err != nil {
		return nil, err
	}
	l.List = data.([]string)
	return l.List, nil
}


func (l *imageList)getImageList(req *artifactoryRequest)(int, interface{}, error)  {
	result := make([]string, 0)
	// 获取image列表
	imageList, err := DockerRegistryApiV2Repository(l.c, req)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// 调用getTagList获取tag列表
	for _, img := range imageList {
		req.params = TAG_LIST
		taglist, err := DockerRegistryApiV2Tags(l.c, req, img)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}

		// 组装数据
		for _, tag := range taglist {
			result = append(result, img + ":" + tag)
		}
	}
	return http.StatusOK, result, nil
}


type tagList struct {
	r *artifactoryRequest
	c Client
	List []string `json:"list"`
}

func (l *tagList)setGetTagList()  {
	l.r.method = http.MethodGet
	l.r.prefix = API_DOCKER
	l.r.params = TAG_LIST
}

func (l *tagList)GetTagList() *tagList {
	l.setGetTagList()
	return l
}

func (l *tagList)GetObject() (interface{}, error) {
	var err error
	var data interface{}

	_, data, err = l.getTagList(l.r)
	if err != nil {
		return nil, err
	}

	l.List = data.([]string)
	return l.List, nil
}

func (l *tagList)getTagList(req *artifactoryRequest) (int, interface{}, error)  {
	tagslist, err := DockerRegistryApiV2Tags(l.c, req, l.r.image)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, tagslist, nil
}
