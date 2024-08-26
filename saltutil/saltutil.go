package saltutil

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/TheCacophonyProject/go-utils/logging"
	"github.com/sirupsen/logrus"
)

const (
	grainsFile    = "/etc/salt/grains"
	nodegroupFile = "/etc/cacophony/salt-nodegroup"
	minionIdFile  = "/etc/salt/minion_id"
)

// GetSaltGrains returns the salt grains. Optional to pass a logger, pass nil to use default.
func GetSaltGrains(log *logrus.Logger) (map[string]string, error) {
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
		return grains, nil
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

	return grains, nil
}

func GetNodegroupFromFile() (string, error) {
	nodegroup, err := os.ReadFile(nodegroupFile)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(nodegroup)), nil
}
func GetMinionID(log *logrus.Logger) (string, error) {
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
