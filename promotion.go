package artifactory

import (
	"net/http"
	"encoding/json"
)

type promotion struct{
	r *artifactoryRequest
	c Client
	data interface{}
}

const (
	PROMOTE_PARAMS = "promote"
)

const (
	// KIND_PROMOTION int = iota
)

func (p *promotion)TargetRepo(targetRepo string) *promotion  {
	p.setPromotion()
	p.data = &ArtifactoryPromote{
		TargetRepo: targetRepo,
		DockerRepository: p.r.image,
		// TargetDockerRepository: "",
		Tag: p.r.tag,
		// TargetTag: "",
		Copy: true,
	}
	return p
}

func (p *promotion)setPromotion() {
	p.r.method = http.MethodPost
	p.r.prefix = API_DOCKER
	p.r.params = PROMOTE_PARAMS
}

func (p *promotion)Set() (interface{}, error) {
	var err error
	var data string
	_, data, err = p.promote(p.r, p.data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (p *promotion)promote(req *artifactoryRequest, data interface{}) (int, string, error) {
	path := req.host + req.prefix + req.repository + "/" + API_VERSION + "/" + req.params
	params := data.(*ArtifactoryPromote)
	postdata, err := json.Marshal(params)
	if err != nil {
		return http.StatusInternalServerError, "", err
	}
	code, rtndata, err := p.c.Post(path, postdata)
	return code, rtndata, err
}
