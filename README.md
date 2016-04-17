# Prometheus SD plugin for Swarm
These codes contain a service discovery based on [Prometheus](https://github.com/prometheus/prometheus) 0.17.0 for [Swarm](https://github.com/docker/swarm) v1.1.0 and later. The implementation code can be found in commits.
## Prerequisite
Cause Swarm dosen't provide metrics by itself, a CAdvisor service should be running on every masters and nodes. The CAdvisor listening port should be 8070 by default.
## Configuration
```
- job_name: service-swarm
  swarm_sd_configs:
  - masters:
    - 'http://swarm.example.com:8080'

    refresh_interval: 1s
    metrics_port: '8060'
```
- **refresh_interval** instructs how often SD plugin retrieve the Swarm node informations.
- **metrics_port** instructs which port the CAdvisor service is listening on.
- The label policy configurations are the same as Kubernetes SD plugin.