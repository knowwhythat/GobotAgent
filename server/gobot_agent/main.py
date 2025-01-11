import asyncio
import json
import logging
import os
import subprocess
import sys

import psutil
import requests
import websockets
from browser_use.browser.browser import Browser
from browser_use.browser.browser import BrowserConfig
from browser_use.browser.context import BrowserContextConfig
from langchain_openai import ChatOpenAI
from websockets.legacy.server import WebSocketServerProtocol

from gobot_agent.custom_agent import CustomAgent
from gobot_agent.custom_browser_context import CustomerBrowserContext
from gobot_agent.custom_prompts import CustomSystemPrompt
from gobot_agent.utils import get_browser_path


class AgentServer(object):
    def __init__(self, port: int, log_path: str):
        self.port = port
        self.message_id = ""
        self.websocket_conn: WebSocketServerProtocol | None = None
        self.logger = logging.getLogger()
        self.logger.setLevel(logging.INFO)
        formatter = logging.Formatter(
            "%(asctime)s - %(name)s - %(levelname)s - %(message)s"
        )
        console_handler = logging.StreamHandler(stream=sys.stdout)
        console_handler.setFormatter(formatter)
        th = logging.FileHandler(filename=log_path, encoding="utf-8")
        th.setFormatter(formatter)
        self.logger.addHandler(console_handler)
        self.logger.addHandler(th)

    async def start(self) -> None:
        os.environ["ANONYMIZED_TELEMETRY"] = "false"
        async with websockets.serve(self.handler, "localhost", self.port):
            await asyncio.Future()  # run forever

    async def handler(self, websocket: WebSocketServerProtocol):
        async for message in websocket:
            await self.process(message, websocket)

    async def process(self, message, websocket: WebSocketServerProtocol):
        self.logger.info(f"收到请求:{message}")
        client_message = json.loads(message)
        self.message_id = client_message["message_id"]
        self.websocket_conn = websocket
        try:
            if "method" not in client_message:
                raise Exception("No method in client message")
            method = getattr(self, client_message["method"])
            if method:
                await method(client_message["data"])
            else:
                raise Exception(f"{client_message['method']}方法未定义")
        except Exception as e:
            self.logger.exception(e)
            await websocket.send(
                json.dumps(
                    {
                        "message_id": client_message["message_id"],
                        "result": "error",
                        "reason": str(e),
                    }
                )
            )

    async def execute_task(self, data):
        # 初始化AI代理，并设置任务描述和LLM模型
        if "browser_path" in data and data["browser_path"]:
            browser_path = data["browser_path"]
        else:
            browser_path = get_browser_path()
        try:
            # Check if browser is already running
            response = requests.get("http://localhost:9222/json/version", timeout=2)
            if response.status_code == 200:
                pass
        except requests.ConnectionError:
            file_name = os.path.basename(browser_path)
            pids = psutil.pids()
            for pid in pids:
                try:
                    pro = psutil.Process(pid)
                    if file_name == pro.name() or pro.name() + ".exe" == file_name:
                        pro.kill()
                except:
                    pass
            subprocess.Popen(
                [
                    browser_path,
                    "--remote-debugging-port=9222",
                ],
                stdout=subprocess.DEVNULL,
                stderr=subprocess.DEVNULL,
            )
        browser = Browser(
            BrowserConfig(headless=False, chrome_instance_path=browser_path)
        )
        agent = CustomAgent(
            self.websocket_conn,
            browser_context=CustomerBrowserContext(
                browser=browser, config=BrowserContextConfig()
            ),
            task=data["task"],
            llm=ChatOpenAI(
                model_name=data["modelName"],
                openai_api_base=data["apiBase"],
                api_key=data["apiKey"],
            ),
            system_prompt_class=CustomSystemPrompt,
            use_vision=False,
        )
        # 运行AI代理并获取结果
        result = await agent.run()


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("参数错误", file=sys.stderr)
    else:
        inspect_server = AgentServer(int(sys.argv[1]), sys.argv[2])
        asyncio.run(inspect_server.start())
