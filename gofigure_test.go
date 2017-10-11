package gofigure

import (
	"io/ioutil"
	"os"
	"testing"
)

func createTempFile(dir string) (*os.File, func()) {
	file, _ := ioutil.TempFile(dir, "configurationtesting")

	return file, func() {
		os.Remove(file.Name())
	}
}

func createFile(dir, name string) (*os.File, func()) {
	file, _ := createTempFile(dir)
	os.Rename(file.Name(), name)

	return file, func() {
		os.Remove(name)
	}
}

func assertMissingLocation(t *testing.T, returnedLoc *FileLocation) {
	if returnedLoc != nil {
		t.Errorf("Foung configuration file that should not exists: %v", returnedLoc)
	}
}

func assertLocation(t *testing.T, returnedLoc, expectedLoc *FileLocation) {
	if returnedLoc == nil {
		t.Errorf("Expecting to find '%v' but did not", expectedLoc)
	} else if returnedLoc.Path != expectedLoc.Path {
		t.Errorf("Expecting '%v' as the Path but found '%v' instead",
			expectedLoc.Path, returnedLoc.Path)
	} else if returnedLoc.Parser != expectedLoc.Parser {
		t.Errorf("Expecting '%v' as the Parser but found '%v' instead",
			expectedLoc.Parser, returnedLoc.Parser)
	}
}

func TestMissingConfigurationFileLocation(t *testing.T) {
	loc := LocateConfigurationFile("testing", []string{})
	assertMissingLocation(t, loc)
}

func TestConfigurationFileLocation(t *testing.T) {
	_, rm := createFile(".", "testing.yml")
	defer rm()

	loc := LocateConfigurationFile("testing", []string{})
	assertLocation(t, loc, &FileLocation{
		Path:   "./testing.yml",
		Parser: "yaml",
	})
}

func TestConfigurationFileLocationWithASingleVariant(t *testing.T) {
	_, rm := createFile(".", "testing.one.yml")
	defer rm()

	loc := LocateConfigurationFile("testing", []string{"one"})

	assertLocation(t, loc, &FileLocation{
		Path:   "./testing.one.yml",
		Parser: "yaml",
	})
}

func TestConfigurationFileLocationWithASingleVariantButAMissingMatch(t *testing.T) {
	_, rm := createFile(".", "testing.two.yml")
	defer rm()

	loc := LocateConfigurationFile("testing", []string{"one"})
	assertMissingLocation(t, loc)
}

func TestConfigurationFileLocationWithMultipleVariants(t *testing.T) {
	_, rm := createFile(".", "testing.two.yml")
	defer rm()

	loc := LocateConfigurationFile("testing", []string{"one", "two"})

	assertLocation(t, loc, &FileLocation{
		Path:   "./testing.two.yml",
		Parser: "yaml",
	})
}
