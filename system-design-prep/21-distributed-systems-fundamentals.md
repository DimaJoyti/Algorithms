# 🌐 Distributed Systems Fundamentals

## 🎯 Core Concepts Overview

Distributed systems are collections of independent computers that appear to users as a single coherent system. Understanding these fundamentals is crucial for designing scalable, reliable, and maintainable systems.

### Key Characteristics
```
✅ Distributed System Properties:
- Multiple autonomous computers
- Connected by a network
- Coordinate activities through message passing
- No shared memory or global clock
- Independent failure modes
- Concurrent execution

✅ Goals:
- Scalability (handle increased load)
- Reliability (continue operating despite failures)
- Availability (system remains operational)
- Performance (low latency, high throughput)
- Consistency (data remains coherent)
```

## 🔄 CAP Theorem

### The Fundamental Trade-off
```
CAP Theorem States: In a distributed system, you can only guarantee 2 out of 3:

C - Consistency: All nodes see the same data simultaneously
A - Availability: System remains operational
P - Partition Tolerance: System continues despite network failures

Real-world implications:
- CP Systems: Sacrifice availability for consistency (traditional RDBMS)
- AP Systems: Sacrifice consistency for availability (DNS, web caches)
- CA Systems: Only possible without network partitions (single-node systems)
```

### CAP Theorem in Practice
```go
// Example: Implementing different CAP trade-offs

// CP System (Consistency + Partition Tolerance)
type CPSystem struct {
    nodes []Node
    quorum int // Majority required for operations
}

func (cp *CPSystem) Write(key, value string) error {
    successCount := 0
    errors := make([]error, 0)
    
    // Write to all nodes
    for _, node := range cp.nodes {
        if err := node.Write(key, value); err != nil {
            errors = append(errors, err)
        } else {
            successCount++
        }
    }
    
    // Require quorum for consistency
    if successCount >= cp.quorum {
        return nil // Success - consistent write
    }
    
    // Failed to achieve quorum - system becomes unavailable
    return fmt.Errorf("failed to achieve quorum: %d/%d nodes succeeded", 
        successCount, len(cp.nodes))
}

// AP System (Availability + Partition Tolerance)
type APSystem struct {
    nodes []Node
}

func (ap *APSystem) Write(key, value string) error {
    // Write to any available node
    for _, node := range ap.nodes {
        if node.IsAvailable() {
            return node.Write(key, value) // Always available, eventual consistency
        }
    }
    return fmt.Errorf("no nodes available")
}

func (ap *APSystem) Read(key string) (string, error) {
    // Read from any available node
    for _, node := range ap.nodes {
        if node.IsAvailable() {
            return node.Read(key) // May return stale data
        }
    }
    return "", fmt.Errorf("no nodes available")
}
```

## 🎭 Consistency Models

### Strong Consistency
```go
// Linearizability - Strongest consistency model
type LinearizableStore struct {
    data    map[string]string
    version map[string]int64
    mutex   sync.RWMutex
    clock   LogicalClock
}

func (ls *LinearizableStore) Write(key, value string) error {
    ls.mutex.Lock()
    defer ls.mutex.Unlock()
    
    // Increment logical clock
    ls.clock.Tick()
    
    // Write with timestamp
    ls.data[key] = value
    ls.version[key] = ls.clock.Time()
    
    return nil
}

func (ls *LinearizableStore) Read(key string) (string, error) {
    ls.mutex.RLock()
    defer ls.mutex.RUnlock()
    
    value, exists := ls.data[key]
    if !exists {
        return "", fmt.Errorf("key not found")
    }
    
    return value, nil
}

// Sequential Consistency
type SequentialStore struct {
    operations []Operation
    mutex      sync.Mutex
}

type Operation struct {
    Type      string // "read" or "write"
    Key       string
    Value     string
    Timestamp time.Time
    NodeID    string
}

func (ss *SequentialStore) Write(key, value, nodeID string) error {
    ss.mutex.Lock()
    defer ss.mutex.Unlock()
    
    op := Operation{
        Type:      "write",
        Key:       key,
        Value:     value,
        Timestamp: time.Now(),
        NodeID:    nodeID,
    }
    
    ss.operations = append(ss.operations, op)
    return nil
}
```

