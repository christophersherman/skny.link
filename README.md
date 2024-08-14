# skny.link

Welcome to **skny.link** – a **free, hyperscalable URL shortener** designed to be fast, reliable, and easy to use. Whether you need a short link for social media, marketing campaigns, or just to share with friends, skny.link has got you covered.

## 🚀 Features

- **Free to Use**: No sign-up required, simply shorten your links and share them instantly.
- **Hyperscalable**: Built with modern technologies, skny.link can handle millions of requests per day, ensuring high availability and performance even under heavy traffic.
- **Customizable**: Choose your own custom aliases for shortened URLs.
- **Analytics (Coming Soon)**: Track how your shortened links are performing with detailed analytics and insights.

## 🛠️ Technology Stack

skny.link is powered by:

- **Go**: The core backend is written in Go, ensuring fast performance and efficient concurrency.
- **Redis**: Used for quick read and write operations, providing ultra-fast access to frequently used URLs.
 **Docker & Docker Compose**: All services are containerized, making deployment and scaling a breeze.
- **Nginx**: Serving as the load balancer, it efficiently distributes traffic across multiple instances.
- **Kubernetes (Optional)**: For large-scale deployments, skny.link can be deployed on Kubernetes, allowing you to easily manage and scale the service.

