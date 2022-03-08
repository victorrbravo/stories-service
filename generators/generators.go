package generators

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
)

type FileListHandler interface {
	Write() error

	HandleFile(file os.FileInfo) error

	GetTargetPath() string
	GetSourcePath() string
	SetSourcePath(string)
}

func Generate(handler FileListHandler) error {
	var files []os.FileInfo
	var err error
	var outfile *os.File
	var info os.FileInfo
	files, err = ioutil.ReadDir(handler.GetSourcePath())
	if err != nil {
		return err
	}
	info, err = os.Stat(handler.GetTargetPath())
	if err == nil && info.IsDir() == false {
		log.Fatalf("error `%s' is not a directory", handler.GetTargetPath())
	} else if err != nil {
		err = os.Mkdir(handler.GetTargetPath(), os.ModeDir|0700)
		if err != nil {
			return err
		}
	}
	defer func() { err = outfile.Close() }()
	// Now visit all files
	for _, file := range files {
		err = handler.HandleFile(file)
		if err != nil {
			log.Fatalf(
				"generator failed for file `%s/%s'\n* \033[1;31m%s\033[0m",
				handler.GetSourcePath(),
				file.Name(),
				err.Error(),
			)
		}
	}
	return nil
}

func GetSourceFiles(handler FileListHandler) ([]os.FileInfo, error) {
	var files []os.FileInfo
	var err error
	files, err = ioutil.ReadDir(handler.GetSourcePath())
	if err != nil {
		return files, err
	}
	return files, nil
}
func GenerateFile(file os.FileInfo, handler FileListHandler) error {
	var err error
	var info os.FileInfo
	info, err = os.Stat(handler.GetTargetPath())
	if err == nil && info.IsDir() == false {
		log.Fatalf("error `%s' is not a directory", handler.GetTargetPath())
	} else if err != nil {
		err = os.Mkdir(handler.GetTargetPath(), os.ModeDir|0700)
		if err != nil {
			return err
		}
	}
	err = handler.HandleFile(file)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func StripExtension(input string) string {
	var length = strings.LastIndexByte(input, '.')
	return input[:length]
}

func GetTypeName(name string) (string, error) {
	var result string
	if len(name) <= 4 {
		return "", fmt.Errorf("invalid input string %s", name)
	}
	var items = strings.Split(StripExtension(name), "-")
	for _, item := range items {
		result += strings.ToUpper(item[0:1]) + item[1:]
	}
	return strings.Replace(result, "Id", "ID", 1), nil
}

func GetFileName(name string) string {
	var result string
	re := regexp.MustCompile(`[A-Z][^A-Z]*`)
	subMatchAll := re.FindAllString(name, -1)
	result = strings.Join(subMatchAll, "-")
	return strings.ToLower(result)
}

func GetContent(dir string, name string) (string, error) {
	var pathString string
	var err error
	var file *os.File
	var info os.FileInfo
	var data []byte
	pathString = path.Join(dir, name)
	info, err = os.Stat(pathString)
	if err != nil {
		return "", err
	}
	file, err = os.Open(pathString)
	if err != nil {
		return "", err
	}
	defer func() {
		err = file.Close()
	}()
	data = make([]byte, info.Size())
	_, err = io.ReadFull(file, data)
	if err != nil {
		return "", err
	}
	// One extra byte or else, there's no content actually
	if len(data) >= 4 {
		var bom = binary.BigEndian.Uint32(data[:4])
		if bom&0xffffff00 == 0xefbbbf00 {
			return string(data[3:]), nil
		}
	}
	return string(data), nil
}
