# RiskAnalysisService Component Documentation

## 1. OVERVIEW

### Purpose and Primary Functionality

The `RiskAnalysisService` is the core business logic component responsible for assessing transaction risk levels based on fraud scoring results. It evaluates various scoring factors to determine whether a transaction should be approved or rejected, providing a comprehensive risk assessment system for financial transactions.

**Key responsibilities:**
- Process fraud scoring results from multiple scoring cards
- Calculate aggregate risk scores using weighted scoring algorithms
- Classify transactions into risk levels (Low, Medium, High)
- Make approval/rejection decisions based on risk thresholds
- Persist risk analysis results for audit and tracking purposes

### When to Use This Component vs. Alternatives

**Use RiskAnalysisService when:**
- You need automated risk assessment for financial transactions
- You have fraud scoring data from multiple sources (value, seller, currency, etc.)
- You require configurable risk thresholds and decision logic
- You need audit trails for risk decisions

**Alternative approaches:**
- **Rule-based engines**: For more complex, dynamic rule sets
- **ML-based scoring**: For advanced pattern recognition and adaptive scoring
- **Third-party risk services**: For outsourced risk assessment

### Architectural Context

The `RiskAnalysisService` sits at the core of a hexagonal architecture:

```
┌─────────────────────────────────────────────────────────────┐
│                    External Systems                          │
├─────────────────────────────────────────────────────────────┤
│  Kafka (Fraud Scoring) → FraudScoringReceiver (Adapter)    │
│                                ↓                            │
│                    RiskAnalysisService (Domain)            │
│                                ↓                            │
│  RiskAnalysisRepository → KafkaRiskAnalysisRepository      │
│                                ↓                            │
│                    Kafka (Risk Decisions)                   │
└─────────────────────────────────────────────────────────────┘
```

**Position in the system:**
- **Input**: Receives `ScoringResult` objects from fraud detection systems
- **Processing**: Applies business rules for risk assessment
- **Output**: Generates `RiskAnalysis` objects with approval/rejection decisions
- **Integration**: Works with repository pattern for persistence

## 2. TECHNICAL SPECIFICATION

### API Reference

#### Constructor

```go
func NewRiskAnalysisService(repo RiskAnalysisRepository, log *zap.Logger) *RiskAnalysisService
```

**Parameters:**
- `repo RiskAnalysisRepository`: Repository interface for storing risk analysis results
- `log *zap.Logger`: Logger instance for error and info logging

**Returns:**
- `*RiskAnalysisService`: New instance of the service

#### Methods

##### Assessment

```go
func (ra *RiskAnalysisService) Assessment(scoring *ScoringResult) error
```

**Purpose**: Performs risk assessment on a transaction based on scoring results.

**Parameters:**
- `scoring *ScoringResult`: Fraud scoring results containing various score cards

**Returns:**
- `error`: Returns nil on success, error on failure

**Processing Logic:**
1. Calculates aggregate score from all scoring cards
2. Determines risk level based on score thresholds
3. Makes approval/rejection decision
4. Creates RiskAnalysis object with timestamp
5. Persists result via repository

### Data Types

#### ScoringResult

```go
type ScoringResult struct {
    Score       ScoreCard           `json:"score"`
    Transaction TransactionAnalysis `json:"transaction"`
}
```

#### ScoreCard

```go
type ScoreCard struct {
    ValueScore        ValueScoreCard        `json:"valueScore"`
    SellerScore       SellerScoreCard       `json:"sellerScore"`
    AverageValueScore AverageValueScoreCard `json:"averageValueScore"`
    CurrencyScore     CurrencyScoreCard     `json:"currencyScore"`
}
```

#### Scoring Weights

| Score Card | Weight Multiplier | Impact |
|------------|-------------------|---------|
| ValueScore | 10x | High impact on final score |
| SellerScore | 5x | Medium-high impact |
| AverageValueScore | 3x | Medium impact |
| CurrencyScore | 1x | Low impact |

### Risk Level Thresholds

| Risk Level | Score Range | Decision |
|------------|-------------|----------|
| Low | 0-30 | Approved |
| Medium | 31-60 | Approved |
| High | 61-100 | Rejected |

### State Management

The service is **stateless** - it does not maintain any internal state between assessments. Each assessment is independent and relies on:
- Input scoring data
- Configured risk thresholds (compile-time constants)
- Repository for persistence

