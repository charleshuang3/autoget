from typing import override, AsyncGenerator

from google.adk.events import Event
from google.adk.agents import BaseAgent, Agent, InvocationContext

from .categorizer import agent as categorizer_agent, CategoryResponse
from .models import PlanRequest, Category, category_list
from .utils.utils import simple_move_plan_event

simple_move_categories = [
  Category.photobook.name,
  Category.audio_book.name,
  Category.book.name,
  Category.music.name,
  Category.music_video.name,
]


class OrganizerAgent(BaseAgent):
  categorizer: Agent

  def __init__(self):
    categorizer_agent_ = categorizer_agent()
    sub_agents_list = [categorizer_agent_]

    super().__init__(
      name="organizer",
      description="this agent creates the organization plan",
      categorizer=categorizer_agent_,
      sub_agents=sub_agents_list,
    )

  @override
  async def _run_async_impl(self, ctx: InvocationContext) -> AsyncGenerator[Event, None]:
    # the caller should put files to state
    if "file" not in ctx.session.state:
      # to allow run with adk web, parse files from user_content.
      if ctx.user_content and ctx.user_content.parts and ctx.user_content.parts[0].text:
        req = PlanRequest.model_validate_json(ctx.user_content.parts[0].text)
        ctx.session.state["files"] = req.files

    async for event in self.categorizer.run_async(ctx):
      yield event

    cat = CategoryResponse.model_validate(ctx.session.state["category"])
    if cat.category not in category_list:
      raise Exception(f"Unknown category: {cat.category}")

    if cat.category in simple_move_categories:
      event = simple_move_plan_event(self.name, Category[cat.category], ctx.session.state["files"])
      yield event
      return

    return
