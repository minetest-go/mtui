package mod

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
)

//go:embed *
var modFS embed.FS

func Install(target string) error {
	walkfn := func(path string, d fs.DirEntry, err error) error {
		if path == "." {
			return nil
		}

		if d.Type().IsDir() {
			err := os.MkdirAll(target+"/"+d.Name(), 0755)
			if err != nil {
				return err
			}
		}

		if d.Type().IsRegular() {
			in, err := modFS.Open(path)
			if err != nil {
				return err
			}
			defer in.Close()

			out, err := os.Create(target + "/" + path)
			if err != nil {
				return err
			}
			defer out.Close()

			_, err = io.Copy(out, in)
			return err
		}
		return nil
	}

	_, err := os.Stat(target)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Creating target directory: %s\n", target)
		err = os.MkdirAll(target, 0755)
		if err != nil {
			return err
		}
	}

	return fs.WalkDir(modFS, ".", walkfn)
}
