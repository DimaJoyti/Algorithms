# 📨 Message Queues & Event Streaming

## 🎯 Why Asynchronous Communication?

### Synchronous vs Asynchronous Trade-offs
```go
// Synchronous - Blocking, immediate response
func ProcessOrderSync(order *Order) error {
    // All operations must complete before returning
    if err := validateOrder(order); err != nil {
        return err
    }
    
    if err := chargePayment(order); err != nil {
        return err
    }
    
    if err := updateInventory(order); err != nil {
        return err
    }
    
    if err := sendConfirmationEmail(order); err != nil {
        return err // Email failure blocks entire operation
    }
    
    return nil
}

// Asynchronous - Non-blocking, eventual consistency
func ProcessOrderAsync(order *Order) error {
    // Validate immediately (critical path)
    if err := validateOrder(order); err != nil {
        return err
    }
    
    // Publish events for async processing
    eventBus.Publish("order.created", OrderCreatedEvent{
        OrderID: order.ID,
        UserID:  order.UserID,
        Items:   order.Items,
    })
    
    return nil // Return immediately
}

// Event handlers process asynchronously
func HandleOrderCreated(event OrderCreatedEvent) {
    // These can fail independently without affecting order creation
    go chargePayment(event.OrderID)
    go updateInventory(event.Items)
    go sendConfirmationEmail(event.UserID, event.OrderID)
}
```

### Benefits of Message Queues
```
1. Decoupling: Services don't need to know about each other
2. Scalability: Process messages at different rates
3. Reliability: Messages persist until processed
4. Fault Tolerance: Failed messages can be retried
5. Load Leveling: Smooth out traffic spikes
```

## 📬 Message Queue Patterns

### 1. Point-to-Point (Queue)
```go
type MessageQueue struct {
    messages chan Message
    workers  int
}

type Message struct {
    ID      string
    Payload []byte
    Retry   int
    MaxRetry int
}

func NewMessageQueue(bufferSize, workers int) *MessageQueue {
    mq := &MessageQueue{
        messages: make(chan Message, bufferSize),
        workers:  workers,
    }
    
    // Start worker goroutines
    for i := 0; i < workers; i++ {
        go mq.worker(i)
    }
    
    return mq
}

func (mq *MessageQueue) Publish(message Message) error {
    select {
    case mq.messages <- message:
        return nil
    default:
        return ErrQueueFull
    }
}

func (mq *MessageQueue) worker(id int) {
    for message := range mq.messages {
        if err := mq.processMessage(message); err != nil {
            // Retry logic
            if message.Retry < message.MaxRetry {
                message.Retry++
                time.Sleep(time.Duration(message.Retry) * time.Second) // Exponential backoff
                mq.Publish(message)
            } else {
                // Send to dead letter queue
                mq.sendToDeadLetterQueue(message)
            }
        }
    }
}

// Use Case: Task processing, job queues
// Characteristics: One consumer per message
```

### 2. Publish-Subscribe (Topic)
```go
type PubSubBroker struct {
    topics map[string]*Topic
    mutex  sync.RWMutex
}

type Topic struct {
    name        string
    subscribers []chan Message
    mutex       sync.RWMutex
}

func (psb *PubSubBroker) Subscribe(topicName string) <-chan Message {
    psb.mutex.Lock()
    defer psb.mutex.Unlock()
    
    topic, exists := psb.topics[topicName]
    if !exists {
        topic = &Topic{
            name:        topicName,
            subscribers: make([]chan Message, 0),
        }
        psb.topics[topicName] = topic
    }
    
    subscriber := make(chan Message, 100)
    topic.mutex.Lock()
    topic.subscribers = append(topic.subscribers, subscriber)
    topic.mutex.Unlock()
    
    return subscriber
}

func (psb *PubSubBroker) Publish(topicName string, message Message) {
    psb.mutex.RLock()
    topic, exists := psb.topics[topicName]
    psb.mutex.RUnlock()
    
    if !exists {
        return
    }
    
    topic.mutex.RLock()
    defer topic.mutex.RUnlock()
    
    // Send to all subscribers
    for _, subscriber := range topic.subscribers {
        select {
        case subscriber <- message:
        default:
            // Subscriber buffer full, skip or handle
        }
    }
}

// Use Case: Event notifications, real-time updates
// Characteristics: Multiple consumers per message
```

## 🛠️ Message Queue Technologies

