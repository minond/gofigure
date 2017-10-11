package gofigure

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Parse func([]byte, interface{}) error

type Parser struct {
	Extensions []string
	Unmarshal  Parse
}

type FileLocation struct {
	Path   string
	Parser string
}

var (
	directories = []string{
		"./",
		"./config/",
	}

	variants = []string{
		"",
	}

	parsers = map[string]*Parser{
		"yaml": &Parser{
			Extensions: []string{".yml", ".yaml"},
			Unmarshal:  yaml.Unmarshal,
		},
	}
)

func AddVariants(list ...string) {
	variants = append(variants, list...)
}

func AddDirectories(list ...string) {
	directories = append(directories, list...)
}

func AddParser(name string, extensions []string, handler Parse) {
	parsers[name] = &Parser{
		Extensions: extensions,
		Unmarshal:  handler,
	}
}

func GetConfigurationFileContents(name string) ([]byte, error) {
	return ioutil.ReadFile(name)
}

func LocateConfigurationFile(configFileName string, activeVariants []string) *FileLocation {
	activeVariants = append(activeVariants, variants...)

	for _, currDir := range directories {
		dir, err := ioutil.ReadDir(currDir)

		if err != nil {
			continue
		}

		for _, file := range dir {
			for parserName, parser := range parsers {
				for _, ext := range parser.Extensions {
					currFile := currDir + file.Name()

					for _, variant := range activeVariants {
						var proposedFile string

						if len(variant) == 0 {
							proposedFile = currDir + configFileName + ext
						} else {
							proposedFile = currDir + configFileName + "." + variant + ext
						}

						if proposedFile == currFile {
							return &FileLocation{proposedFile, parserName}
						}
					}
				}
			}
		}
	}

	return nil
}

func Load(label string, store interface{}) {
	LoadWithVariants(label, []string{}, store)
}

func LoadWithVariants(label string, variants []string, store interface{}) {
	floc := LocateConfigurationFile(label, variants)

	if floc != nil {
		file, err := GetConfigurationFileContents(floc.Path)

		if err == nil {
			parsers[floc.Parser].Unmarshal(file, store)
		}
	}
}
