package swarm

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/common/log"
)

const (
	ROLE     = "Role"
	STRATEGY = "Strategy"
	FILTER   = "Filter"
	NODES    = "Nodes"

	NODE_SYSTEM_STATUS_NUM    = 8
	NODE_SYSTEM_STATUS_PREFIX = "  └ "

	NODE_SYSTEM_STATUS_KEY_STATUS     = "Status"
	NODE_SYSTEM_STATUS_KEY_CONTAINERS = "Containers"
	NODE_SYSTEM_STATUS_KEY_CPU        = "Reserved CPUs"
	NODE_SYSTEM_STATUS_KEY_MEMORY     = "Reserved Memory"
	NODE_SYSTEM_STATUS_KEY_LABELS     = "Labels"
	NODE_SYSTEM_STATUS_KEY_ERROR      = "Error"
	NODE_SYSTEM_STATUS_KEY_UPDATE     = "UpdatedAt"
)

type InfoEntity struct {
	Id                string `json:"ID,omitempty"`
	Containers        int    `json:"Containers,omitempty"`
	ContainersRunning int    `json:"ContainersRunning,omitempty"`
	ContainersPaused  int    `json:"ContainersPaused,omitempty"`
	ContainersStopped int    `json:"ContainersStopped,omitempty"`

	SystemStatus [][]string `json:"SystemStatus,omitempty"`
}

type Node struct {
	Host           string
	Addr           string
	Status         string
	Containers     int
	ReservedCPUs   string
	ReservedMemory string
	Labels         map[string]string
	Err            string
	UpdatedAt      time.Time
}

type Info struct {
	Id                string
	Containers        int
	ContainersRunning int
	ContainersPaused  int
	ContainersStopped int

	Role     string
	Strategy string
	Filters  []string
	Nodes    []*Node
}

func (e *InfoEntity) toInfo() (info *Info, err error) {
	info = &Info{
		Id:                e.Id,
		Containers:        e.Containers,
		ContainersRunning: e.ContainersRunning,
		ContainersPaused:  e.ContainersPaused,
		ContainersStopped: e.ContainersStopped,
	}

	// parse the first 4 infos in system status
	var i int
	nodeNum := 0
	for i = 0; i < 4; i = i + 1 {
		tokens := e.SystemStatus[i]
		if nodeNum, err = e.injectGeneralInfo(tokens, info); err != nil {
			return
		}
	}

	if nodeNum == 0 {
		log.Debugf("Swarm has no node.")
		return
	}

	if len(e.SystemStatus) != 4+NODE_SYSTEM_STATUS_NUM*nodeNum {
		return nil, fmt.Errorf("Fail to parse response info: no compatible response content. The system status is:\n%v", e.SystemStatus)
	}

	log.Debugf("Begin to parse swarm response info: %v", e.SystemStatus)
	if info.Nodes == nil {
		info.Nodes = []*Node{}
	}
	// parse rest of system status
	for i = 0; i < nodeNum; i = i + 1 {
		offset := 4 + NODE_SYSTEM_STATUS_NUM*i
		j := 0
		node := &Node{}
		for ; j < NODE_SYSTEM_STATUS_NUM; j = j + 1 {
			tokens := e.SystemStatus[offset+j]
			if err = e.injectNodeInfo(tokens, node); err != nil {
				return
			}
		}
		info.Nodes = append(info.Nodes, node)
	}

	return info, nil
}

func (e *InfoEntity) injectGeneralInfo(tokens []string, info *Info) (int, error) {
	nodeNum := 0
	switch tokens[0] {
	case ROLE:
		info.Role = tokens[1]
	case STRATEGY:
		info.Strategy = tokens[1]
	case FILTER:
		if info.Filters == nil {
			info.Filters = []string{}
		}
		filters := strings.Split(tokens[1], ",")
		for _, f := range filters {
			f = strings.TrimSpace(f)
			if f == "" {
				continue
			}
			info.Filters = append(info.Filters, f)
		}
	case NODES:
		n, _err := strconv.Atoi(tokens[1])
		if _err != nil {
			return 0, fmt.Errorf("Fail to parse node num in response info: %s.", _err.Error())
		}
		nodeNum = n
	}
	return nodeNum, nil
}

func (e *InfoEntity) injectNodeInfo(tokens []string, node *Node) (err error) {
	key := strings.TrimSpace(strings.Trim(tokens[0], NODE_SYSTEM_STATUS_PREFIX))
	switch key {
	case NODE_SYSTEM_STATUS_KEY_STATUS:
		node.Status = strings.TrimSpace(tokens[1])
	case NODE_SYSTEM_STATUS_KEY_CONTAINERS:
		node.Containers, err = strconv.Atoi(tokens[1])
		if err != nil {
			return fmt.Errorf("Fail to parse containers of node info in response info: %s", err.Error())
		}
	case NODE_SYSTEM_STATUS_KEY_CPU:
		node.ReservedCPUs = strings.TrimSpace(tokens[1])
	case NODE_SYSTEM_STATUS_KEY_MEMORY:
		node.ReservedMemory = strings.TrimSpace(tokens[1])
	case NODE_SYSTEM_STATUS_KEY_LABELS:
		if node.Labels == nil {
			node.Labels = map[string]string{}
		}
		kvs := strings.Split(tokens[1], ",")
		for _, kv := range kvs {
			kv = strings.TrimSpace(kv)
			if kv == "" || strings.Index(kv, "=") == -1 {
				continue
			}
			pair := strings.Split(kv, "=")
			node.Labels[strings.TrimSpace(pair[0])] = strings.TrimSpace(pair[1])
		}
	case NODE_SYSTEM_STATUS_KEY_ERROR:
		node.Err = strings.TrimSpace(tokens[1])
	case NODE_SYSTEM_STATUS_KEY_UPDATE:
		node.UpdatedAt, err = time.Parse(time.RFC3339Nano, strings.TrimSpace(tokens[1]))
		if err != nil {
			return fmt.Errorf("Fail to parse update date of node info in response info: %s", err.Error())
		}
	default:
		host := strings.TrimSpace(tokens[0])
		addr := strings.TrimSpace(tokens[1])
		node.Host = host
		node.Addr = addr
	}
	return nil
}
