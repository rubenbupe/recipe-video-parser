package gallery

import (
	"github.com/kelseyhightower/envconfig"
)

func CreateConfig() (*Galleryconfig, error) {
	var cfg Galleryconfig
	err := envconfig.Process("GALLERY", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

type Galleryconfig struct {
	DownloadDir string `default:"./tmp"`
	PublicUrl   string `default:"http://localhost:8080"`
}
