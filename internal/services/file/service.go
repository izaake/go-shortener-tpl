package file

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/izaake/go-shortener-tpl/internal/models"
)

type Writer struct {
	file    *os.File
	encoder *json.Encoder
}

type Reader struct {
	file    *os.File
	scanner *bufio.Scanner
}

func NewWriter(fileName string) (*Writer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &Writer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *Writer) WriteEvent(event *models.User) error {
	return p.encoder.Encode(&event)
}

func (p *Writer) Close() error {
	return p.file.Close()
}

func NewReader(fileName string) (*Reader, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	return &Reader{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (c *Reader) ReadEvent() (*models.URL, error) {
	if !c.scanner.Scan() {
		return nil, c.scanner.Err()
	}
	data := c.scanner.Bytes()

	event := models.URL{}
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (c *Reader) Close() error {
	return c.file.Close()
}

func WriteToFile(fileName string, user *models.User) error {
	writer, err := NewWriter(fileName)
	if err != nil {
		return err
	}
	defer writer.Close()

	if err := writer.WriteEvent(user); err != nil {
		return err
	}
	return nil
}

func ReadLines(fileName string) ([]models.User, error) {
	reader, err := NewReader(fileName)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var Users []models.User
	for reader.scanner.Scan() {
		data := reader.scanner.Bytes()

		user := models.User{}
		err := json.Unmarshal(data, &user)
		if err != nil {
			return nil, err
		}

		Users = append(Users, user)
	}
	return Users, reader.scanner.Err()
}
