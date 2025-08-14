# 🛡️ Disaster Recovery & Reliability

## 🎯 Disaster Recovery Overview

Disaster Recovery (DR) and Reliability engineering focus on building systems that can withstand failures, recover from disasters, and maintain high availability even under adverse conditions.

### Key Concepts
```
✅ Reliability Metrics:
- Mean Time Between Failures (MTBF)
- Mean Time To Recovery (MTTR)
- Recovery Time Objective (RTO)
- Recovery Point Objective (RPO)
- Service Level Objectives (SLOs)

✅ Failure Types:
- Hardware failures (server, disk, network)
- Software failures (bugs, memory leaks)
- Human errors (misconfigurations, accidents)
- Natural disasters (earthquakes, floods)
- Cyber attacks (DDoS, data breaches)

✅ Availability Levels:
- 99.9% (8.77 hours downtime/year)
- 99.95% (4.38 hours downtime/year)
- 99.99% (52.6 minutes downtime/year)
- 99.999% (5.26 minutes downtime/year)
```

## 🔄 Backup and Recovery Strategies

### Automated Backup Systems
```go
// Comprehensive backup management system
type BackupManager struct {
    strategies    []BackupStrategy
    scheduler     BackupScheduler
    storage       BackupStorage
    encryption    EncryptionService
    verification  BackupVerifier
    retention     RetentionPolicy
    monitoring    BackupMonitor
}

type BackupStrategy interface {
    CreateBackup(ctx context.Context, source DataSource) (*Backup, error)
    RestoreBackup(ctx context.Context, backup *Backup, target DataTarget) error
    ValidateBackup(ctx context.Context, backup *Backup) error
}

// Full backup strategy
type FullBackupStrategy struct {
    compression CompressionLevel
    encryption  bool
    checksum    bool
}

func (fbs *FullBackupStrategy) CreateBackup(ctx context.Context, source DataSource) (*Backup, error) {
    backup := &Backup{
        ID:        generateBackupID(),
        Type:      BackupTypeFull,
        Source:    source.GetIdentifier(),
        StartTime: time.Now(),
        Status:    BackupStatusInProgress,
    }
    
    // Create backup stream
    reader, err := source.CreateReader(ctx)
    if err != nil {
        return nil, err
    }
    defer reader.Close()
    
    // Apply compression if enabled
    var dataStream io.Reader = reader
    if fbs.compression != CompressionNone {
        dataStream = fbs.compressStream(dataStream)
    }
    
    // Apply encryption if enabled
    if fbs.encryption {
        dataStream, err = fbs.encryptStream(dataStream)
        if err != nil {
            return nil, err
        }
    }
    
    // Calculate checksum while writing
    var checksum string
    if fbs.checksum {
        hasher := sha256.New()
        dataStream = io.TeeReader(dataStream, hasher)
        defer func() {
            checksum = hex.EncodeToString(hasher.Sum(nil))
            backup.Checksum = checksum
        }()
    }
    
    // Write to backup storage
    writer, err := fbs.storage.CreateWriter(backup.ID)
    if err != nil {
        return nil, err
    }
    defer writer.Close()
    
    bytesWritten, err := io.Copy(writer, dataStream)
    if err != nil {
        backup.Status = BackupStatusFailed
        backup.Error = err.Error()
        return backup, err
    }
    
    backup.Size = bytesWritten
    backup.EndTime = time.Now()
    backup.Duration = backup.EndTime.Sub(backup.StartTime)
    backup.Status = BackupStatusCompleted
    
    return backup, nil
}

// Incremental backup strategy
type IncrementalBackupStrategy struct {
    baseBackup     *Backup
    changeTracker  ChangeTracker
    compression    CompressionLevel
}

func (ibs *IncrementalBackupStrategy) CreateBackup(ctx context.Context, source DataSource) (*Backup, error) {
    // Get changes since last backup
    lastBackupTime := ibs.getLastBackupTime()
    changes, err := ibs.changeTracker.GetChangesSince(lastBackupTime)
    if err != nil {
        return nil, err
    }
    
    backup := &Backup{
        ID:         generateBackupID(),
        Type:       BackupTypeIncremental,
        Source:     source.GetIdentifier(),
        BaseBackup: ibs.baseBackup.ID,
        StartTime:  time.Now(),
        Status:     BackupStatusInProgress,
    }
    
    // Create incremental backup data
    incrementalData := &IncrementalBackupData{
        Changes:   changes,
        Timestamp: time.Now(),
    }
    
    // Serialize and compress
    data, err := json.Marshal(incrementalData)
    if err != nil {
        return nil, err
    }
    
    if ibs.compression != CompressionNone {
        data = ibs.compress(data)
    }
    
    // Write to storage
    writer, err := ibs.storage.CreateWriter(backup.ID)
    if err != nil {
        return nil, err
    }
    defer writer.Close()
    
    _, err = writer.Write(data)
    if err != nil {
        backup.Status = BackupStatusFailed
        return backup, err
    }
    
    backup.Size = int64(len(data))
    backup.EndTime = time.Now()
    backup.Duration = backup.EndTime.Sub(backup.StartTime)
    backup.Status = BackupStatusCompleted
    
    return backup, nil
}

// Point-in-time recovery
type PointInTimeRecovery struct {
    backupManager   *BackupManager
    transactionLog  TransactionLog
    recoveryEngine  RecoveryEngine
}

func (pitr *PointInTimeRecovery) RecoverToPoint(ctx context.Context, targetTime time.Time) error {
    // Find the latest full backup before target time
    fullBackup, err := pitr.backupManager.FindLatestFullBackup(targetTime)
    if err != nil {
        return err
    }
    
    // Find all incremental backups between full backup and target time
    incrementalBackups, err := pitr.backupManager.FindIncrementalBackups(
        fullBackup.EndTime, targetTime)
    if err != nil {
        return err
    }
    
    // Restore full backup
    if err := pitr.recoveryEngine.RestoreFullBackup(ctx, fullBackup); err != nil {
        return err
    }
    
    // Apply incremental backups in order
    for _, backup := range incrementalBackups {
        if err := pitr.recoveryEngine.ApplyIncrementalBackup(ctx, backup); err != nil {
            return err
        }
    }
    
    // Apply transaction log entries from last backup to target time
    logEntries, err := pitr.transactionLog.GetEntriesBetween(
        incrementalBackups[len(incrementalBackups)-1].EndTime, targetTime)
    if err != nil {
        return err
    }
    
    for _, entry := range logEntries {
        if err := pitr.recoveryEngine.ApplyLogEntry(ctx, entry); err != nil {
            return err
        }
    }
    
    return nil
}

// Cross-region backup replication
type CrossRegionReplication struct {
    primaryRegion   string
    backupRegions   []string
    replicator      BackupReplicator
    consistency     ConsistencyChecker
    encryption      RegionalEncryption
}

func (crr *CrossRegionReplication) ReplicateBackup(ctx context.Context, backup *Backup) error {
    var errors []error
    
    for _, region := range crr.backupRegions {
        go func(targetRegion string) {
            if err := crr.replicateToRegion(ctx, backup, targetRegion); err != nil {
                errors = append(errors, fmt.Errorf("replication to %s failed: %w", targetRegion, err))
            }
        }(region)
    }
    
    // Wait for all replications to complete
    time.Sleep(30 * time.Second) // TODO: Use proper synchronization
    
    if len(errors) > 0 {
        return fmt.Errorf("replication errors: %v", errors)
    }
    
    return nil
}

func (crr *CrossRegionReplication) replicateToRegion(ctx context.Context, backup *Backup, region string) error {
    // Encrypt for target region
    encryptedBackup, err := crr.encryption.EncryptForRegion(backup, region)
    if err != nil {
        return err
    }
    
    // Transfer to target region
    if err := crr.replicator.Transfer(ctx, encryptedBackup, region); err != nil {
        return err
    }
    
    // Verify consistency
    return crr.consistency.VerifyReplication(ctx, backup.ID, region)
}
```

