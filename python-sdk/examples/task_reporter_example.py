"""Example of using the task reporter"""

import time
import uuid
from blog_sdk import TaskReporter, TaskStatus, Config


def simulate_crawling_task(reporter: TaskReporter, task_id: str, task_name: str):
    """Simulate a crawling task with progress updates"""
    
    # Register the task
    print(f"Registering task: {task_id}")
    task = reporter.register_task(
        task_id=task_id,
        task_name=task_name,
        metadata={"type": "example", "simulation": True},
    )
    print(f"Task registered successfully: {task.get('task_id')}")
    print()
    
    # Simulate task progress
    steps = [
        ("Initializing crawler", 10),
        ("Fetching URLs", 30),
        ("Crawling pages", 60),
        ("Processing content", 80),
        ("Saving results", 95),
    ]
    
    try:
        for step_name, progress in steps:
            print(f"Step: {step_name} ({progress}%)")
            
            # Update task status
            reporter.update_status(
                task_id=task_id,
                status=TaskStatus.RUNNING,
                progress=progress,
                message=step_name,
            )
            
            # Simulate work
            time.sleep(1)
        
        # Complete the task
        print("\nCompleting task...")
        reporter.complete_task(
            task_id=task_id,
            message="Task completed successfully",
            metadata={
                "total_pages": 100,
                "success_rate": 0.95,
                "duration_seconds": len(steps),
            },
        )
        print("Task completed successfully!")
        
    except Exception as e:
        print(f"\nError during task execution: {e}")
        # Mark task as failed
        try:
            reporter.fail_task(
                task_id=task_id,
                message="Task failed due to error",
                error=str(e),
                metadata={"error_type": type(e).__name__},
            )
            print("Task marked as failed")
        except Exception as fail_error:
            print(f"Failed to report task failure: {fail_error}")


def main():
    # Load configuration from environment
    try:
        config = Config.from_env()
    except Exception as e:
        print(f"Error loading configuration: {e}")
        print("\nPlease create a .env file with the following variables:")
        print("  BLOG_API_BASE_URL=http://localhost:8080")
        print("  BLOG_CRAWLER_TOKEN=your-crawler-token")
        return
    
    # Validate configuration
    if not config.crawler_token:
        print("Error: BLOG_CRAWLER_TOKEN not set")
        print("Please set BLOG_CRAWLER_TOKEN in your .env file or environment variables")
        return
    
    try:
        config.validate()
    except ValueError as e:
        print(f"Configuration error: {e}")
        return
    
    # Initialize task reporter
    reporter = TaskReporter(
        api_base_url=config.api_base_url,
        token=config.crawler_token,
        timeout=config.timeout,
    )
    
    # Generate a unique task ID
    task_id = f"example-task-{uuid.uuid4().hex[:8]}"
    task_name = "Example Crawling Task"
    
    print(f"API Base URL: {config.api_base_url}")
    print(f"Task ID: {task_id}")
    print(f"Task Name: {task_name}")
    print("-" * 50)
    print()
    
    # Run the simulation
    simulate_crawling_task(reporter, task_id, task_name)


if __name__ == "__main__":
    main()

