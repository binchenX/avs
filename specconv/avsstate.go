package specconv

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pierrchen/avs/utils"
)

// AvsState record the state accross the avs commands.
// It will be persistend to $genDir/.avsstate
// avs s will create this file.
// avs u will update it.
type AvsState struct {
	GenDir          string   `json:"gen_dir"`
	GenereatedFiles []string `json:"generated_files"`
}

// Update persistent it self to $GenDir/.avsstae
func (s *AvsState) Update() error {
	if s.GenDir == "" {
		s.GenDir, _ = os.Getwd()
	}
	o := filepath.Join(s.GenDir, ".avsstate")
	return SaveSpecToJSON(s, o)
}

// LoadAvsState load the avs state, or error
func LoadAvsState(genDir string) (state *AvsState, err error) {
	stateFile := filepath.Join(genDir, ".avsstate")
	if r, err := utils.FileExists(stateFile); r != true {
		return nil, err
	}

	stateData, err := os.Open(stateFile)
	if err = json.NewDecoder(stateData).Decode(&state); err != nil {
		fmt.Printf("%#v", err)
		return nil, err
	}

	return state, nil
}
