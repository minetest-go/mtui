package cdb

import (
	"fmt"
)

type ResolvedDependency struct {
	Name      string   `json:"name"`
	Choices   []string `json:"choices"`
	Installed bool     `json:"installed"`
}

func ResolveDependencies(cc *CachedCDBClient, required_pkg string, selected_pkgs, installed_pkgs []string) ([]*ResolvedDependency, error) {
	rd := []*ResolvedDependency{}

	// already processed dependencies
	processed_deps := map[string]bool{}

	// convert to lookup maps
	installed_pkg_map := map[string]bool{}
	for _, pkg := range installed_pkgs {
		installed_pkg_map[pkg] = true
	}

	selected_pkg_map := map[string]bool{}
	for _, pkg := range selected_pkgs {
		selected_pkg_map[pkg] = true
	}

	mod_list, err := cc.SearchPackages(&PackageQuery{Type: []PackageType{PackageTypeMod}})
	if err != nil {
		return nil, fmt.Errorf("failed to query mods: %v", err)
	}

	mod_map := map[string]*Package{}
	for _, pkg := range mod_list {
		mod_map[GetPackagename(pkg.Author, pkg.Name)] = pkg
	}

	resolved_dep_infos := map[string][]*DependencyInfo{}

	// recursive resolver
	var resolve func(string) error
	resolve = func(pkgname string) error {
		if processed_deps[pkgname] {
			return nil
		}

		dep := resolved_dep_infos[pkgname]
		author, name := GetAuthorName(pkgname)

		if dep == nil {
			// fetch dep infos
			deps, err := cc.GetDependencies(author, name)
			if err != nil {
				return fmt.Errorf("failed to resolve deps for mod '%s': %v", pkgname, err)
			}

			for n, dep := range deps {
				resolved_dep_infos[n] = dep
			}
			dep = resolved_dep_infos[pkgname]
		}

		if dep == nil {
			// should not happen but check anyway
			return fmt.Errorf("dep unresolved: '%s'", pkgname)
		}

		for _, di := range dep {
			if di.IsOptional {
				continue
			}

			processed_deps[di.Name] = true

			d := &ResolvedDependency{
				Name:    di.Name,
				Choices: []string{},
			}

			if installed_pkg_map[di.Name] {
				// already installed
				d.Installed = true
				rd = append(rd, d)
				continue
			}

			for _, dep_pkgname := range di.Packages {
				if mod_map[dep_pkgname] == nil {
					// not of "mod"-type
					continue
				}

				err = resolve(dep_pkgname)
				if err != nil {
					return err
				}
			}

			rd = append(rd, d)
		}

		return nil
	}

	err = resolve(required_pkg)
	if err != nil {
		return nil, err
	}

	return rd, nil
}
