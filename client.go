package artifactory

import (
	"net/http"
)

var _ Client = &HttpsClient{}
type Client interface {
	BasicAuth(user, pass string) Client
	Certificate(filepath string) Client
	Put(path string, data []byte)(int, string ,error)
	Delete(path string) (int, string, error)
	Get(path string) (int, string, error)
	Post(path string, data []byte)(int, string, error)
	RawGet(path string)(int, *http.Response, error)
}

type ArtifactoryClient struct{
	core Client
	r *artifactoryRequest
}

//TODO: host port需要指定
func NewArtifactoryClient(host string) *ArtifactoryClient {
	return &ArtifactoryClient{
		core: &HttpsClient{
			Client: &http.Client{},
		},
		r : &artifactoryRequest{
			host: host,
		},
	}
}

func (c *ArtifactoryClient)BasicAuth(user, pass string) *ArtifactoryClient {
	c.core.BasicAuth(user, pass)
	return c
}

func (c *ArtifactoryClient)Certificate(certpath string) *ArtifactoryClient {
	c.core.Certificate(certpath)
	return c
}

func (c *ArtifactoryClient) Repository(repo string) *ArtifactoryClient {
	c.r.repository = repo
	return c
}

func (c *ArtifactoryClient) Image(image string) *ArtifactoryClient {
	c.r.image = image
	return c
}

func (c *ArtifactoryClient) Tag(tag string) *ArtifactoryClient {
	c.r.tag = tag
	return c
}

func (c *ArtifactoryClient)ImageList() *imageList  {
	l := &imageList{
		r: c.r,
		c: c.core,
		//TODO: 支持map[string][]string
		List: make([]string, 0),
	}
	return l
}

func (c *ArtifactoryClient)TagList() *tagList  {
	l := &tagList{
		r: c.r,
		c: c.core,
		//TODO: 支持map[string][]string
		List: make([]string, 0),
	}
	return l
}

func (c *ArtifactoryClient)Property() *properties {
	p := &properties{
		r: c.r,
		c: c.core,
		//TODO: 支持map[string][]string
		data: make(map[string]string),
	}
	return p
}

func (c *ArtifactoryClient)Storage() *storage {
	return &storage{
		r: c.r,
		c: c.core,
	}
}

func (c *ArtifactoryClient)BuildInfo() *buildInfo {
	return &buildInfo{
		r: c.r,
		c: c.core,
	}
}

func (c *ArtifactoryClient)Promotion() *promotion {
	return &promotion {
		r: c.r,
		c: c.core,
	}
}
