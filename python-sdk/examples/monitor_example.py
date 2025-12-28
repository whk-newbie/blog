"""Example of using the monitor"""

from blog_sdk import Monitor
import time


def main():
    # Initialize monitor
    monitor = Monitor(api_base_url="http://localhost:8080")
    
    # Continuous health monitoring
    while True:
        health = monitor.check_health()
        print(f"[{time.strftime('%Y-%m-%d %H:%M:%S')}] Status: {health['status']}")
        
        if health['status'] == 'healthy':
            print(f"  Response time: {health['response_time']}")
        else:
            print(f"  Error: {health.get('error', 'Unknown')}")
            
        time.sleep(10)  # Check every 10 seconds


if __name__ == "__main__":
    main()

