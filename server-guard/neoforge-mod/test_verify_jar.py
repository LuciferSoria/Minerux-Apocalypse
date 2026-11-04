import importlib.util
import tempfile
import unittest
import warnings
import zipfile
from pathlib import Path

spec = importlib.util.spec_from_file_location("verify_jar", Path(__file__).with_name("verify_jar.py"))
verifier = importlib.util.module_from_spec(spec)
spec.loader.exec_module(verifier)

class VerifyJarTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.jar = Path(self.temp.name) / "test.jar"

    def build(self, omit=(), metadata=b'modId="mineruxbetaguard"', extra=None, duplicate=False):
        with zipfile.ZipFile(self.jar, "w") as archive:
            for name in sorted(verifier.REQUIRED - set(omit)):
                archive.writestr(name, metadata if name.endswith(".toml") else b"class")
            if extra:
                archive.writestr(extra, b"dummy")
            if duplicate:
                with warnings.catch_warnings():
                    warnings.simplefilter("ignore", UserWarning)
                    archive.writestr(sorted(verifier.REQUIRED)[0], b"duplicate")

    def test_valid(self):
        self.build()
        verifier.check(self.jar)

    def test_missing_entrypoint(self):
        self.build(omit={"dev/minerux/guard/neoforge/MineruxBetaGuardMod.class"})
        with self.assertRaisesRegex(ValueError, "Missing"):
            verifier.check(self.jar)

    def test_missing_client(self):
        self.build(omit={"dev/minerux/guard/BetaAccessClient.class"})
        with self.assertRaisesRegex(ValueError, "Missing"):
            verifier.check(self.jar)

    def test_wrong_mod_id(self):
        self.build(metadata=b'modId="other"')
        with self.assertRaisesRegex(ValueError, "Incorrect"):
            verifier.check(self.jar)

    def test_private_file(self):
        self.build(extra="private/.env")
        with self.assertRaisesRegex(ValueError, "Private"):
            verifier.check(self.jar)

    def test_duplicate(self):
        self.build(duplicate=True)
        with self.assertRaisesRegex(ValueError, "Duplicate"):
            verifier.check(self.jar)

    def test_invalid_zip(self):
        self.jar.write_bytes(b"not zip")
        with self.assertRaises(zipfile.BadZipFile):
            verifier.check(self.jar)

if __name__ == "__main__":
    unittest.main()