### RabbitMQ Implementation
```go
type RabbitMQClient struct {
    connection *amqp.Connection
    channel    *amqp.Channel
}

func NewRabbitMQClient(url string) (*RabbitMQClient, error) {
    conn, err := amqp.Dial(url)
    if err != nil {
        return nil, err
    }
    
    ch, err := conn.Channel()
    if err != nil {
        return nil, err
    }
    
    return &RabbitMQClient{
        connection: conn,
        channel:    ch,
    }, nil
}

// Work Queue Pattern
func (r *RabbitMQClient) SetupWorkQueue(queueName string) error {
    _, err := r.channel.QueueDeclare(
        queueName, // name
        true,      // durable
        false,     // delete when unused
        false,     // exclusive
        false,     // no-wait
        nil,       // arguments
    )
    return err
}

func (r *RabbitMQClient) PublishToQueue(queueName string, message []byte) error {
    return r.channel.Publish(
        "",        // exchange
        queueName, // routing key
        false,     // mandatory
        false,     // immediate
        amqp.Publishing{
            DeliveryMode: amqp.Persistent, // Make message persistent
            ContentType:  "application/json",
            Body:         message,
        })
}

func (r *RabbitMQClient) ConsumeFromQueue(queueName string, handler func([]byte) error) error {
    // Set QoS to process one message at a time
    r.channel.Qos(1, 0, false)
    
    msgs, err := r.channel.Consume(
        queueName, // queue
        "",        // consumer
        false,     // auto-ack (manual ack for reliability)
        false,     // exclusive
        false,     // no-local
        false,     // no-wait
        nil,       // args
    )
    if err != nil {
        return err
    }
    
    go func() {
        for msg := range msgs {
            if err := handler(msg.Body); err != nil {
                // Reject and requeue
                msg.Nack(false, true)
            } else {
                // Acknowledge successful processing
                msg.Ack(false)
            }
        }
    }()
    
    return nil
}

// Pub/Sub Pattern with Exchanges
func (r *RabbitMQClient) SetupPubSub(exchangeName string) error {
    return r.channel.ExchangeDeclare(
        exchangeName, // name
        "fanout",     // type
        true,         // durable
        false,        // auto-deleted
        false,        // internal
        false,        // no-wait
        nil,          // arguments
    )
}

func (r *RabbitMQClient) PublishToExchange(exchangeName string, message []byte) error {
    return r.channel.Publish(
        exchangeName, // exchange
        "",           // routing key (ignored for fanout)
        false,        // mandatory
        false,        // immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:        message,
        })
}
```

### Apache Kafka Implementation
```go
type KafkaProducer struct {
    producer sarama.SyncProducer
}

func NewKafkaProducer(brokers []string) (*KafkaProducer, error) {
    config := sarama.NewConfig()
    config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all replicas
    config.Producer.Retry.Max = 5
    config.Producer.Return.Successes = true
    
    producer, err := sarama.NewSyncProducer(brokers, config)
    if err != nil {
        return nil, err
    }
    
    return &KafkaProducer{producer: producer}, nil
}

func (kp *KafkaProducer) PublishMessage(topic string, key, value []byte) error {
    message := &sarama.ProducerMessage{
        Topic: topic,
        Key:   sarama.ByteEncoder(key),
        Value: sarama.ByteEncoder(value),
    }
    
    partition, offset, err := kp.producer.SendMessage(message)
    if err != nil {
        return err
    }
    
    log.Printf("Message sent to partition %d at offset %d", partition, offset)
    return nil
}

type KafkaConsumer struct {
    consumer sarama.ConsumerGroup
}

func NewKafkaConsumer(brokers []string, groupID string) (*KafkaConsumer, error) {
    config := sarama.NewConfig()
    config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
    config.Consumer.Offsets.Initial = sarama.OffsetOldest
    
    consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
    if err != nil {
        return nil, err
    }
    
    return &KafkaConsumer{consumer: consumer}, nil
}

func (kc *KafkaConsumer) ConsumeMessages(topics []string, handler func(string, []byte) error) error {
    ctx := context.Background()
    
    for {
        err := kc.consumer.Consume(ctx, topics, &ConsumerGroupHandler{
            handler: handler,
        })
        if err != nil {
            return err
        }
    }
}

type ConsumerGroupHandler struct {
    handler func(string, []byte) error
}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for message := range claim.Messages() {
        if err := h.handler(message.Topic, message.Value); err != nil {
            // Log error but continue processing
            log.Printf("Error processing message: %v", err)
        }
        
        // Mark message as processed
        session.MarkMessage(message, "")
    }
    
    return nil
}
```

