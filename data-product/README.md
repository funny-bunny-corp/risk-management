# Risk Analysis Data Product

## Overview

### Owner & Domain
- **Owner**: Risk Management Team (risk-management@funny-bunny.xyz)
- **Domain**: Risk Management for Financial Transactions
- **Product Name**: Risk Analysis Decisions
- **Version**: 1.0.0

### Description & Purpose

The Risk Analysis Data Product provides real-time and historical analytical access to risk decisions made on financial transactions. This data product aggregates fraud detection scoring results with final risk assessment decisions, enabling comprehensive analysis of transaction risk patterns, fraud detection effectiveness, and business impact.

**Key Business Value:**
- **Risk Monitoring**: Real-time visibility into transaction risk levels and approval rates
- **Fraud Pattern Analysis**: Historical analysis of fraud scoring effectiveness and emerging patterns
- **Business Intelligence**: Decision support for risk policy optimization and business growth
- **Compliance Reporting**: Audit trail and regulatory reporting for risk management practices
- **Performance Optimization**: Data-driven insights for improving fraud detection accuracy

### Intended Consumers

1. **Risk Management Team**: Real-time monitoring dashboards and policy optimization
2. **Business Intelligence Analysts**: Historical trend analysis and business impact assessment
3. **Compliance Officers**: Regulatory reporting and audit trail generation
4. **Data Scientists**: Model performance evaluation and pattern discovery
5. **Product Managers**: Transaction approval rate optimization and user experience insights

## Technical Specification

### Data Source
- **Event Stream**: Apache Kafka
- **Input Topics**: 
  - `risk-decision-approved` (CloudEvents format)
  - `risk-decision-rejected` (CloudEvents format)
- **Processing**: Real-time stream processing with Apache Pinot
- **Update Frequency**: Real-time (sub-second latency)

### Schema Definition

| Field Name | Data Type | Description | Example | Business Context |
|------------|-----------|-------------|---------|------------------|
| `decision_status` | STRING | Risk decision outcome | "approved", "rejected" | Final risk assessment decision for the transaction |
| `risk_level` | STRING | Risk level classification | "Low", "Medium", "High" | Risk category based on calculated score thresholds |
| `transaction_id` | STRING | Unique transaction identifier | "ORDER-789" | Primary business key for transaction tracking |
| `order_id` | STRING | Unique order identifier | "ORDER-789" | Order reference for business operations |
| `payment_id` | STRING | Unique payment identifier | "PAY-456" | Payment processing reference |
| `buyer_document` | STRING | Buyer's identity document | "12345678901" | Buyer identification for fraud analysis |
| `buyer_name` | STRING | Buyer's full name | "John Doe" | Customer identification |
| `seller_id` | STRING | Seller identifier | "SELLER123" | Merchant/seller reference |
| `payment_currency` | STRING | Payment currency code | "USD", "EUR", "BRL" | Currency for regional risk analysis |
| `payment_status` | STRING | Payment status at analysis | "pending", "processing" | Payment lifecycle context |
| `card_info` | STRING | Masked card information | "****1234" | Card pattern analysis (PCI compliant) |
| `card_token` | STRING | Card tokenization reference | "tok_abc123" | Card usage pattern tracking |
| `event_type` | STRING | Type of risk event | "risk.decision.approved" | Event classification |
| `event_source` | STRING | Source system | "risk-management" | System provenance |
| `value_score` | INT | Raw value-based fraud score | 0-10 | Transaction amount risk assessment |
| `seller_score` | INT | Raw seller-based fraud score | 0-10 | Merchant reputation risk assessment |
| `average_value_score` | INT | Raw average value fraud score | 0-10 | User spending pattern risk assessment |
| `currency_score` | INT | Raw currency-based fraud score | 0-10 | Currency exchange risk assessment |
| `weighted_value_score` | INT | Weighted value score | 0-100 | Value score × 10 (high impact) |
| `weighted_seller_score` | INT | Weighted seller score | 0-50 | Seller score × 5 (medium impact) |
| `weighted_currency_score` | INT | Weighted currency score | 0-10 | Currency score × 1 (low impact) |
| `total_risk_score` | INT | Total calculated risk score | 0-170 | Sum of all weighted scores |
| `risk_score_min` | INT | Risk level minimum threshold | 0, 31, 61 | Lower bound for risk classification |
| `risk_score_max` | INT | Risk level maximum threshold | 30, 60, 100 | Upper bound for risk classification |
| `payment_amount` | DECIMAL | Transaction amount | 1500.00 | Financial value for risk analysis |
| `analysis_timestamp` | TIMESTAMP | Risk analysis timestamp | 2024-01-15T10:30:15.123Z | When risk decision was made |
| `order_timestamp` | TIMESTAMP | Order creation timestamp | 2024-01-15T10:30:00Z | When transaction was initiated |
| `analysis_date` | STRING | Analysis date (partitioning) | "2024-01-15" | Date-based data partitioning |
| `analysis_hour` | INT | Analysis hour (0-23) | 10 | Time-based pattern analysis |

### Data Quality Guarantees (SLA/SLO)

#### Freshness
- **Real-time Latency**: < 1 second from event generation to query availability
- **Data Availability**: 99.9% uptime during business hours
- **Update Frequency**: Continuous real-time ingestion

#### Completeness
- **Data Coverage**: 100% of risk decisions are captured and stored
- **Required Fields**: All schema fields marked as `notNull` are guaranteed to be present
- **Validation**: Input events are validated against schema before ingestion

