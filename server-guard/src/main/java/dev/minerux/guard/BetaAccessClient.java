package dev.minerux.guard;

import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.time.Duration;
import java.util.Locale;
import java.util.UUID;
import java.util.regex.Pattern;

/**
 * Private server-side adapter for the Minerux beta backend.
 * The caller MUST supply a UUID verified by Minecraft online-mode login.
 * Never put MINERUX_SERVER_TOKEN in a launcher or client-side mod.
 */
public final class BetaAccessClient {
    private static final Pattern RESPONSE = Pattern.compile(
            "^\\{\\s*\\\"allowed\\\"\\s*:\\s*(true|false)\\s*,\\s*\\\"status\\\"\\s*:\\s*\\\"(online|beta|maintenance)\\\"\\s*\\}$");
    private final HttpClient http;
    private final URI base;
    private final String token;

    public BetaAccessClient(String baseUrl, String serverToken) {
        if (baseUrl == null || baseUrl.isBlank()) {
            throw new IllegalArgumentException("HTTPS backend URL required");
        }
        URI u = URI.create(baseUrl);
        if (!"https".equalsIgnoreCase(u.getScheme()) || u.getHost() == null
                || u.getUserInfo() != null || u.getQuery() != null || u.getFragment() != null
                || (u.getRawPath() != null && !u.getRawPath().matches("/?"))
                || u.getPort() == 0) {
            throw new IllegalArgumentException("HTTPS backend URL required");
        }
        if (serverToken == null || serverToken.length() < 32 || serverToken.length() > 4096
                || !serverToken.matches("[A-Za-z0-9._~+/-]+")) {
            throw new IllegalArgumentException("server token must be 32-4096 safe ASCII characters");
        }
        this.base = URI.create(u.getScheme() + "://" + u.getRawAuthority() + "/");
        this.token = serverToken;
        this.http = HttpClient.newBuilder()
                .connectTimeout(Duration.ofSeconds(4))
                .followRedirects(HttpClient.Redirect.NEVER)
                .build();
    }

    /** Fail closed on all errors. Call asynchronously, never on the main server tick thread. */
    public boolean isAllowed(UUID verifiedUuid) {
        if (verifiedUuid == null) return false;
        String uuid = verifiedUuid.toString().replace("-", "").toLowerCase(Locale.ROOT);
        URI endpoint = base.resolve("api/v1/internal/player-access/" + uuid);
        try {
            HttpRequest request = HttpRequest.newBuilder(endpoint)
                    .header("Authorization", "Bearer " + token)
                    .timeout(Duration.ofSeconds(6))
                    .GET().build();
            HttpResponse<java.io.InputStream> response = http.send(request, HttpResponse.BodyHandlers.ofInputStream());
            try (java.io.InputStream stream = response.body()) {
                if (response.statusCode() != 200) return false;
                // Bound response consumption before decoding: a compromised backend must not exhaust server memory.
                byte[] bytes = stream.readNBytes(4097);
                if (bytes.length > 4096) return false;
                return allowedFromJson(new String(bytes, java.nio.charset.StandardCharsets.UTF_8));
            }
        } catch (IOException | InterruptedException | IllegalArgumentException e) {
            if (e instanceof InterruptedException) Thread.currentThread().interrupt();
            return false;
        }
    }

    static boolean allowedFromJson(String json) {
        if (json == null || json.length() > 4096) return false;
        var matcher = RESPONSE.matcher(json.trim());
        return matcher.matches() && "true".equals(matcher.group(1))
                && !"maintenance".equals(matcher.group(2));
    }
}
