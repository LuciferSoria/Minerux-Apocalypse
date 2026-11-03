package dev.minerux.guard;

import java.util.UUID;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.Executor;

/**
 * Integration boundary for NeoForge login handlers.
 *
 * The caller must obtain the UUID from a fully authenticated online-mode
 * Minecraft profile. Do not use usernames, client-provided UUIDs or offline-mode UUIDs.
 *
 * IMPORTANT: This is a framework-neutral helper, not a registered NeoForge event.
 */
public final class LoginAccessGate {
    private final BetaAccessClient client;
    private final Executor worker;

    public LoginAccessGate(BetaAccessClient client, Executor worker) {
        this.client = java.util.Objects.requireNonNull(client);
        this.worker = java.util.Objects.requireNonNull(worker);
    }

    /**
     * Always completes with false when authorization fails. Do not call join()
     * or get() on the Minecraft main server thread.
     */
    public CompletableFuture<Boolean> check(UUID authenticatedUuid) {
        if (authenticatedUuid == null) return CompletableFuture.completedFuture(false);
        try {
            return CompletableFuture.supplyAsync(
                    () -> client.isAllowed(authenticatedUuid), worker)
                    .exceptionally(error -> false);
        } catch (RuntimeException error) {
            return CompletableFuture.completedFuture(false);
        }
    }
}
