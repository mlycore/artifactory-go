package artifactory
import (
	"net/http"
	"strings"
	"io/ioutil"
	"encoding/json"

	myregexp "github.com/maxwell92/gokits/regexp"
)


func DockerRegistryApiV2Repository(c Client, req *artifactoryRequest) ([]string, error) {
	// 检查响应头
	repositories, err := DockerRegistryApiV2CheckHeaderLink(c, req.host, req.prefix, req.repository, req.params)
	if err != nil {
		return []string{}, err
	}
	return repositories, nil
}

func DockerRegistryApiV2CheckHeaderLink(c Client, host, prefix, repository, params string) ([]string, error) {
	baseUrl := host + prefix + repository
	catalog := baseUrl + "/" + params

	_, resp, err := c.RawGet(catalog)
	if resp == nil || err != nil {
		return nil, err
	}

	repos := make([]string, 0)
	link, r := DockerRegistryApiV2ProcessResponse(resp)
	repos = append(repos, r...)

	// 如果条数超过100，则循环处理
	// TODO: 需要验证
	for {
		if strings.EqualFold(link, "") {
			break
		}

		// 请求下一页
		next := DockerRegistryApiV2ProcessRegexp(link)
		if strings.EqualFold(next, "") {
			break
		}

		url := baseUrl + next
		_, resp, err := c.RawGet(url)
		if resp == nil || err != nil {
			break
		}

		// 处理响应体
		link, r = DockerRegistryApiV2ProcessResponse(resp)
		repos = append(repos, r...)
	}

	return repos, nil
}

func DockerRegistryApiV2ProcessResponse(resp *http.Response) (string, []string){
	link := ""
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}

	repository := new(DockerRepository)
	err = json.Unmarshal(body, repository)
	if err != nil {
		return "", nil
	}

	link = resp.Header.Get("Link")
	return link, repository.Repositories
}


const Content   = `<(.*)>; `
func DockerRegistryApiV2ProcessRegexp(link string) string  {
	contentExp := myregexp.Match(Content)
	ss := contentExp.FindStringSubmatch(link)
	if len(ss) == 2 {
		return ss[1]
	}
	return ""
}

func DockerRegistryApiV2Tags(c Client, req *artifactoryRequest, image string)([]string, error)  {
	baseUrl := req.host + req.prefix + req.repository
	url := baseUrl + "/" + API_VERSION + "/" + image + req.params
	_, resp, err := c.RawGet(url)
	if resp == nil || err != nil {
		return nil, nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	i := new(DockerImage)
	err = json.Unmarshal(body, i)
	if err != nil {
		return nil, err
	}

	return i.Tags, nil
}
