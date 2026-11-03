package dev.minerux.guard;

import java.util.Objects;

/**
 * Server-only configuration. Requires environment variables at startup.
 * Do not print the token or serialize this object into logs.
 */
public final class ServerGuardConfig {
    private final String apiUrl;
    private final String serverToken;
    private final int workers;
    private final int queueCapacity;

    public ServerGuardConfig(String apiUrl, String serverToken, int workers, int queueCapacity) {
        this.apiUrl = Objects.requireNonNull(apiUrl, "apiUrl");
        this.serverToken = Objects.requireNonNull(serverToken, "serverToken");
        if (workers < 1 || workers > 32 || queueCapacity < 1 || queueCapacity > 4096) {
            throw new IllegalArgumentException("Invalid worker limits");
        }
        // Validate credentials and URL immediately; do not contact the network.
        new BetaAccessClient(apiUrl, serverToken);
        this.workers = workers;
        this.queueCapacity = queueCapacity;
    }

    public static ServerGuardConfig fromEnvironment() {
        return new ServerGuardConfig(
                required("MINERUX_API_URL"),
                required("MINERUX_SERVER_TOKEN"),
                positiveInt("MINERUX_AUTH_WORKERS", 4),
                positiveInt("MINERUX_AUTH_QUEUE", 64));
    }

    private static String required(String name) {
        String value = System.getenv(name);
        if (value == null || value.isBlank()) {
            throw new IllegalStateException("Missing required server configuration: " + name);
        }
        return value;
    }

    private static int positiveInt(String name, int fallback) {
        String value = System.getenv(name);
        if (value == null || value.isBlank()) return fallback;
        try { return Integer.parseInt(value); }
        catch (NumberFormatException ex) {
            throw new IllegalStateException("Invalid server configuration: " + name);
        }
    }

    public BetaAccessClient client() { return new BetaAccessClient(apiUrl, serverToken); }
    public LoginCheckWorkers workers() { return new LoginCheckWorkers(workers, queueCapacity); }

    @Override public String toString() {
        return "ServerGuardConfig{apiUrl=[redacted],serverToken=[redacted],workers="
                + workers + ",queueCapacity=" + queueCapacity + "}";
    }
}