### Backup Verification and Testing
```go
// Automated backup verification system
type BackupVerifier struct {
    testEnvironment TestEnvironment
    validator       BackupValidator
    scheduler       VerificationScheduler
    reporter        VerificationReporter
}

func (bv *BackupVerifier) VerifyBackup(ctx context.Context, backup *Backup) (*VerificationResult, error) {
    result := &VerificationResult{
        BackupID:  backup.ID,
        StartTime: time.Now(),
        Tests:     make([]VerificationTest, 0),
    }
    
    // Test 1: Integrity check
    integrityTest := bv.performIntegrityCheck(ctx, backup)
    result.Tests = append(result.Tests, integrityTest)
    
    // Test 2: Restore test in isolated environment
    restoreTest := bv.performRestoreTest(ctx, backup)
    result.Tests = append(result.Tests, restoreTest)
    
    // Test 3: Data validation
    dataValidationTest := bv.performDataValidation(ctx, backup)
    result.Tests = append(result.Tests, dataValidationTest)
    
    // Test 4: Performance test
    performanceTest := bv.performPerformanceTest(ctx, backup)
    result.Tests = append(result.Tests, performanceTest)
    
    result.EndTime = time.Now()
    result.Duration = result.EndTime.Sub(result.StartTime)
    result.Success = bv.allTestsPassed(result.Tests)
    
    return result, nil
}

func (bv *BackupVerifier) performRestoreTest(ctx context.Context, backup *Backup) VerificationTest {
    test := VerificationTest{
        Name:      "Restore Test",
        StartTime: time.Now(),
    }
    
    // Create isolated test environment
    testEnv, err := bv.testEnvironment.CreateIsolatedEnvironment()
    if err != nil {
        test.Error = err.Error()
        test.Success = false
        return test
    }
    defer testEnv.Cleanup()
    
    // Attempt restore
    restoreStart := time.Now()
    err = testEnv.RestoreBackup(ctx, backup)
    restoreTime := time.Since(restoreStart)
    
    if err != nil {
        test.Error = err.Error()
        test.Success = false
        return test
    }
    
    // Validate restored data
    validationErr := testEnv.ValidateRestoredData()
    if validationErr != nil {
        test.Error = validationErr.Error()
        test.Success = false
        return test
    }
    
    test.Success = true
    test.Duration = restoreTime
    test.EndTime = time.Now()
    test.Metadata = map[string]interface{}{
        "restore_time_seconds": restoreTime.Seconds(),
        "data_size_bytes":      backup.Size,
    }
    
    return test
}

// Backup retention and lifecycle management
type BackupRetentionManager struct {
    policies    []RetentionPolicy
    storage     BackupStorage
    scheduler   RetentionScheduler
    auditor     RetentionAuditor
}

type RetentionPolicy struct {
    Name            string
    BackupType      BackupType
    RetentionPeriod time.Duration
    ArchiveAfter    time.Duration
    DeleteAfter     time.Duration
    MinimumCopies   int
}

func (brm *BackupRetentionManager) ApplyRetentionPolicies(ctx context.Context) error {
    for _, policy := range brm.policies {
        if err := brm.applyPolicy(ctx, policy); err != nil {
            log.Printf("Failed to apply retention policy %s: %v", policy.Name, err)
        }
    }
    return nil
}

func (brm *BackupRetentionManager) applyPolicy(ctx context.Context, policy RetentionPolicy) error {
    // Get backups matching this policy
    backups, err := brm.storage.GetBackupsByType(policy.BackupType)
    if err != nil {
        return err
    }
    
    now := time.Now()
    
    for _, backup := range backups {
        age := now.Sub(backup.CreatedAt)
        
        // Archive old backups
        if age > policy.ArchiveAfter && backup.Status != BackupStatusArchived {
            if err := brm.archiveBackup(ctx, backup); err != nil {
                log.Printf("Failed to archive backup %s: %v", backup.ID, err)
            }
        }
        
        // Delete very old backups (but keep minimum copies)
        if age > policy.DeleteAfter {
            activeBackups := brm.countActiveBackups(backups, policy.BackupType)
            if activeBackups > policy.MinimumCopies {
                if err := brm.deleteBackup(ctx, backup); err != nil {
                    log.Printf("Failed to delete backup %s: %v", backup.ID, err)
                }
            }
        }
    }
    
    return nil
}
```

