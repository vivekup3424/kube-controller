# Kubernetes Resource Monitor Controller

The **Kubernetes Resource Monitor Controller** is a custom Kubernetes controller designed to monitor CPU and Memory usage for every Deployment within a cluster. It aggregates usage metrics based on namespaces and exposes this data through both a `/metrics` endpoint and an API route. The controller is containerized using Docker for easy deployment.

## Features

- Monitors CPU and Memory usage for Deployments.
- Aggregates usage metrics based on namespaces.
- Exposes metrics through a `/metrics` endpoint for kubernetes controller.
- Provides an API route for accessing aggregated usage data.
- Implemented **Concurrency using Go routines** to improve performance and efficiency. By leveraging Go routines, the controller can handle multiple monitoring tasks concurrently, leading to faster data collection and processing.
## Working
+---------------------------+
|        main.go            |
|---------------------------|
| - Initializes Clientsets  |
| - Starts HTTP Server      |
| - Starts Controller       |
| - Sets up Cron Job        |
+-----+---------------------+
      |
      v
+----------------------------------+
|       StartController()          |
|----------------------------------|
| - Sets up Informer               |
| - Handles Deployment Events      |
| - Calls handleDeploymentChange() |
+-----+----------------------------+
      |                  |
      v                  |
+----------------------------------+  |
| handleDeploymentChange()         |  |
|----------------------------------|  |
| - Calls Compute()                |  |
| - Logs Errors                    |  |
+-----+----------------------------+  |
      |                              |
      v                              v
+----------------------------------+  |
|       Compute()                  |  |
|----------------------------------|  |
| - Iterates through Namespaces    |  |
| - Calls getNamespaceMetrics()    |  |
| - Aggregates Metrics             |  |
+-----+----------------------------+  |
      |                              |
      v                              |
+----------------------------------+  |
| getNamespaceMetrics()            |  |
|----------------------------------|  |
| - Fetches Deployments            |  |
| - Collects Pod Metrics           |  |
| - Aggregates CPU and Memory Usage|  |
+----------------------------------+  |
      |                              |
      v                              |
+----------------------------------+  |
|      HTTP Server                 |  |
|----------------------------------|  |
| - /metrics Endpoint              |  |
| - /update Endpoint               |  |
+----------------------------------+  |
      |                              |
      v                              |
+----------------------------------+  |
|       Cron Job                   |  |
|----------------------------------|  |
| - Triggers UpdateMetrics()       |  |
| - Calls Compute()                |  |
+----------------------------------+--+

## Getting Started

These instructions will help you get the Kubernetes Resource Monitor Controller up and running on your Kubernetes cluster.

### Prerequisites

- Kubernetes cluster up and running.
- `kubectl` configured to manage your cluster.
- Docker installed on your local machine.

### Installation

1. Clone this repository to your local machine:

   ```bash
   git clone https://github.com/gurpreet-legend/kube-controller.git
   cd kube-controller
   ```

2. Build the Docker container:

   ```bash
   docker build -t kube-controller:v1 .
   ```

3. Deploy the controller to your Kubernetes cluster:

   ```bash
   kubectl apply -f resources/controller.yaml
   ```

### Usage

Once the controller is deployed, it will start monitoring the CPU and Memory usage of Deployments in your cluster. You can access the metrics through the following endpoints:

- **Metrics Endpoint**: The `/metrics` endpoint is designed for usage statistics. It provides aggregated usage metrics for each namespace.

  Example: `http://your-controller-ip:8080/metrics`


### Cleanup

To remove the Kubernetes Resource Monitor Controller and all associated resources, run:

```bash
kubectl delete -f resources/controller.yaml
```

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request. Adding some new stuffs here.
