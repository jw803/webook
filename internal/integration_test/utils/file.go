package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/xuri/excelize/v2"
)

func ReadJSONFile[T any](filePath string) (T, error) {
	var data T
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func ReadExcelFile(filePath, sheetName string) ([][]string, error) {
	rows := make([][]string, 0)
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		fmt.Println("无法打开文件:", err)
		return rows, err
	}
	rows, err = f.GetRows(sheetName)
	return rows, err
}

func RemoveTestCaseFilesByExtension(fileExtension string) error {
	files, err := filepath.Glob(fmt.Sprintf("./*.%s", fileExtension))
	if err != nil {
		return err
	}
	for _, f := range files {
		if err = os.Remove(f); err != nil {
			return err
		}
	}
	return nil
}