## ⚡ Failover and High Availability

### Automatic Failover Systems
```go
// High availability cluster manager
type HAClusterManager struct {
    nodes           []ClusterNode
    leader          string
    failureDetector FailureDetector
    consensus       ConsensusManager
    loadBalancer    LoadBalancer
    healthChecker   HealthChecker
    failoverPolicy  FailoverPolicy
}

type ClusterNode struct {
    ID              string
    Address         string
    Role            NodeRole
    Status          NodeStatus
    LastHeartbeat   time.Time
    FailureCount    int
    LoadMetrics     LoadMetrics
}

type NodeRole int

const (
    NodeRoleLeader NodeRole = iota
    NodeRoleFollower
    NodeRoleCandidate
)

func (ha *HAClusterManager) StartCluster(ctx context.Context) error {
    // Initialize cluster
    if err := ha.initializeCluster(); err != nil {
        return err
    }

    // Start failure detection
    go ha.monitorNodeHealth(ctx)

    // Start leader election if needed
    if ha.leader == "" {
        go ha.electLeader(ctx)
    }

    // Start load balancing
    go ha.manageLoadBalancing(ctx)

    return nil
}

func (ha *HAClusterManager) monitorNodeHealth(ctx context.Context) {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            ha.checkNodeHealth()
        }
    }
}

func (ha *HAClusterManager) checkNodeHealth() {
    for i := range ha.nodes {
        node := &ha.nodes[i]

        // Perform health check
        healthy := ha.healthChecker.CheckNode(node)

        if !healthy {
            node.FailureCount++
            node.Status = NodeStatusUnhealthy

            // Trigger failover if this is the leader
            if node.ID == ha.leader {
                go ha.handleLeaderFailure(node)
            }
        } else {
            node.FailureCount = 0
            node.Status = NodeStatusHealthy
            node.LastHeartbeat = time.Now()
        }
    }
}

func (ha *HAClusterManager) handleLeaderFailure(failedLeader *ClusterNode) {
    log.Printf("Leader %s failed, initiating failover", failedLeader.ID)

    // Remove failed leader from load balancer
    ha.loadBalancer.RemoveNode(failedLeader.ID)

    // Start leader election
    newLeader, err := ha.electNewLeader()
    if err != nil {
        log.Printf("Leader election failed: %v", err)
        return
    }

    // Promote new leader
    ha.promoteToLeader(newLeader)

    // Update load balancer
    ha.loadBalancer.UpdateLeader(newLeader.ID)

    log.Printf("Failover completed, new leader: %s", newLeader.ID)
}

func (ha *HAClusterManager) electNewLeader() (*ClusterNode, error) {
    // Get healthy nodes
    healthyNodes := ha.getHealthyNodes()
    if len(healthyNodes) == 0 {
        return nil, fmt.Errorf("no healthy nodes available")
    }

    // Select leader based on policy
    switch ha.failoverPolicy.SelectionStrategy {
    case SelectionStrategyLowestLoad:
        return ha.selectNodeByLowestLoad(healthyNodes), nil
    case SelectionStrategyHighestCapacity:
        return ha.selectNodeByHighestCapacity(healthyNodes), nil
    case SelectionStrategyRoundRobin:
        return ha.selectNodeByRoundRobin(healthyNodes), nil
    default:
        return healthyNodes[0], nil
    }
}

// Database failover with read replicas
type DatabaseFailoverManager struct {
    master          DatabaseNode
    readReplicas    []DatabaseNode
    writeReplicas   []DatabaseNode
    healthChecker   DatabaseHealthChecker
    replicationLag  ReplicationLagMonitor
    failoverPolicy  DatabaseFailoverPolicy
}

func (dfm *DatabaseFailoverManager) HandleMasterFailure(ctx context.Context) error {
    log.Printf("Master database failure detected, starting failover")

    // Stop accepting writes
    dfm.master.SetReadOnly(true)

    // Find best replica to promote
    bestReplica, err := dfm.selectBestReplica()
    if err != nil {
        return err
    }

    // Wait for replica to catch up
    if err := dfm.waitForReplicationCatchup(ctx, bestReplica); err != nil {
        return err
    }

    // Promote replica to master
    if err := dfm.promoteReplicaToMaster(bestReplica); err != nil {
        return err
    }

    // Update connection strings
    if err := dfm.updateConnectionStrings(bestReplica); err != nil {
        return err
    }

    // Reconfigure remaining replicas
    if err := dfm.reconfigureReplicas(bestReplica); err != nil {
        return err
    }

    log.Printf("Database failover completed, new master: %s", bestReplica.ID)
    return nil
}

func (dfm *DatabaseFailoverManager) selectBestReplica() (*DatabaseNode, error) {
    var bestReplica *DatabaseNode
    var minLag time.Duration = time.Hour // Start with high value

    for _, replica := range dfm.readReplicas {
        if !dfm.healthChecker.IsHealthy(&replica) {
            continue
        }

        lag := dfm.replicationLag.GetLag(replica.ID)
        if lag < minLag {
            minLag = lag
            bestReplica = &replica
        }
    }

    if bestReplica == nil {
        return nil, fmt.Errorf("no healthy replicas available for promotion")
    }

    return bestReplica, nil
}

// Application-level circuit breaker
type CircuitBreaker struct {
    name            string
    maxFailures     int
    resetTimeout    time.Duration
    halfOpenMaxCalls int
    state           CircuitState
    failureCount    int
    lastFailureTime time.Time
    halfOpenCalls   int
    mutex           sync.RWMutex
    onStateChange   func(name string, from, to CircuitState)
}

type CircuitState int

const (
    CircuitStateClosed CircuitState = iota
    CircuitStateOpen
    CircuitStateHalfOpen
)

func (cb *CircuitBreaker) Execute(operation func() (interface{}, error)) (interface{}, error) {
    cb.mutex.RLock()
    state := cb.state
    cb.mutex.RUnlock()

    switch state {
    case CircuitStateClosed:
        return cb.executeInClosedState(operation)
    case CircuitStateOpen:
        return cb.executeInOpenState(operation)
    case CircuitStateHalfOpen:
        return cb.executeInHalfOpenState(operation)
    default:
        return nil, fmt.Errorf("unknown circuit state")
    }
}

func (cb *CircuitBreaker) executeInClosedState(operation func() (interface{}, error)) (interface{}, error) {
    result, err := operation()

    cb.mutex.Lock()
    defer cb.mutex.Unlock()

    if err != nil {
        cb.failureCount++
        cb.lastFailureTime = time.Now()

        if cb.failureCount >= cb.maxFailures {
            cb.setState(CircuitStateOpen)
        }
    } else {
        cb.failureCount = 0
    }

    return result, err
}

func (cb *CircuitBreaker) executeInOpenState(operation func() (interface{}, error)) (interface{}, error) {
    cb.mutex.Lock()
    defer cb.mutex.Unlock()

    if time.Since(cb.lastFailureTime) > cb.resetTimeout {
        cb.setState(CircuitStateHalfOpen)
        cb.halfOpenCalls = 0

        // Execute operation in half-open state
        cb.mutex.Unlock()
        return cb.executeInHalfOpenState(operation)
    }

    return nil, fmt.Errorf("circuit breaker is open")
}

func (cb *CircuitBreaker) executeInHalfOpenState(operation func() (interface{}, error)) (interface{}, error) {
    cb.mutex.Lock()
    if cb.halfOpenCalls >= cb.halfOpenMaxCalls {
        cb.mutex.Unlock()
        return nil, fmt.Errorf("circuit breaker half-open call limit exceeded")
    }
    cb.halfOpenCalls++
    cb.mutex.Unlock()

    result, err := operation()

    cb.mutex.Lock()
    defer cb.mutex.Unlock()

    if err != nil {
        cb.setState(CircuitStateOpen)
        cb.lastFailureTime = time.Now()
    } else {
        if cb.halfOpenCalls >= cb.halfOpenMaxCalls {
            cb.setState(CircuitStateClosed)
            cb.failureCount = 0
        }
    }

    return result, err
}

func (cb *CircuitBreaker) setState(newState CircuitState) {
    if cb.state == newState {
        return
    }

    oldState := cb.state
    cb.state = newState

    if cb.onStateChange != nil {
        go cb.onStateChange(cb.name, oldState, newState)
    }
}
```