### Eventual Consistency
```go
// Vector Clocks for tracking causality
type VectorClock map[string]int

func (vc VectorClock) Increment(nodeID string) {
    vc[nodeID]++
}

func (vc VectorClock) Update(other VectorClock) {
    for nodeID, timestamp := range other {
        if vc[nodeID] < timestamp {
            vc[nodeID] = timestamp
        }
    }
}

func (vc VectorClock) HappensBefore(other VectorClock) bool {
    lessThanOrEqual := true
    strictlyLess := false
    
    for nodeID := range vc {
        if vc[nodeID] > other[nodeID] {
            lessThanOrEqual = false
            break
        }
        if vc[nodeID] < other[nodeID] {
            strictlyLess = true
        }
    }
    
    return lessThanOrEqual && strictlyLess
}

// Eventually Consistent Store
type EventuallyConsistentStore struct {
    nodeID    string
    data      map[string]VersionedValue
    vectorClock VectorClock
    peers     []Peer
}

type VersionedValue struct {
    Value string
    Clock VectorClock
}

func (ecs *EventuallyConsistentStore) Write(key, value string) error {
    // Increment local clock
    ecs.vectorClock.Increment(ecs.nodeID)
    
    // Store with vector clock
    ecs.data[key] = VersionedValue{
        Value: value,
        Clock: ecs.vectorClock.Copy(),
    }
    
    // Propagate to peers asynchronously
    go ecs.propagateUpdate(key, value, ecs.vectorClock.Copy())
    
    return nil
}

func (ecs *EventuallyConsistentStore) Read(key string) (string, error) {
    versionedValue, exists := ecs.data[key]
    if !exists {
        return "", fmt.Errorf("key not found")
    }
    
    return versionedValue.Value, nil
}

func (ecs *EventuallyConsistentStore) propagateUpdate(key, value string, clock VectorClock) {
    update := Update{
        Key:   key,
        Value: value,
        Clock: clock,
    }
    
    for _, peer := range ecs.peers {
        go peer.SendUpdate(update)
    }
}

func (ecs *EventuallyConsistentStore) ReceiveUpdate(update Update) {
    existing, exists := ecs.data[update.Key]
    
    if !exists {
        // New key, accept update
        ecs.data[update.Key] = VersionedValue{
            Value: update.Value,
            Clock: update.Clock,
        }
        ecs.vectorClock.Update(update.Clock)
        return
    }
    
    // Resolve conflicts using vector clocks
    if update.Clock.HappensBefore(existing.Clock) {
        // Update is older, ignore
        return
    } else if existing.Clock.HappensBefore(update.Clock) {
        // Update is newer, accept
        ecs.data[update.Key] = VersionedValue{
            Value: update.Value,
            Clock: update.Clock,
        }
        ecs.vectorClock.Update(update.Clock)
    } else {
        // Concurrent updates, need conflict resolution
        ecs.resolveConflict(update.Key, existing, VersionedValue{
            Value: update.Value,
            Clock: update.Clock,
        })
    }
}
```

## 🔐 Distributed Consensus

