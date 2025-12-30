"""Example of using the monitor"""

import time
from blog_sdk import Monitor, Config


def main():
    # Load configuration from environment
    config = Config.from_env()
    
    # Initialize monitor
    monitor = Monitor(
        api_base_url=config.api_base_url,
        api_key=config.crawler_token,  # Optional: use crawler token if available
    )
    
    print(f"Monitoring API: {config.api_base_url}")
    print("Press Ctrl+C to stop\n")
    
    # Continuous health monitoring
    try:
        while True:
            health = monitor.check_health()
            timestamp = time.strftime('%Y-%m-%d %H:%M:%S')
            print(f"[{timestamp}] Status: {health['status']}")
            
            if health['status'] == 'healthy':
                print(f"  Response time: {health['response_time']}")
            else:
                print(f"  Error: {health.get('error', 'Unknown')}")
            
            # Try to get metrics if available
            if health['status'] == 'healthy':
                metrics = monitor.get_metrics()
                if metrics:
                    print(f"  Metrics available: {len(metrics)} items")
            
            print()
            time.sleep(10)  # Check every 10 seconds
            
    except KeyboardInterrupt:
        print("\nMonitoring stopped")


if __name__ == "__main__":
    main()