### Chaos Engineering Implementation
```go
// Chaos engineering framework
type ChaosEngineer struct {
    experiments     []ChaosExperiment
    scheduler       ExperimentScheduler
    safetyChecks    []SafetyCheck
    metrics         ChaosMetrics
    rollbackManager RollbackManager
}

type ChaosExperiment interface {
    Name() string
    Description() string
    Execute(ctx context.Context, target Target) error
    Rollback(ctx context.Context, target Target) error
    GetSafetyChecks() []SafetyCheck
}

// Network partition experiment
type NetworkPartitionExperiment struct {
    duration    time.Duration
    partitions  [][]string // Groups of nodes that can communicate
    rollbackCmd string
}

func (npe *NetworkPartitionExperiment) Execute(ctx context.Context, target Target) error {
    log.Printf("Starting network partition experiment for %v", npe.duration)

    // Create network partitions using iptables
    for i, partition := range npe.partitions {
        for j, otherPartition := range npe.partitions {
            if i != j {
                if err := npe.blockTrafficBetween(partition, otherPartition); err != nil {
                    return err
                }
            }
        }
    }

    // Wait for experiment duration
    select {
    case <-time.After(npe.duration):
        return npe.Rollback(ctx, target)
    case <-ctx.Done():
        return npe.Rollback(ctx, target)
    }
}

func (npe *NetworkPartitionExperiment) blockTrafficBetween(partition1, partition2 []string) error {
    for _, node1 := range partition1 {
        for _, node2 := range partition2 {
            cmd := fmt.Sprintf("iptables -A INPUT -s %s -j DROP", node2)
            if err := npe.executeOnNode(node1, cmd); err != nil {
                return err
            }
        }
    }
    return nil
}

func (npe *NetworkPartitionExperiment) Rollback(ctx context.Context, target Target) error {
    log.Printf("Rolling back network partition experiment")

    // Remove all iptables rules
    for _, partition := range npe.partitions {
        for _, node := range partition {
            cmd := "iptables -F INPUT"
            if err := npe.executeOnNode(node, cmd); err != nil {
                log.Printf("Failed to rollback on node %s: %v", node, err)
            }
        }
    }

    return nil
}

// CPU stress experiment
type CPUStressExperiment struct {
    cpuPercent int
    duration   time.Duration
    processes  []int // PIDs of stress processes
}

func (cse *CPUStressExperiment) Execute(ctx context.Context, target Target) error {
    log.Printf("Starting CPU stress experiment: %d%% for %v", cse.cpuPercent, cse.duration)

    // Start stress processes
    numCPU := runtime.NumCPU()
    for i := 0; i < numCPU; i++ {
        cmd := exec.CommandContext(ctx, "stress-ng", "--cpu", "1",
            "--cpu-load", strconv.Itoa(cse.cpuPercent),
            "--timeout", cse.duration.String())

        if err := cmd.Start(); err != nil {
            return err
        }

        cse.processes = append(cse.processes, cmd.Process.Pid)
    }

    // Wait for completion
    time.Sleep(cse.duration)
    return nil
}

func (cse *CPUStressExperiment) Rollback(ctx context.Context, target Target) error {
    log.Printf("Rolling back CPU stress experiment")

    // Kill all stress processes
    for _, pid := range cse.processes {
        if process, err := os.FindProcess(pid); err == nil {
            process.Kill()
        }
    }

    cse.processes = nil
    return nil
}

// Memory pressure experiment
type MemoryPressureExperiment struct {
    memoryMB  int
    duration  time.Duration
    allocatedMemory [][]byte
}

func (mpe *MemoryPressureExperiment) Execute(ctx context.Context, target Target) error {
    log.Printf("Starting memory pressure experiment: %dMB for %v", mpe.memoryMB, mpe.duration)

    // Allocate memory in chunks
    chunkSize := 1024 * 1024 // 1MB chunks
    numChunks := mpe.memoryMB

    for i := 0; i < numChunks; i++ {
        chunk := make([]byte, chunkSize)
        // Write to memory to ensure it's actually allocated
        for j := range chunk {
            chunk[j] = byte(j % 256)
        }
        mpe.allocatedMemory = append(mpe.allocatedMemory, chunk)

        // Small delay to avoid overwhelming the system
        time.Sleep(10 * time.Millisecond)
    }

    // Hold memory for duration
    time.Sleep(mpe.duration)

    return mpe.Rollback(ctx, target)
}

func (mpe *MemoryPressureExperiment) Rollback(ctx context.Context, target Target) error {
    log.Printf("Rolling back memory pressure experiment")

    // Release allocated memory
    mpe.allocatedMemory = nil
    runtime.GC() // Force garbage collection

    return nil
}

// Chaos experiment scheduler
type ChaosScheduler struct {
    experiments []ScheduledExperiment
    safetyNet   SafetyNet
    metrics     ChaosMetrics
}

type ScheduledExperiment struct {
    Experiment ChaosExperiment
    Schedule   string // Cron expression
    Enabled    bool
    LastRun    time.Time
    NextRun    time.Time
}

func (cs *ChaosScheduler) Start(ctx context.Context) {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            cs.checkAndRunExperiments(ctx)
        }
    }
}

func (cs *ChaosScheduler) checkAndRunExperiments(ctx context.Context) {
    now := time.Now()

    for i := range cs.experiments {
        experiment := &cs.experiments[i]

        if !experiment.Enabled {
            continue
        }

        if now.After(experiment.NextRun) {
            go cs.runExperiment(ctx, experiment)
        }
    }
}

func (cs *ChaosScheduler) runExperiment(ctx context.Context, scheduled *ScheduledExperiment) {
    // Check safety conditions
    if !cs.safetyNet.IsSafeToRun(scheduled.Experiment) {
        log.Printf("Safety check failed for experiment %s", scheduled.Experiment.Name())
        return
    }

    log.Printf("Running chaos experiment: %s", scheduled.Experiment.Name())

    start := time.Now()
    err := scheduled.Experiment.Execute(ctx, Target{})
    duration := time.Since(start)

    // Record metrics
    cs.metrics.RecordExperiment(scheduled.Experiment.Name(), duration, err)

    if err != nil {
        log.Printf("Chaos experiment %s failed: %v", scheduled.Experiment.Name(), err)
    }

    // Update schedule
    scheduled.LastRun = start
    scheduled.NextRun = cs.calculateNextRun(scheduled.Schedule, start)
}
```

