package controller

import (
	"archive/zip"
	"fmt"
	"github.com/csby/gom/config"
	"github.com/csby/gwsf/gtype"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type base struct {
	gtype.Base

	cfg *config.Config
	wsc gtype.SocketChannelCollection
}

func (s *base) createCatalog(doc gtype.Doc, names ...string) gtype.Catalog {
	root := doc.AddCatalog("管理平台接口")

	count := len(names)
	if count < 1 {
		return root
	}

	child := root
	for i := 0; i < count; i++ {
		name := names[i]
		child = child.AddChild(name)
	}

	return child
}

func (s *base) writeOptMessage(id int, data interface{}) bool {
	if s.wsc == nil {
		return false
	}
	msg := &gtype.SocketMessage{
		ID:   id,
		Data: data,
	}

	s.wsc.Write(msg, nil)

	return true
}

func (s *base) sizeToText(v float64) string {
	kb := float64(1024)
	mb := 1024 * kb
	gb := 1024 * mb

	if v >= gb {
		return fmt.Sprintf("%.1fGB", v/gb)
	} else if v >= mb {
		return fmt.Sprintf("%.1fMB", v/mb)
	} else if v >= kb {
		return fmt.Sprintf("%.1fKB", v/kb)
	} else {
		return fmt.Sprintf("%.0fB", v)
	}
}

func (s *base) compressFolder(fileWriter io.Writer, folderPath, folderName string, ignore func(name string) bool) error {
	zw := zip.NewWriter(fileWriter)
	defer zw.Close()

	return s.createSubFolder(zw, folderPath, folderName, ignore)
}

func (s *base) createSubFolder(zw *zip.Writer, folderPath, folderName string, ignore func(name string) bool) error {
	paths, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return err
	}
	for _, path := range paths {
		if ignore != nil {
			if ignore(path.Name()) {
				continue
			}
		}

		fp := filepath.Join(folderPath, path.Name())
		if path.IsDir() {
			subFolderName := path.Name()
			if folderName != "" {
				subFolderName = fmt.Sprintf("%s/%s", folderName, path.Name())
			}
			err = s.createSubFolder(zw, fp, subFolderName, nil)
			if err != nil {
				return err
			}
		} else {
			fi, err := os.Stat(fp)
			if err != nil {
				return err
			}

			fr, err := os.Open(fp)
			if err != nil {
				return err
			}
			defer fr.Close()

			fn := fi.Name()
			if folderName != "" {
				fn = fmt.Sprintf("%s/%s", folderName, fi.Name())
			}
			fh, err := zip.FileInfoHeader(fi)
			if err != nil {
				return err
			}
			fh.Name = fn
			fh.Method = zip.Deflate
			fw, err := zw.CreateHeader(fh)
			if err != nil {
				return err
			}
			_, err = io.Copy(fw, fr)
			if err != nil {
				return err
			}
			zw.Flush()

			fr.Close()
		}
	}

	return nil
}

func (s *base) getFilePath(folderPath, fileName string) (string, error) {
	fs, e := ioutil.ReadDir(folderPath)
	if e != nil {
		return "", e
	}

	for _, f := range fs {
		name := f.Name()
		path := filepath.Join(folderPath, name)

		if !f.IsDir() {
			if name == fileName {
				return path, nil
			}
		} else {
			p, e := s.getFilePath(path, fileName)
			if e == nil {
				return p, nil
			}
		}
	}

	return "", fmt.Errorf("not found")
}
