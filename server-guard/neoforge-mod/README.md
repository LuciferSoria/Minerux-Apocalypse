# NeoForge mod packaging — incomplete

This is a **Java 25 / Minecraft 26.1.2 / NeoForge 26.1.2.125** Gradle mod project. The source now contains a **pre-join PlayerNegotiationEvent gate** that waits for an asynchronous backend authorization decision and disconnects denied players. **Compilation and behavior on NeoForge 26.1.2.125 are not yet verified; do not deploy as production security.**

Run with a JDK 25 installation and Gradle 9.1+:

```sh
cd server-guard/neoforge-mod
gradle build
```

Do not deploy the resulting JAR as an access control solution until its build succeeds, denied/authorized/timeout scenarios are tested on a real NeoForge 26.1.2.125 server, and Minecraft `online-mode=true` is verified.

The framework-neutral Java adapter is reused from `../src/main/java`. This avoids duplicating authorization logic. **Never** package `MINERUX_SERVER_TOKEN` in the JAR; inject it only into the server process environment.

The standalone Java 21 CI does not validate this Gradle mod build. Java 25 is required for Minecraft 26.1.

## Lifecycle entrypoint

`MineruxBetaGuardMod` registers NeoForge server starting/stopped listeners. It initializes the bounded authorization service from environment variables on server startup and closes it on server shutdown. **It now subscribes to PlayerNegotiationEvent, enqueues the authorization future, and disconnects denied connections**. The startup listener throws if mandatory configuration is absent, but this behavior has not been tested in an actual NeoForge server.

## Required server acceptance tests

1. Ensure `online-mode=true` in `server.properties` and inject the backend token only into the server process.
2. Test a verified UUID without a beta grant: the connection must terminate **before world join**.
3. Grant the same UUID: the player must join successfully.
4. Revoke the UUID and reconnect: the connection must be denied.
5. Disable the backend or supply an invalid token: the connection must fail closed.
6. Saturate the worker queue: excess logins must fail closed, not hang indefinitely.
7. Inspect server logs and built JAR: no private token or beta codes may be included.

**Note:** the current gate is designed around the documented NeoForge negotiation event, but the precise runtime disconnect behavior and Gradle build have not been verified in this repository yet.

## CI build artifact

The experimental workflow `.github/workflows/neoforge-mod-build.yml` runs with Java 25 and Gradle 9.2.1. On a **successful** build, GitHub Actions publishes `minerux-neoforge-guard-development` as a workflow artifact; this is **not** a signed release or a production-ready server plugin. Build success, loader compatibility, and actual pre-join denial must all be verified independently.

In GitHub: **Actions → NeoForge mod build (experimental) → Run workflow**. Inspect the build log and artifact before considering a deployment. Do not distribute any artifact that has not passed the acceptance tests above.

## JAR packaging verification

Before the workflow uploads the development JAR, `verify_jar.py` checks that the mod entrypoint, authorization classes, and NeoForge metadata are packaged and rejects a small set of known private-file names. The workflow fails if no JAR or more than one JAR is produced.

This is a **structural** check only. It cannot detect every possible secret, confirm that the Minecraft loader accepts the mod, or establish that denied players are disconnected before world join. Those require build logs, secret review, and live integration tests.

## Read-only server deployment preflight

Before exposing the beta server, set `online-mode=true`, `white-list=true`, and `enforce-whitelist=true` in `server.properties`. Add only approved beta users to the vanilla Minecraft whitelist as a temporary, independently enforced barrier.

With `MINERUX_API_URL` and `MINERUX_SERVER_TOKEN` supplied to the server process, run:

```bash
python3 server-guard/neoforge-mod/server_preflight.py /path/to/server.properties
```

The script only reads configuration, does not print token values, and does not contact the backend. A passing result **does not** establish that the token is valid, that the NeoForge mod builds, or that the negotiation gate denies unauthorized players. The vanilla whitelist and backend beta grants must both be maintained until the full gate is verified.

### Ambiguous properties are rejected

The deployment preflight also rejects duplicate keys in `server.properties` (including duplicate `online-mode` and `white-list` settings). Remove duplicates rather than relying on a particular parser's last-value behavior. The check is conservative and does not replace a live server test.
