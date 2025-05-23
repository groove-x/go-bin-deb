// Package go-bin-deb creates binary package for debian system
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mh-cbon/verbose"
	"github.com/urfave/cli"

	"github.com/groove-x/go-bin-deb/debian"
)

// VERSION is the last build number.
var VERSION = "0.0.0"
var logger = verbose.Auto()

func main() {
	app := cli.NewApp()
	app.Name = "go-bin-deb"
	app.Version = VERSION
	app.Usage = "Generate a binary debian package"
	app.UsageText = "go-bin-deb <cmd> <options>"
	app.Commands = []cli.Command{
		{
			Name:   "generate",
			Usage:  "Generate the contents of the package",
			Action: generateContents,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "wd, w",
					Value: "pkg-build",
					Usage: "Working directory to prepare the package",
				},
				cli.StringFlag{
					Name:  "output, o",
					Value: "",
					Usage: "Output directory for the debian package files",
				},
				cli.StringFlag{
					Name:  "file, f",
					Value: "deb.json",
					Usage: "Path to the deb.json,deb.yml,deb.yaml file",
				},
				cli.StringFlag{
					Name:  "version",
					Value: "",
					Usage: "Version of the package",
				},
				cli.StringFlag{
					Name:  "arch, a",
					Value: "",
					Usage: "Arch of the package",
				},
			},
		},
		{
			Name:   "test",
			Usage:  "Test the package json file",
			Action: testPkg,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file, f",
					Value: "deb.json",
					Usage: "Path to the deb.json,deb.yml,deb.yaml file",
				},
			},
		},
	}

	app.Run(os.Args)
}

func generateContents(c *cli.Context) error {
	output := c.String("output")
	wd := c.String("wd")
	file := c.String("file")
	version := c.String("version")
	arch := c.String("arch")

	pkgDir := filepath.Join(wd)

	if o, err := filepath.Abs(output); err != nil {
		return cli.NewExitError(err.Error(), 1)
	} else {
		output = o
	}

	// Check if the file has a supported extension
	if ext := filepath.Ext(file); ext != ".json" && ext != ".yaml" && ext != ".yml" {
		return cli.NewExitError(fmt.Sprintf("unsupported file format: %s. only .json, .yaml, and .yml are supported.", ext), 1)
	}

	packageInfo := debian.Package{}

	// load the deb.json,deb.yml,deb.yaml file
	if err := packageInfo.Load(file); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	logger.Println(fmt.Sprintf("%s loaded", file))

	// normalize data
	packageInfo.Normalize(pkgDir, version, arch)
	logger.Println("pkg data normalized")

	logger.Printf("Generating files in %s", pkgDir)
	if err := packageInfo.GenerateFiles(filepath.Dir(file), pkgDir); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	logger.Printf("Building package in %s to %s", wd, output)
	if err := buildPackage(pkgDir, output); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	logger.Printf("linting package in %s to %s", wd, output)
	lintPackage(pkgDir, output) // it does not need to fail.

	return nil
}

func testPkg(c *cli.Context) error {
	file := c.String("file")

	// Check if the file has a supported extension
	if ext := filepath.Ext(file); ext != ".json" && ext != ".yaml" && ext != ".yml" {
		return cli.NewExitError(fmt.Sprintf("unsupported file format: %s. only .json, .yaml, and .yml are supported.", ext), 1)
	}

	packageInfo := debian.Package{}

	if err := packageInfo.Load(file); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Println("File is correct")

	return nil
}

func buildPackage(wd string, output string) error {
	oCmd := exec.Command("fakeroot", "dpkg-deb", "--build", "debian", output)
	oCmd.Dir = wd
	oCmd.Stdout = os.Stdout
	oCmd.Stderr = os.Stderr
	return oCmd.Run()
}

func lintPackage(wd string, output string) error {
	oCmd := exec.Command("lintian", output)
	oCmd.Dir = wd
	oCmd.Stdout = os.Stdout
	oCmd.Stderr = os.Stderr
	return oCmd.Run()
}
