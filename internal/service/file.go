package service

import (
	"encoding/json"
	"errors"
	"os"
)

type FileService interface {
	Save(gauges map[string]float64, counters map[string]int64) error
	Restore() (map[string]float64, map[string]int64, error)
}

type FileHanlder struct {
	Gauge   map[string]float64
	Counter map[string]int64
	path    string
}

func NewFileHanlder(path string) *FileHanlder {
	return &FileHanlder{path: path}
}

func (f *FileHanlder) Save(gauges map[string]float64, counters map[string]int64) error {
	if f.path == "" {
		return errors.New("can`t save to file. no file path specify")
	}
	f.Gauge = gauges
	f.Counter = counters
	data, err := json.MarshalIndent(f, "", "   ")
	if err != nil {
		return err
	}
	return os.WriteFile(f.path, data, 0666)
}

func (f *FileHanlder) Restore() (map[string]float64, map[string]int64, error) {
	if f.path == "" {
		return nil, nil, errors.New("can`t restore from file. no file path specify")
	}
	data, err := os.ReadFile(f.path)
	if err != nil {
		return nil, nil, err
	}
	if err := json.Unmarshal(data, f); err != nil {
		return nil, nil, err
	}
	return f.Gauge, f.Counter, nil
}
