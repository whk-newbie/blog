"""Configuration Management"""

import os
from typing import Optional
from pathlib import Path
from dotenv import load_dotenv
from pydantic import BaseModel, Field


class Config(BaseModel):
    """SDK Configuration"""
    
    api_base_url: str = Field(..., description="Base URL of the blog API")
    crawler_token: Optional[str] = Field(None, description="Bearer token for crawler authentication")
    timeout: int = Field(30, description="Request timeout in seconds")
    
    @classmethod
    def from_env(cls, env_file: Optional[str] = None) -> "Config":
        """
        Load configuration from environment variables
        
        Args:
            env_file: Path to .env file (optional)
            
        Returns:
            Config instance
        """
        if env_file:
            load_dotenv(env_file)
        else:
            # Try to load from current directory or parent directories
            load_dotenv()
        
        return cls(
            api_base_url=os.getenv("BLOG_API_BASE_URL", "http://localhost:8080"),
            crawler_token=os.getenv("BLOG_CRAWLER_TOKEN"),
            timeout=int(os.getenv("BLOG_TIMEOUT", "30")),
        )
    
    @classmethod
    def from_file(cls, config_path: str) -> "Config":
        """
        Load configuration from file
        
        Args:
            config_path: Path to configuration file
            
        Returns:
            Config instance
        """
        path = Path(config_path)
        if not path.exists():
            raise FileNotFoundError(f"Config file not found: {config_path}")
        
        # Load as .env file
        return cls.from_env(str(path))
    
    def validate(self) -> None:
        """
        Validate configuration
        
        Raises:
            ValueError: If configuration is invalid
        """
        if not self.api_base_url:
            raise ValueError("api_base_url is required")
        
        if not self.api_base_url.startswith(("http://", "https://")):
            raise ValueError("api_base_url must start with http:// or https://")