#### Accuracy
- **Score Calculation**: Risk scores are calculated using validated business logic
- **Data Integrity**: Primary key constraints ensure no duplicate risk decisions
- **Transformation Accuracy**: All data transformations are tested and validated

#### Retention
- **Hot Tier**: 7 days of data stored on high-performance storage
- **Cold Tier**: 30 days of data stored on standard storage
- **Total Retention**: 90 days of historical data available for analysis
- **Backup**: Daily backups with 30-day retention for disaster recovery

### Data Governance

#### Access Control
- **Read Access**: Granted to Risk Management, BI, and Compliance teams
- **Write Access**: Restricted to risk-management service only
- **Audit Trail**: All data access is logged and monitored

#### Privacy & Compliance
- **PII Protection**: Buyer names and documents are access-controlled
- **Card Data**: Only tokenized/masked card information is stored
- **Regulatory Compliance**: Data retention and access patterns comply with financial regulations

## Example Queries

### 1. Real-time Risk Decision Monitoring

```sql
-- Monitor risk decisions in the last hour
SELECT 
    decision_status,
    risk_level,
    COUNT(*) as transaction_count,
    AVG(total_risk_score) as avg_risk_score,
    AVG(payment_amount) as avg_payment_amount
FROM risk_analysis 
WHERE analysis_timestamp >= NOW() - INTERVAL '1' HOUR
GROUP BY decision_status, risk_level
ORDER BY transaction_count DESC;
```

**Business Value**: Real-time monitoring of risk decision patterns and transaction volumes for operational awareness.

### 2. Fraud Detection Effectiveness Analysis

```sql
-- Analyze fraud detection score distribution by decision outcome
SELECT 
    decision_status,
    risk_level,
    COUNT(*) as transactions,
    AVG(value_score) as avg_value_score,
    AVG(seller_score) as avg_seller_score,
    AVG(currency_score) as avg_currency_score,
    AVG(total_risk_score) as avg_total_score
FROM risk_analysis 
WHERE analysis_date >= '2024-01-01'
GROUP BY decision_status, risk_level
ORDER BY decision_status, avg_total_score DESC;
```

**Business Value**: Evaluate the effectiveness of different fraud detection components and optimize scoring weights.

### 3. High-Risk Transaction Pattern Analysis

```sql
-- Identify patterns in high-risk rejected transactions
SELECT 
    seller_id,
    payment_currency,
    COUNT(*) as rejected_count,
    AVG(payment_amount) as avg_amount,
    AVG(total_risk_score) as avg_risk_score,
    COUNT(DISTINCT buyer_document) as unique_buyers
FROM risk_analysis 
WHERE decision_status = 'rejected' 
    AND risk_level = 'High'
    AND analysis_date >= '2024-01-01'
GROUP BY seller_id, payment_currency
HAVING rejected_count >= 5
ORDER BY rejected_count DESC, avg_risk_score DESC
LIMIT 20;
```

**Business Value**: Identify high-risk merchants and currency combinations for targeted risk policy adjustments.

### 4. Time-based Risk Pattern Analysis

```sql
-- Analyze risk patterns by hour of day
SELECT 
    analysis_hour,
    decision_status,
    COUNT(*) as transaction_count,
    AVG(total_risk_score) as avg_risk_score,
    COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (PARTITION BY analysis_hour) as percentage
FROM risk_analysis 
WHERE analysis_date >= '2024-01-01'
GROUP BY analysis_hour, decision_status
ORDER BY analysis_hour, decision_status;
```

**Business Value**: Understand time-based risk patterns for optimizing fraud detection sensitivity and business operations.

### 5. Business Impact Assessment

```sql
-- Calculate business impact of risk decisions
SELECT 
    analysis_date,
    decision_status,
    COUNT(*) as transaction_count,
    SUM(payment_amount) as total_amount,
    AVG(payment_amount) as avg_amount,
    SUM(CASE WHEN decision_status = 'approved' THEN payment_amount ELSE 0 END) as approved_amount,
    SUM(CASE WHEN decision_status = 'rejected' THEN payment_amount ELSE 0 END) as rejected_amount
FROM risk_analysis 
WHERE analysis_date >= '2024-01-01'
GROUP BY analysis_date, decision_status
ORDER BY analysis_date DESC, decision_status;
```

**Business Value**: Measure the financial impact of risk decisions on business revenue and assess the balance between security and growth.

## Getting Started

### Prerequisites
- Apache Pinot cluster access
- Kafka cluster with risk decision topics
- Appropriate permissions for data access

### Setup Instructions

1. **Deploy Pinot Schema**:
   ```bash
   curl -X POST "http://pinot-controller:9000/schemas" \
     -H "Content-Type: application/json" \
     -d @risk-analysis-schema.json
   ```

2. **Deploy Pinot Table**:
   ```bash
   curl -X POST "http://pinot-controller:9000/tables" \
     -H "Content-Type: application/json" \
     -d @risk-analysis-table-config.json
   ```

3. **Verify Setup**:
   ```bash
   curl "http://pinot-controller:9000/tables/risk_analysis/schema"
   ```

### Data Access

- **Pinot Broker**: `pinot-broker:8099`
- **Query Console**: `http://pinot-console:9000`
- **Grafana Dashboard**: `http://grafana:3000/dashboard/risk-analysis`

### Support & Contact

- **Technical Support**: risk-management-team@funny-bunny.xyz
- **Business Questions**: product-management@funny-bunny.xyz
- **Data Governance**: data-governance@funny-bunny.xyz

### Change Log

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2024-01-15 | Initial release with real-time risk analysis data product |

---

*This data product is part of the Risk Management domain and follows Data Mesh principles for decentralized data ownership and management.*