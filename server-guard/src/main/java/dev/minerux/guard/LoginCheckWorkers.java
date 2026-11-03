package dev.minerux.guard;

import java.util.concurrent.ArrayBlockingQueue;
import java.util.concurrent.ThreadPoolExecutor;
import java.util.concurrent.TimeUnit;

/**
 * Bounded worker pool for remote login checks. Rejected checks fail closed
 * through LoginAccessGate's exception handler.
 */
public final class LoginCheckWorkers implements AutoCloseable {
    private final ThreadPoolExecutor executor;

    public LoginCheckWorkers(int threads, int queueCapacity) {
        if (threads < 1 || threads > 32 || queueCapacity < 1 || queueCapacity > 4096) {
            throw new IllegalArgumentException("Invalid worker pool limits");
        }
        executor = new ThreadPoolExecutor(
                threads, threads, 0L, TimeUnit.MILLISECONDS,
                new ArrayBlockingQueue<>(queueCapacity),
                task -> {
                    Thread thread = new Thread(task, "minerux-beta-auth");
                    thread.setDaemon(true);
                    return thread;
                },
                new ThreadPoolExecutor.AbortPolicy());
    }

    public java.util.concurrent.Executor executor() {
        return executor;
    }

    @Override
    public void close() {
        executor.shutdownNow();
    }
}
