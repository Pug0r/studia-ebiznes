import os
import json
import socket
import threading
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from urllib.error import HTTPError
from urllib.parse import urlparse
from urllib.request import Request, urlopen

import pytest
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.chrome.service import Service


@pytest.fixture
def base_url() -> str:
    return os.getenv("APP_URL", "http://localhost:5173")


@pytest.fixture
def wait_timeout() -> int:
    return int(os.getenv("WAIT_TIMEOUT", "20"))


def _pick_free_port() -> int:
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.bind(("127.0.0.1", 0))
    _, port = sock.getsockname()
    sock.close()
    return port


def _create_proxy_handler(upstream_base_url: str):
    upstream = urlparse(upstream_base_url)
    host_header = "localhost:5173"

    class ProxyHandler(BaseHTTPRequestHandler):
        def _forward(self) -> None:
            target_url = f"{upstream_base_url}{self.path}"
            req = Request(target_url, method=self.command)
            req.add_header("Host", host_header)

            for header_name, header_value in self.headers.items():
                header_name_lower = header_name.lower()
                if header_name_lower in {"host", "connection", "content-length"}:
                    continue
                req.add_header(header_name, header_value)

            body = None
            content_length = self.headers.get("Content-Length")
            if content_length:
                body = self.rfile.read(int(content_length))

            try:
                with urlopen(req, data=body, timeout=20) as response:
                    response_body = response.read()
                    self.send_response(response.status)
                    for key, value in response.headers.items():
                        key_lower = key.lower()
                        if key_lower in {"transfer-encoding", "connection", "content-length"}:
                            continue
                        self.send_header(key, value)
                    self.send_header("Content-Length", str(len(response_body)))
                    self.end_headers()
                    self.wfile.write(response_body)
            except HTTPError as error:
                response_body = error.read()
                self.send_response(error.code)
                self.send_header("Content-Length", str(len(response_body)))
                self.end_headers()
                self.wfile.write(response_body)

        def do_GET(self) -> None:
            self._forward()

        def do_POST(self) -> None:
            self._forward()

        def do_PUT(self) -> None:
            self._forward()

        def do_DELETE(self) -> None:
            self._forward()

        def do_OPTIONS(self) -> None:
            self._forward()

        def log_message(self, _format: str, *_args) -> None:
            return

    return ProxyHandler


@pytest.fixture
def effective_base_url(base_url: str):
    parsed = urlparse(base_url)
    if parsed.hostname != "host.docker.internal":
        yield base_url
        return

    proxy_port = _pick_free_port()
    proxy_server = ThreadingHTTPServer(("127.0.0.1", proxy_port), _create_proxy_handler(base_url))
    server_thread = threading.Thread(target=proxy_server.serve_forever, daemon=True)
    server_thread.start()

    try:
        yield f"http://127.0.0.1:{proxy_port}"
    finally:
        proxy_server.shutdown()
        proxy_server.server_close()


@pytest.fixture
def driver(effective_base_url: str):
    remote_url = os.getenv("SELENIUM_REMOTE_URL")

    if remote_url:
        options = Options()
        options.set_capability("browserName", os.getenv("BROWSER_NAME", "chrome"))
        options.set_capability("browserVersion", os.getenv("BROWSER_VERSION", "latest"))
        bstack_options_raw = os.getenv("BROWSERSTACK_OPTIONS")
        if bstack_options_raw:
            options.set_capability("bstack:options", json.loads(bstack_options_raw))
        driver_instance = webdriver.Remote(command_executor=remote_url, options=options)
    else:
        options = Options()
        chrome_binary = os.getenv("CHROME_BINARY")
        if chrome_binary:
            options.binary_location = chrome_binary
        if os.getenv("HEADLESS", "true").lower() == "true":
            options.add_argument("--headless=new")
        options.add_argument("--no-sandbox")
        options.add_argument("--disable-dev-shm-usage")
        options.add_argument("--disable-gpu")
        options.add_argument("--remote-debugging-port=9222")
        options.add_argument("--window-size=1366,900")

        chromedriver_path = os.getenv("CHROMEDRIVER_PATH")
        if chromedriver_path:
            driver_instance = webdriver.Chrome(service=Service(chromedriver_path), options=options)
        else:
            driver_instance = webdriver.Chrome(options=options)

    driver_instance.get(effective_base_url)
    yield driver_instance
    driver_instance.quit()
