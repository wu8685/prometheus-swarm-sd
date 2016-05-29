package swarm

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	NODE_SYSTEM_STATUS_PREFIX = "└ "

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

	SystemStatus *SystemStatus
}

type SystemStatus struct {
	Nodes []*Node
}

type Status struct {
	Key      string
	Value    string
	Children []*Status
	deep     int
}

func (s *Status) addChild(c *Status) {
	s.Children = append(s.Children, c)
}

func (e *InfoEntity) toInfo() (info *Info, err error) {
	info = &Info{
		Id:                e.Id,
		Containers:        e.Containers,
		ContainersRunning: e.ContainersRunning,
		ContainersPaused:  e.ContainersPaused,
		ContainersStopped: e.ContainersStopped,
	}

	ss, err := e.toSystemStatus()
	if err != nil {
		return
	}
	info.SystemStatus = ss

	return
}

func (e *InfoEntity) toSystemStatus() (ss *SystemStatus, err error) {
	systemStatus, err := e.processSystemStatus()
	if err != nil {
		return
	}

	ss = &SystemStatus{
		Nodes: []*Node{},
	}
	for _, s := range systemStatus.Children {
		if s.Key == "Nodes" {
			for _, n := range s.Children {
				node := &Node{}
				ss.Nodes = append(ss.Nodes, node)
				node.inject(n)
			}
			break
		}
	}
	return
}

func (e *InfoEntity) processSystemStatus() (root *Status, err error) {
	root = &Status{
		Key:      "root",
		Value:    "",
		Children: []*Status{},
		deep:     0,
	}
	parents := []*Status{root}
	for _, ssToken := range e.SystemStatus {
		s, _err := e.toStatus(ssToken)
		if _err != nil {
			return nil, _err
		}

		parent := parents[len(parents)-1]
		if parent.deep < s.deep {
			parent.addChild(s)
		} else {
			for {
				if parent.deep < s.deep {
					break
				}
				parents = parents[:len(parents)-1]
				parent = parents[len(parents)-1]
			}
			parent.addChild(s)
		}
		parents = append(parents, s)
	}
	return
}

func (e *InfoEntity) toStatus(tokens []string) (s *Status, err error) {
	if len(tokens) != 2 {
		return nil, fmt.Errorf("Incompatible system status content: %v\n", tokens)
	}
	s = &Status{
		Value:    tokens[1],
		Children: []*Status{},
	}

	deep := 1
	for _, c := range tokens[0] {
		if string(c) != " " {
			break
		}
		deep++
	}
	s.deep = deep
	s.Key = strings.TrimSpace(tokens[0])
	if strings.HasPrefix(s.Key, NODE_SYSTEM_STATUS_PREFIX) {
		s.Key = s.Key[2:]
	}
	return
}

func (n *Node) inject(s *Status) (err error) {
	n.Host = s.Key
	n.Addr = s.Value

	for _, attr := range s.Children {
		switch attr.Key {
		case NODE_SYSTEM_STATUS_KEY_CONTAINERS:
			n.Containers, err = strconv.Atoi(attr.Value)
			if err != nil {
				return fmt.Errorf("Fail to parse containers of node info in response info: %s", err.Error())
			}
		case NODE_SYSTEM_STATUS_KEY_CPU:
			n.ReservedCPUs = attr.Value
		case NODE_SYSTEM_STATUS_KEY_ERROR:
			n.Err = attr.Value
		case NODE_SYSTEM_STATUS_KEY_MEMORY:
			n.ReservedMemory = attr.Value
		case NODE_SYSTEM_STATUS_KEY_STATUS:
			n.Status = attr.Value
		case NODE_SYSTEM_STATUS_KEY_LABELS:
			n.Labels = map[string]string{}
			kvs := strings.Split(attr.Value, ",")
			for _, kv := range kvs {
				kv = strings.TrimSpace(kv)
				if kv == "" || strings.Index(kv, "=") == -1 {
					continue
				}
				pair := strings.Split(kv, "=")
				n.Labels[strings.TrimSpace(pair[0])] = strings.TrimSpace(pair[1])
			}
		case NODE_SYSTEM_STATUS_KEY_UPDATE:
			n.UpdatedAt, err = time.Parse(time.RFC3339Nano, strings.TrimSpace(attr.Value))
			if err != nil {
				return fmt.Errorf("Fail to parse update date of node info in response info: %s", err.Error())
			}
		}
	}
	return
}
