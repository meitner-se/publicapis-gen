# GitHub Actions Cache Management

This document explains the cache optimization strategies implemented to resolve the cache storage limit issue (INF-298).

## Problem Statement

The repository was approaching GitHub's 10GB cache storage limit with 11.8GB usage and 60 active caches. This was caused by:

- Matrix strategy creating multiple cache entries per job
- Duplicate cache configurations across jobs
- No automated cache cleanup
- Inefficient cache key strategies

## Implemented Solutions

### 1. Optimized CI Workflow (`.github/workflows/ci.yml`)

**Key Improvements:**
- **Centralized Cache Setup**: Added a `setup` job that creates and manages the primary cache
- **Reduced Cache Duplication**: Only the primary Go version (1.22) uses caching, other versions download fresh
- **Optimized Cache Keys**: Improved cache key strategy with version control
- **Cache Restoration**: Uses `actions/cache/restore` for read-only cache access in dependent jobs

**Cache Strategy:**
```yaml
# Primary cache key includes Go version, cache version, and dependency hashes
key: ${{ runner.os }}-go-${{ env.GO_VERSION }}-${{ env.GO_CACHE_VERSION }}-${{ hashFiles('**/go.sum', '**/go.mod') }}
```

### 2. Automated Cache Cleanup (`.github/workflows/cache-cleanup.yml`)

**Features:**
- **Scheduled Cleanup**: Runs daily at 2 AM UTC to remove caches older than 7 days
- **Manual Triggers**: Supports on-demand cleanup with customizable parameters
- **Dry Run Mode**: Test cleanup operations without actually deleting caches
- **Configurable Age Limits**: Adjust maximum age of caches to keep
- **Detailed Reporting**: Provides comprehensive cleanup summaries

**Usage Examples:**
```bash
# Manual cleanup with default settings (7 days)
gh workflow run cache-cleanup.yml

# Dry run to see what would be deleted
gh workflow run cache-cleanup.yml -f dry_run=true

# Custom cleanup (keep only 3 days)
gh workflow run cache-cleanup.yml -f max_age_days=3
```

### 3. Cache Monitoring (`.github/workflows/cache-monitoring.yml`)

**Features:**
- **Usage Reports**: Regular cache usage analysis and statistics
- **Pattern Analysis**: Identifies cache key patterns and potential optimizations
- **Age Distribution**: Shows how old caches are
- **Recommendations**: Provides actionable insights based on usage patterns
- **Health Checks**: Weekly automated monitoring

## Expected Impact

### Cache Reduction
- **From 6 caches per push** (2 Go versions × 3 jobs) **to 1 cache per push**
- **~83% reduction** in cache creation rate
- **Estimated 70-80% reduction** in total cache storage usage

### Storage Management
- Automated daily cleanup of caches older than 7 days
- Manual cleanup capabilities for immediate space recovery
- Proactive monitoring to prevent future limit issues

## Cache Management Best Practices

### 1. Cache Key Design
- Include dependency file hashes (`go.sum`, `go.mod`)
- Use cache version variables for controlled cache busting
- Avoid overly specific keys that prevent cache reuse

### 2. Cache Lifecycle
- **Creation**: Only create caches for the primary build configuration
- **Usage**: Restore caches in read-only mode for dependent jobs
- **Cleanup**: Automatically remove caches older than 7 days
- **Monitoring**: Weekly health checks and usage reports

### 3. Manual Cache Management

**Check current usage:**
```bash
gh cache list --repo meitner-se/publicapis-gen
```

**Delete specific cache:**
```bash
gh cache delete <cache-id> --repo meitner-se/publicapis-gen
```

**Run immediate cleanup:**
```bash
gh workflow run cache-cleanup.yml
```

### 4. Troubleshooting

**If cache limit is reached again:**
1. Run immediate cleanup: `gh workflow run cache-cleanup.yml`
2. Check cache monitoring report for patterns
3. Consider reducing cleanup interval (from 7 to 3-5 days)
4. Review cache key strategies for optimization

**If builds are slower after optimization:**
1. Check that the primary Go version matches your main development version
2. Verify cache hits in the setup job logs
3. Consider caching additional paths if needed

## Monitoring and Alerts

### Automated Monitoring
- **Daily**: Cleanup runs automatically
- **Weekly**: Cache health check reports
- **After CI**: Usage monitoring after each CI run

### Manual Monitoring
```bash
# View current cache status
gh cache list --repo meitner-se/publicapis-gen

# Generate detailed report
gh workflow run cache-monitoring.yml

# Check workflow run logs for detailed analysis
gh run list --workflow=cache-monitoring.yml
```

## Configuration Variables

### Environment Variables (CI)
- `GO_VERSION`: Primary Go version for caching (default: '1.22')
- `GO_CACHE_VERSION`: Cache version for controlled busting (default: 'v1')

### Cleanup Configuration
- `MAX_AGE_DAYS`: Maximum age of caches to keep (default: 7 days)
- `DRY_RUN`: Test mode without actual deletion (default: false)

## Future Considerations

1. **Cache Size Monitoring**: Consider adding cache size tracking and alerts
2. **Advanced Cleanup**: Implement smart cleanup based on cache usage patterns
3. **Cross-Job Dependencies**: Optimize caching for complex workflow dependencies
4. **Performance Metrics**: Track build time impact of cache optimizations

This optimization should resolve the immediate cache storage issue and provide sustainable cache management going forward.