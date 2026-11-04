package dev.minerux.guard;

import java.util.Objects;
import java.util.UUID;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.atomic.AtomicBoolean;

/**
 * Owns the worker pool and the asynchronous beta authorization gate.
 *
 * Lifecycle boundary for a future NeoForge integration. The Minecraft login
 * handler must reject or hold connections until the returned future resolves.
 * A true pre-join gate has NOT been implemented here.
 */
public final class AuthorizationService implements AutoCloseable {
    private final LoginCheckWorkers workers;
    private final LoginAccessGate gate;
    private final AtomicBoolean closed = new AtomicBoolean();

    public AuthorizationService(ServerGuardConfig config) {
        Objects.requireNonNull(config, "config");
        this.workers = config.workers();
        this.gate = new LoginAccessGate(config.client(), workers.executor());
    }

    /** Always fail closed if this service has already shut down. */
    public CompletableFuture<Boolean> check(UUID verifiedPlayerUuid) {
        if (closed.get() || verifiedPlayerUuid == null) {
            return CompletableFuture.completedFuture(false);
        }
        return gate.check(verifiedPlayerUuid).thenApply(allowed -> !closed.get() && allowed);
    }

    @Override
    public void close() {
        if (closed.compareAndSet(false, true)) {
            workers.close();
        }
    }
}