### Events and Dependencies

**Dependencies:**
- `RiskAnalysisRepository`: For storing assessment results
- `*zap.Logger`: For error logging and monitoring

**Events:**
- **Input**: Triggered by fraud scoring events from Kafka
- **Output**: Publishes risk decision events via repository implementation

## 3. IMPLEMENTATION EXAMPLES

### Basic Usage Example

```go
package main

import (
    "go.uber.org/zap"
    "risk-management/internal/domain"
)

func main() {
    // Initialize logger
    logger, _ := zap.NewDevelopment()
    
    // Initialize repository (implementation depends on your needs)
    repo := NewInMemoryRiskAnalysisRepository()
    
    // Create service
    service := domain.NewRiskAnalysisService(repo, logger)
    
    // Create scoring result
    scoring := &domain.ScoringResult{
        Score: domain.ScoreCard{
            ValueScore:        domain.ValueScoreCard{Score: 3},
            SellerScore:       domain.SellerScoreCard{Score: 2},
            AverageValueScore: domain.AverageValueScoreCard{Score: 1},
            CurrencyScore:     domain.CurrencyScoreCard{Score: 1},
        },
        Transaction: domain.TransactionAnalysis{
            Order: domain.Checkout{Id: "txn-123"},
            // ... other transaction details
        },
    }
    
    // Perform assessment
    err := service.Assessment(scoring)
    if err != nil {
        logger.Error("Assessment failed", zap.Error(err))
    }
}
```

### Advanced Configuration Example

```go
// Custom repository implementation with database
type DatabaseRiskAnalysisRepository struct {
    db *sql.DB
}

func (repo *DatabaseRiskAnalysisRepository) Store(analysis *domain.RiskAnalysis) error {
    query := `
        INSERT INTO risk_analyses (status, level, transaction_id, created_at)
        VALUES (?, ?, ?, ?)
    `
    _, err := repo.db.Exec(query, 
        analysis.Status.String(), 
        analysis.Level.Name, 
        analysis.Transaction.Order.Id, 
        analysis.At)
    return err
}

// Usage with custom repository
func setupWithDatabase() *domain.RiskAnalysisService {
    db, _ := sql.Open("postgres", "connection-string")
    repo := &DatabaseRiskAnalysisRepository{db: db}
    logger, _ := zap.NewProduction()
    
    return domain.NewRiskAnalysisService(repo, logger)
}
```

### Customization Scenarios

#### Custom Risk Thresholds

```go
// To customize risk thresholds, you would need to modify the service
// or create a configurable version:

type ConfigurableRiskAnalysisService struct {
    repo RiskAnalysisRepository
    log  *zap.Logger
    lowThreshold    int
    mediumThreshold int
    highThreshold   int
}

func (ra *ConfigurableRiskAnalysisService) Assessment(scoring *ScoringResult) error {
    total := ra.calculateTotal(scoring)
    
    var level *RiskLevel
    if total <= ra.lowThreshold {
        level = &RiskLevel{Name: "Low", From: 0, To: ra.lowThreshold}
    } else if total <= ra.mediumThreshold {
        level = &RiskLevel{Name: "Medium", From: ra.lowThreshold + 1, To: ra.mediumThreshold}
    } else {
        level = &RiskLevel{Name: "High", From: ra.mediumThreshold + 1, To: ra.highThreshold}
    }
    
    // ... rest of the assessment logic
}
```

### Common Patterns and Best Practices

#### 1. Error Handling Pattern

```go
func handleAssessment(service *domain.RiskAnalysisService, scoring *domain.ScoringResult) {
    if err := service.Assessment(scoring); err != nil {
        // Log error but don't fail the entire process
        log.Error("Risk assessment failed", 
            zap.Error(err),
            zap.String("transaction_id", scoring.Transaction.Order.Id))
        
        // Implement fallback strategy
        // e.g., default to rejection for safety
    }
}
```

#### 2. Batch Processing Pattern

```go
func processBatch(service *domain.RiskAnalysisService, scorings []*domain.ScoringResult) {
    for _, scoring := range scorings {
        go func(s *domain.ScoringResult) {
            defer func() {
                if r := recover(); r != nil {
                    log.Error("Panic in assessment", zap.Any("panic", r))
                }
            }()
            
            service.Assessment(s)
        }(scoring)
    }
}
```

#### 3. Testing Pattern

