package dev.minerux.guard;

import java.util.concurrent.CountDownLatch;
import java.util.concurrent.RejectedExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;

public final class LoginCheckWorkersTest {
    public static void main(String[] args) throws Exception {
        for (int[] limits : new int[][] {{0,1},{33,1},{1,0},{1,4097}}) {
            try {
                new LoginCheckWorkers(limits[0], limits[1]);
                throw new AssertionError("Unsafe limits accepted");
            } catch (IllegalArgumentException expected) {
                // Correctly rejected.
            }
        }
        CountDownLatch started = new CountDownLatch(1);
        CountDownLatch release = new CountDownLatch(1);
        AtomicInteger finished = new AtomicInteger();
        try (LoginCheckWorkers workers = new LoginCheckWorkers(1,1)) {
            workers.executor().execute(() -> {
                started.countDown();
                try { release.await(3, TimeUnit.SECONDS); }
                catch (InterruptedException e) { Thread.currentThread().interrupt(); }
            });
            if (!started.await(2, TimeUnit.SECONDS)) throw new AssertionError("Worker did not start");
            workers.executor().execute(finished::incrementAndGet);
            try {
                workers.executor().execute(finished::incrementAndGet);
                throw new AssertionError("Full queue accepted another task");
            } catch (RejectedExecutionException expected) {
                // Overload must not create unlimited threads or queue entries.
            }
            release.countDown();
            long deadline = System.nanoTime() + TimeUnit.SECONDS.toNanos(3);
            while (finished.get() == 0 && System.nanoTime() < deadline) Thread.sleep(10);
            if (finished.get() != 1) throw new AssertionError("Queued task not executed");
        } finally {
            release.countDown();
        }
        System.out.println("LoginCheckWorkers tests passed");
    }
}
