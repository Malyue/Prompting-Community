package storage

import (
	"go.uber.org/zap"
	"io"
	"os"
	"path"
)

/*
	Author:Malyue
	Description:本地存储
	CreatedAt:2023/6/7
*/

type LocalStorage struct {
	RootPath string
}

func NewLocalStorage(rootPath string) *LocalStorage {
	return &LocalStorage{
		RootPath: rootPath,
	}
}

func (s *LocalStorage) MakeBucket(bucketName string) error {
	dirName := path.Join(s.RootPath, bucketName)
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		if err := os.Mkdir(dirName, 0755); err != nil {
			return err
		}
	}
	return nil
}

/*
@Params:
bucketName 存储桶名称
objectName 对象名
offset 偏移量
length 读取长度
*/
func (s *LocalStorage) GetObject(bucketName, objectName string, offset, length int64) ([]byte, error) {
	objectPath := path.Join(s.RootPath, bucketName, objectName)
	file, err := os.Open(objectPath)
	if err != nil {
		zap.L().Info("Failed to open file:" + err.Error())
		return nil, err
	}
	defer file.Close()
	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		zap.L().Info("Seek File Error:" + err.Error())
		return nil, err
	}
	buffer := make([]byte, length)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return buffer, nil
}

func (s *LocalStorage) PutObject(bucketName, objectName string, fileContent *[]byte, contentType string) error {
	//sourceFile, err := os.Open(filePath)
	//if err != nil {
	//	return err
	//}
	//defer sourceFile.Close()

	objectPath := path.Join(s.RootPath, bucketName, objectName)
	file, err := os.Create(objectPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(*fileContent)
	//err = os.WriteFile(objectPath, *fileContent, 0644)
	if err != nil {
		return err
	}

	return nil
}