### Two-Phase Commit (2PC)
```go
// Two-Phase Commit Coordinator
type TwoPhaseCommitCoordinator struct {
    participants []Participant
    transactionID string
    state        TransactionState
}

type TransactionState int

const (
    Preparing TransactionState = iota
    Prepared
    Committed
    Aborted
)

func (tpc *TwoPhaseCommitCoordinator) Commit(transaction Transaction) error {
    tpc.transactionID = transaction.ID
    tpc.state = Preparing
    
    // Phase 1: Prepare
    prepareResponses := make(chan bool, len(tpc.participants))
    
    for _, participant := range tpc.participants {
        go func(p Participant) {
            response := p.Prepare(transaction)
            prepareResponses <- response
        }(participant)
    }
    
    // Collect prepare responses
    preparedCount := 0
    for i := 0; i < len(tpc.participants); i++ {
        if <-prepareResponses {
            preparedCount++
        }
    }
    
    // Phase 2: Commit or Abort
    if preparedCount == len(tpc.participants) {
        // All participants prepared, commit
        tpc.state = Committed
        return tpc.sendCommitToAll(transaction)
    } else {
        // Some participants failed to prepare, abort
        tpc.state = Aborted
        return tpc.sendAbortToAll(transaction)
    }
}

func (tpc *TwoPhaseCommitCoordinator) sendCommitToAll(transaction Transaction) error {
    var errors []error
    
    for _, participant := range tpc.participants {
        if err := participant.Commit(transaction); err != nil {
            errors = append(errors, err)
        }
    }
    
    if len(errors) > 0 {
        return fmt.Errorf("commit failed on some participants: %v", errors)
    }
    
    return nil
}

// Participant in 2PC
type TwoPhaseCommitParticipant struct {
    id           string
    transactions map[string]Transaction
    state        map[string]TransactionState
    storage      Storage
}

func (tpcp *TwoPhaseCommitParticipant) Prepare(transaction Transaction) bool {
    // Check if transaction can be committed
    if !tpcp.canCommit(transaction) {
        return false
    }
    
    // Write transaction to stable storage
    if err := tpcp.storage.WriteTransaction(transaction); err != nil {
        return false
    }
    
    // Mark as prepared
    tpcp.transactions[transaction.ID] = transaction
    tpcp.state[transaction.ID] = Prepared
    
    return true
}

func (tpcp *TwoPhaseCommitParticipant) Commit(transaction Transaction) error {
    // Apply transaction changes
    if err := tpcp.applyTransaction(transaction); err != nil {
        return err
    }
    
    // Mark as committed
    tpcp.state[transaction.ID] = Committed
    
    // Clean up transaction log
    delete(tpcp.transactions, transaction.ID)
    delete(tpcp.state, transaction.ID)
    
    return nil
}
```

## 💥 Failure Detection & Recovery

