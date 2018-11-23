package artifactory

import "time"

type artifactoryRequest struct {
	protocol   string
	host       string
	method     string
	repository string
	image      string
	tag        string
	prefix     string
	params     string

	kind int
}

type QueryImagesByProperties struct {
	Results []URI `json:"results"`
}

type URI struct {
	URI string `json:"uri"`
}

type QueryPropertiesByImage struct {
	Properties map[string][]string `json:"properties"`
	URI        string              `json:"uri"`
}

type FolderInfo struct {
	URI          string    `json:"uri"`
	Repo         string    `json:"repo"`
	Path         string    `json:"path"`
	Created      time.Time `json:"created"`
	CreatedBy    string    `json:"createdBy"`
	LastModified time.Time `json:"lastModified"`
	Children     []struct {
		URI    string `json:"uri"`
		Folder bool   `json:"folder"`
	} `json:"children"`
}

type FileInfo struct {
	// General
	URI          string    `json:"uri"`
	Size         int32     `json:"size"`
	LastModified time.Time `json:"lastModified"`

	// File
	DownloadURI       string    `json:"downloadUri"`
	Repo              string    `json:"repo"`
	Path              string    `json:"path"`
	RemoteURL         string    `json:"remoteUrl"`
	Created           time.Time `json:"created"`
	CreatedBy         string    `json:"createdBy"`
	ModifiedBy        string    `json:"modifiedBy"`
	LastUpdated       time.Time `json:"lastUpdated"`
	MimeType          string    `json:"mimeType"`
	Checksums         Checksums `json:"checkSums"`
	OriginalChecksums Checksums `json:"originalChecksums"`

	// File list
	Folder bool   `json:"folder"`
	SHA1   string `json:"sha1"`
	SHA2   string `json:"sha2"`
}

type Checksums struct {
	MD5    string `json:"md5"`
	SHA1   string `json:"sha1"`
	SHA256 string `json:"sha256"`
}

type FileList struct {
	URI     string      `json:"uri"`
	Created time.Time   `json:"created"`
	Files   []*FileInfo `json:"files"`
}

type BuildInfo struct {
	BuildInfo buildInfoSpec `json:"buildInfo"`
}

type buildInfoSpec struct {
	Version              string     `json:"version"`
	Name                 string     `json:"name"`
	Number               string     `json:"number"`
	BuildAgent           buildAgent `json:"buildAgent"`
	Agent                buildAgent `json:"agent"`
	Started              string     `json:"started"`
	DurationMillis       int32      `json:"durationMills"`
	ArtifactoryPrinciple string     `json:"artifactoryPrinciple"`
	Modules              []Module   `json:"modules"`
	URI                  string     `json:"uri"`
}

type buildAgent struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Module struct {
	Properties map[string]string `json:"properties"`
	Id         string            `json:"id"`
	Artifacts  []Artifact        `json:"artifacts"`
}

type Artifact struct {
	Name   string `json:"name"`
	SHA1   string `json:"sha1"`
	SHA256 string `json:"sha256"`
	MD5    string `json:"md5"`
}

type DockerRepository struct {
	Repositories []string `json:"repositories"`
}


type DockerImage struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type ArtifactoryPromote struct {
	TargetRepo             string `json:"targetRepo"`        // The target repository for the move or copy
	DockerRepository       string `json:"dockerRepository"`  // The docker repository name to promote
	TargetDockerRepository string `json:"targetDockerRepository"` // An optional docker repository name, if null, will use the same name as 'dockerRepository'
	Tag                    string `json:"tag"`               // An optional tag name to promote, if null - the entire docker repository will be promoted. Available from v4.10.
	TargetTag              string `json:"targetTag"`         // An optional target tag to assign the image after promotion, if null - will use the same tag
	Copy                   bool   `json:"copy"`              // An optional value to set whether to copy instead of move. Default: false
}
