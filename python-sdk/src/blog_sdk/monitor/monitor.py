"""System Monitoring Implementation"""

import time
import requests
from loguru import logger
from typing import Dict, Any, Optional


class Monitor:
    """System and API Monitoring"""
    
    def __init__(self, api_base_url: str, api_key: Optional[str] = None):
        """
        Initialize monitor
        
        Args:
            api_base_url: Base URL of the blog API
            api_key: Optional API key for authentication
        """
        self.api_base_url = api_base_url.rstrip("/")
        self.api_key = api_key
        self.session = requests.Session()
        
        if api_key:
            self.session.headers.update({"Authorization": f"Bearer {api_key}"})
            
    def check_health(self) -> Dict[str, Any]:
        """
        Check API health status
        
        Returns:
            Health check result
        """
        try:
            start_time = time.time()
            response = self.session.get(f"{self.api_base_url}/health", timeout=5)
            elapsed = time.time() - start_time
            
            return {
                "status": "healthy" if response.status_code == 200 else "unhealthy",
                "status_code": response.status_code,
                "response_time": f"{elapsed:.3f}s",
                "timestamp": time.time(),
            }
        except requests.RequestException as e:
            logger.error(f"Health check failed: {e}")
            return {
                "status": "unhealthy",
                "error": str(e),
                "timestamp": time.time(),
            }
            
    def get_metrics(self) -> Optional[Dict[str, Any]]:
        """
        Get system metrics
        
        Returns:
            Metrics data or None
        """
        try:
            response = self.session.get(f"{self.api_base_url}/api/v1/admin/metrics")
            response.raise_for_status()
            return response.json()
        except requests.RequestException as e:
            logger.error(f"Failed to get metrics: {e}")
            return None

