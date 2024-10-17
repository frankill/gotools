package file

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SplitFileByLines 按文件行数拆分文件
// 参数:
//
//   - inputFile - 文件路径
//   - linesPerFile - 每个文件的行数
//   - addHeader - 是否添加文件头
//   - deleteSource - 是否删除源文件
//
// 返回:
//
//   - 错误信息
func SplitFileByLines(inputFile string, linesPerFile int, addHeader, deleteSource bool) error {
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
	header, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

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

	if addHeader {
		_, err = writer.WriteString(header)
		if err != nil {
			return err
		}
	}

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
			if addHeader {
				_, err = writer.WriteString(header)
				if err != nil {
					return err
				}
			}
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

// SplitFileBySize 按文件大小拆分文件
// 参数:
//
//   - inputFile - 文件路径
//   - chunkSizeMB - 每个文件的大小（单位 MB）
//   - addHeader - 是否添加文件头
//   - deleteSource - 是否删除源文件
//
// 返回:
//
//   - 错误信息
func SplitFileBySize(inputFile string, chunkSizeMB int, addHeader, deleteSource bool) error {
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
	header, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

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

	if addHeader {
		_, err = writer.WriteString(header)
		if err != nil {
			return err
		}
	}

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
			if addHeader {
				_, err = writer.WriteString(header)
				if err != nil {
					return err
				}
			}
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
