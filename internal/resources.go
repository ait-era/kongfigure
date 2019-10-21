package kongfigure

import (
	"bytes"
	"fmt"
	"gopkg.in/resty.v1"
	"io/ioutil"
	"strings"
)

func ApplyResources(resourceName string, kongSettings AppSettings, restyClient *resty.Client) error {
	configs := KongResourceConfig{Resource: resourceName}
	if err := configs.loadFiles(kongSettings.KongConfPath); err != nil {
		return err
	}
	if err := configs.apply(restyClient, kongSettings); err != nil {
		return err
	}
	return nil

}

func ApplyFileResource(resourceName string, kongSettings AppSettings, restyClient *resty.Client) error {
	configs := KongResourceConfig{Resource: resourceName}

	if err := configs.loadFile(kongSettings.KongConfPath); err != nil {
		return err
	}
	if err := configs.apply(restyClient, kongSettings); err != nil {
		return err
	}
	return nil
}

func pushResource(action string, resourcePath string, restyClient *resty.Client, kongData *KongData, fileData []byte) error {

	var err error
	var response *resty.Response

	if action == "post" {
		response, err = restyClient.R().SetBody(bytes.NewBuffer(fileData)).Post(resourcePath)
	} else if action == "patch" {

		response, err = restyClient.R().SetBody(bytes.NewBuffer(fileData)).Patch(fmt.Sprintf("%s/%s", resourcePath, kongData.Id))
	}
	if err != nil {
		return err
	}

	if response.StatusCode() >= 400 {
		return &KongfigureHttpError{fmt.Sprintf("Failed to apply resource `%v`", kongData), response.Request.URL, string(response.Body())}
	}

	return nil
}

func listJsonFiles(directory string) ([]string, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	var jsonFiles []string
	for _, f := range files {
		if f.IsDir() != true && strings.HasSuffix(f.Name(), ".json") {
			jsonFiles = append(jsonFiles, f.Name())
		}
	}
	return jsonFiles, nil
}
