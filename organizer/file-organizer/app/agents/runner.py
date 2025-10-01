import uuid

from google.adk.sessions import InMemorySessionService
from google.adk.runners import Runner
from google.genai import types

from .models import PlanRequest, PlanResponse

APP_NAME = "autoget_organizer"
USER_ID = "user"

session_service = InMemorySessionService()
runner = Runner(app_name=APP_NAME, session_service=session_service)


async def run(request: PlanRequest) -> PlanResponse:
  query = request.model_dump_json()

  content = types.Content(role="user", parts=[types.Part(text=query)])

  final_response_text = "Agent did not produce a final response."

  session_id = uuid.uuid4().hex

  await session_service.create_session(app_name=APP_NAME, user_id=USER_ID, session_id=session_id)

  # Execute the agent and find the final response
  async for event in runner.run_async(user_id=USER_ID, session_id=session_id, new_message=content):
    if event.is_final_response():
      if event.content and event.content.parts:
        final_response_text = event.content.parts[0].text
      break

  await session_service.delete_session(app_name=APP_NAME, user_id=USER_ID, session_id=session_id)

  return PlanResponse.model_validate_json(final_response_text)
