package main

import (
	"errors"
	"flag"
	"fmt"
	"goFolderCleaner/internal"
	"os"
	"path"
	"sort"
	"time"

	"github.com/dustin/go-humanize"
)

type Backup struct {
	name string
	path string
	size uint64
	time time.Time
}

func main() {
	var place, maxSizeStr string
	var minCount, maxCount, maxDays int
	var useModTime, isDryRun, isVerbose bool

	flag.StringVar(&place, "place", "", "Backup place")
	flag.StringVar(&maxSizeStr, "maxSize", "", "Place max size")
	flag.IntVar(&minCount, "minCount", 3, "Min file count")
	flag.IntVar(&maxCount, "maxCount", 0, "Max file count")
	flag.IntVar(&maxDays, "maxDays", 0, "Max day count")
	flag.BoolVar(&useModTime, "useModTime", false, "Use modification time")
	flag.BoolVar(&isDryRun, "dry", false, "Dry run")
	flag.BoolVar(&isVerbose, "v", false, "Verbose")
	flag.Parse()

	var maxSize uint64
	if maxSizeStr != "" {
		s, err := humanize.ParseBytes(maxSizeStr)
		if err != nil {
			panic(err)
		}
		maxSize = s
	}
	if maxSize == 0 && maxCount == 0 && maxDays == 0 {
		panic("maxSize or maxCount or maxDays should be set")
	}

	placeFile, err := os.Open(place)
	if err != nil {
		panic(err)
	}
	defer placeFile.Close()

	stat, err := placeFile.Stat()
	if err == nil && !stat.IsDir() {
		err = errors.New("Place is not dir")
	}
	if err != nil {
		panic(err)
	}

	files, err := placeFile.ReadDir(-1)
	if err != nil {
		panic(err)
	}

	var backups []Backup
	for _, file := range files {
		name := file.Name()
		filePath := path.Join(place, name)

		info, err := file.Info()
		if err != nil {
			fmt.Printf("Unable get file info `%s`, skip: %v", name, err)
			continue
		}

		var date time.Time
		if useModTime {
			date = info.ModTime()
		} else {
			d, err := internal.GetDate(name)
			if err != nil {
				fmt.Printf("Unable get date `%s`, skip: %v\n", name, err)
				continue
			}
			date = d
		}

		var size int64
		if file.IsDir() {
			dirSize, err := internal.DirSize(filePath)
			if err != nil {
				fmt.Printf("Unable get dir size `%s`: %v, skip\n", name, err)
				continue
			}
			size = dirSize
		} else {
			size = info.Size()
		}

		backup := Backup{
			name: name,
			path: filePath,
			size: uint64(size),
			time: date,
		}
		backups = append(backups, backup)
	}

	sort.Slice(backups, func(i, j int) bool {
		return backups[i].time.After(backups[j].time)
	})

	dayMap := make(map[string]bool)
	newDayMap := make(map[string]bool)
	var sum uint64
	var newSize uint64
	var newCount int
	for idx, backup := range backups {
		if isVerbose {
			fmt.Printf("%s: %s\n", backup.name, humanize.Bytes(backup.size))
		}
		sum += backup.size
		date := backup.time.Format("2006-01-02")
		dayMap[date] = true
		isSizeLimit := maxSize > 0 && sum > maxSize
		isCountLimit := maxCount > 0 && idx >= maxCount
		isDaysLimit := maxDays > 0 && len(dayMap) > maxDays
		if idx >= minCount && (isSizeLimit || isCountLimit || isDaysLimit) {
			var reason string
			if isSizeLimit {
				reason = "Size limit"
			}
			if isCountLimit {
				reason = "Count limit"
			}
			if isDaysLimit {
				reason = "Days limit"
			}
			fmt.Printf("Delete %s, cause: %s\n", backup.name, reason)
			if !isDryRun {
				if err := os.RemoveAll(backup.path); err != nil {
					fmt.Printf("Unable remove %s: %v\n", backup.path, err)
				}
			}
		} else {
			newSize += backup.size
			newCount += 1
			newDayMap[date] = true
		}
	}

	fmt.Printf("%s: size: %s/%s, count: %v/%v, days: %v/%v\n", place,
		humanize.Bytes(newSize), humanize.Bytes(maxSize),
		newCount, maxCount,
		len(newDayMap), maxDays,
	)
}
