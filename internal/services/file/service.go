package file

import (
	"bufio"
	"encoding/json"
	"log"
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

func (p *Writer) WriteEvent(event *models.URL) error {
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

func WriteToFile(fileName string, url *models.URL) {
	writer, err := NewWriter(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer writer.Close()

	if err := writer.WriteEvent(url); err != nil {
		log.Fatal(err)
	}
}

func ReadLines(fileName string) ([]models.URL, error) {
	reader, err := NewReader(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	var URLs []models.URL
	for reader.scanner.Scan() {
		data := reader.scanner.Bytes()

		url := models.URL{}
		err := json.Unmarshal(data, &url)
		if err != nil {
			return nil, err
		}

		URLs = append(URLs, url)
	}
	return URLs, reader.scanner.Err()
}
