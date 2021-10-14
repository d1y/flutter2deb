package main

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type pubspecParse struct {
	Path string
	data []byte
	Pub  PubYamlModel
}

type PubYamlModel struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
}

func NewParsePubFile(path string) (*pubspecParse, error) {
	var pub = pubspecParse{
		Path: path,
	}
	e := pub.Load()
	return &pub, e
}

// Load load pub file
func (pub *pubspecParse) Load() error {
	var f, e = ioutil.ReadFile(pub.Path)
	if e != nil {
		return errors.New("load pubspec.yaml file faild")
	}
	y := PubYamlModel{}
	err := yaml.Unmarshal(f, &y)
	if err != nil {
		return errors.New("parse pubspec.yaml faild")
	}
	pub.Pub = y
	pub.data = f
	return nil
}
