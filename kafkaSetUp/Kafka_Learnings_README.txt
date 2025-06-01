
# Kafka Learnings Summary

This document summarizes key learnings about Kafka based on today's hands-on experience with the `segmentio/kafka-go` library in Go.

---

## ✅ Kafka Consumer Basics

### Explicit Partition vs Consumer Group

| Config                            | Behavior                                              |
|----------------------------------|-------------------------------------------------------|
| `Partition` set, no `GroupID`    | Reads only from that partition, no offset tracking   |
| `GroupID` set, no `Partition`    | Kafka manages partition assignment & offset tracking |

### Reading From a Partition

```go
reader := kafka.NewReader(kafka.ReaderConfig{
    Brokers:   []string{"localhost:9092"},
    Topic:     "my-topic",
    Partition: 0,
})
reader.SetOffset(kafka.LastOffset) // Start from latest messages
```

### Reading with Consumer Group (Offset Persistence)

```go
reader := kafka.NewReader(kafka.ReaderConfig{
    Brokers: []string{"localhost:9092"},
    Topic:   "my-topic",
    GroupID: "my-consumer-group",
})
```

> Using a `GroupID` enables Kafka to resume consumption from the last committed offset.

---

## 🔄 Goroutine & Fault Recovery

Wrap partition readers in goroutines with panic recovery and restart logic:

```go
go func() {
    for {
        func() {
            defer func() {
                if r := recover(); r != nil {
                    log.Printf("Recovered from panic: %v", r)
                }
            }()

            reader := kafka.NewReader(...) // with Partition or GroupID
            defer reader.Close()

            for {
                msg, err := reader.ReadMessage(ctx)
                if err != nil {
                    log.Printf("Error: %v", err)
                    return // restart reader
                }

                fmt.Printf("Received from actual partition %d: %s\n", msg.Partition, msg.Value)
            }
        }()

        time.Sleep(3 * time.Second) // restart delay
    }
}()
```

---

## 🧠 Common Pitfalls

- **Don't assume the partition from config equals the message partition.**  
  Always log `msg.Partition`, not just the configured value.

- **Without `GroupID`, no offset is persisted.**  
  Use `reader.SetOffset(kafka.LastOffset)` if you want to start at the tail.

---

## 📤 Kafka Producer Learnings

### Hash Balancing + Partition Confusion

If you set both a custom key and an explicit `Partition`, the `Partition` **takes priority**.

```go
writer := &kafka.Writer{
    Addr:     kafka.TCP("localhost:9092"),
    Topic:    "my-topic",
    Balancer: &kafka.Hash{}, // Used only if Partition is not set
}
```

### To Route Messages Based on Key (e.g., `glid` ending with 0/1):

```go
glid := "1230"
partition := 0
if glid[len(glid)-1] == '1' {
    partition = 1
}

message := kafka.Message{
    Key:       []byte(glid),
    Value:     []byte(`{"glid":"1230","name":"yashwant"}`),
    Partition: partition, // Manually route to correct partition
}
```

---

## 🧪 Testing Observations

- Kafka consumers will block on `ReadMessage()` until a new message arrives.
- Manual partition assignment ignores consumer group load balancing.
- Offsets are persisted **only when using `GroupID`**.
- Partition numbers do **not** guarantee FIFO between partitions — only *within* a partition.

---

## ✅ Summary of Best Practices

- Use consumer groups (`GroupID`) for scalable, offset-managed consumption.
- Use manual partitions only for advanced scenarios (e.g., custom partition routing).
- Always log `msg.Partition` to trace true data flow.
- Implement retries and recoveries for production consumers.
- Use message keys or payload-based logic (e.g., `glid`) for custom routing.

---

## 🔗 Tools Used

- [segmentio/kafka-go](https://github.com/segmentio/kafka-go)
- Go 1.21+
- Local Kafka via Docker

---
