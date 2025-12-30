"""Task Reporter for Crawler Tasks"""

import requests
from enum import Enum
from typing import Optional, Dict, Any
from loguru import logger
from ..utils.http_client import HTTPClient


class TaskStatus(str, Enum):
    """Crawl Task Status"""
    RUNNING = "running"
    COMPLETED = "completed"
    FAILED = "failed"


class TaskReporter:
    """Task Reporter for reporting crawler task status to the blog API"""
    
    def __init__(
        self,
        api_base_url: str,
        token: str,
        timeout: int = 30,
    ):
        """
        Initialize task reporter
        
        Args:
            api_base_url: Base URL of the blog API
            token: Bearer token for authentication
            timeout: Request timeout in seconds
        """
        self.client = HTTPClient(
            base_url=api_base_url,
            token=token,
            timeout=timeout,
        )
    
    def register_task(
        self,
        task_id: str,
        task_name: str,
        metadata: Optional[Dict[str, Any]] = None,
    ) -> Dict[str, Any]:
        """
        Register a new crawler task
        
        Args:
            task_id: Unique task identifier
            task_name: Task name
            metadata: Optional metadata dictionary
            
        Returns:
            Task data from API response
            
        Raises:
            requests.RequestException: If request fails
        """
        url = "/api/v1/crawler/tasks"
        payload = {
            "task_id": task_id,
            "task_name": task_name,
        }
        if metadata:
            payload["metadata"] = metadata
        
        try:
            response = self.client.post(url, json=payload)
            data = response.json()
            if data.get("code") == 201:
                logger.info(f"Task registered: {task_id}")
                return data.get("data", {})
            else:
                raise Exception(f"Failed to register task: {data.get('message', 'Unknown error')}")
        except requests.RequestException as e:
            logger.error(f"Failed to register task {task_id}: {e}")
            raise
    
    def update_status(
        self,
        task_id: str,
        status: TaskStatus,
        progress: int = 0,
        message: Optional[str] = None,
    ) -> Dict[str, Any]:
        """
        Update task status
        
        Args:
            task_id: Task identifier
            status: Task status
            progress: Progress percentage (0-100)
            message: Optional status message
            
        Returns:
            Updated task data from API response
            
        Raises:
            requests.RequestException: If request fails
            ValueError: If progress is not between 0 and 100
        """
        if not 0 <= progress <= 100:
            raise ValueError("Progress must be between 0 and 100")
        
        url = f"/api/v1/crawler/tasks/{task_id}"
        payload = {
            "status": status.value,
            "progress": progress,
        }
        if message:
            payload["message"] = message
        
        try:
            response = self.client.put(url, json=payload)
            data = response.json()
            if data.get("code") == 200:
                logger.debug(f"Task status updated: {task_id} -> {status.value} ({progress}%)")
                return data.get("data", {})
            else:
                raise Exception(f"Failed to update task status: {data.get('message', 'Unknown error')}")
        except requests.RequestException as e:
            logger.error(f"Failed to update task status for {task_id}: {e}")
            raise
    
    def complete_task(
        self,
        task_id: str,
        message: Optional[str] = None,
        metadata: Optional[Dict[str, Any]] = None,
    ) -> Dict[str, Any]:
        """
        Mark task as completed
        
        Args:
            task_id: Task identifier
            message: Optional completion message
            metadata: Optional metadata dictionary
            
        Returns:
            Completed task data from API response
            
        Raises:
            requests.RequestException: If request fails
        """
        url = f"/api/v1/crawler/tasks/{task_id}/complete"
        payload = {}
        if message:
            payload["message"] = message
        if metadata:
            payload["metadata"] = metadata
        
        try:
            response = self.client.put(url, json=payload)
            data = response.json()
            if data.get("code") == 200:
                logger.info(f"Task completed: {task_id}")
                return data.get("data", {})
            else:
                raise Exception(f"Failed to complete task: {data.get('message', 'Unknown error')}")
        except requests.RequestException as e:
            logger.error(f"Failed to complete task {task_id}: {e}")
            raise
    
    def fail_task(
        self,
        task_id: str,
        message: Optional[str] = None,
        error: Optional[str] = None,
        metadata: Optional[Dict[str, Any]] = None,
    ) -> Dict[str, Any]:
        """
        Mark task as failed
        
        Args:
            task_id: Task identifier
            message: Optional failure message
            error: Optional error message
            metadata: Optional metadata dictionary
            
        Returns:
            Failed task data from API response
            
        Raises:
            requests.RequestException: If request fails
        """
        url = f"/api/v1/crawler/tasks/{task_id}/fail"
        payload = {}
        if message:
            payload["message"] = message
        if error:
            payload["error"] = error
        if metadata:
            payload["metadata"] = metadata
        
        try:
            response = self.client.put(url, json=payload)
            data = response.json()
            if data.get("code") == 200:
                logger.warning(f"Task failed: {task_id}")
                return data.get("data", {})
            else:
                raise Exception(f"Failed to mark task as failed: {data.get('message', 'Unknown error')}")
        except requests.RequestException as e:
            logger.error(f"Failed to mark task as failed {task_id}: {e}")
            raise

