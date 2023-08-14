package modmanager

import (
	"fmt"
	"mtui/api/cdb"
	"strings"
)

type ResolvedDependency struct {
	Name   string `json:"name"`
	Author string `json:"author"`
}

func resolve_deps(cli *cdb.CDBClient, game_map map[string]*bool, author, name string) ([]*ResolvedDependency, error) {
	harddeps := []*ResolvedDependency{}

	deps, err := cli.GetDependencies(author, name)
	if err != nil {
		return nil, fmt.Errorf("error resolving deps for %s/%s: %v", author, name, err)
	}

	key := fmt.Sprintf("%s/%s", author, name)
	dep := deps[key]
	if dep == nil {
		return nil, fmt.Errorf("no deps found for %s/%s", author, name)
	}

	for _, di := range dep {
		if !di.IsOptional {
			for _, pkg := range di.Packages {
				is_game := game_map[pkg]
				if is_game == nil {
					// no info, fetch data
					parts := strings.Split(pkg, "/")

				}

				if *is_game {
					// dependency is a game, not a valid choice
					continue
				} else {
					// valid dependency
					//TODO
				}
			}
		}
	}

	return harddeps, nil
}

func ResolveDependencies(cli *cdb.CDBClient, installed []string, author, name string) ([]*ResolvedDependency, error) {

	rd := []*ResolvedDependency{}
	// games
	game_map := map[string]*bool{}

	installed_depmap := map[string]bool{}
	for _, i := range installed {
		installed_depmap[i] = true
	}

	di, err := resolve_deps(cli, game_map, author, name)
	if err != nil {
		return nil, fmt.Errorf("resolve error with '%s/%s': %v", author, name, err)
	}

	return rd, nil
}
