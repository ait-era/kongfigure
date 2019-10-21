package kongfigure

import (
	"fmt"
	"github.com/urfave/cli"
	"gopkg.in/resty.v1"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

func Run(c *cli.Context, settings AppSettings) error {
	if settings.KongConfPath == "" {
		log.Fatal("Missing `kong-configs`")
	}

	servicesExist := true
	if _, err := os.Stat(path.Join(settings.KongConfPath, "services")); err != nil {
		if settings.DryRun {
			log.Print("No Service resource")
		}
		servicesExist = false
	}

	routesExist := true
	if _, err := os.Stat(path.Join(settings.KongConfPath, "routes")); err != nil {
		if settings.DryRun {
			log.Print("No Route resource")
		}
		routesExist = false
	}

	consumersExist := true
	if _, err := os.Stat(path.Join(settings.KongConfPath, "consumers")); err != nil {
		if settings.DryRun {
			log.Print("No Consumer resource")
		}
		consumersExist = false
	}

	pluginsExist := true
	if _, err := os.Stat(path.Join(settings.KongConfPath, "plugins")); err != nil {
		if settings.DryRun {
			log.Print("No Plugin resource")
		}
		pluginsExist = false
	}

	restyClient := *resty.New()
	restyClient.SetHostURL(settings.KongUrl)

	if servicesExist == true {
		if err := ApplyResources("services", settings, &restyClient); err != nil {
			log.Fatal(err)
		}
	}

	if routesExist == true {
		if err := ApplyResources("routes", settings, &restyClient); err != nil {
			log.Fatal(err)
		}
	}

	if consumersExist == true {
		if err := ApplyResources("consumers", settings, &restyClient); err != nil {
			log.Fatal(err)
		}
	}

	if pluginsExist == true {
		if err := ApplyResources("plugins", settings, &restyClient); err != nil {
			log.Fatal(err)
		}
	}

	if consumersExist == true {
		consumerFiles, err := ioutil.ReadDir(path.Join(settings.KongConfPath, "consumers"))
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range consumerFiles {
			if f.IsDir() == true {
				credentialFiles, err := ioutil.ReadDir(path.Join(settings.KongConfPath, fmt.Sprintf("consumers/%s", f.Name())))
				if err != nil {
					log.Fatal(err)
				}

				for _, cf := range credentialFiles {
					if cf.IsDir() != true && strings.HasSuffix(cf.Name(), ".json") {
						resourceNameSlice := strings.Split(cf.Name(), ".")
						resourceName := resourceNameSlice[len(resourceNameSlice)-2]
						if err := ApplyFileResource(fmt.Sprintf("consumers/%s/%s", f.Name(), resourceName), settings, &restyClient); err != nil {
							log.Fatal(err)
						}
					}
				}
			}
		}
	}

	return nil
}
