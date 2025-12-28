# Blog SDK

Python SDK for blog system crawler and monitoring.

## Features

- Web crawling and content extraction
- API health monitoring
- System metrics collection

## Installation

Using uv (recommended):

```bash
uv pip install -e .
```

Using pip:

```bash
pip install -e .
```

## Usage

### Crawler

```python
from blog_sdk import Crawler

# Initialize crawler
crawler = Crawler(base_url="https://example.com")

# Crawl a URL
result = crawler.crawl("https://example.com/article")
print(result)
```

### Monitor

```python
from blog_sdk import Monitor

# Initialize monitor
monitor = Monitor(api_base_url="http://localhost:8080", api_key="your-api-key")

# Check health
health = monitor.check_health()
print(health)

# Get metrics
metrics = monitor.get_metrics()
print(metrics)
```

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

