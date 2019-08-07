package main

import (
	"encoding/json"
	"fmt"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"os"
)

type depProjects map[string]dependencyInfo

// `Revision` is always present on a project stanza. Additionally, either `branch` or `version` can be present.
// Check the [dep docs](https://golang.github.io/dep/docs/Gopkg.lock.html#projects) for full reference.
type dependencyInfo struct {
	Name     string `json:"name"`
	Version  string `json:"version,omitempty"`
	Revision string `json:"revision"`
	Branch   string `json:"branch,omitempty"`
}

type dependencyInfoPair struct {
	Plugin dependencyInfo `json:"pluginDependencies"`
	GlooE  dependencyInfo `json:"glooeDependencies"`
}

func (d dependencyInfo) matches(that dependencyInfo) bool {
	// `Revision` is the ultimate source of truth, `version` or `branch` are potentially floating references
	if d.Revision == that.Revision {
		return true
	}
	return false
}

// TODO(marco): would be nice if this script would generate the snippets that the user needs to add to their Gopkg.toml
func main() {

	pluginDependencies := getPluginDependencies()
	glooeDependencies := getGlooeDependencies()

	var nonMatchingDeps []dependencyInfoPair
	for name, depInfo := range pluginDependencies {

		// Just check libraries that are shared with GlooE
		if glooeEquivalent, ok := glooeDependencies[name]; ok {
			if !glooeEquivalent.matches(depInfo) {
				nonMatchingDeps = append(nonMatchingDeps, dependencyInfoPair{
					Plugin: depInfo,
					GlooE:  glooeEquivalent,
				})
			}
		}
	}

	if len(nonMatchingDeps) > 0 {

		reportBytes, err := json.Marshal(nonMatchingDeps)
		if err != nil {
			panic(err)
		}

		if err := ioutil.WriteFile("mismatched_dependencies.json", reportBytes, 0644); err != nil {
			panic(err)
		}

		os.Exit(1)
	}

	os.Exit(0)
}

func getPluginDependencies() depProjects {
	pluginDependencies, err := parseGoPkgLock("Gopkg.lock")
	if err != nil {
		panic(err)
	}
	return collectDependencyInfo(pluginDependencies)
}

func getGlooeDependencies() depProjects {
	pluginDependencies, err := parseGoPkgLock("glooe/Gopkg.lock")
	if err != nil {
		panic(err)
	}
	return collectDependencyInfo(pluginDependencies)
}

func parseGoPkgLock(path string) ([]*toml.Tree, error) {
	config, err := toml.LoadFile(path)
	if err != nil {
		return nil, err
	}

	tomlTree := config.Get("projects")

	switch typedTree := tomlTree.(type) {
	case []*toml.Tree:
		return typedTree, nil
	default:
		return nil, fmt.Errorf("unable to parse toml tree")
	}
}

func collectDependencyInfo(deps []*toml.Tree) depProjects {
	dependencies := make(depProjects)

	for _, t := range deps {

		name := t.Get("name").(string)

		info := dependencyInfo{
			Name:     name,
			Revision: t.Get("revision").(string),
		}

		if version, ok := t.Get("version").(string); ok {
			info.Version = version
		}
		if branch, ok := t.Get("branch").(string); ok {
			info.Branch = branch
		}

		dependencies[name] = info
	}
	return dependencies
}
