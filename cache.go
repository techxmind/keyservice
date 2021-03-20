package keyservice

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/sync/singleflight"

	"github.com/techxmind/go-utils/fileutil"
)

var (
	NoCache = &nilCache{}
)

type Cache interface {
	Load(id string) (cipherData []byte, err error)
	Store(id string, cipherData []byte) error
}

type nilCache struct {
}
func (c *nilCache) Load(id string) ([]byte, error) { return nil, nil }
func (c *nilCache) Store(id string, data []byte) error { return nil}

type fileCache struct {
	dir string
	sf  singleflight.Group
}

func (c *fileCache) Load(id string) (content []byte, err error) {
	if c.dir == "" {
		return nil, ErrNotFound
	}

	cacheFile := c.getCacheFile(id)

	if !fileutil.Exist(cacheFile) {
		return nil, ErrNotFound
	}

	val, err, _ := c.sf.Do("r-"+id, func() (interface{}, error) {
		return ioutil.ReadFile(cacheFile)
	})

	if err == nil {
		content = val.([]byte)
	}

	return
}

func (c *fileCache) Store(id string, content []byte) error {
	if c.dir == "" {
		return nil
	}

	_, err, _ := c.sf.Do("w-"+id, func() (interface{}, error) {
		cacheFile := c.getCacheFile(id)
		err := ioutil.WriteFile(cacheFile, content, fileutil.PrivateFileMode)
		return nil, err
	})

	return err
}

func (c *fileCache) getCacheFile(id string) string {
	filename := fmt.Sprintf("%x", md5.Sum([]byte("keymanager."+id)))

	return filepath.Join(c.dir, filename)
}

func NewCache() Cache {
	// 优先使用内存文件
	dir := "/dev/shm"
	if err := fileutil.IsDirWriteable(dir); err != nil {
		dir = os.TempDir()
		if err := fileutil.IsDirWriteable(dir); err != nil {
			dir = ""
		}
	}

	if dir != "" {
		dir = filepath.Join(dir, fmt.Sprintf("keymanager_%d", os.Getuid()))
		if fileutil.Exist(dir) {
			if err := fileutil.IsDirWriteable(dir); err != nil {
				dir = ""
			}
		} else if err := fileutil.CreateDirAll(dir); err != nil {
			dir = ""
		}
	}

	if dir == "" {
		logger.Errorf("no fileCache dir")
	}

	return &fileCache{
		dir: dir,
	}
}