### Failure Detection Mechanisms
```go
// Heartbeat-based failure detection
type HeartbeatFailureDetector struct {
    nodes           map[string]*NodeStatus
    heartbeatInterval time.Duration
    timeoutThreshold  time.Duration
    suspicionLevel    int
    callbacks         []FailureCallback
    mutex            sync.RWMutex
}

type NodeStatus struct {
    NodeID        string
    LastHeartbeat time.Time
    IsAlive       bool
    SuspicionLevel int
    ConsecutiveMisses int
}

func (hfd *HeartbeatFailureDetector) StartMonitoring(nodeID string) {
    hfd.mutex.Lock()
    hfd.nodes[nodeID] = &NodeStatus{
        NodeID:        nodeID,
        LastHeartbeat: time.Now(),
        IsAlive:       true,
    }
    hfd.mutex.Unlock()

    // Start heartbeat monitoring
    go hfd.monitorNode(nodeID)
}

func (hfd *HeartbeatFailureDetector) monitorNode(nodeID string) {
    ticker := time.NewTicker(hfd.heartbeatInterval)
    defer ticker.Stop()

    for range ticker.C {
        hfd.mutex.Lock()
        node, exists := hfd.nodes[nodeID]
        if !exists {
            hfd.mutex.Unlock()
            return
        }

        timeSinceLastHeartbeat := time.Since(node.LastHeartbeat)

        if timeSinceLastHeartbeat > hfd.timeoutThreshold {
            node.ConsecutiveMisses++
            node.SuspicionLevel++

            if node.ConsecutiveMisses >= 3 && node.IsAlive {
                // Mark node as failed
                node.IsAlive = false
                hfd.mutex.Unlock()

                // Notify failure callbacks
                for _, callback := range hfd.callbacks {
                    go callback.OnNodeFailure(nodeID)
                }
                continue
            }
        } else {
            // Reset counters on successful heartbeat
            node.ConsecutiveMisses = 0
            node.SuspicionLevel = 0
            if !node.IsAlive {
                node.IsAlive = true
                // Notify recovery callbacks
                for _, callback := range hfd.callbacks {
                    go callback.OnNodeRecovery(nodeID)
                }
            }
        }
        hfd.mutex.Unlock()
    }
}

func (hfd *HeartbeatFailureDetector) ReceiveHeartbeat(nodeID string) {
    hfd.mutex.Lock()
    defer hfd.mutex.Unlock()

    if node, exists := hfd.nodes[nodeID]; exists {
        node.LastHeartbeat = time.Now()
    }
}

// Phi Accrual Failure Detector (more sophisticated)
type PhiAccrualFailureDetector struct {
    nodes           map[string]*PhiNodeStatus
    threshold       float64
    maxSampleSize   int
    minStdDeviation float64
}

type PhiNodeStatus struct {
    NodeID           string
    ArrivalIntervals []time.Duration
    LastArrival      time.Time
    Mean             float64
    Variance         float64
}

func (pafd *PhiAccrualFailureDetector) ReceiveHeartbeat(nodeID string) {
    now := time.Now()

    node, exists := pafd.nodes[nodeID]
    if !exists {
        node = &PhiNodeStatus{
            NodeID:      nodeID,
            LastArrival: now,
        }
        pafd.nodes[nodeID] = node
        return
    }

    // Calculate interval since last heartbeat
    interval := now.Sub(node.LastArrival)
    node.LastArrival = now

    // Add to arrival intervals
    node.ArrivalIntervals = append(node.ArrivalIntervals, interval)
    if len(node.ArrivalIntervals) > pafd.maxSampleSize {
        node.ArrivalIntervals = node.ArrivalIntervals[1:]
    }

    // Update statistics
    pafd.updateStatistics(node)
}

func (pafd *PhiAccrualFailureDetector) GetPhi(nodeID string) float64 {
    node, exists := pafd.nodes[nodeID]
    if !exists {
        return math.Inf(1) // Infinite phi for unknown nodes
    }

    if len(node.ArrivalIntervals) < 2 {
        return 0.0 // Not enough data
    }

    timeSinceLastHeartbeat := time.Since(node.LastArrival)

    // Calculate phi value
    phi := pafd.calculatePhi(timeSinceLastHeartbeat, node.Mean, node.Variance)
    return phi
}

func (pafd *PhiAccrualFailureDetector) IsNodeSuspected(nodeID string) bool {
    phi := pafd.GetPhi(nodeID)
    return phi > pafd.threshold
}

func (pafd *PhiAccrualFailureDetector) calculatePhi(interval time.Duration, mean, variance float64) float64 {
    if variance < pafd.minStdDeviation*pafd.minStdDeviation {
        variance = pafd.minStdDeviation * pafd.minStdDeviation
    }

    stdDev := math.Sqrt(variance)
    y := (interval.Seconds() - mean) / stdDev
    e := math.Exp(-y * (1.5976 + 0.070566*y*y))

    if interval.Seconds() > mean {
        return -math.Log10(e / (1.0 + e))
    }
    return -math.Log10(1.0 - 1.0/(1.0+e))
}
```