```go
func TestRiskAnalysisService_Assessment(t *testing.T) {
    // Mock repository
    mockRepo := &MockRiskAnalysisRepository{}
    logger, _ := zap.NewDevelopment()
    
    service := domain.NewRiskAnalysisService(mockRepo, logger)
    
    // Test case: Low risk scenario
    scoring := &domain.ScoringResult{
        Score: domain.ScoreCard{
            ValueScore:        domain.ValueScoreCard{Score: 1},
            SellerScore:       domain.SellerScoreCard{Score: 1},
            AverageValueScore: domain.AverageValueScoreCard{Score: 1},
            CurrencyScore:     domain.CurrencyScoreCard{Score: 1},
        },
    }
    
    err := service.Assessment(scoring)
    
    assert.NoError(t, err)
    assert.Equal(t, domain.Approved, mockRepo.StoredAnalysis.Status)
    assert.Equal(t, "Low", mockRepo.StoredAnalysis.Level.Name)
}
```

## 4. TROUBLESHOOTING

### Common Errors and Solutions

#### 1. Repository Storage Failures

**Error**: `error to store risk analysis: connection refused`

**Cause**: Database/Kafka connection issues

**Solution**:
```go
// Implement retry logic in repository
func (repo *DatabaseRepo) Store(analysis *domain.RiskAnalysis) error {
    for i := 0; i < 3; i++ {
        if err := repo.doStore(analysis); err != nil {
            if i == 2 {
                return err
            }
            time.Sleep(time.Second * time.Duration(i+1))
            continue
        }
        return nil
    }
    return fmt.Errorf("failed after 3 retries")
}
```

#### 2. Nil Pointer Exceptions

**Error**: `panic: runtime error: invalid memory address or nil pointer dereference`

**Cause**: Missing validation of input data

**Solution**:
```go
func (ra *RiskAnalysisService) Assessment(scoring *ScoringResult) error {
    if scoring == nil {
        return errors.New("scoring result cannot be nil")
    }
    
    if scoring.Transaction.Order.Id == "" {
        return errors.New("transaction ID is required")
    }
    
    // ... rest of validation
}
```

#### 3. Inconsistent Scoring Calculations

**Error**: Unexpected risk levels for known transaction patterns

**Cause**: Incorrect weight calculations or missing score card data

**Solution**:
```go
func (ra *RiskAnalysisService) calculateTotal(scoring *ScoringResult) int {
    // Add validation and default handling
    valueScore := scoring.Score.ValueScore.Value()
    sellerScore := scoring.Score.SellerScore.Value()
    avgScore := scoring.Score.AverageValueScore.Score // Note: no multiplier
    currencyScore := scoring.Score.CurrencyScore.Value()
    
    total := valueScore + sellerScore + avgScore + currencyScore
    
    // Log for debugging
    ra.log.Debug("Score calculation",
        zap.Int("value_score", valueScore),
        zap.Int("seller_score", sellerScore),
        zap.Int("avg_score", avgScore),
        zap.Int("currency_score", currencyScore),
        zap.Int("total", total))
    
    return total
}
```

### Debugging Strategies

#### 1. Enable Debug Logging

```go
// Add detailed logging in Assessment method
func (ra *RiskAnalysisService) Assessment(scoring *ScoringResult) error {
    ra.log.Debug("Starting assessment",
        zap.String("transaction_id", scoring.Transaction.Order.Id))
    
    total := ra.calculateTotal(scoring)
    
    ra.log.Debug("Risk assessment completed",
        zap.Int("total_score", total),
        zap.String("risk_level", level.Name),
        zap.String("decision", status.String()))
    
    // ... rest of method
}
```

#### 2. Add Metrics Collection

```go
import "github.com/prometheus/client_golang/prometheus"

var (
    assessmentCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "risk_assessments_total",
            Help: "Total number of risk assessments",
        },
        []string{"risk_level", "decision"},
    )
)

func (ra *RiskAnalysisService) Assessment(scoring *ScoringResult) error {
    // ... assessment logic
    
    assessmentCounter.WithLabelValues(level.Name, status.String()).Inc()
    
    return nil
}
```

### Performance Considerations

#### 1. Memory Usage

- **Current**: Service is stateless, minimal memory footprint
- **Optimization**: Consider object pooling for high-throughput scenarios

#### 2. Processing Speed

