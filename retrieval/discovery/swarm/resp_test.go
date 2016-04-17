package swarm

const (
	oneNode = `
{
  "ID": "",
  "Containers": 16,
  "ContainersRunning": 10,
  "ContainersPaused": 0,
  "ContainersStopped": 6,
  "Images": 30,
  "Driver": "",
  "DriverStatus": null,
  "SystemStatus": [
    [
      "Role",
      "primary"
    ],
    [
      "Strategy",
      "spread"
    ],
    [
      "Filters",
      "health, port, dependency, affinity, constraint"
    ],
    [
      "Nodes",
      "1"
    ],
    [
      " hh-yun-k8s-128049.vclound.com",
      "10.199.128.49:2375"
    ],
    [
      "  └ Status",
      "Healthy"
    ],
    [
      "  └ Containers",
      "8"
    ],
    [
      "  └ Reserved CPUs",
      "12 / 25"
    ],
    [
      "  └ Reserved Memory",
      "3.75 GiB / 132 GiB"
    ],
    [
      "  └ Labels",
      "executiondriver=native-0.2, kernelversion=3.10.0-229.4.2.el7.x86_64, operatingsystem=CentOS Linux 7 (Core), storagedriver=devicemapper"
    ],
    [
      "  └ Error",
      "(none)"
    ],
    [
      "  └ UpdatedAt",
      "2016-04-05T09:22:57Z"
    ]
  ],
  "Plugins": {
    "Volume": null,
    "Network": null,
    "Authorization": null
  },
  "MemoryLimit": true,
  "SwapLimit": true,
  "CpuCfsPeriod": true,
  "CpuCfsQuota": true,
  "CPUShares": true,
  "CPUSet": true,
  "IPv4Forwarding": true,
  "BridgeNfIptables": true,
  "BridgeNfIp6tables": true,
  "Debug": false,
  "NFd": 0,
  "OomKillDisable": true,
  "NGoroutines": 0,
  "SystemTime": "2016-04-05T17:23:50.830465718+08:00",
  "ExecutionDriver": "",
  "LoggingDriver": "",
  "NEventsListener": 0,
  "KernelVersion": "3.10.0-229.4.2.el7.x86_64",
  "OperatingSystem": "linux",
  "OSType": "",
  "Architecture": "amd64",
  "IndexServerAddress": "",
  "RegistryConfig": null,
  "NCPU": 50,
  "MemTotal": 283536760012,
  "DockerRootDir": "",
  "HttpProxy": "",
  "HttpsProxy": "",
  "NoProxy": "",
  "Name": "hh-yun-k8s-128050.vclound.com",
  "Labels": null,
  "ExperimentalBuild": false,
  "ServerVersion": "",
  "ClusterStore": "",
  "ClusterAdvertise": ""
}
	`

	twoNode = `
{
  "ID": "",
  "Containers": 16,
  "ContainersRunning": 10,
  "ContainersPaused": 0,
  "ContainersStopped": 6,
  "Images": 30,
  "Driver": "",
  "DriverStatus": null,
  "SystemStatus": [
    [
      "Role",
      "primary"
    ],
    [
      "Strategy",
      "spread"
    ],
    [
      "Filters",
      "health, port, dependency, affinity, constraint"
    ],
    [
      "Nodes",
      "2"
    ],
    [
      " hh-yun-k8s-128049.vclound.com",
      "10.199.128.49:2375"
    ],
    [
      "  └ Status",
      "Healthy"
    ],
    [
      "  └ Containers",
      "8"
    ],
    [
      "  └ Reserved CPUs",
      "12 / 25"
    ],
    [
      "  └ Reserved Memory",
      "3.75 GiB / 132 GiB"
    ],
    [
      "  └ Labels",
      "executiondriver=native-0.2, kernelversion=3.10.0-229.4.2.el7.x86_64, operatingsystem=CentOS Linux 7 (Core), storagedriver=devicemapper"
    ],
    [
      "  └ Error",
      "(none)"
    ],
    [
      "  └ UpdatedAt",
      "2016-04-05T09:22:57Z"
    ],
    [
      " hh-yun-k8s-128050.vclound.com",
      "10.199.128.50:2375"
    ],
    [
      "  └ Status",
      "Healthy"
    ],
    [
      "  └ Containers",
      "8"
    ],
    [
      "  └ Reserved CPUs",
      "12 / 25"
    ],
    [
      "  └ Reserved Memory",
      "3 GiB / 132 GiB"
    ],
    [
      "  └ Labels",
      "executiondriver=native-0.2, kernelversion=3.10.0-229.4.2.el7.x86_64, operatingsystem=CentOS Linux 7 (Core), storagedriver=devicemapper"
    ],
    [
      "  └ Error",
      "(none)"
    ],
    [
      "  └ UpdatedAt",
      "2016-04-05T09:23:29Z"
    ]
  ],
  "Plugins": {
    "Volume": null,
    "Network": null,
    "Authorization": null
  },
  "MemoryLimit": true,
  "SwapLimit": true,
  "CpuCfsPeriod": true,
  "CpuCfsQuota": true,
  "CPUShares": true,
  "CPUSet": true,
  "IPv4Forwarding": true,
  "BridgeNfIptables": true,
  "BridgeNfIp6tables": true,
  "Debug": false,
  "NFd": 0,
  "OomKillDisable": true,
  "NGoroutines": 0,
  "SystemTime": "2016-04-05T17:23:50.830465718+08:00",
  "ExecutionDriver": "",
  "LoggingDriver": "",
  "NEventsListener": 0,
  "KernelVersion": "3.10.0-229.4.2.el7.x86_64",
  "OperatingSystem": "linux",
  "OSType": "",
  "Architecture": "amd64",
  "IndexServerAddress": "",
  "RegistryConfig": null,
  "NCPU": 50,
  "MemTotal": 283536760012,
  "DockerRootDir": "",
  "HttpProxy": "",
  "HttpsProxy": "",
  "NoProxy": "",
  "Name": "hh-yun-k8s-128050.vclound.com",
  "Labels": null,
  "ExperimentalBuild": false,
  "ServerVersion": "",
  "ClusterStore": "",
  "ClusterAdvertise": ""
}
	`

	threeNode = `
{
  "ID": "",
  "Containers": 16,
  "ContainersRunning": 10,
  "ContainersPaused": 0,
  "ContainersStopped": 6,
  "Images": 30,
  "Driver": "",
  "DriverStatus": null,
  "SystemStatus": [
    [
      "Role",
      "primary"
    ],
    [
      "Strategy",
      "spread"
    ],
    [
      "Filters",
      "health, port, dependency, affinity, constraint"
    ],
    [
      "Nodes",
      "3"
    ],
    [
      " hh-yun-k8s-128049.vclound.com",
      "10.199.128.49:2375"
    ],
    [
      "  └ Status",
      "Healthy"
    ],
    [
      "  └ Containers",
      "8"
    ],
    [
      "  └ Reserved CPUs",
      "12 / 25"
    ],
    [
      "  └ Reserved Memory",
      "3.75 GiB / 132 GiB"
    ],
    [
      "  └ Labels",
      "executiondriver=native-0.2, kernelversion=3.10.0-229.4.2.el7.x86_64, operatingsystem=CentOS Linux 7 (Core), storagedriver=devicemapper"
    ],
    [
      "  └ Error",
      "(none)"
    ],
    [
      "  └ UpdatedAt",
      "2016-04-05T09:22:57Z"
    ],
    [
      " hh-yun-k8s-128050.vclound.com",
      "10.199.128.50:2375"
    ],
    [
      "  └ Status",
      "Healthy"
    ],
    [
      "  └ Containers",
      "8"
    ],
    [
      "  └ Reserved CPUs",
      "12 / 25"
    ],
    [
      "  └ Reserved Memory",
      "3 GiB / 132 GiB"
    ],
    [
      "  └ Labels",
      "executiondriver=native-0.2, kernelversion=3.10.0-229.4.2.el7.x86_64, operatingsystem=CentOS Linux 7 (Core), storagedriver=devicemapper"
    ],
    [
      "  └ Error",
      "(none)"
    ],
    [
      "  └ UpdatedAt",
      "2016-04-05T09:23:29Z"
    ],
    [
      " hh-yun-k8s-128051.vclound.com",
      "10.199.128.51:2375"
    ],
    [
      "  └ Status",
      "Healthy"
    ],
    [
      "  └ Containers",
      "8"
    ],
    [
      "  └ Reserved CPUs",
      "12 / 25"
    ],
    [
      "  └ Reserved Memory",
      "3 GiB / 132 GiB"
    ],
    [
      "  └ Labels",
      "executiondriver=native-0.2, kernelversion=3.10.0-229.4.2.el7.x86_64, operatingsystem=CentOS Linux 7 (Core), storagedriver=devicemapper"
    ],
    [
      "  └ Error",
      "(none)"
    ],
    [
      "  └ UpdatedAt",
      "2016-04-05T09:23:29Z"
    ]
  ],
  "Plugins": {
    "Volume": null,
    "Network": null,
    "Authorization": null
  },
  "MemoryLimit": true,
  "SwapLimit": true,
  "CpuCfsPeriod": true,
  "CpuCfsQuota": true,
  "CPUShares": true,
  "CPUSet": true,
  "IPv4Forwarding": true,
  "BridgeNfIptables": true,
  "BridgeNfIp6tables": true,
  "Debug": false,
  "NFd": 0,
  "OomKillDisable": true,
  "NGoroutines": 0,
  "SystemTime": "2016-04-05T17:23:50.830465718+08:00",
  "ExecutionDriver": "",
  "LoggingDriver": "",
  "NEventsListener": 0,
  "KernelVersion": "3.10.0-229.4.2.el7.x86_64",
  "OperatingSystem": "linux",
  "OSType": "",
  "Architecture": "amd64",
  "IndexServerAddress": "",
  "RegistryConfig": null,
  "NCPU": 50,
  "MemTotal": 283536760012,
  "DockerRootDir": "",
  "HttpProxy": "",
  "HttpsProxy": "",
  "NoProxy": "",
  "Name": "hh-yun-k8s-128050.vclound.com",
  "Labels": null,
  "ExperimentalBuild": false,
  "ServerVersion": "",
  "ClusterStore": "",
  "ClusterAdvertise": ""
}
	`
)