### Recovery Strategies
```go
// Automatic failover with leader election
type FailoverManager struct {
    nodes          []Node
    currentLeader  string
    electionInProgress bool
    callbacks      []LeadershipCallback
    mutex          sync.RWMutex
}

func (fm *FailoverManager) OnNodeFailure(nodeID string) {
    fm.mutex.Lock()
    defer fm.mutex.Unlock()

    if nodeID == fm.currentLeader {
        // Leader failed, start election
        fm.currentLeader = ""
        if !fm.electionInProgress {
            go fm.startLeaderElection()
        }
    }

    // Remove failed node from active nodes
    fm.removeNode(nodeID)
}

func (fm *FailoverManager) startLeaderElection() {
    fm.mutex.Lock()
    fm.electionInProgress = true
    fm.mutex.Unlock()

    defer func() {
        fm.mutex.Lock()
        fm.electionInProgress = false
        fm.mutex.Unlock()
    }()

    // Simple bully algorithm
    candidates := fm.getHealthyNodes()
    if len(candidates) == 0 {
        return
    }

    // Select node with highest ID as leader
    sort.Strings(candidates)
    newLeader := candidates[len(candidates)-1]

    fm.mutex.Lock()
    fm.currentLeader = newLeader
    fm.mutex.Unlock()

    // Notify callbacks
    for _, callback := range fm.callbacks {
        go callback.OnLeaderElected(newLeader)
    }
}

// Circuit breaker for handling cascading failures
type CircuitBreaker struct {
    name           string
    maxFailures    int
    resetTimeout   time.Duration
    state          CircuitState
    failureCount   int
    lastFailureTime time.Time
    mutex          sync.RWMutex
}

type CircuitState int

const (
    CircuitClosed CircuitState = iota
    CircuitOpen
    CircuitHalfOpen
)

func (cb *CircuitBreaker) Call(operation func() error) error {
    cb.mutex.Lock()
    state := cb.state
    cb.mutex.Unlock()

    switch state {
    case CircuitOpen:
        // Check if we should transition to half-open
        if time.Since(cb.lastFailureTime) > cb.resetTimeout {
            cb.mutex.Lock()
            cb.state = CircuitHalfOpen
            cb.mutex.Unlock()
            return cb.callInHalfOpenState(operation)
        }
        return fmt.Errorf("circuit breaker is open")

    case CircuitHalfOpen:
        return cb.callInHalfOpenState(operation)

    default: // CircuitClosed
        return cb.callInClosedState(operation)
    }
}

func (cb *CircuitBreaker) callInClosedState(operation func() error) error {
    err := operation()

    cb.mutex.Lock()
    defer cb.mutex.Unlock()

    if err != nil {
        cb.failureCount++
        cb.lastFailureTime = time.Now()

        if cb.failureCount >= cb.maxFailures {
            cb.state = CircuitOpen
        }
    } else {
        cb.failureCount = 0
    }

    return err
}

func (cb *CircuitBreaker) callInHalfOpenState(operation func() error) error {
    err := operation()

    cb.mutex.Lock()
    defer cb.mutex.Unlock()

    if err != nil {
        cb.state = CircuitOpen
        cb.failureCount++
        cb.lastFailureTime = time.Now()
    } else {
        cb.state = CircuitClosed
        cb.failureCount = 0
    }

    return err
}
```

## 🔄 Replication Strategies

### Master-Slave Replication
```go
// Master-Slave replication system
type MasterSlaveReplicator struct {
    master   *MasterNode
    slaves   []*SlaveNode
    binlog   BinaryLog
    health   HealthChecker
}

type MasterNode struct {
    id       string
    storage  Storage
    binlog   BinaryLog
    slaves   []*SlaveNode
}

func (mn *MasterNode) Write(key, value string) error {
    // Write to local storage
    if err := mn.storage.Set(key, value); err != nil {
        return err
    }

    // Write to binary log
    logEntry := BinaryLogEntry{
        Operation: "SET",
        Key:       key,
        Value:     value,
        Timestamp: time.Now(),
        LSN:       mn.binlog.GetNextLSN(),
    }

    if err := mn.binlog.Append(logEntry); err != nil {
        // Rollback local write
        mn.storage.Delete(key)
        return err
    }

    // Replicate to slaves asynchronously
    go mn.replicateToSlaves(logEntry)

    return nil
}

func (mn *MasterNode) replicateToSlaves(entry BinaryLogEntry) {
    for _, slave := range mn.slaves {
        go func(s *SlaveNode) {
            if err := s.ApplyLogEntry(entry); err != nil {
                log.Printf("Failed to replicate to slave %s: %v", s.ID, err)
                // Mark slave as lagging
                s.MarkAsLagging()
            }
        }(slave)
    }
}

type SlaveNode struct {
    ID           string
    storage      Storage
    lastAppliedLSN int64
    isLagging    bool
    master       *MasterNode
}

func (sn *SlaveNode) ApplyLogEntry(entry BinaryLogEntry) error {
    // Check LSN ordering
    if entry.LSN != sn.lastAppliedLSN+1 {
        return fmt.Errorf("LSN gap detected: expected %d, got %d",
            sn.lastAppliedLSN+1, entry.LSN)
    }

    // Apply operation
    switch entry.Operation {
    case "SET":
        if err := sn.storage.Set(entry.Key, entry.Value); err != nil {
            return err
        }
    case "DELETE":
        if err := sn.storage.Delete(entry.Key); err != nil {
            return err
        }
    }

    sn.lastAppliedLSN = entry.LSN
    sn.isLagging = false

    return nil
}

func (sn *SlaveNode) CatchUp() error {
    // Request missing log entries from master
    missingEntries, err := sn.master.GetLogEntriesSince(sn.lastAppliedLSN)
    if err != nil {
        return err
    }

    // Apply missing entries in order
    for _, entry := range missingEntries {
        if err := sn.ApplyLogEntry(entry); err != nil {
            return err
        }
    }

    return nil
}
```

