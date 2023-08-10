package depanalyzer

import (
	"fmt"
	"os"
	"path"
)

type dependencycontext struct {
	installed map[string]bool
	missing   map[string]bool
}

type DependencyAnalysis struct {
	Installed []string `json:"installed"`
	Missing   []string `json:"missing"`
}

func scanDir(dir string, ctx *dependencycontext) error {
	stat, _ := os.Stat(path.Join(dir, "modpack.conf"))
	if stat != nil && !stat.IsDir() {
		return scanDirs(dir, ctx)
	}

	modconf := path.Join(dir, "mod.conf")
	stat, _ = os.Stat(modconf)
	if stat != nil && !stat.IsDir() {
		data, err := os.ReadFile(modconf)
		if err != nil {
			return fmt.Errorf("error reading '%s': %v", modconf, err)
		}
		mc, err := ParseModConf(data)
		if err != nil {
			return fmt.Errorf("error parsing '%s': %v", modconf, err)
		}

		modname := mc.Name
		if modname == "" {
			// fall back to dirname if name not specified
			modname = path.Base(dir)
		}

		ctx.installed[modname] = true
		ctx.missing[modname] = false

		for _, dep := range mc.Depends {
			if !ctx.installed[dep] {
				ctx.missing[dep] = true
			}
		}
	}

	depends := path.Join(dir, "depends.txt")
	stat, _ = os.Stat(depends)
	if stat != nil && !stat.IsDir() {
		data, err := os.ReadFile(depends)
		if err != nil {
			return fmt.Errorf("error reading '%s': %v", depends, err)
		}
		di, err := ParseDependsTXT(data)
		if err != nil {
			return fmt.Errorf("error in depends.txt: %v", err)
		}

		modname := path.Base(dir)
		ctx.installed[modname] = true
		ctx.missing[modname] = false

		for _, dep := range di.Depends {
			if !ctx.installed[dep] {
				ctx.missing[dep] = true
			}
		}
	}

	return nil
}

func scanDirs(dir string, ctx *dependencycontext) error {
	fi, _ := os.Stat(dir)
	if fi == nil {
		// nonexistent dir
		return nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			err = scanDir(path.Join(dir, entry.Name()), ctx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func AnalyzeDeps(dirs ...string) (*DependencyAnalysis, error) {
	ctx := &dependencycontext{
		installed: map[string]bool{},
		missing:   map[string]bool{},
	}
	for _, dir := range dirs {
		err := scanDirs(dir, ctx)
		if err != nil {
			return nil, err
		}
	}

	da := &DependencyAnalysis{
		Installed: make([]string, 0),
		Missing:   make([]string, 0),
	}

	for name, ok := range ctx.installed {
		if ok {
			da.Installed = append(da.Installed, name)
		}
	}

	for name, ok := range ctx.missing {
		if ok {
			da.Missing = append(da.Missing, name)
		}
	}

	return da, nil
}
