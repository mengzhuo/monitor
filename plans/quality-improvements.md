# Quality Improvements TODO

> **Note**: Primary testing strategy is now in `testing-strategy-unified.md` - this document tracks specific quality initiatives and technical debt.

## Overview

This document consolidates outstanding quality improvement tasks including testing enhancements, bug fixes, and code quality initiatives across the NTP Pool monitoring system.

## Test Coverage Improvement Plan (High Priority)

### Current Status: 53.6% Coverage
**Target**: 80%+ overall with 100% coverage for safety-critical functions

### Critical Functions Requiring Coverage

#### Phase 1: Safety Logic Testing (Critical)
**Target**: 100% coverage of safety and emergency conditions

**Functions with Low Coverage**:
- `applyRule6BootstrapPromotion`: **13.6%** ❌ (very low coverage)
- `calculateSafetyLimits`: 84.6% ⚠️ (needs improvement)
- `applyRule1ImmediateBlocking`: 75.0% ⚠️ (could be better)
- `applyRule2GradualConstraintRemoval`: 87.5% ⚠️ (good but could be complete)

**Emergency Condition Test Scenarios**:
```go
func TestCalculateSafetyLimits_EmergencyConditions(t *testing.T) {
    tests := []struct {
        name                string
        targetNumber        int
        totalMonitors       int
        healthyActive       int
        activeCount         int
        expectedMaxRemovals int
        expectEmergency     bool
    }{
        {
            name: "emergency_not_enough_monitors",
            targetNumber: 7, totalMonitors: 4,
            healthyActive: 2, activeCount: 3,
            expectedMaxRemovals: 0, expectEmergency: true,
        },
        {
            name: "safety_below_target_unhealthy",
            targetNumber: 7, totalMonitors: 10,
            healthyActive: 2, activeCount: 3,
            expectedMaxRemovals: 0, expectEmergency: false,
        },
        {
            name: "normal_operation_above_target",
            targetNumber: 7, totalMonitors: 10,
            healthyActive: 8, activeCount: 8,
            expectedMaxRemovals: 2, expectEmergency: false,
        },
    }
}
```

#### Phase 2: Bootstrap and Edge Case Coverage (High Priority)
**Target**: 90%+ coverage for bootstrap logic

**Bootstrap Test Scenarios**:
```go
func TestApplyRule6BootstrapPromotion(t *testing.T) {
    tests := []struct {
        name                string
        testingMonitors     []evaluatedMonitor
        candidateMonitors   []evaluatedMonitor
        expectedPromotions  int
        expectedReason      string
    }{
        {
            name: "zero_testing_promotes_candidates",
            testingMonitors: []evaluatedMonitor{},
            candidateMonitors: createHealthyCandidates(3),
            expectedPromotions: 3,
        },
        {
            name: "bootstrap_with_account_constraints",
            testingMonitors: []evaluatedMonitor{},
            candidateMonitors: createCandidatesWithAccountLimits(),
            expectedPromotions: 2, // limited by account
        },
        {
            name: "bootstrap_with_network_constraints",
            testingMonitors: []evaluatedMonitor{},
            candidateMonitors: createCandidatesWithNetworkConflicts(),
            expectedPromotions: 1, // limited by network diversity
        },
    }
}
```

#### Phase 3: Untested Infrastructure (Medium Priority)
**Target**: 80%+ coverage for database and metrics functions

**Currently Untested (0% coverage)**:
- `loadServerInfo`: Database loading logic
- `applyStatusChange`: Actual status change application
- `buildAccountLimitsFromMonitors`: Account limit calculation
- All metrics tracking functions
- `isGrandfathered`: Constraint grandfathering logic

**Database Integration Testing**:
```go
func TestDatabaseIntegration(t *testing.T) {
    // Use CI tools for proper database setup
    db := setupTestDatabase(t) // ./scripts/test-db.sh start
    defer cleanupTestDatabase(t, db)

    // Test constraint tracking persistence
    testConstraintViolationPersistence(t, db)

    // Test status change application
    testStatusChangeApplication(t, db)

    // Test concurrent operations
    testConcurrentSelectorRuns(t, db)
}
```

### Testing Infrastructure Improvements

#### CI Tool Integration
**Always use CI tools instead of manual testing**:
- `./scripts/test-ci-local.sh`: Full CI replication for final validation
- `./scripts/test-scorer-integration.sh`: Component testing for focused debugging
- `./scripts/test-minimal-ci.sh`: Quick isolation testing
- `./scripts/test-db.sh`: Direct database management

