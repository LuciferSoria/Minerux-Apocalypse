package dev.minerux.guard;

public final class ServerGuardConfigTest {
    public static void main(String[] args) {
        String token = "0123456789abcdef0123456789abcdef";
        ServerGuardConfig config = new ServerGuardConfig("https://api.example.org", token, 4, 64);
        if (config.toString().contains(token) || config.toString().contains("api.example.org")) {
            throw new AssertionError("Configuration must not reveal credentials or backend address");
        }
        try (LoginCheckWorkers workers = config.workers()) {
            if (workers.executor() == null || config.client() == null) {
                throw new AssertionError("Configuration did not create required components");
            }
        }
        reject("http://api.example.org", token, 4, 64);
        reject("https://api.example.org/unexpected/path", token, 4, 64);
        reject("https://api.example.org/api/", token, 4, 64);
        reject("https://api.example.org:0", token, 4, 64);
        reject("https://api.example.org", "short", 4, 64);
        reject("https://api.example.org", token, 0, 64);
        reject("https://api.example.org", token, 33, 64);
        reject("https://api.example.org", token, 4, 0);
        reject("https://api.example.org", token, 4, 4097);
        System.out.println("ServerGuardConfig tests passed");
    }

    private static void reject(String url, String token, int workers, int queue) {
        try {
            new ServerGuardConfig(url, token, workers, queue);
            throw new AssertionError("Unsafe configuration accepted");
        } catch (IllegalArgumentException expected) {
            // Expected.
        }
    }
}
