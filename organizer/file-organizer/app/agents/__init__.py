import os

if os.getenv("XAI_API_KEY"):
  from .organizer import OrganizerAgent

  root_agent = OrganizerAgent()
