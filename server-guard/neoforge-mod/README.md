# NeoForge mod packaging — incomplete

This is a **Java 25 / Minecraft 26.1.2 / NeoForge 26.1.2.125** Gradle project scaffold. It is **not** a functional beta-enforcement mod yet: an `@Mod` entrypoint now exists, but no pre-login hook or verified server join rejection has been added.

Run with a JDK 25 installation and Gradle 9.1+:

```sh
cd server-guard/neoforge-mod
gradle build
```

Do not deploy the resulting JAR as an access control solution until the NeoForge lifecycle is wired, a true pre-join denial is implemented and tested, and Minecraft `online-mode=true` is verified.

The framework-neutral Java adapter is reused from `../src/main/java`. This avoids duplicating authorization logic. **Never** package `MINERUX_SERVER_TOKEN` in the JAR; inject it only into the server process environment.

The standalone Java 21 CI does not validate this Gradle mod build. Java 25 is required for Minecraft 26.1.

## Lifecycle entrypoint

`MineruxBetaGuardMod` registers NeoForge server starting/stopped listeners. It initializes the bounded authorization service from environment variables on server startup and closes it on server shutdown. **It does not yet subscribe to a login event or deny players**. The startup listener throws if mandatory configuration is absent, but this behavior has not been tested in an actual NeoForge server.
