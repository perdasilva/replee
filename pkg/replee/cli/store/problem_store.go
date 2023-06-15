package store

import (
	"encoding/json"
	"fmt"
	"github.com/perdasilva/replee/pkg/deppy"
	"github.com/perdasilva/replee/pkg/deppy/resolution"
	"github.com/perdasilva/replee/pkg/deppy/resolver"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type ResolutionProblemStore interface {
	New(resolutionProblemID deppy.Identifier) error
	Delete(resolutionProblemID deppy.Identifier) error
	Get(resolutionProblemID deppy.Identifier) (deppy.MutableResolutionProblem, error)
	Clean() error
	Save(resolutionProblem deppy.MutableResolutionProblem) error
	Solve(resolutionProblemID deppy.Identifier) (resolver.Solution, error)
}

var _ ResolutionProblemStore = &FSResolutionProblemStore{}

type FSResolutionProblemStore struct {
	rootDir string
}

func NewFSResolutionProblemStore(rootDir string) (*FSResolutionProblemStore, error) {
	if rootDir == "" {
		rootDir = "."
	}
	if err := createDirectoryIfNotExist(rootDir); err != nil {
		return nil, err
	}
	return &FSResolutionProblemStore{rootDir: rootDir}, nil
}

func (f FSResolutionProblemStore) New(resolutionProblemID deppy.Identifier) error {
	problemPath := path.Join(f.rootDir, resolutionProblemID.String()+".json")
	if exists, err := pathAlreadyExists(problemPath); err != nil {
		return err
	} else if exists {
		return deppy.ConflictErrorf("resolution problem %s already exists", resolutionProblemID)
	}
	return writeJSONToFile(resolution.NewMutableResolutionProblem(resolutionProblemID), problemPath)
}

func (f FSResolutionProblemStore) Delete(resolutionProblemID deppy.Identifier) error {
	problemPath := path.Join(f.rootDir, resolutionProblemID.String()+".json")
	if exists, err := pathAlreadyExists(problemPath); err != nil {
		return err
	} else if !exists {
		return deppy.ConflictErrorf("resolution problem %s does not exist", resolutionProblemID)
	}
	return os.Remove(problemPath)
}

func (f FSResolutionProblemStore) Get(resolutionProblemID deppy.Identifier) (deppy.MutableResolutionProblem, error) {
	problemPath := path.Join(f.rootDir, resolutionProblemID.String()+".json")
	if exists, err := pathAlreadyExists(problemPath); err != nil {
		return nil, err
	} else if !exists {
		return nil, deppy.ConflictErrorf("resolution problem %s does not exist", resolutionProblemID)
	}
	m := &resolution.MutableResolutionProblem{}
	if err := unmarshalJSONFile(problemPath, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (f FSResolutionProblemStore) Clean() error {
	err := filepath.Walk(f.rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}

		if !info.IsDir() && filepath.Ext(path) == ".json" {
			err := os.Remove(path)
			if err != nil {
				return fmt.Errorf("error deleting file %s: %w", path, err)
			}
			fmt.Printf("Deleted file: %s\n", path)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking directory: %w", err)
	}

	return nil
}

func (f FSResolutionProblemStore) Save(resolutionProblem deppy.MutableResolutionProblem) error {
	problemPath := path.Join(f.rootDir, resolutionProblem.ResolutionProblemID().String()+".json")
	if exists, err := pathAlreadyExists(problemPath); err != nil {
		return err
	} else if !exists {
		return deppy.ConflictErrorf("resolution problem %s does not exist", resolutionProblem.ResolutionProblemID())
	}
	return writeJSONToFile(resolutionProblem, problemPath)
}

func (f FSResolutionProblemStore) Solve(resolutionProblemID deppy.Identifier) (resolver.Solution, error) {
	//TODO implement me
	panic("implement me")
}

func pathAlreadyExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Path does not exist
			return false, nil
		}
		// Other error occurred
		return false, err
	}
	return true, nil
}

func writeJSONToFile(data interface{}, filename string) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write JSON to file: %w", err)
	}

	return nil
}

func unmarshalJSONFile(filename string, v interface{}) error {
	jsonData, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read JSON file: %w", err)
	}

	err = json.Unmarshal(jsonData, &v)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

func createDirectoryIfNotExist(dirPath string) error {
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
		fmt.Printf("Created directory: %s\n", dirPath)
	} else if err != nil {
		return fmt.Errorf("failed to access directory: %w", err)
	}

	return nil
}