## 🌍 Multi-Region Deployment Strategies

### Global Disaster Recovery Architecture
```go
// Multi-region disaster recovery coordinator
type MultiRegionDRCoordinator struct {
    regions          []Region
    primaryRegion    string
    backupRegions    []string
    replicationMgr   CrossRegionReplication
    failoverMgr      RegionalFailoverManager
    consistencyMgr   GlobalConsistencyManager
    monitoringMgr    GlobalMonitoringManager
}

type Region struct {
    ID               string
    Name             string
    Location         GeoLocation
    Status           RegionStatus
    Capacity         RegionCapacity
    HealthScore      float64
    LastHealthCheck  time.Time
    DataCenters      []DataCenter
    NetworkLatency   map[string]time.Duration // Latency to other regions
}

type RegionStatus int

const (
    RegionStatusActive RegionStatus = iota
    RegionStatusDegraded
    RegionStatusFailed
    RegionStatusMaintenance
)

func (mdr *MultiRegionDRCoordinator) InitializeGlobalDR() error {
    // Initialize each region
    for _, region := range mdr.regions {
        if err := mdr.initializeRegion(region); err != nil {
            return fmt.Errorf("failed to initialize region %s: %w", region.ID, err)
        }
    }

    // Setup cross-region replication
    if err := mdr.replicationMgr.SetupReplication(mdr.regions); err != nil {
        return err
    }

    // Start global monitoring
    go mdr.monitoringMgr.StartGlobalMonitoring()

    // Start health checks
    go mdr.startRegionalHealthChecks()

    return nil
}

func (mdr *MultiRegionDRCoordinator) HandleRegionalFailure(failedRegion string) error {
    log.Printf("Regional failure detected in %s, initiating disaster recovery", failedRegion)

    // Mark region as failed
    if err := mdr.markRegionAsFailed(failedRegion); err != nil {
        return err
    }

    // If primary region failed, promote backup region
    if failedRegion == mdr.primaryRegion {
        return mdr.promoteBackupRegion()
    }

    // If backup region failed, redistribute load
    return mdr.redistributeLoad(failedRegion)
}

func (mdr *MultiRegionDRCoordinator) promoteBackupRegion() error {
    // Select best backup region based on health and capacity
    bestBackup := mdr.selectBestBackupRegion()
    if bestBackup == "" {
        return fmt.Errorf("no healthy backup regions available")
    }

    log.Printf("Promoting backup region %s to primary", bestBackup)

    // Update DNS to point to new primary
    if err := mdr.updateGlobalDNS(bestBackup); err != nil {
        return err
    }

    // Promote region to primary
    mdr.primaryRegion = bestBackup

    // Update load balancer configuration
    if err := mdr.updateGlobalLoadBalancer(bestBackup); err != nil {
        return err
    }

    // Notify all services of the change
    return mdr.notifyServicesOfFailover(bestBackup)
}

// Cross-region data replication
type CrossRegionReplication struct {
    replicationStrategy ReplicationStrategy
    consistencyLevel    ConsistencyLevel
    conflictResolver    ConflictResolver
    bandwidthManager    BandwidthManager
}

func (crr *CrossRegionReplication) ReplicateData(data ReplicationData) error {
    switch crr.replicationStrategy {
    case ReplicationStrategyAsynchronous:
        return crr.replicateAsync(data)
    case ReplicationStrategySynchronous:
        return crr.replicateSync(data)
    case ReplicationStrategyEventual:
        return crr.replicateEventual(data)
    default:
        return fmt.Errorf("unknown replication strategy")
    }
}

func (crr *CrossRegionReplication) replicateAsync(data ReplicationData) error {
    // Replicate to all regions asynchronously
    var wg sync.WaitGroup
    errors := make(chan error, len(data.TargetRegions))

    for _, region := range data.TargetRegions {
        wg.Add(1)
        go func(targetRegion string) {
            defer wg.Done()

            // Apply bandwidth throttling
            if err := crr.bandwidthManager.AcquireBandwidth(targetRegion, data.Size); err != nil {
                errors <- err
                return
            }
            defer crr.bandwidthManager.ReleaseBandwidth(targetRegion, data.Size)

            // Replicate data
            if err := crr.replicateToRegion(data, targetRegion); err != nil {
                errors <- err
            }
        }(region)
    }

    wg.Wait()
    close(errors)

    // Check for errors
    var replicationErrors []error
    for err := range errors {
        replicationErrors = append(replicationErrors, err)
    }

    if len(replicationErrors) > 0 {
        return fmt.Errorf("replication errors: %v", replicationErrors)
    }

    return nil
}

// Global load balancing with disaster recovery
type GlobalLoadBalancer struct {
    regions          []Region
    routingPolicy    RoutingPolicy
    healthChecker    GlobalHealthChecker
    trafficSplitter  TrafficSplitter
    geoResolver      GeoResolver
}

type RoutingPolicy int

const (
    RoutingPolicyGeographic RoutingPolicy = iota
    RoutingPolicyLatency
    RoutingPolicyCapacity
    RoutingPolicyFailover
)

func (glb *GlobalLoadBalancer) RouteRequest(request *Request) (*Region, error) {
    // Get client location
    clientLocation := glb.geoResolver.GetLocation(request.ClientIP)

    // Filter healthy regions
    healthyRegions := glb.getHealthyRegions()
    if len(healthyRegions) == 0 {
        return nil, fmt.Errorf("no healthy regions available")
    }

    // Apply routing policy
    switch glb.routingPolicy {
    case RoutingPolicyGeographic:
        return glb.routeByGeography(clientLocation, healthyRegions), nil
    case RoutingPolicyLatency:
        return glb.routeByLatency(clientLocation, healthyRegions), nil
    case RoutingPolicyCapacity:
        return glb.routeByCapacity(healthyRegions), nil
    case RoutingPolicyFailover:
        return glb.routeByFailover(healthyRegions), nil
    default:
        return healthyRegions[0], nil
    }
}

func (glb *GlobalLoadBalancer) routeByGeography(clientLocation GeoLocation, regions []*Region) *Region {
    var closestRegion *Region
    var minDistance float64 = math.MaxFloat64

    for _, region := range regions {
        distance := calculateDistance(clientLocation, region.Location)
        if distance < minDistance {
            minDistance = distance
            closestRegion = region
        }
    }

    return closestRegion
}

func (glb *GlobalLoadBalancer) routeByLatency(clientLocation GeoLocation, regions []*Region) *Region {
    var bestRegion *Region
    var minLatency time.Duration = time.Hour

    for _, region := range regions {
        // Use historical latency data or perform real-time measurement
        latency := glb.getLatencyToRegion(clientLocation, region)
        if latency < minLatency {
            minLatency = latency
            bestRegion = region
        }
    }

    return bestRegion
}
```