- **Current**: Synchronous processing, O(1) complexity
- **Bottleneck**: Repository storage operations
- **Optimization**: Implement async repository operations

```go
type AsyncRiskAnalysisRepository struct {
    queue chan *domain.RiskAnalysis
    repo  domain.RiskAnalysisRepository
}

func (ar *AsyncRiskAnalysisRepository) Store(analysis *domain.RiskAnalysis) error {
    select {
    case ar.queue <- analysis:
        return nil
    default:
        return errors.New("queue full")
    }
}
```

#### 3. Concurrency

- **Current**: Thread-safe (stateless)
- **Consideration**: Repository implementations must be thread-safe
- **Testing**: Use race detector: `go test -race`

## 5. RELATED COMPONENTS

### Dependencies

#### Direct Dependencies

1. **RiskAnalysisRepository Interface**
   - **Purpose**: Persistence abstraction
   - **Location**: `internal/domain/risk_analysis_repo.go`
   - **Implementations**: `KafkaRiskAnalysisRepository`

2. **Zap Logger**
   - **Purpose**: Structured logging
   - **Library**: `go.uber.org/zap`
   - **Usage**: Error logging and debugging

#### Domain Models

1. **ScoringResult**
   - **Location**: `internal/domain/scoring_result.go`
   - **Purpose**: Input data structure for assessments

2. **RiskAnalysis**
   - **Location**: `internal/domain/risk_analysis.go`
   - **Purpose**: Output data structure with risk decisions

3. **TransactionAnalysis**
   - **Location**: `internal/domain/transaction_analysis.go`
   - **Purpose**: Transaction context and metadata

### Components Commonly Used Alongside

#### 1. FraudScoringReceiver
- **Location**: `internal/adapter/kafka/in/fraud_scoring_receiver.go`
- **Relationship**: Upstream component that calls `RiskAnalysisService`
- **Usage Pattern**: Event-driven processing

```go
// FraudScoringReceiver.Handle calls RiskAnalysisService.Assessment
func (fsr *FraudScoringReceiver) Handle(ctx context.Context, event cloudevents.Event) error {
    data := &domain.ScoringResult{}
    event.DataAs(data)
    return fsr.rs.Assessment(data) // Calls our service
}
```

#### 2. KafkaRiskAnalysisRepository
- **Location**: `internal/adapter/kafka/out/kafka_risk_analysis.go`
- **Relationship**: Downstream component that publishes risk decisions
- **Usage Pattern**: Repository pattern implementation

#### 3. Manager
- **Location**: `cmd/manager.go`
- **Relationship**: Orchestrates the entire risk assessment flow
- **Usage Pattern**: Dependency injection and lifecycle management

### Alternative Approaches

#### 1. Rule Engine Based Approach

```go
type RuleEngine interface {
    Evaluate(ctx context.Context, facts map[string]interface{}) (Decision, error)
}

type RuleBasedRiskService struct {
    engine RuleEngine
}

// More flexible but complex rule management
```

#### 2. Machine Learning Based Approach

```go
type MLRiskService struct {
    model MLModel
}

func (ml *MLRiskService) Assessment(scoring *ScoringResult) (*RiskAnalysis, error) {
    features := extractFeatures(scoring)
    prediction := ml.model.Predict(features)
    return convertToRiskAnalysis(prediction), nil
}
```

#### 3. Microservice Based Approach

```go
type RemoteRiskService struct {
    client *http.Client
    endpoint string
}

func (rs *RemoteRiskService) Assessment(scoring *ScoringResult) error {
    // HTTP call to dedicated risk service
    resp, err := rs.client.Post(rs.endpoint, "application/json", scoring)
    // ... handle response
}
```

### Integration Patterns

#### 1. Event-Driven Architecture
- **Current**: Kafka-based event processing
- **Benefits**: Loose coupling, scalability
- **Considerations**: Eventual consistency, error handling

#### 2. Hexagonal Architecture
- **Current**: Domain-driven design with adapters
- **Benefits**: Testability, technology independence
- **Components**: Domain (business logic), Adapters (I/O), Infrastructure

#### 3. Repository Pattern
- **Current**: Abstracted persistence layer
- **Benefits**: Testability, multiple storage backends
- **Implementations**: Kafka, Database, In-memory

---

*This documentation covers the RiskAnalysisService component version 1.0. For updates and additional examples, refer to the project repository and test files.*