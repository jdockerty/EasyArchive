package zipdir

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// PathSeparator returns the appropriate separator, either / or \, depending on Windows or Linux OS.
func PathSeparator() string {
	return string(filepath.Separator)
}

// ZipFiles compresses one or many files into a single zip archive file.
// Param 1: filename is the output zip file's name.
// Param 2: files is a list of files to add to the zip.
// Param 3: Location of the archive to attach onto the param 2.
func zipUp(filename string, files []string, archivePath string) error {

	newZipFile, err := os.Create("archives" + PathSeparator() + filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)

	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = addFileToZip(zipWriter, file, archivePath); err != nil {
			return err
		}
	}
	return nil
}

func addFileToZip(zipWriter *zip.Writer, filename string, archivePath string) error {
	fileString := fmt.Sprintf("%s%s%s", archivePath, PathSeparator(), filename)
	fileToZip, err := os.Open(fileString)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = filename

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

func timeStamp() string {
	now := time.Now()
	currentDate := fmt.Sprintf("%02d-%02d-%d", now.Day(), now.Month(), now.Year())
	return currentDate
}

// ZipFiles will compress and zip the files contained within the set archive path from config.json.
func ZipFiles(files []string, archivePath string) string {
	outputFile := fmt.Sprintf("%s.zip", timeStamp())

	if err := zipUp(outputFile, files, archivePath); err != nil {
		panic(err)
	}
	log.Println("Zip file successfully created:", outputFile)

	return outputFile
}
