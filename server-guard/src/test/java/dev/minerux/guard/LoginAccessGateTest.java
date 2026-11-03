package dev.minerux.guard;

import java.net.URI;
import java.util.UUID;
import java.util.concurrent.Executor;
import java.util.concurrent.TimeUnit;

/**
 * Executable smoke tests for the asynchronous login integration boundary.
 * Does not contact the backend or require NeoForge.
 */
public final class LoginAccessGateTest {
    private static final String TOKEN = "0123456789abcdef0123456789abcdef";
    public static void main(String[] args) throws Exception {
        BetaAccessClient client = new BetaAccessClient("https://127.0.0.1:1", TOKEN);
        Executor inline = Runnable::run;
        LoginAccessGate gate = new LoginAccessGate(client, inline);
        if (gate.check(null).get(1, TimeUnit.SECONDS)) {
            throw new AssertionError("Missing authenticated UUID must be denied");
        }
        Executor rejected = task -> { throw new java.util.concurrent.RejectedExecutionException("offline"); };
        LoginAccessGate unavailable = new LoginAccessGate(client, rejected);
        if (unavailable.check(UUID.randomUUID()).get(1, TimeUnit.SECONDS)) {
            throw new AssertionError("Executor rejection must deny login");
        }
        if (gate.check(UUID.randomUUID()).get(10, TimeUnit.SECONDS)) {
            throw new AssertionError("Unreachable backend must deny login");
        }
        System.out.println("LoginAccessGate tests passed");
    }
}
