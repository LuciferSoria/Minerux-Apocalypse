package dev.minerux.guard;

import java.util.UUID;

/**
 * Server-side connectivity preflight (not a NeoForge mod).
 *
 * Usage: MINERUX_API_URL=https://... MINERUX_SERVER_TOKEN=... \
 *        java -cp out dev.minerux.guard.ServerGuardPreflight <verified-player-uuid>
 *
 * Exit codes: 0=authorized; 2=denied/backend unavailable; 3=invalid configuration.
 * Never logs secrets or the complete backend URL.
 */
public final class ServerGuardPreflight {
    public static void main(String[] args) {
        if (args.length != 1) {
            System.err.println("Usage: ServerGuardPreflight <verified-player-uuid>");
            System.exit(3);
        }
        final UUID uuid;
        try {
            uuid = UUID.fromString(args[0]);
            if (!uuid.toString().equalsIgnoreCase(args[0])) {
                throw new IllegalArgumentException("UUID must be in canonical format");
            }
        } catch (IllegalArgumentException e) {
            System.err.println("Invalid player UUID");
            System.exit(3);
            return;
        }
        final BetaAccessClient client;
        try {
            client = ServerGuardConfig.fromEnvironment().client();
        } catch (RuntimeException e) {
            System.err.println("Invalid server guard configuration");
            System.exit(3);
            return;
        }
        if (client.isAllowed(uuid)) {
            System.out.println("Access granted");
            return;
        }
        System.err.println("Access denied (or backend unavailable)");
        System.exit(2);
    }
}
