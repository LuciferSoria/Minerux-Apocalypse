import importlib.util
import unittest
from pathlib import Path

spec = importlib.util.spec_from_file_location("server_preflight", Path(__file__).with_name("server_preflight.py"))
preflight = importlib.util.module_from_spec(spec)
spec.loader.exec_module(preflight)

PROPERTIES = "online-mode=true\nwhite-list=true\nenforce-whitelist=true\n"
ENV = {"MINERUX_API_URL": "https://beta.example.org", "MINERUX_SERVER_TOKEN": "x" * 40}

class ServerPreflightTests(unittest.TestCase):
    def test_valid(self):
        self.assertEqual(preflight.inspect(PROPERTIES, ENV), [])

    def test_offline_mode_denied(self):
        self.assertTrue(any("online-mode" in e for e in preflight.inspect(PROPERTIES.replace("online-mode=true", "online-mode=false"), ENV)))

    def test_missing_whitelist_denied(self):
        self.assertTrue(any("white-list" in e for e in preflight.inspect("online-mode=true\n", ENV)))

    def test_enforce_whitelist_denied(self):
        self.assertTrue(any("enforce-whitelist" in e for e in preflight.inspect(PROPERTIES.replace("enforce-whitelist=true", "enforce-whitelist=false"), ENV)))

    def test_http_backend_denied(self):
        self.assertTrue(preflight.inspect(PROPERTIES, {**ENV, "MINERUX_API_URL": "http://beta.example.org"}))

    def test_missing_token_denied(self):
        self.assertTrue(preflight.inspect(PROPERTIES, {**ENV, "MINERUX_SERVER_TOKEN": ""}))

    def test_invalid_port_denied(self):
        self.assertTrue(preflight.inspect(PROPERTIES, {**ENV, "MINERUX_API_URL": "https://beta.example.org:bad"}))

if __name__ == "__main__":
    unittest.main()
