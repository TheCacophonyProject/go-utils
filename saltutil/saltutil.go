package saltutil

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/TheCacophonyProject/go-utils/logging"
)

const (
	grainsFile    = "/etc/salt/grains"
	nodegroupFile = "/etc/cacophony/salt-nodegroup"
	minionIdFile  = "/etc/salt/minion_id"
)

type Grains struct {
	DeviceName  string `json:"device_name,omitempty"`
	Environment string `json:"environment,omitempty"`
	Group       string `json:"group,omitempty"`
}

func SetGrains(grains Grains, log *logging.Logger) error {
	// Setup logger if none is provided
	if log == nil {
		log = logging.NewLogger("info")
	}

	grainsJSON, err := json.Marshal(grains)
	if err != nil {
		log.Errorf("Failed to marshal grains: %v", err)
		return err
	}

	command := []string{"salt-call", "grains.setvals", string(grainsJSON)}
	log.Debugf("Running command: %s", strings.Join(command, " "))
	out, err := exec.Command(command[0], command[1:]...).CombinedOutput()
	if err != nil {
		log.Errorf("Failed to set grains: %s, error: %v, when running: %s", string(out), err, strings.Join(command, " "))
		return err
	}

	return nil
}

// GetSaltGrains returns the salt grains. Optional to pass a logger, pass nil to use default.
func GetSaltGrains(log *logging.Logger) (*Grains, error) {
	// Setup logger if none is provided
	if log == nil {
		log = logging.NewLogger("info")
	}

	// Open up grains file
	grains := make(map[string]string)
	file, err := os.Open(grainsFile)
	if os.IsNotExist(err) {
		// Some devices will not have any grains set so this is not an error
		log.Debugf("No grains file found: %v", err)
		return &Grains{}, nil
	} else if err != nil {
		log.Errorf("Failed to open grains file: %v", err)
		return nil, err
	}
	defer file.Close()

	// Parse the grains file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			grains[key] = value
		} else if line != "" {
			err = fmt.Errorf("failed to parse line in grains file: '%s'", line)
			log.Errorf(err.Error())
			return nil, err
		}
	}
	if err := scanner.Err(); err != nil {
		log.Errorf("Error reading grains file: %v", err)
		return nil, err
	}

	// Return the grains
	grainsJSON, err := json.Marshal(grains)
	if err != nil {
		log.Errorf("Failed to marshal grains: %v", err)
		return nil, err
	}
	log.Debugf("Grains: %s", string(grainsJSON))

	grainsStruct := &Grains{}
	err = json.Unmarshal(grainsJSON, grainsStruct)
	if err != nil {
		log.Errorf("Failed to unmarshal grains: %v", err)
		return nil, err
	}

	return grainsStruct, nil
}

func GetNodegroupFromFile() (string, error) {
	nodegroup, err := os.ReadFile(nodegroupFile)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(nodegroup)), nil
}

func GetMinionID(log *logging.Logger) (string, error) {
	// Setup logger if none is provided
	if log == nil {
		log = logging.NewLogger("info")
	}

	// Read salt minion ID
	idRaw, err := os.ReadFile(minionIdFile)
	if err != nil {
		log.Error("Error reading minion ID: " + err.Error())
		return "", err
	}
	id := strings.TrimSpace(string(idRaw))
	log.Debugf("Minion ID: '%s'", id)
	return id, nil
}