### Multi-Master Replication
```go
// Multi-master replication with conflict resolution
type MultiMasterReplicator struct {
    nodes           []*MasterNode
    conflictResolver ConflictResolver
    vectorClock     VectorClock
    nodeID          string
}

func (mmr *MultiMasterReplicator) Write(key, value string) error {
    // Increment vector clock
    mmr.vectorClock.Increment(mmr.nodeID)

    // Create versioned write
    versionedWrite := VersionedWrite{
        Key:         key,
        Value:       value,
        VectorClock: mmr.vectorClock.Copy(),
        NodeID:      mmr.nodeID,
        Timestamp:   time.Now(),
    }

    // Write locally
    localNode := mmr.getLocalNode()
    if err := localNode.WriteVersioned(versionedWrite); err != nil {
        return err
    }

    // Propagate to other masters
    go mmr.propagateWrite(versionedWrite)

    return nil
}

func (mmr *MultiMasterReplicator) propagateWrite(write VersionedWrite) {
    for _, node := range mmr.nodes {
        if node.ID != mmr.nodeID {
            go func(n *MasterNode) {
                if err := n.ReceiveWrite(write); err != nil {
                    log.Printf("Failed to propagate write to node %s: %v", n.ID, err)
                }
            }(node)
        }
    }
}

func (mn *MasterNode) ReceiveWrite(write VersionedWrite) error {
    // Check if we already have this write
    existing, exists := mn.getVersionedValue(write.Key)

    if !exists {
        // New key, accept write
        return mn.WriteVersioned(write)
    }

    // Resolve conflict using vector clocks
    if write.VectorClock.HappensBefore(existing.VectorClock) {
        // Write is older, ignore
        return nil
    } else if existing.VectorClock.HappensBefore(write.VectorClock) {
        // Write is newer, accept
        return mn.WriteVersioned(write)
    } else {
        // Concurrent writes, need conflict resolution
        resolved := mn.conflictResolver.Resolve(existing, write)
        return mn.WriteVersioned(resolved)
    }
}

// Last-Writer-Wins conflict resolution
type LWWConflictResolver struct{}

func (lww *LWWConflictResolver) Resolve(existing, incoming VersionedWrite) VersionedWrite {
    if incoming.Timestamp.After(existing.Timestamp) {
        return incoming
    }
    return existing
}

// Application-specific conflict resolution
type ApplicationConflictResolver struct {
    mergeFunc func(existing, incoming VersionedWrite) VersionedWrite
}

func (acr *ApplicationConflictResolver) Resolve(existing, incoming VersionedWrite) VersionedWrite {
    return acr.mergeFunc(existing, incoming)
}
```

## 🔀 Partitioning Strategies

