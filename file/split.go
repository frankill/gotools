package file

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// 按行拆分文件
func SplitFileByLines(inputFile string, linesPerFile int, deleteSource bool) error {
	fileInfo, err := os.Stat(inputFile)
	if err != nil {
		return fmt.Errorf("file %s does not exist", inputFile)
	}

	fileSize := fileInfo.Size()
	if fileSize == 0 {
		fmt.Println("File is empty, no need to split.")
		return nil
	}

	fmt.Println("Splitting file by lines.")

	file, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	baseName := strings.TrimSuffix(inputFile, filepath.Ext(inputFile))
	ext := filepath.Ext(inputFile)
	chunkIndex := 1

	chunkFileName := fmt.Sprintf("%s_%d%s", baseName, chunkIndex, ext)
	chunkFile, err := os.Create(chunkFileName)
	if err != nil {
		return err
	}
	defer chunkFile.Close()

	writer := bufio.NewWriter(chunkFile)
	defer writer.Flush()

	lineCount := 0

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() != "EOF" {
				return err
			}
			break
		}

		_, err = writer.WriteString(line)
		if err != nil {
			return err
		}

		lineCount++
		if lineCount >= linesPerFile {
			writer.Flush()
			chunkFile.Close()

			chunkIndex++
			chunkFileName = fmt.Sprintf("%s_%d%s", baseName, chunkIndex, ext)
			chunkFile, err = os.Create(chunkFileName)
			if err != nil {
				return err
			}
			defer chunkFile.Close()

			writer = bufio.NewWriter(chunkFile)
			lineCount = 0
		}
	}

	writer.Flush()

	if deleteSource {
		err = os.Remove(inputFile)
		if err != nil {
			return fmt.Errorf("error removing original file: %v", err)
		}
		fmt.Println("File split completed successfully and original file removed.")
	} else {
		fmt.Println("File split completed successfully. Original file retained.")
	}

	return nil
}

// 按文件大小拆分文件
func SplitFileBySize(inputFile string, chunkSizeMB int, deleteSource bool) error {
	fileInfo, err := os.Stat(inputFile)
	if err != nil {
		return fmt.Errorf("file %s does not exist", inputFile)
	}

	fileSize := fileInfo.Size()
	if fileSize == 0 {
		fmt.Println("File is empty, no need to split.")
		return nil
	}

	chunkSizeBytes := int64(chunkSizeMB * 1024 * 1024)

	if fileSize <= chunkSizeBytes {
		fmt.Println("File size is not larger than the specified chunk size, no need to split.")
		return nil
	}

	fmt.Println("Splitting file by size.")

	file, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	baseName := strings.TrimSuffix(inputFile, filepath.Ext(inputFile))
	ext := filepath.Ext(inputFile)
	chunkIndex := 1

	chunkFileName := fmt.Sprintf("%s_%d%s", baseName, chunkIndex, ext)
	chunkFile, err := os.Create(chunkFileName)
	if err != nil {
		return err
	}
	defer chunkFile.Close()

	writer := bufio.NewWriter(chunkFile)
	defer writer.Flush()

	var currentSize int64

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() != "EOF" {
				return err
			}
			break
		}

		lineBytes := int64(len(line))
		if currentSize+lineBytes > chunkSizeBytes {
			writer.Flush()
			chunkFile.Close()

			chunkIndex++
			chunkFileName = fmt.Sprintf("%s_%d%s", baseName, chunkIndex, ext)
			chunkFile, err = os.Create(chunkFileName)
			if err != nil {
				return err
			}
			defer chunkFile.Close()

			writer = bufio.NewWriter(chunkFile)
			currentSize = 0
		}

		_, err = writer.WriteString(line)
		if err != nil {
			return err
		}
		currentSize += lineBytes
	}

	writer.Flush()

	if deleteSource {
		err = os.Remove(inputFile)
		if err != nil {
			return fmt.Errorf("error removing original file: %v", err)
		}
		fmt.Println("File split completed successfully and original file removed.")
	} else {
		fmt.Println("File split completed successfully. Original file retained.")
	}

	return nil
}
