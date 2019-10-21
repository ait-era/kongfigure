package kongfigure

import (
	"encoding/json"
	"fmt"
	"gopkg.in/resty.v1"
	"io/ioutil"
	"log"
	"path"
)

type KongData struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type KongResourceConfig struct {
	Files    [][]byte
	Resource string
}

func (k *KongResourceConfig) loadFiles(confPath string) error {
	resourcePath := path.Join(confPath, k.Resource)

	resourceFilenames, err := listJsonFiles(resourcePath)
	if err != nil {
		return err
	}
	for _, f := range resourceFilenames {
		file, _ := ioutil.ReadFile(path.Join(confPath, k.Resource, f))
		k.Files = append(k.Files, file)
	}
	return nil
}

func (k *KongResourceConfig) loadFile(confPath string) error {
	file, err := ioutil.ReadFile(path.Join(confPath, fmt.Sprintf("%s.json", k.Resource)))
	if err != nil {
		return err
	}
	k.Files = append(k.Files, file)

	return nil
}

func (k *KongResourceConfig) apply(restyClient *resty.Client, kongSettings AppSettings) error {
	for _, fileData := range k.Files {
		kongData := KongData{}
		_ = json.Unmarshal(fileData, &kongData)

		response, err := restyClient.R().Get(fmt.Sprintf("%s/%s", k.Resource, kongData.Id))
		if err != nil {
			return err
		}

		if response.IsError() && response.StatusCode() != 404 {
			return &KongfigureHttpError{fmt.Sprintf("Failed to fetch status for resource `%v`", kongData), response.Request.URL, string(response.Body())}
		}

		var action string
		if response.StatusCode() == 404 {
			action = "post"
		} else {
			action = "patch"
		}

		if kongSettings.DryRun {
			log.Printf("Kongfigure would %s resource %s with id %s", action, k.Resource, kongData.Id)
		} else {
			err = pushResource(action, k.Resource, restyClient, &kongData, fileData)
			if err != nil {
				return err
			}
			log.Printf("Successfully applied %v resource %v", k.Resource, kongData)
		}
	}
	return nil
}