### Horizontal Partitioning (Sharding)
```go
// Hash-based sharding
type HashSharding struct {
    shards    []Shard
    hashFunc  hash.Hash64
    replicas  int
}

func (hs *HashSharding) GetShard(key string) Shard {
    hs.hashFunc.Reset()
    hs.hashFunc.Write([]byte(key))
    hash := hs.hashFunc.Sum64()

    shardIndex := hash % uint64(len(hs.shards))
    return hs.shards[shardIndex]
}

func (hs *HashSharding) Write(key, value string) error {
    shard := hs.GetShard(key)
    return shard.Write(key, value)
}

func (hs *HashSharding) Read(key string) (string, error) {
    shard := hs.GetShard(key)
    return shard.Read(key)
}

// Consistent hashing for better load distribution
type ConsistentHashing struct {
    ring     map[uint64]string // hash -> node
    nodes    map[string]bool   // node -> active
    replicas int
    hashFunc hash.Hash64
}

func (ch *ConsistentHashing) AddNode(nodeID string) {
    ch.nodes[nodeID] = true

    // Add virtual nodes for better distribution
    for i := 0; i < ch.replicas; i++ {
        virtualNodeID := fmt.Sprintf("%s:%d", nodeID, i)
        hash := ch.hash(virtualNodeID)
        ch.ring[hash] = nodeID
    }
}

func (ch *ConsistentHashing) RemoveNode(nodeID string) {
    delete(ch.nodes, nodeID)

    // Remove virtual nodes
    for i := 0; i < ch.replicas; i++ {
        virtualNodeID := fmt.Sprintf("%s:%d", nodeID, i)
        hash := ch.hash(virtualNodeID)
        delete(ch.ring, hash)
    }
}

func (ch *ConsistentHashing) GetNode(key string) string {
    if len(ch.ring) == 0 {
        return ""
    }

    hash := ch.hash(key)

    // Find first node clockwise from hash
    var keys []uint64
    for k := range ch.ring {
        keys = append(keys, k)
    }
    sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

    for _, k := range keys {
        if k >= hash {
            return ch.ring[k]
        }
    }

    // Wrap around to first node
    return ch.ring[keys[0]]
}

func (ch *ConsistentHashing) hash(key string) uint64 {
    ch.hashFunc.Reset()
    ch.hashFunc.Write([]byte(key))
    return ch.hashFunc.Sum64()
}

// Range-based partitioning
type RangePartitioning struct {
    partitions []RangePartition
}

type RangePartition struct {
    StartKey string
    EndKey   string
    NodeID   string
}

func (rp *RangePartitioning) GetPartition(key string) *RangePartition {
    for _, partition := range rp.partitions {
        if key >= partition.StartKey && key < partition.EndKey {
            return &partition
        }
    }
    return nil
}

func (rp *RangePartitioning) SplitPartition(partitionIndex int, splitKey string) error {
    if partitionIndex >= len(rp.partitions) {
        return fmt.Errorf("invalid partition index")
    }

    oldPartition := rp.partitions[partitionIndex]

    // Create two new partitions
    leftPartition := RangePartition{
        StartKey: oldPartition.StartKey,
        EndKey:   splitKey,
        NodeID:   oldPartition.NodeID,
    }

    rightPartition := RangePartition{
        StartKey: splitKey,
        EndKey:   oldPartition.EndKey,
        NodeID:   rp.allocateNewNode(),
    }

    // Replace old partition with new ones
    rp.partitions = append(rp.partitions[:partitionIndex],
        append([]RangePartition{leftPartition, rightPartition},
            rp.partitions[partitionIndex+1:]...)...)

    return nil
}
```

## 🎯 Summary: Distributed Systems Fundamentals

### Core Concepts Mastered
- **CAP Theorem** - Understanding trade-offs between Consistency, Availability, and Partition tolerance
- **Consistency Models** - Strong consistency (linearizability) vs eventual consistency with vector clocks
- **Consensus Algorithms** - Two-phase commit for distributed transactions
- **Failure Detection** - Heartbeat and Phi accrual failure detectors
- **Replication Strategies** - Master-slave and multi-master replication with conflict resolution
- **Partitioning** - Hash-based, consistent hashing, and range-based partitioning

### Key Patterns Implemented
- **Vector Clocks** - Tracking causality in distributed systems
- **Circuit Breakers** - Preventing cascading failures
- **Leader Election** - Automatic failover and coordination
- **Conflict Resolution** - Last-writer-wins and application-specific strategies
- **Consistent Hashing** - Load balancing with minimal reshuffling

You now have a solid foundation in distributed systems theory and practice! 🚀
