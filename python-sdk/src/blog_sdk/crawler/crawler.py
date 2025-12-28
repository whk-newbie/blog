"""Web Crawler Implementation"""

import requests
from bs4 import BeautifulSoup
from loguru import logger
from typing import Optional, Dict, Any


class Crawler:
    """Web Crawler for fetching and parsing web content"""
    
    def __init__(self, base_url: str, headers: Optional[Dict[str, str]] = None):
        """
        Initialize crawler
        
        Args:
            base_url: Base URL for crawling
            headers: Custom HTTP headers
        """
        self.base_url = base_url
        self.headers = headers or {
            "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
        }
        self.session = requests.Session()
        self.session.headers.update(self.headers)
        
    def fetch(self, url: str, **kwargs) -> Optional[requests.Response]:
        """
        Fetch content from URL
        
        Args:
            url: URL to fetch
            **kwargs: Additional arguments for requests
            
        Returns:
            Response object or None if failed
        """
        try:
            response = self.session.get(url, **kwargs)
            response.raise_for_status()
            return response
        except requests.RequestException as e:
            logger.error(f"Failed to fetch {url}: {e}")
            return None
            
    def parse_html(self, html: str) -> BeautifulSoup:
        """
        Parse HTML content
        
        Args:
            html: HTML string
            
        Returns:
            BeautifulSoup object
        """
        return BeautifulSoup(html, "lxml")
        
    def crawl(self, url: str) -> Optional[Dict[str, Any]]:
        """
        Crawl and parse a URL
        
        Args:
            url: URL to crawl
            
        Returns:
            Parsed data dictionary or None
        """
        response = self.fetch(url)
        if not response:
            return None
            
        soup = self.parse_html(response.text)
        
        return {
            "url": url,
            "status_code": response.status_code,
            "title": soup.title.string if soup.title else None,
            "content": soup.get_text(strip=True),
        }