### Monitoring and Alerting Systems
```go
// Comprehensive disaster recovery monitoring
type DRMonitoringSystem struct {
    healthCheckers    []HealthChecker
    alertManager      AlertManager
    dashboardManager  DashboardManager
    incidentManager   IncidentManager
    metricsCollector  MetricsCollector
    thresholds        MonitoringThresholds
}

type MonitoringThresholds struct {
    RPOThreshold         time.Duration // Recovery Point Objective
    RTOThreshold         time.Duration // Recovery Time Objective
    ReplicationLagMax    time.Duration
    FailoverTimeMax      time.Duration
    DataLossMax          int64 // Maximum acceptable data loss in bytes
    AvailabilityMin      float64 // Minimum availability percentage
}

func (drm *DRMonitoringSystem) StartMonitoring(ctx context.Context) {
    // Start health checkers
    for _, checker := range drm.healthCheckers {
        go checker.Start(ctx)
    }

    // Start metrics collection
    go drm.metricsCollector.Start(ctx)

    // Start threshold monitoring
    go drm.monitorThresholds(ctx)

    // Start dashboard updates
    go drm.dashboardManager.Start(ctx)
}

func (drm *DRMonitoringSystem) monitorThresholds(ctx context.Context) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            drm.checkAllThresholds()
        }
    }
}

func (drm *DRMonitoringSystem) checkAllThresholds() {
    metrics := drm.metricsCollector.GetCurrentMetrics()

    // Check RPO threshold
    if metrics.ReplicationLag > drm.thresholds.RPOThreshold {
        drm.alertManager.SendAlert(Alert{
            Level:       AlertLevelCritical,
            Type:        AlertTypeRPOViolation,
            Message:     fmt.Sprintf("RPO threshold violated: %v > %v", metrics.ReplicationLag, drm.thresholds.RPOThreshold),
            Timestamp:   time.Now(),
            Metadata:    map[string]interface{}{"replication_lag": metrics.ReplicationLag},
        })
    }

    // Check availability threshold
    if metrics.Availability < drm.thresholds.AvailabilityMin {
        drm.alertManager.SendAlert(Alert{
            Level:       AlertLevelCritical,
            Type:        AlertTypeAvailabilityDrop,
            Message:     fmt.Sprintf("Availability below threshold: %.2f%% < %.2f%%", metrics.Availability*100, drm.thresholds.AvailabilityMin*100),
            Timestamp:   time.Now(),
            Metadata:    map[string]interface{}{"availability": metrics.Availability},
        })
    }

    // Check data loss
    if metrics.DataLoss > drm.thresholds.DataLossMax {
        drm.alertManager.SendAlert(Alert{
            Level:       AlertLevelCritical,
            Type:        AlertTypeDataLoss,
            Message:     fmt.Sprintf("Data loss detected: %d bytes", metrics.DataLoss),
            Timestamp:   time.Now(),
            Metadata:    map[string]interface{}{"data_loss_bytes": metrics.DataLoss},
        })
    }
}

// Automated incident response system
type IncidentResponseSystem struct {
    playbooks       map[IncidentType]Playbook
    escalationRules []EscalationRule
    notificationMgr NotificationManager
    automationMgr   AutomationManager
    statusPage      StatusPageManager
}

type Playbook struct {
    Name        string
    Steps       []ResponseStep
    Automation  []AutomatedAction
    Escalation  EscalationPolicy
}

type ResponseStep struct {
    Name        string
    Description string
    Action      func(ctx context.Context, incident *Incident) error
    Timeout     time.Duration
    Required    bool
}

func (irs *IncidentResponseSystem) HandleIncident(incident *Incident) error {
    log.Printf("Handling incident: %s", incident.ID)

    // Get appropriate playbook
    playbook, exists := irs.playbooks[incident.Type]
    if !exists {
        return fmt.Errorf("no playbook found for incident type: %s", incident.Type)
    }

    // Update status page
    irs.statusPage.UpdateIncidentStatus(incident.ID, "investigating")

    // Execute automated actions first
    for _, action := range playbook.Automation {
        if err := irs.automationMgr.ExecuteAction(action, incident); err != nil {
            log.Printf("Automated action failed: %v", err)
        }
    }

    // Execute manual steps
    for _, step := range playbook.Steps {
        if err := irs.executeStep(step, incident); err != nil {
            log.Printf("Response step failed: %v", err)

            // Escalate if required step fails
            if step.Required {
                return irs.escalateIncident(incident, playbook.Escalation)
            }
        }
    }

    // Mark incident as resolved
    incident.Status = IncidentStatusResolved
    incident.ResolvedAt = time.Now()

    // Update status page
    irs.statusPage.UpdateIncidentStatus(incident.ID, "resolved")

    return nil
}

func (irs *IncidentResponseSystem) executeStep(step ResponseStep, incident *Incident) error {
    ctx, cancel := context.WithTimeout(context.Background(), step.Timeout)
    defer cancel()

    log.Printf("Executing response step: %s", step.Name)

    return step.Action(ctx, incident)
}

// Business continuity planning
type BusinessContinuityManager struct {
    continuityPlans   map[DisasterType]ContinuityPlan
    riskAssessment    RiskAssessmentEngine
    impactAnalysis    BusinessImpactAnalysis
    recoveryPlanner   RecoveryPlanner
    testScheduler     BCPTestScheduler
}

type ContinuityPlan struct {
    DisasterType     DisasterType
    ImpactLevel      ImpactLevel
    RecoveryStrategy RecoveryStrategy
    ResourceRequirements []Resource
    Dependencies     []Dependency
    TestSchedule     TestSchedule
}

type DisasterType int

const (
    DisasterTypeNatural DisasterType = iota
    DisasterTypeCyber
    DisasterTypeHardware
    DisasterTypeHuman
    DisasterTypePandemic
)

func (bcm *BusinessContinuityManager) ActivateContinuityPlan(disasterType DisasterType) error {
    plan, exists := bcm.continuityPlans[disasterType]
    if !exists {
        return fmt.Errorf("no continuity plan for disaster type: %d", disasterType)
    }

    log.Printf("Activating business continuity plan for disaster type: %d", disasterType)

    // Assess current impact
    impact := bcm.impactAnalysis.AssessCurrentImpact(disasterType)

    // Execute recovery strategy
    switch plan.RecoveryStrategy {
    case RecoveryStrategyHotSite:
        return bcm.activateHotSite(plan)
    case RecoveryStrategyColdSite:
        return bcm.activateColdSite(plan)
    case RecoveryStrategyCloud:
        return bcm.activateCloudRecovery(plan)
    case RecoveryStrategyMobile:
        return bcm.activateMobileRecovery(plan)
    default:
        return fmt.Errorf("unknown recovery strategy: %d", plan.RecoveryStrategy)
    }
}

func (bcm *BusinessContinuityManager) TestContinuityPlan(planType DisasterType) (*TestResult, error) {
    plan := bcm.continuityPlans[planType]

    testResult := &TestResult{
        PlanType:  planType,
        StartTime: time.Now(),
        Tests:     make([]TestCase, 0),
    }

    // Test each component of the plan
    for _, test := range plan.TestSchedule.Tests {
        result := bcm.executeTest(test)
        testResult.Tests = append(testResult.Tests, result)
    }

    testResult.EndTime = time.Now()
    testResult.Duration = testResult.EndTime.Sub(testResult.StartTime)
    testResult.Success = bcm.allTestsPassed(testResult.Tests)

    return testResult, nil
}
```