### AWS SQS Implementation
```go
type SQSClient struct {
    client   *sqs.Client
    queueURL string
}

func NewSQSClient(region, queueURL string) *SQSClient {
    cfg, _ := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
    client := sqs.NewFromConfig(cfg)
    
    return &SQSClient{
        client:   client,
        queueURL: queueURL,
    }
}

func (s *SQSClient) SendMessage(messageBody string) error {
    _, err := s.client.SendMessage(context.TODO(), &sqs.SendMessageInput{
        QueueUrl:    aws.String(s.queueURL),
        MessageBody: aws.String(messageBody),
        DelaySeconds: 0,
    })
    
    return err
}

func (s *SQSClient) ReceiveMessages(handler func(string) error) error {
    for {
        result, err := s.client.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
            QueueUrl:            aws.String(s.queueURL),
            MaxNumberOfMessages: 10,
            WaitTimeSeconds:     20, // Long polling
        })
        
        if err != nil {
            return err
        }
        
        for _, message := range result.Messages {
            if err := handler(*message.Body); err != nil {
                // Message will become visible again after visibility timeout
                log.Printf("Error processing message: %v", err)
                continue
            }
            
            // Delete message after successful processing
            s.client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
                QueueUrl:      aws.String(s.queueURL),
                ReceiptHandle: message.ReceiptHandle,
            })
        }
    }
}

// Dead Letter Queue setup
func (s *SQSClient) SetupDLQ(mainQueueURL, dlqURL string, maxReceiveCount int) error {
    // Configure main queue to send failed messages to DLQ
    policy := map[string]interface{}{
        "deadLetterTargetArn": dlqURL,
        "maxReceiveCount":     maxReceiveCount,
    }
    
    policyJSON, _ := json.Marshal(policy)
    
    _, err := s.client.SetQueueAttributes(context.TODO(), &sqs.SetQueueAttributesInput{
        QueueUrl: aws.String(mainQueueURL),
        Attributes: map[string]string{
            "RedrivePolicy": string(policyJSON),
        },
    })
    
    return err
}
```

## 🌊 Event Streaming Patterns

### Event Sourcing
```go
type EventStore struct {
    events []Event
    mutex  sync.RWMutex
}

type Event struct {
    ID          string    `json:"id"`
    AggregateID string    `json:"aggregate_id"`
    Type        string    `json:"type"`
    Data        []byte    `json:"data"`
    Version     int       `json:"version"`
    Timestamp   time.Time `json:"timestamp"`
}

func (es *EventStore) AppendEvent(aggregateID string, eventType string, data interface{}) error {
    es.mutex.Lock()
    defer es.mutex.Unlock()
    
    eventData, err := json.Marshal(data)
    if err != nil {
        return err
    }
    
    event := Event{
        ID:          generateUUID(),
        AggregateID: aggregateID,
        Type:        eventType,
        Data:        eventData,
        Version:     es.getNextVersion(aggregateID),
        Timestamp:   time.Now(),
    }
    
    es.events = append(es.events, event)
    
    // Publish event to message queue
    eventBus.Publish(eventType, event)
    
    return nil
}

func (es *EventStore) GetEvents(aggregateID string) []Event {
    es.mutex.RLock()
    defer es.mutex.RUnlock()
    
    var result []Event
    for _, event := range es.events {
        if event.AggregateID == aggregateID {
            result = append(result, event)
        }
    }
    
    return result
}

// Rebuild aggregate state from events
func (es *EventStore) ReplayEvents(aggregateID string) (*UserAggregate, error) {
    events := es.GetEvents(aggregateID)
    
    user := &UserAggregate{ID: aggregateID}
    for _, event := range events {
        user.Apply(event)
    }
    
    return user, nil
}
```

### CQRS with Event Streaming
```go
// Command side - writes events
type UserCommandHandler struct {
    eventStore EventStore
}

func (uch *UserCommandHandler) CreateUser(cmd CreateUserCommand) error {
    // Validate command
    if err := cmd.Validate(); err != nil {
        return err
    }
    
    // Store event
    return uch.eventStore.AppendEvent(cmd.UserID, "UserCreated", UserCreatedEvent{
        UserID: cmd.UserID,
        Email:  cmd.Email,
        Name:   cmd.Name,
    })
}

// Query side - reads from projections
type UserQueryHandler struct {
    readDB Database
}

func (uqh *UserQueryHandler) GetUser(userID string) (*UserView, error) {
    return uqh.readDB.GetUserView(userID)
}

// Event handler updates read models
type UserProjectionHandler struct {
    readDB Database
}

func (uph *UserProjectionHandler) HandleUserCreated(event UserCreatedEvent) error {
    userView := &UserView{
        ID:        event.UserID,
        Email:     event.Email,
        Name:      event.Name,
        CreatedAt: time.Now(),
    }
    
    return uph.readDB.SaveUserView(userView)
}

func (uph *UserProjectionHandler) HandleUserUpdated(event UserUpdatedEvent) error {
    return uph.readDB.UpdateUserView(event.UserID, map[string]interface{}{
        "email": event.Email,
        "name":  event.Name,
    })
}
```

