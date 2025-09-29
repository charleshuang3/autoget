import os
import sys
from contextlib import asynccontextmanager
from fastapi import FastAPI


@asynccontextmanager
async def lifespan(app: FastAPI):
  # Startup
  grok_key = os.getenv("GROK_KEY")
  search_mcp = os.getenv("SEARCH_MCP")

  if not grok_key:
    print("Error: GROK_KEY environment variable is not set or is empty.", file=sys.stderr)
    sys.exit(1)

  if not search_mcp:
    print("Error: SEARCH_MCP environment variable is not set or is empty.", file=sys.stderr)
    sys.exit(1)

  yield
  # Shutdown


app = FastAPI(lifespan=lifespan)
