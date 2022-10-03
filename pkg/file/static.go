package file

import (
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"github.com/igevin/sepweb/pkg/context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type StaticResourceHandlerOption func(handler *StaticResourceHandler)

type StaticResourceHandler struct {
	dir                     string
	pathPrefix              string
	extensionContentTypeMap map[string]string

	cache       *lru.Cache
	maxFileSize int
}

type fileCacheItem struct {
	fileName    string
	fileSize    int
	contentType string
	data        []byte
}

func NewStaticResourceHandler(dir string, pathPrefix string,
	options ...StaticResourceHandlerOption) *StaticResourceHandler {
	res := &StaticResourceHandler{
		dir:        dir,
		pathPrefix: pathPrefix,
		extensionContentTypeMap: map[string]string{
			"jpeg": "image/jpeg",
			"jpe":  "image/jpeg",
			"jpg":  "image/jpeg",
			"png":  "image/png",
			"pdf":  "image/pdf",
		},
	}
	for _, opt := range options {
		opt(res)
	}
	return res
}

func WithFileCache(maxFileSizeThreshold int, maxCacheFileCnt int) StaticResourceHandlerOption {
	return func(h *StaticResourceHandler) {
		c, err := lru.New(maxCacheFileCnt)
		if err != nil {
			log.Printf("创建缓存失败，将不会缓存静态资源")
		}
		h.maxFileSize = maxFileSizeThreshold
		h.cache = c
	}
}

func WithMoreExtension(extMap map[string]string) StaticResourceHandlerOption {
	return func(h *StaticResourceHandler) {
		for ext, contentType := range extMap {
			h.extensionContentTypeMap[ext] = contentType
		}
	}
}

func (s *StaticResourceHandler) Handler(ctx *context.Context) {
	req, _ := ctx.PathValue("file").ToString()
	if item, ok := s.readFileFromData(req); ok {
		log.Printf("从缓存中读取数据...")
		s.writeFileAsResponse(item, ctx.Resp)
		return
	}
	path := filepath.Join(s.dir, req)
	f, err := os.Open(path)
	if err != nil {
		ctx.Resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	ext := getFileExt(f.Name())
	t, ok := s.extensionContentTypeMap[ext]
	if !ok {
		ctx.Resp.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		ctx.Resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	item := &fileCacheItem{
		fileSize:    len(data),
		data:        data,
		contentType: t,
		fileName:    req,
	}
	s.cacheFile(item)
	s.writeFileAsResponse(item, ctx.Resp)
}

func (s *StaticResourceHandler) cacheFile(item *fileCacheItem) {
	if s.cache != nil && item.fileSize < s.maxFileSize {
		s.cache.Add(item.fileName, item)
	}
}

func (s *StaticResourceHandler) readFileFromData(fileName string) (*fileCacheItem, bool) {
	if s.cache != nil {
		if item, ok := s.cache.Get(fileName); ok {
			return item.(*fileCacheItem), true
		}
	}
	return nil, false
}

func (s *StaticResourceHandler) writeFileAsResponse(item *fileCacheItem, writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", item.contentType)
	writer.Header().Set("Content-Length", fmt.Sprintf("%d", item.fileSize))
	_, _ = writer.Write(item.data)
}

func getFileExt(name string) string {
	index := strings.LastIndex(name, ".")
	if index == len(name)-1 {
		return ""
	}
	return name[index+1:]
}