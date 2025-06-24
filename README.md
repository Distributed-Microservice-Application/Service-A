
**Service A** is a gRPC-based microservice designed to **perform addition operations** reliably and efficiently. When it receives a request containing two numbers:

1. **Computes the Sum**: Adds the two input numbers.
2. **Implements the Outbox Pattern**:

   * Instead of sending the result directly to Kafka, it **stores the sum in an outbox table** in its database.
   * This ensures **message durability and fault-tolerance**, even if Kafka or downstream systems are temporarily unavailable.
3. **Background Job**:

   * A separate **background worker** continuously scans the outbox table and **forwards the stored messages to Kafka**.
   * This decouples the core logic from external messaging, increasing reliability and **ensuring at-least-once delivery**.