#### Test Data Validation Strategy
**Mathematical Test Correctness**:
```go
// Example: Validate test conditions make expected outcome possible
// MaxPerServer=2, ActiveCount=1, TestingCount=2
// Total limit = MaxPerServer + 1 = 3
// Current total = 1 + 2 = 3 (at limit)
// Promoting testing→active: would become 2 active + 1 testing = 3 (still valid)
accountLimits := map[uint32]*accountLimit{
    1: {AccountID: 1, MaxPerServer: 2, ActiveCount: 1, TestingCount: 2},
}
```

#### Test Helper Consolidation
```go
// Centralized test data creation for consistency
func createHealthyMonitors(count int, status ServerScoreStatus) []evaluatedMonitor
func createMonitorsWithConstraintViolations(violationType ConstraintType) []evaluatedMonitor
func createAccountLimitsScenario(limits map[uint32]int) map[uint32]*accountLimit

// Validation helpers
func validateWorkingCounts(result RuleResult, expected workingCounts) error
func validateConstraintViolations(changes []statusChange, expected []violationType) error
func validateStateTransitions(changes []statusChange, expected []transition) error
```

## Outstanding Bug Fixes (Medium Priority)

### Critical Bug #1: Emergency Override Coverage Gap
**Location**: `selector/constraints.go:611-629` in `canPromoteToTesting`
**Issue**: Emergency override only applies to testing→active promotions, not candidate→testing
**Impact**: System could get stuck with zero monitors if candidates can't be promoted

**Required Fix**:
```go
// Current signature - lacks emergency override
func (sl *Selector) canPromoteToTesting(
    monitor *monitorCandidate,
    server *serverInfo,
    accountLimits map[uint32]*accountLimit,
    existingMonitors []ntpdb.GetMonitorPriorityRow,
) bool

// Required signature - add emergency override parameter
func (sl *Selector) canPromoteToTesting(
    monitor *monitorCandidate,
    server *serverInfo,
    accountLimits map[uint32]*accountLimit,
    existingMonitors []ntpdb.GetMonitorPriorityRow,
    emergencyOverride bool, // NEW
) bool
```

**Implementation Steps**:
1. Add `emergencyOverride` parameter to function signature
2. Update all call sites to pass emergency override status
3. Apply same emergency logic as `canPromoteToActive()`
4. Add comprehensive test coverage for emergency scenarios

### Critical Bug #2: Non-Functional Grandfathering Logic
**Location**: `selector/state.go:104-113`
**Issue**: Grandfathered and non-grandfathered violations have identical behavior

**Current Broken Logic**:
```go
if violation.Type != violationNone {
    if violation.IsGrandfathered {
        return candidateOut  // Same as non-grandfathered
    }
    return candidateOut      // Same behavior
}
```

**Required Fix**: Implement different behavior for grandfathered violations
```go
if violation.Type != violationNone {
    if violation.IsGrandfathered {
        // Slower removal rate for grandfathered violations
        return candidateOut // But with different priority/timing
    }
    // Immediate removal for new violations
    return candidateBlock
}
```

### Logic Issue: Bootstrap Constraint Inconsistency
**Location**: `selector/process.go:484, 504` in bootstrap scenarios
**Issue**: Bootstrap checks constraints but emergency override doesn't apply consistently

**Current Problem**:
```go
// Bootstrap calls canPromoteToTesting without emergency override
if sl.canPromoteToTesting(&em.monitor, server, workingAccountLimits, assignedMonitors) {
```

**Required Fix**: Apply consistent emergency override logic across all promotion paths
```go
// Pass emergency override to bootstrap promotion attempts
emergencyOverride := len(activeMonitors) == 0
if sl.canPromoteToTesting(&em.monitor, server, workingAccountLimits, assignedMonitors, emergencyOverride) {
```

## Code Quality Initiatives (Low Priority)

### Common Bug Pattern Prevention

#### Safety Variable Scope Creep Prevention
**Pattern**: Variables intended for specific safety checks shouldn't be repurposed for broader blocking
**Implementation**:
```go
// Good: Specific purpose variables
maxActiveRemovals := calculateActiveRemovalLimit()
maxTestingRemovals := calculateTestingRemovalLimit()

// Bad: Generic variable used for multiple purposes
maxRemovals := 0 // Don't use for all removal types
```

#### Mathematical Consistency Enforcement
**Pattern**: Ensure counts track all state changes accurately
**Implementation**:
```go
// Always update working counts immediately after each change decision
func (wc *workingCounts) applyChange(change statusChange) {
    // Update counts synchronously
    switch change.fromStatus {
    case ServerScoresStatusActive: wc.active--
    case ServerScoresStatusTesting: wc.testing--
    }
    switch change.toStatus {
    case ServerScoresStatusActive: wc.active++
    case ServerScoresStatusTesting: wc.testing++
    }
}
```

