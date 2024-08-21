# skny.link

Welcome to **skny.link** – a **free, hyperscalable URL shortener** designed to be fast, reliable, and easy to use. 

## Purpose

Emulating managed complexity to practice practical production engineering. This project is mostly for practice desigining a highly available, robust application.  

This app can, and will, become more complex as I try to build something that can scale to millions of simulated users. For example, I picked Redis for speed but this will not scale yet until I configure Redis replication (master-slave) or Redis cluster to achieve meaningful horizontal scalability. 


## 🛠️ Technology Stack

skny.link is powered by:

- **Go**: The core backend is written in Go, ensuring fast performance and efficient concurrency(coming soon). I also personally love Go's build simplicity for CI/CD purposes. 
- **Redis**: Used for quick read and write operations, providing ultra-fast access to frequently used URLs.
 **Docker & Docker Compose**: All services are containerized, making deployment and scaling a breeze.
- **Nginx**: Serving as the load balancer, it efficiently distributes traffic across multiple instances.
- **Kubernetes (Optional)**: The scalability is mostly derived from K8s horizontal pod autoscaling in response to increased loads. I will host this app on K8s until my spot user node pools get evicted.

