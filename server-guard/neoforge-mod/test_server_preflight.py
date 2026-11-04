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

    def test_duplicate_online_mode_denied(self):
        errors = preflight.inspect(PROPERTIES + "online-mode=false\n", ENV)
        self.assertTrue(any("duplicate setting: online-mode" in e for e in errors))

    def test_duplicate_whitelist_denied(self):
        errors = preflight.inspect(PROPERTIES + "white-list=true\n", ENV)
        self.assertTrue(any("duplicate setting: white-list" in e for e in errors))

    def test_noncanonical_true_denied(self):
        for setting in ("online-mode", "white-list", "enforce-whitelist"):
            with self.subTest(setting=setting):
                props = PROPERTIES.replace(setting + "=true", setting + "=TRUE")
                self.assertTrue(preflight.inspect(props, ENV))

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

    def test_out_of_range_port_denied(self):
        self.assertTrue(preflight.inspect(PROPERTIES, {**ENV, "MINERUX_API_URL": "https://beta.example.org:99999"}))

    def test_invalid_port_denied(self):
        self.assertTrue(preflight.inspect(PROPERTIES, {**ENV, "MINERUX_API_URL": "https://beta.example.org:bad"}))

if __name__ == "__main__":
    unittest.main()
