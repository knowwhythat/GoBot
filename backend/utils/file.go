package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

func SaveFile(src io.Reader, dst string) error {
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			//
		}
	}(out)

	_, err = io.Copy(out, src)
	return err
}

func GetFileMD5(file io.Reader) string {
	h := md5.New()
	_, err := io.Copy(h, file)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}

func PathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func ContainsOnlyOneDir(path string) string {
	fileList, err := os.ReadDir(path)
	if err == nil && len(fileList) == 1 && fileList[0].IsDir() {
		return fileList[0].Name()
	}
	return ""
}
