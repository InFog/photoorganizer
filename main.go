package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/InFog/photoorganizer/utils"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

type pictureData struct {
	Year      int
	Month     int
	Day       int
	Hour      int
	Minute    int
	Second    int
	Latitude  string
	Longitude string
}

var (
	keepFiles bool = true
)

func output(s ...interface{}) {
	fmt.Println(s...)
}

func parseInput() ([]string, error) {
	var filesToMove []string

	if len(os.Args) < 2 {
		return filesToMove, errors.New("Missing parameters")
	}

	for i := 1; i < len(os.Args); i++ {
		if utils.IsFile(os.Args[i]) {
			filesToMove = append(filesToMove, os.Args[i])
		}
	}

	return filesToMove, nil
}

func getPicturesData(file *os.File) *pictureData {
	exif.RegisterParsers(mknote.All...)
	exifData, err := exif.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	d, _ := exifData.DateTime()
	year, month, day := d.Date()
	hour := d.Hour()
	minute := d.Minute()
	second := d.Second()

	latitude, longitude, _ := exifData.LatLong()

	return &pictureData{
		Year:      int(year),
		Month:     int(month),
		Day:       int(day),
		Hour:      int(hour),
		Minute:    int(minute),
		Second:    int(second),
		Latitude:  fmt.Sprintf("%f", latitude),
		Longitude: fmt.Sprintf("%f", longitude),
	}
}

func getMonthName(month int) string {
	monthName := "Unknown"

	switch month {
	case 1:
		monthName = "January"
	case 2:
		monthName = "February"
	case 3:
		monthName = "March"
	case 4:
		monthName = "April"
	case 5:
		monthName = "May"
	case 6:
		monthName = "June"
	case 7:
		monthName = "July"
	case 8:
		monthName = "August"
	case 9:
		monthName = "September"
	case 10:
		monthName = "October"
	case 11:
		monthName = "November"
	case 12:
		monthName = "December"
	}

	return monthName
}

func getFinalDestinationFile(destination string, i int) string {
	if utils.IsFile(destination) {
		if i == 1 {
			destination = destination[0:len(destination)-4] + "_" + strconv.Itoa(i) + destination[len(destination)-4:]
		} else if i > 1 {
			destination = destination[0:len(destination)-6] + "_" + strconv.Itoa(i) + destination[len(destination)-4:]
		}

		destination = getFinalDestinationFile(destination, i+1)
	}

	return destination
}

func moveFile(origin string, destFilename string, destDirectory string) {
	if _, err := os.Stat(destDirectory); err != nil {
		err := os.MkdirAll(destDirectory, 0777)
		if err != nil {
			log.Fatal(fmt.Sprintf("It was not possible to create the directory: '%s'", destDirectory))
			return
		}
	}

	destCompleteFilename := getFinalDestinationFile(destDirectory+destFilename, 0)
	dest, err := os.Create(destCompleteFilename)
	if err != nil {
		log.Fatal("It was not possible to create the destination file")
	}
	defer dest.Close()

	orig, _ := os.Open(origin)
	defer orig.Close()

	if keepFiles {
		io.Copy(dest, orig)
	}
}

func moveImage(fileName string) {
	fileExtension := fileName[len(fileName)-3:]

	if !utils.IsFile(fileName) {
		log.Fatal(errors.New(fmt.Sprintf("The file '%s' does not exist", fileName)))
		return
	}

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	pData := getPicturesData(file)
	monthName := getMonthName(pData.Month)
	destDirectory := fmt.Sprintf("/tmp/dest/%d/%s/%d/", pData.Year, monthName, pData.Day)
	destFilename := fmt.Sprintf("%02d_%02d_%02d.%s", pData.Hour, pData.Minute, pData.Second, fileExtension)

	moveFile(fileName, destFilename, destDirectory)

	/** TODO check if GPS information is available */

	/** Destination should be DEST/YYYY/MM/DD/HH_MM_SS.EXT when no location available */
	/** Destination should be DEST/Country/City/YYYY/MM/DD/HH_MM_SS.EXT when no location available */
}

func main() {
	output("Starting")

	fileNames, err := parseInput()
	if err != nil {
		log.Fatal(err)
		return
	}

	for i := 0; i < len(fileNames); i++ {
		moveImage(fileNames[i])
	}

	output(fmt.Sprintf("Moved %d files", len(fileNames)))
}
