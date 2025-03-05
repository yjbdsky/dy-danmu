package router

import (
	"danmu-http/logger"
	"embed"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
)

type Resource struct {
	fs   embed.FS
	path string
}

func NewResource() *Resource {
	return &Resource{
		fs:   Static,
		path: "dist",
	}
}

func (r *Resource) Open(name string) (fs.File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("http: invalid character in file path")
	}
	//fullName := filepath.Join(r.path, name)
	fullName := r.path + "/" + "assets" + "/" + name
	file, err := r.fs.Open(fullName)
	//file, err := r.fs.Open(name)
	logger.Info().Msg(fmt.Sprintf("open file: %s  path: %s  name: %s", fullName, r.path, name))
	if err != nil {
		logger.Error().Err(err).Msg("failed to open file")
	}
	return file, err
}

func InitResourceServer(engine *gin.Engine) *gin.Engine {
	engine.StaticFS("/assets", http.FS(NewResource()))
	return engine
}