### Saga Pattern (Distributed Transactions)
```go
type OrderSaga struct {
    orderID string
    state   SagaState
    steps   []SagaStep
}

type SagaStep struct {
    Name            string
    Execute         func() error
    Compensate      func() error
    Executed        bool
    CompensationExecuted bool
}

func NewOrderSaga(orderID string) *OrderSaga {
    return &OrderSaga{
        orderID: orderID,
        state:   SagaStateStarted,
        steps: []SagaStep{
            {
                Name:       "ReserveInventory",
                Execute:    func() error { return reserveInventory(orderID) },
                Compensate: func() error { return releaseInventory(orderID) },
            },
            {
                Name:       "ProcessPayment",
                Execute:    func() error { return processPayment(orderID) },
                Compensate: func() error { return refundPayment(orderID) },
            },
            {
                Name:       "ShipOrder",
                Execute:    func() error { return shipOrder(orderID) },
                Compensate: func() error { return cancelShipment(orderID) },
            },
        },
    }
}

func (saga *OrderSaga) Execute() error {
    for i, step := range saga.steps {
        if err := step.Execute(); err != nil {
            // Compensate all executed steps in reverse order
            for j := i - 1; j >= 0; j-- {
                if saga.steps[j].Executed {
                    saga.steps[j].Compensate()
                    saga.steps[j].CompensationExecuted = true
                }
            }
            
            saga.state = SagaStateFailed
            return err
        }
        
        saga.steps[i].Executed = true
    }
    
    saga.state = SagaStateCompleted
    return nil
}

// Event-driven saga
func (saga *OrderSaga) HandleEvent(event Event) {
    switch event.Type {
    case "InventoryReserved":
        saga.processNextStep()
    case "PaymentProcessed":
        saga.processNextStep()
    case "OrderShipped":
        saga.complete()
    case "InventoryReservationFailed":
        saga.compensate()
    case "PaymentFailed":
        saga.compensate()
    }
}
```

## 📊 Message Queue Comparison

### Technology Comparison
```
Feature          | RabbitMQ | Apache Kafka | AWS SQS | Redis Pub/Sub
-----------------|----------|--------------|---------|---------------
Throughput       | Medium   | Very High    | High    | Very High
Latency          | Low      | Low          | Medium  | Very Low
Persistence      | Yes      | Yes          | Yes     | No
Ordering         | Yes      | Yes          | FIFO    | No
Scalability      | Medium   | Very High    | High    | High
Complexity       | Medium   | High         | Low     | Low
Use Case         | General  | Streaming    | Simple  | Real-time
```

### When to Use What
```go
// RabbitMQ: Complex routing, reliable delivery
type RabbitMQUseCase struct {
    UseCases []string
    Pros     []string
    Cons     []string
}

var RabbitMQ = RabbitMQUseCase{
    UseCases: []string{
        "Task queues with complex routing",
        "Request-response patterns",
        "Priority queues",
        "Delayed message delivery",
    },
    Pros: []string{
        "Flexible routing",
        "Management UI",
        "Multiple protocols",
        "Clustering support",
    },
    Cons: []string{
        "Lower throughput than Kafka",
        "More complex setup",
        "Memory usage can be high",
    },
}

// Kafka: High-throughput streaming
var Kafka = struct {
    UseCases []string
    Pros     []string
    Cons     []string
}{
    UseCases: []string{
        "Event streaming",
        "Log aggregation",
        "Real-time analytics",
        "Event sourcing",
    },
    Pros: []string{
        "Very high throughput",
        "Horizontal scaling",
        "Message replay",
        "Strong ordering guarantees",
    },
    Cons: []string{
        "Complex setup and operations",
        "Higher latency than Redis",
        "Requires ZooKeeper (older versions)",
    },
}

// SQS: Simple, managed queuing
var SQS = struct {
    UseCases []string
    Pros     []string
    Cons     []string
}{
    UseCases: []string{
        "Simple task queues",
        "Decoupling microservices",
        "Batch processing",
        "Dead letter queues",
    },
    Pros: []string{
        "Fully managed",
        "Auto-scaling",
        "Integration with AWS",
        "No operational overhead",
    },
    Cons: []string{
        "AWS vendor lock-in",
        "Limited message size",
        "No message ordering (standard queues)",
    },
}
```

## 🔗 Next Steps
- Study API design and communication patterns
- Learn about microservices architecture
- Explore event-driven architecture patterns
