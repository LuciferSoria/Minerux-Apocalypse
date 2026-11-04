package dev.minerux.guard;

public final class BetaAccessClientTest {
    public static void main(String[] args) {
        check("{\"allowed\":true,\"status\":\"beta\"}", true);
        check("{\"allowed\":true,\"status\":\"online\"}", true);
        check("{\"allowed\":false,\"status\":\"beta\"}", false);
        check("{\"allowed\":false,\"status\":\"maintenance\"}", false);
        check("{\"allowed\":true,\"status\":\"maintenance\"}", false);
        check("{\"allowed\":true}", false);
        check("{\"allowed\":true,\"status\":\"beta\",\"allowed\":false}", false);
        check("error: \"allowed\":true", false);
        check(null, false);
        check("{\"allowed\":true,\"status\":\"unknown\"}", false);
        check("{\"status\":\"beta\",\"allowed\":true}", false);
        check("{\"allowed\":true,\"status\":\"beta\"} trailing", false);
        invalidConfiguration();
        System.out.println("BetaAccessClient tests passed");
    }
    private static void invalidConfiguration() {
        reject("http://example.com", "12345678901234567890123456789012");
        reject("https://example.com", "short");
        reject("https://example.com", "a".repeat(40) + " ");
        reject("https://example.com", "a".repeat(4097));
        reject("https://user:password@example.com", "12345678901234567890123456789012");
    }
    private static void reject(String url, String token) {
        try {
            new BetaAccessClient(url, token);
            throw new AssertionError("Unsafe configuration accepted: " + url);
        } catch (IllegalArgumentException expected) {
            // Invalid configuration must be rejected before a request is made.
        }
    }
    private static void check(String json, boolean expected) {
        boolean actual = BetaAccessClient.allowedFromJson(json);
        if (actual != expected) throw new AssertionError("expected " + expected + " for " + json);
    }
}
