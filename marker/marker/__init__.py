"""Marker — Python script grader that runs inside Docker containers."""

__version__ = "0.1.0"

from .grader import get_result

__all__ = ["get_result"]
