# Blog SDK

Python SDK for blog system crawler and monitoring.

## Features

- Web crawling and content extraction
- Task reporting and status tracking
- API health monitoring
- System metrics collection
- Configuration management
- HTTP client with authentication

## Installation

Using uv (recommended):

```bash
uv pip install -e .
```

Using pip:

```bash
pip install -e .
```

## Configuration

Create a `.env` file in your project root:

```env
BLOG_API_BASE_URL=http://localhost:8080
BLOG_CRAWLER_TOKEN=your-crawler-token-here
BLOG_TIMEOUT=30
```

Or set environment variables:

```bash
export BLOG_API_BASE_URL=http://localhost:8080
export BLOG_CRAWLER_TOKEN=your-crawler-token-here
export BLOG_TIMEOUT=30
```

## Usage

### Configuration

```python
from blog_sdk import Config

# Load from environment variables
config = Config.from_env()

# Or load from a specific .env file
config = Config.from_env(".env")

# Or load from a config file
config = Config.from_file("config.env")

# Access configuration
print(config.api_base_url)
print(config.crawler_token)
```

### Crawler

```python
from blog_sdk import Crawler

# Initialize crawler
crawler = Crawler(base_url="https://example.com")

# Crawl a URL
result = crawler.crawl("https://example.com/article")
if result:
    print(f"Title: {result['title']}")
    print(f"Status: {result['status_code']}")
    print(f"Content: {result['content']}")
```

### Task Reporter

The TaskReporter allows you to report crawler task status to the blog API:

```python
from blog_sdk import TaskReporter, TaskStatus, Config

# Load configuration
config = Config.from_env()

# Initialize task reporter
reporter = TaskReporter(
    api_base_url=config.api_base_url,
    token=config.crawler_token,
)

# Register a new task
task_id = "my-task-123"
task = reporter.register_task(
    task_id=task_id,
    task_name="My Crawling Task",
    metadata={"source": "my_script"},
)

# Update task progress
reporter.update_status(
    task_id=task_id,
    status=TaskStatus.RUNNING,
    progress=50,
    message="Processing pages...",
)

# Complete the task
reporter.complete_task(
    task_id=task_id,
    message="Task completed successfully",
    metadata={"total_pages": 100},
)

# Or mark as failed
reporter.fail_task(
    task_id=task_id,
    message="Task failed",
    error="Connection timeout",
)
```

### Monitor

```python
from blog_sdk import Monitor, Config

# Load configuration
config = Config.from_env()

# Initialize monitor
monitor = Monitor(
    api_base_url=config.api_base_url,
    api_key=config.crawler_token,  # Optional
)

# Check health
health = monitor.check_health()
print(f"Status: {health['status']}")
print(f"Response time: {health['response_time']}")

# Get metrics
metrics = monitor.get_metrics()
if metrics:
    print(metrics)
```

### HTTP Client

For custom API requests:

```python
from blog_sdk import HTTPClient

# Initialize HTTP client
client = HTTPClient(
    base_url="http://localhost:8080",
    token="your-token",
)

# Make requests
response = client.get("/api/v1/some-endpoint")
data = response.json()

# Or with custom headers
response = client.post(
    "/api/v1/some-endpoint",
    json={"key": "value"},
)
```

## Examples

See the `examples/` directory for complete examples:

- `crawler_example.py` - Basic crawler usage with task reporting
- `monitor_example.py` - API health monitoring
- `task_reporter_example.py` - Task reporting workflow

Run an example:

```bash
python examples/task_reporter_example.py
```

## API Reference

### TaskReporter

- `register_task(task_id, task_name, metadata=None)` - Register a new task
- `update_status(task_id, status, progress, message=None)` - Update task status
- `complete_task(task_id, message=None, metadata=None)` - Mark task as completed
- `fail_task(task_id, message=None, error=None, metadata=None)` - Mark task as failed

### TaskStatus Enum

- `TaskStatus.RUNNING` - Task is running
- `TaskStatus.COMPLETED` - Task completed successfully
- `TaskStatus.FAILED` - Task failed

### Config

- `Config.from_env(env_file=None)` - Load from environment variables
- `Config.from_file(config_path)` - Load from file
- `config.validate()` - Validate configuration

## Development

Install development dependencies:

```bash
uv pip install -e ".[dev]"
```

Run tests:

```bash
pytest
```

Format code:

```bash
black src/
```

Lint code:

```bash
ruff check src/
```

## License

MIT

