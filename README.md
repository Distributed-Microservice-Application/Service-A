# Service A: The Calculator & Outbox Publisher

**Service A** is a high-performance microservice built with Go. Its primary role is to accept gRPC requests, perform a simple calculation, and reliably hand off the result for asynchronous processing using the **Outbox Pattern**. This ensures that every calculation is durably stored and eventually delivered, even in the face of downstream service failures.

## 🚀 Core Functionality

- **gRPC Server**: Exposes a `CalculateSum` RPC endpoint that accepts two numbers and computes their sum.
- **Outbox Pattern for Guaranteed Delivery**: Instead of publishing results directly to a message broker, Service A writes the result to an `outbox` table within its own PostgreSQL database. This operation is atomic with the primary business logic.
- **Change Data Capture (CDC)**: The service relies on **Debezium** to monitor the `outbox` table. Debezium captures any new rows inserted into the table and automatically publishes them as events to an Apache Kafka topic (`user-events`). This decouples the service from the messaging system and guarantees message delivery.
- **Horizontal Scalability**: The architecture supports running multiple instances of Service A, which are `load-balanced by Nginx` for high availability and throughput.
- **Observability**: Exposes a `/metrics` endpoint for Prometheus to scrape performance metrics.

## 🛠️ Technology Stack

| Category      | Technology                                      |
|---------------|-------------------------------------------------|
| **Language**  | Go                                              |
| **Framework** | Standard Library, gRPC-Go                       |
| **API**       | gRPC                                            |
| **Database**  | PostgreSQL                                      |
| **Messaging** | Kafka (via Debezium for Change Data Capture)    |
| **Metrics**   | Prometheus Client for Go                        |

## ⚙️ How It Works

1.  **gRPC Request**: An external client (or the K6 load test) sends a `CalculateSum` request via the Nginx gRPC load balancer.
2.  **Calculation**: One of the Service A instances receives the request and calculates the sum of the two numbers.
3.  **Atomic Write**: The service writes the result into the `outbox` table in its PostgreSQL database. This is the end of its synchronous work for the request.
4.  **CDC with Debezium**: The Debezium connector, configured to watch the `outbox` table, detects the new row.
5.  **Kafka Message**: Debezium creates a JSON message containing the data from the new row and publishes it to the `user-events` Kafka topic.
6.  **Downstream Consumption**: Service B (or any other consumer) can now consume this event from Kafka for further processing.

This approach ensures that the result is captured durably and will be sent to Kafka as soon as the CDC platform processes the change, providing a highly reliable and resilient system.

## ⚡ Load Testing

The service includes a K6 script to test the performance of its HTTP endpoint via the Nginx load balancer.

- **Test Script**: `internal/test/k6_http.js`
- **To run the test**:
    1.  Ensure the full application stack is running (`make start` in the `Docker-Compose` directory).
    2.  [Install K6](https://k6.io/docs/getting-started/installation/).
    3.  Navigate to the `Service-A` directory and execute the script:
    ```sh
    k6 run internal/test/k6_http.js
    ```