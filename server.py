#!/usr/bin/env python3
"""本机 Grok Bot 记忆管理。只监听 127.0.0.1。"""

import json
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from pathlib import Path
from urllib.parse import parse_qs, urlparse

import grok_client

ROOT = Path(__file__).resolve().parent
STATIC = ROOT / "static"
HOST = "127.0.0.1"
PORT = 8765


def session():
    return grok_client.load_session()


class Handler(BaseHTTPRequestHandler):
    def log_message(self, fmt, *args):
        print(f"[memory] {self.address_string()} {fmt % args}")

    def _send(self, status, payload=None, content_type="application/json; charset=utf-8", body=None):
        if body is None:
            body = json.dumps(payload, ensure_ascii=False).encode()
        self.send_response(status)
        self.send_header("Content-Type", content_type)
        self.send_header("Content-Length", str(len(body)))
        self.send_header("Cache-Control", "no-store")
        self.end_headers()
        self.wfile.write(body)

    def _read_json(self):
        length = int(self.headers.get("Content-Length") or 0)
        if length <= 0:
            return {}
        return json.loads(self.rfile.read(length).decode())

    def _error(self, status, message):
        self._send(status, {"error": message})

    def do_GET(self):
        parsed = urlparse(self.path)
        if parsed.path == "/":
            page = (STATIC / "index.html").read_bytes()
            self._send(200, content_type="text/html; charset=utf-8", body=page)
            return
        if parsed.path == "/api/bots":
            try:
                agents = grok_client.list_agents(session())
            except grok_client.GrokError as error:
                self._error(400, str(error))
                return
            self._send(200, {"bots": agents})
            return
        if parsed.path == "/api/memories":
            agent_id = (parse_qs(parsed.query).get("agentId") or [""])[0]
            if not agent_id:
                self._error(400, "缺少 agentId")
                return
            try:
                result = grok_client.pull_memories(session(), agent_id)
            except grok_client.GrokError as error:
                self._error(400, str(error))
                return
            self._send(200, result)
            return
        self._error(404, "没有这个地址")

    def do_POST(self):
        parsed = urlparse(self.path)
        try:
            payload = self._read_json()
        except json.JSONDecodeError:
            self._error(400, "请求不是 JSON")
            return
        agent_id = payload.get("agentId") or ""
        if parsed.path == "/api/memories":
            if not agent_id:
                self._error(400, "缺少 agentId")
                return
            try:
                result = grok_client.write_memory(
                    session(), agent_id, payload.get("text") or "", payload.get("tier") or "log"
                )
            except grok_client.GrokError as error:
                self._error(400, str(error))
                return
            self._send(200, result)
            return
        if parsed.path == "/api/memories/forget":
            if not agent_id:
                self._error(400, "缺少 agentId")
                return
            try:
                result = grok_client.forget_memory(session(), agent_id, payload.get("text") or "")
            except grok_client.GrokError as error:
                self._error(400, str(error))
                return
            self._send(200, result)
            return
        self._error(404, "没有这个地址")


def main():
    server = ThreadingHTTPServer((HOST, PORT), Handler)
    print(f"记忆管理打开 http://{HOST}:{PORT}")
    server.serve_forever()


if __name__ == "__main__":
    main()
