package dirs

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type Folder interface {
	Exists(path string) bool
	Create(path string) error
	Delete(path string) error
	Rename(path string, newName string) error
	WalkUpTree(path string) ([]string, error)
	WalkDownTree(path string) ([]string, error)
}

func NewFolder(permissions os.FileMode) Folder {
	return &folderImpl{}
}

type folderImpl struct {
	permissions os.FileMode
}

func (f *folderImpl) Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (f *folderImpl) Create(path string) error {
	err := os.MkdirAll(path, f.permissions)
	return err
}

func (f *folderImpl) Delete(path string) error {
	err := os.RemoveAll(path)
	return err
}

func (f *folderImpl) Rename(path string, newName string) error {
	parentDir := filepath.Dir(path)
	newFolderPath := filepath.Join(parentDir, newName)
	err := os.Rename(path, newFolderPath)
	if err != nil {
		return err
	}
	return nil
}

// WalkUpTree walks up the folder tree from the current folder.
func (f *folderImpl) WalkUpTree(path string) ([]string, error) {
	var walkedPaths []string
	currentDir := path
	for currentDir != "" {
		walkedPaths = append(walkedPaths, currentDir)
		currentDir = filepath.Dir(currentDir)
	}
	return walkedPaths, nil
}

// WalkDownTree walks down the folder tree starting from the current folder.
func (f *folderImpl) WalkDownTree(path string) ([]string, error) {
	var walkedPaths []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		walkedPaths = append(walkedPaths, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return walkedPaths, nil
}

type File interface {
	Write(path string, data []byte) (int, error)
	Read(path string, buf []byte) (int, error)
	Rename(path string, newName string) error
	Delete(path string) error
	ListAllFiles(path string) ([]string, error)
}

func NewFile() File {
	return &fileImpl{}
}

type fileImpl struct {
}

func (f *fileImpl) Write(path string, data []byte) (int, error) {
	file, err := os.Create(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	n, err := file.Write(data)
	return n, err
}

func (f *fileImpl) Read(path string, buf []byte) (int, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}

	copy(buf, data)
	return len(data), nil
}

func (f *fileImpl) Rename(path string, newName string) error {
	err := os.Rename(path, newName)
	if err != nil {
		return err
	}
	return nil
}

func (f *fileImpl) Delete(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

func (f *fileImpl) ListAllFiles(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}

	return fileNames, nil
}