#### Self-Reference Bug Prevention
**Pattern**: Always exclude the entity being evaluated from conflict detection
**Implementation**:
```go
// Always check for self-exclusion in constraint functions
func checkNetworkConstraints(currentMonitor *Monitor, existingMonitors []Monitor) error {
    for _, existing := range existingMonitors {
        if existing.ID == currentMonitor.ID {
            continue // Skip self
        }
        // Check constraint against other monitors
    }
}
```

### Code Quality Metrics

#### Static Analysis Integration
```yaml
# Add to CI pipeline
- name: Code Quality Analysis
  run: |
    # Go linting
    golangci-lint run ./...

    # Cyclomatic complexity check
    gocyclo -over 10 ./...

    # Dead code detection
    deadcode ./...

    # Security scanning
    gosec ./...
```

#### Quality Gates
- **Function Complexity**: Maximum cyclomatic complexity of 10
- **File Size**: Maximum 500 lines per file (except generated files)
- **Test Coverage**: Minimum 80% overall, 100% for safety functions
- **Documentation**: All exported functions must have godoc comments

### Performance Quality Assurance

#### Benchmark Integration
```go
func BenchmarkSelectorPerformance(b *testing.B) {
    scenarios := []struct {
        name         string
        monitorCount int
        complexity   string
    }{
        {"small_simple", 10, "basic"},
        {"medium_complex", 100, "constraints"},
        {"large_realistic", 1000, "production"},
    }

    for _, scenario := range scenarios {
        b.Run(scenario.name, func(b *testing.B) {
            data := generateBenchmarkData(scenario)
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                processSelector(data)
            }
        })
    }
}
```

#### Performance Regression Prevention
```yaml
# CI performance validation
- name: Performance Regression Check
  run: |
    # Run benchmarks
    go test -bench=. -benchmem ./selector

    # Compare with baseline
    benchcmp baseline.txt current.txt

    # Fail if regression > 10%
    if [ $REGRESSION_PERCENT -gt 10 ]; then
      echo "Performance regression detected: ${REGRESSION_PERCENT}%"
      exit 1
    fi
```

## Implementation Timeline

### Sprint 1: Critical Safety Testing
**Week 1-2**:
- ✅ Emergency condition test scenarios
- ✅ Bootstrap logic comprehensive testing
- ✅ Safety limit edge case coverage

### Sprint 2: Bug Fixes
**Week 3-4**:
- 🔲 Emergency override consistency fix
- 🔲 Grandfathering logic implementation
- 🔲 Bootstrap constraint consistency fix

### Sprint 3: Infrastructure Testing
**Week 5-6**:
- 🔲 Database integration test framework
- 🔲 Metrics tracking test coverage
- 🔲 Concurrent operation testing

### Sprint 4: Code Quality
**Week 7-8**:
- 🔲 Static analysis integration
- 🔲 Performance regression framework
- 🔲 Documentation completeness audit

## Success Metrics

### Coverage Targets
- **Overall Package Coverage**: 80%+ (from current 53.6%)
- **Critical Function Coverage**: 95%+ for all safety and emergency logic
- **Safety Logic Coverage**: 100% for all emergency/safety paths
- **Database Integration**: 80%+ coverage for all database operations

### Quality Metrics
- **Bug Density**: < 1 critical bug per 1000 lines of code
- **Test Stability**: < 1% flaky test rate
- **Performance Regression**: 0 regressions > 10% in selector performance
- **Code Review Coverage**: 100% of changes reviewed with quality focus

### Operational Improvements
- **Emergency Recovery Time**: < 5 seconds from zero monitors to restored service
- **Bug Detection Rate**: 90%+ of bugs caught in testing vs production
- **Test Execution Time**: < 60 seconds for full test suite
- **CI Reliability**: 99%+ CI success rate (excluding legitimate failures)

## Continuous Quality Process

### Regular Review Schedule
- **Weekly**: Coverage metrics review and gap identification
- **Monthly**: Code quality metrics assessment and improvement planning
- **Quarterly**: Comprehensive quality audit and process improvement

### Quality Gates Integration
```yaml
# Pre-merge quality gates
quality_gates:
  required_checks:
    - test_coverage_80_percent
    - no_critical_bugs
    - performance_regression_check
    - static_analysis_pass
    - security_scan_pass

  blocking_conditions:
    - coverage_below_threshold
    - unresolved_critical_bugs
    - performance_regression_detected
```

This quality improvement plan provides a systematic approach to enhancing code quality, test coverage, and operational reliability while maintaining development velocity.
