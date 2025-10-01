import os
import sys
from contextlib import asynccontextmanager
from typing import List, Literal, Optional
from fastapi import FastAPI
from pydantic import BaseModel


# Pydantic models for /v1/plan
class PlanRequest(BaseModel):
  files: List[str]
  metadata: dict[str, str] = None


class PlanAction(BaseModel):
  file: str
  action: Literal["move", "ignore"]
  target: Optional[str] = None


class PlanResponse(BaseModel):
  plan: List[PlanAction]


# Pydantic models for /v1/execute
class ExecuteRequest(BaseModel):
  plan: List[PlanAction]


def check_env_vars(name: str):
  var = os.getenv(name)
  if not var:
    print(f"Error: {name} environment variable is not set or is empty.", file=sys.stderr)
    sys.exit(1)


@asynccontextmanager
async def lifespan(app: FastAPI):
  # Startup
  check_env_vars("XAI_API_KEY")
  check_env_vars("XAI_MODEL")
  check_env_vars("SEARCH_MCP")

  yield
  # Shutdown


app = FastAPI(lifespan=lifespan)


@app.post("/v1/plan", response_model=PlanResponse)
async def create_plan(request: PlanRequest):
  # For now, return a fake response as the agent caller is not ready
  fake_plan = []
  for file_path in request.files:
    if "document" in file_path.lower():
      fake_plan.append(
        PlanAction(file=file_path, action="move", target=f"documents/{os.path.basename(file_path)}")
      )
    else:
      fake_plan.append(PlanAction(file=file_path, action="ignore"))
  return PlanResponse(plan=fake_plan)


@app.post("/v1/execute")
async def execute_plan(request: ExecuteRequest):
  # For now, just return a 200 OK status
  return {"message": "Plan executed successfully (fake response)"}
