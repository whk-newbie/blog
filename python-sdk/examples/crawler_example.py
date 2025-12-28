"""Example of using the crawler"""

from blog_sdk import Crawler


def main():
    # Initialize crawler
    crawler = Crawler(base_url="https://example.com")
    
    # Crawl a URL
    result = crawler.crawl("https://example.com")
    
    if result:
        print(f"Title: {result['title']}")
        print(f"Status: {result['status_code']}")
        print(f"Content length: {len(result['content'])} chars")
    else:
        print("Failed to crawl URL")


if __name__ == "__main__":
    main()

