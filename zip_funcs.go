package wfimport

import (
	"archive/zip"
	"bufio"
	"path/filepath"

	"github.com/writeas/go-writeas/v2"
)

// TopLevelZipFunc return a poiter to a writeas.PostParams for any parseable
// zip.File that is not a directory. It does not traverse children.
//
// This is an example of a ZipFunc that can be used to filter files parse from
// a zip archive.
func TopLevelZipFunc(f *zip.File) (*writeas.PostParams, error) {
	if !f.FileInfo().IsDir() {
		return openAndParse(f)
	}
	return nil, nil
}

// TextFileZipFunc parses .txt files into PostParams
func TextFileZipFunc(f *zip.File) (*writeas.PostParams, error) {
	if !f.FileInfo().IsDir() && filepath.Ext(f.FileHeader.Name) == ".txt" {
		return openAndParse(f)
	}
	return nil, nil
}

func openAndParse(f *zip.File) (*writeas.PostParams, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	r := bufio.NewReader(rc)
	b := make([]byte, f.FileInfo().Size())
	_, err = r.Read(b)
	if err != nil {
		return nil, err
	}
	p, err := fromBytes(b)
	if err != nil {
		return nil, err
	}
	p.Created = &f.Modified

	return p, nil
}
