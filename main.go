// Create binary rpm package with ease
package main

import (
	"fmt"
	"github.com/mh-cbon/verbose"
	"github.com/urfave/cli"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// VERSION is the build version number.
var VERSION = "0.0.0"
var logger = verbose.Auto()

func main() {
	app := cli.NewApp()
	app.Name = "go-bin-rpm"
	app.Version = VERSION
	app.Usage = "Generate a binary rpm package"
	app.UsageText = "go-bin-rpm <cmd> <options>"
	app.Commands = []cli.Command{
		{
			Name:   "generate-spec",
			Usage:  "Generate the SPEC file",
			Action: generateSpec,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file, f",
					Value: "rpm.json",
					Usage: "Path to the rpm.json file",
				},
				cli.StringFlag{
					Name:   "a, arch",
					Value:  runtime.GOARCH,
					Usage:  "Target architecture of the build",
					EnvVar: "GOARCH",
				},
				cli.StringFlag{
					Name:  "version",
					Value: "0.0.0",
					Usage: "Target version of the build",
				},
			},
		},
		{
			Name:   "generate",
			Usage:  "Generate the package",
			Action: generatePkg,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file, f",
					Value: "rpm.json",
					Usage: "Path to the rpm.json file",
				},
				cli.StringFlag{
					Name:  "b, build-area",
					Value: "pkg-build",
					Usage: "Path to the build area",
				},
				cli.StringFlag{
					Name:   "a, arch",
					Value:  runtime.GOARCH,
					Usage:  "Target architecture of the build",
					EnvVar: "GOARCH",
				},
				cli.StringFlag{
					Name:  "o, output",
					Usage: "File path to the resulting rpm file",
				},
				cli.StringFlag{
					Name:  "version",
					Value: "0.0.0",
					Usage: "Target version of the build",
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
					Value: "rpm.json",
					Usage: "Path to the rpm.json file",
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func generateSpec(c *cli.Context) error {
	file := c.String("file")
	arch := c.String("arch")
	version := c.String("version")

	rpmJSON := Package{}

	if err := rpmJSON.Load(file); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if err := rpmJSON.Normalize(arch, version); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	spec, err := rpmJSON.GenerateSpecFile("")
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	fmt.Printf("%s", spec)

	return nil
}

func generatePkg(c *cli.Context) error {
	var err error

	file := c.String("file")
	arch := c.String("arch")
	version := c.String("version")
	buildArea := c.String("build-area")
	output := c.String("output")

	if output == "" {
		return cli.NewExitError("--output,-o argument is required", 1)
	} else {
		err := os.Mkdir(output, 0777)
		if err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}
	}

	rpmJSON := Package{}

	if err3 := rpmJSON.Load(file); err3 != nil {
		return cli.NewExitError(err3.Error(), 1)
	}

	if buildArea, err = filepath.Abs(buildArea); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if err2 := rpmJSON.Normalize(arch, version); err2 != nil {
		return cli.NewExitError(err2.Error(), 1)
	}

	err = rpmJSON.InitializeBuildArea(buildArea)
	if err != nil {
		log.Fatal(err)
	}

	if err = rpmJSON.WriteSpecFile("", buildArea); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if err = rpmJSON.RunBuild(buildArea, output); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Println("\n\nAll done!")

	return nil
}

func testPkg(c *cli.Context) error {
	file := c.String("file")

	rpmJSON := Package{}

	if err := rpmJSON.Load(file); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Println("File is correct")

	return nil
}
