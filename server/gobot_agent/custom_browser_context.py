from browser_use.browser.browser import Browser
from browser_use.browser.context import BrowserContext, BrowserContextConfig, BrowserSession
from browser_use.browser.views import BrowserState
from browser_use.dom.views import DOMElementNode


class CustomerBrowserContext(BrowserContext):
    def __init__(self, browser: 'Browser', config: BrowserContextConfig = BrowserContextConfig()):
        super().__init__(browser, config)

    async def _initialize_session(self):
        """Initialize the browser session"""

        playwright_browser = await self.browser.get_playwright_browser()

        context = await self._create_context(playwright_browser)
        if len(context.pages) > 0:
            page = context.pages[-1]
            await page.bring_to_front()
        else:
            page = await context.new_page()

        # Instead of calling _update_state(), create an empty initial state
        initial_state = BrowserState(
            element_tree=DOMElementNode(
                tag_name='root',
                is_visible=True,
                parent=None,
                xpath='',
                attributes={},
                children=[],
            ),
            selector_map={},
            url=page.url,
            title=await page.title(),
            screenshot=None,
            tabs=[],
        )

        self.session = BrowserSession(
            context=context,
            current_page=page,
            cached_state=initial_state,
        )

        await self._add_new_page_listener(context)

        return self.session
