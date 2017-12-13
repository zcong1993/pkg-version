package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/zcong1993/utils"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	// Name is cli name
	Name = "pv"
	// Version is cli current version
	Version = "v0.0.5"
	// Date is build date
	Date = ""
)

// GitCommit is cli current git commit hash
var GitCommit string

// TableData is data type for tablewriter
type TableData [][]string

// Package is json struct for package.json
type Package struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

// ApiResp is json struct for npm search response
type ApiResp struct {
	Name string            `json:"name"`
	Tags map[string]string `json:"tags"`
}

func errPanic(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func searchPkg(name string) (*ApiResp, error) {
	apiUrl := fmt.Sprintf("https://ofcncog2cu-dsn.algolia.net/1/indexes/npm-search/%s?x-algolia-agent=Algolia for vanilla JavaScript (lite) 3.24.6&x-algolia-application-id=OFCNCOG2CU&x-algolia-api-key=f54e21fa3a2a0160595bb058179bfb1e", name)
	var apiResp ApiResp
	err := utils.GetJSON(apiUrl, &apiResp)
	if err != nil {
		return nil, err
	}
	return &apiResp, nil
}

func render(data ...TableData) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Local", "Latest", "Next"})
	cyan := color.New(color.FgCyan)
	red := color.New(color.FgRed)
	for _, td := range data {
		for _, v := range td {
			local := strings.Replace(v[1], "^", "", 1)
			if local != v[2] && local != v[3] {
				v = []string{cyan.Sprintf(v[0]), v[1], red.Sprintf(v[2]), red.Sprintf(v[3])}
			}
			table.Append(v)
		}
	}
	table.Render()
}

func searchAll(pkgs map[string]string) TableData {
	var tableData [][]string
	for k, v := range pkgs {
		resp, err := searchPkg(k)
		if err != nil {
			// ignore
			continue
		}
		data := []string{resp.Name, v, resp.Tags["latest"], resp.Tags["next"]}
		tableData = append(tableData, data)
	}
	return tableData
}

// ShowVersion is handler for version command
func ShowVersion() {
	version := fmt.Sprintf("\n%s version %s", Name, Version)
	if len(GitCommit) != 0 {
		version += fmt.Sprintf(" (%s)", GitCommit)
	}
	if len(Date) != 0 {
		version += fmt.Sprintf(" (%v)", Date)
	}
	fmt.Println(version)
}

func run() {
	var (
		help    bool
		version bool
		path    string
		pkg     Package
	)
	flag.BoolVar(&help, "help", false, "show help")
	flag.BoolVar(&help, "h", false, "show help(short)")

	flag.BoolVar(&version, "version", false, "show version")
	flag.BoolVar(&version, "v", false, "show version(short)")

	flag.StringVar(&path, "path", "package.json", "package.json path")
	flag.StringVar(&path, "p", "package.json", "package.json path(short)")

	flag.Parse()

	if help {
		fmt.Println(helpText)
		return
	}

	if version {
		ShowVersion()
		return
	}

	js, err := ioutil.ReadFile(path)
	errPanic(err)
	err = json.Unmarshal(js, &pkg)
	errPanic(err)
	deps := searchAll(pkg.Dependencies)
	devDeps := searchAll(pkg.DevDependencies)
	render(deps, devDeps)
}

func main() {
	run()
}

var helpText = `
pkg-version is a tool helping to check local npm package version
 Usage:
 	pv [options]
 Options:
	-path, -p                      Path of 'package.json', default is '$cwd/package.json'
 Example:
 	pv -path="package.json"
`
