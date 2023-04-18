package internal

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

var RE_FORMAT = [][]string{
	{`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z`, "2006-01-02T15:04:05Z"},
}

func GetDate(str string) (time.Time, error) {
	var date time.Time
	var err error
	for _, patternFormat := range RE_FORMAT {
		pattern := patternFormat[0]
		format := patternFormat[1]
		isoRe := regexp.MustCompile(pattern)
		isoBytes := isoRe.Find([]byte(str))
		if isoBytes == nil {
			err = errors.New("unable match ISO date")
		}
		if err == nil {
			iso := string(isoBytes)
			date, err = time.Parse(format, iso)
			break
		}
	}
	return date, err
}
