"""HTTP Client with Authentication"""

import requests
from typing import Optional, Dict, Any
from loguru import logger


class HTTPClient:
    """HTTP Client with Bearer Token authentication"""
    
    def __init__(
        self,
        base_url: str,
        token: Optional[str] = None,
        timeout: int = 30,
        headers: Optional[Dict[str, str]] = None,
    ):
        """
        Initialize HTTP client
        
        Args:
            base_url: Base URL of the API
            token: Bearer token for authentication
            timeout: Request timeout in seconds
            headers: Additional headers
        """
        self.base_url = base_url.rstrip("/")
        self.token = token
        self.timeout = timeout
        self.session = requests.Session()
        
        # Set default headers
        default_headers = {
            "Content-Type": "application/json",
            "Accept": "application/json",
        }
        if headers:
            default_headers.update(headers)
        self.session.headers.update(default_headers)
        
        # Set authorization header if token provided
        if token:
            self.set_token(token)
    
    def set_token(self, token: str) -> None:
        """
        Set Bearer token for authentication
        
        Args:
            token: Bearer token
        """
        self.token = token
        self.session.headers.update({"Authorization": f"Bearer {token}"})
    
    def get(self, path: str, params: Optional[Dict[str, Any]] = None, **kwargs) -> requests.Response:
        """
        Send GET request
        
        Args:
            path: API path
            params: Query parameters
            **kwargs: Additional arguments for requests
            
        Returns:
            Response object
        """
        url = f"{self.base_url}{path}"
        kwargs.setdefault("timeout", self.timeout)
        try:
            response = self.session.get(url, params=params, **kwargs)
            response.raise_for_status()
            return response
        except requests.RequestException as e:
            logger.error(f"GET request failed: {url}, error: {e}")
            raise
    
    def post(self, path: str, data: Optional[Dict[str, Any]] = None, json: Optional[Dict[str, Any]] = None, **kwargs) -> requests.Response:
        """
        Send POST request
        
        Args:
            path: API path
            data: Form data
            json: JSON data
            **kwargs: Additional arguments for requests
            
        Returns:
            Response object
        """
        url = f"{self.base_url}{path}"
        kwargs.setdefault("timeout", self.timeout)
        try:
            response = self.session.post(url, data=data, json=json, **kwargs)
            response.raise_for_status()
            return response
        except requests.RequestException as e:
            logger.error(f"POST request failed: {url}, error: {e}")
            raise
    
    def put(self, path: str, data: Optional[Dict[str, Any]] = None, json: Optional[Dict[str, Any]] = None, **kwargs) -> requests.Response:
        """
        Send PUT request
        
        Args:
            path: API path
            data: Form data
            json: JSON data
            **kwargs: Additional arguments for requests
            
        Returns:
            Response object
        """
        url = f"{self.base_url}{path}"
        kwargs.setdefault("timeout", self.timeout)
        try:
            response = self.session.put(url, data=data, json=json, **kwargs)
            response.raise_for_status()
            return response
        except requests.RequestException as e:
            logger.error(f"PUT request failed: {url}, error: {e}")
            raise
    
    def delete(self, path: str, **kwargs) -> requests.Response:
        """
        Send DELETE request
        
        Args:
            path: API path
            **kwargs: Additional arguments for requests
            
        Returns:
            Response object
        """
        url = f"{self.base_url}{path}"
        kwargs.setdefault("timeout", self.timeout)
        try:
            response = self.session.delete(url, **kwargs)
            response.raise_for_status()
            return response
        except requests.RequestException as e:
            logger.error(f"DELETE request failed: {url}, error: {e}")
            raise

