package transformer

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/rantav/go-archetype/log"
	"github.com/rantav/go-archetype/operations"
	"github.com/rantav/go-archetype/reader"
	"github.com/rantav/go-archetype/types"
	"github.com/rantav/go-archetype/writer"
)

type Transformer interface {
	GetName() string
	GetFilePatterns() []types.FilePattern
	Template(vars map[string]string) error
	Transform(types.File) types.File
}

func Transform(source, destination string, transformations Transformations) error {
	err := before(transformations)
	if err != nil {
		return err
	}

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "error walking to file")
		}
		sourceFile := path
		isDir, ignored, file, err := reader.ReadFile(sourceFile, info, transformations.IsGloballyIgnored)
		if err != nil {
			return errors.Wrap(err, "error reading file")
		}
		if isDir {
			return nil
		}

		if ignored {
			log.Debugf("Ignoring file %s", path)
		} else {
			file, err = transformations.Transform(file)
			err := writer.WriteFile(destination, file, info.Mode())
			if err != nil {
				return err
			}
		}
		return errors.Wrap(err, "transforming")
	})
	if err != nil {
		return err
	}

	return after(transformations)
}

func before(ts Transformations) error {
	return executeOperators(ts.before)
}

func after(ts Transformations) error {
	return executeOperators(ts.after)
}

func executeOperators(ops []operations.Operator) error {
	for _, op := range ops {
		if err := op.Operate(); err != nil {
			return err
		}
	}
	return nil
}
