import os

from google.adk.tools.mcp_tool.mcp_toolset import McpToolset, StreamableHTTPConnectionParams

search_mcp = os.getenv("SEARCH_MCP")

mcp_search_tool = McpToolset(
  connection_params=StreamableHTTPConnectionParams(
    url=search_mcp,
  )
)
