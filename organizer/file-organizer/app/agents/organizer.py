from typing import override, AsyncGenerator

from google.adk.events import Event
from google.adk.agents import BaseAgent, LlmAgent, InvocationContext
from .categorizer import agent as categorizer_agent, CategoryResponse, Category


class OrganizerAgent(BaseAgent):
  categorizer: LlmAgent

  def __init__(self):
    sub_agents_list = [categorizer_agent]

    super().__init__(
      name="organizer",
      description="this agent creates the organization plan",
      categorizer=categorizer_agent,
      sub_agents=sub_agents_list,
    )

  @override
  async def _run_async_impl(self, ctx: InvocationContext) -> AsyncGenerator[Event, None]:
    async for event in self.categorizer.run_async(ctx):
      yield event

    cat = CategoryResponse.model_validate_json(ctx.session.state["category"])
    match cat.category:
      case Category.movie.name:
        pass
      case Category.tv_series.name:
        pass
      case Category.anim_tv_series.name:
        pass
      case Category.anim_movie.name:
        pass
      case Category.photobook.name:
        pass
      case Category.porn.name:
        pass
      case Category.audio_book.name:
        pass
      case Category.book.name:
        pass
      case Category.music.name:
        pass
      case Category.music_video.name:
        pass
      case _:
        raise Exception(f"Unknown category: {cat.category}")
