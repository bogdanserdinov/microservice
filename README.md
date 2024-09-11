## Golang Microservice Example

This project is a simple microservice built in Go, designed as a template for creating robust, scalable microservices. It comes with various built-in features such as graceful shutdown, readiness probes, metrics collection, and CI/CD integration.

Features:
- Graceful Shutdown: Ensures that the service shuts down properly, finishing in-flight requests before termination.
- Readiness Probe: Kubernetes-ready health check to determine when the service is ready to handle traffic.
- Metrics: Integrated with Prometheus for collecting service-level metrics.
- Dependabot: Automatic dependency management using GitHub Dependabot.
- GitHub Actions: CI/CD pipelines for automated testing and deployment.
- Pull Request Template: Standardized PR template for contributing to the project.
- Dockerfile: Dockerized service for easy deployment.
- Docker Compose for 3rd dependencies.
- DB Connection Pool Management: Efficient connection pooling to manage database connections and improve performance.
