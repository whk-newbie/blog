"""Example of using the crawler and task reporter"""

import time
from blog_sdk import Crawler, TaskReporter, TaskStatus, Config


def main():
    # Load configuration from environment
    config = Config.from_env()
    
    # Initialize task reporter
    if not config.crawler_token:
        print("Error: BLOG_CRAWLER_TOKEN not set in environment")
        print("Please set BLOG_CRAWLER_TOKEN in your .env file or environment variables")
        return
    
    reporter = TaskReporter(
        api_base_url=config.api_base_url,
        token=config.crawler_token,
    )
    
    # Initialize crawler
    crawler = Crawler(base_url="https://example.com")
    
    # Register a new task
    task_id = f"crawl-task-{int(time.time())}"
    task_name = "Example Crawl Task"
    
    try:
        print(f"Registering task: {task_id}")
        task = reporter.register_task(
            task_id=task_id,
            task_name=task_name,
            metadata={"source": "example_script"},
        )
        print(f"Task registered: {task.get('task_id')}")
        
        # Update progress
        urls = [
            "https://example.com/page1",
            "https://example.com/page2",
            "https://example.com/page3",
        ]
        
        for i, url in enumerate(urls):
            print(f"Crawling: {url}")
            result = crawler.crawl(url)
            
            if result:
                print(f"  Title: {result['title']}")
                print(f"  Status: {result['status_code']}")
                print(f"  Content length: {len(result['content'])} chars")
                
                # Update task progress
                progress = int((i + 1) / len(urls) * 100)
                reporter.update_status(
                    task_id=task_id,
                    status=TaskStatus.RUNNING,
                    progress=progress,
                    message=f"Crawled {i + 1}/{len(urls)} URLs",
                )
            else:
                print(f"  Failed to crawl: {url}")
            
            time.sleep(1)  # Rate limiting
        
        # Complete the task
        print("Completing task...")
        reporter.complete_task(
            task_id=task_id,
            message="All URLs crawled successfully",
            metadata={"total_urls": len(urls), "status": "success"},
        )
        print("Task completed!")
        
    except Exception as e:
        print(f"Error: {e}")
        # Mark task as failed
        try:
            reporter.fail_task(
                task_id=task_id,
                message="Task failed due to error",
                error=str(e),
            )
        except Exception as fail_error:
            print(f"Failed to report task failure: {fail_error}")


if __name__ == "__main__":
    main()

