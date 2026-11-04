package dev.minerux.guard;

import java.util.UUID;
import java.util.concurrent.TimeUnit;

/** Deterministic lifecycle checks without a live backend. */
public final class AuthorizationServiceTest {
    public static void main(String[] args) throws Exception {
        ServerGuardConfig config = new ServerGuardConfig(
                "https://api.example.org",
                "0123456789abcdef0123456789abcdef",
                1, 1);
        AuthorizationService service = new AuthorizationService(config);
        if (service.check(null).get(1, TimeUnit.SECONDS)) {
            throw new AssertionError("Missing UUID must be denied");
        }
        service.close();
        service.close(); // Closing twice must be safe.
        UUID uuid = UUID.fromString("01234567-89ab-cdef-0123-456789abcdef");
        if (service.check(uuid).get(1, TimeUnit.SECONDS)) {
            throw new AssertionError("Closed authorization service must deny access");
        }
        System.out.println("AuthorizationService lifecycle tests passed");
    }
}