## 🎯 Summary: Disaster Recovery & Reliability

### Comprehensive DR Coverage
- **Backup Systems** - Automated full, incremental, and point-in-time recovery
- **Failover Mechanisms** - High availability clusters and database failover
- **Chaos Engineering** - Systematic failure injection and resilience testing
- **Multi-Region Deployment** - Global disaster recovery with cross-region replication
- **Monitoring & Alerting** - Real-time DR monitoring with automated incident response
- **Business Continuity** - Comprehensive planning for various disaster scenarios

### Advanced Reliability Patterns
- **Circuit Breakers** - Prevent cascading failures
- **Health Checking** - Multi-level system health monitoring
- **Automated Recovery** - Self-healing systems with minimal human intervention
- **Chaos Testing** - Proactive resilience validation
- **Global Load Balancing** - Intelligent traffic routing during failures
- **Incident Response** - Automated playbooks and escalation procedures

### Real-world Applications
These DR patterns apply to all the systems we've designed:
- **Twitter**: Multi-region timeline replication and failover
- **Netflix**: Global CDN failover and content availability
- **Uber**: Location service redundancy and driver matching continuity
- **WhatsApp**: Message delivery guarantees and cross-region backup
- **TinyURL**: High availability URL resolution and data backup
- **Search Engine**: Index replication and query service failover

You now have **expert-level disaster recovery and reliability knowledge** covering all aspects of building resilient systems! 🚀
