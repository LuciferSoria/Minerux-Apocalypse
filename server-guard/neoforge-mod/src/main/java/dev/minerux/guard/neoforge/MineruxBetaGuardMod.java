package dev.minerux.guard.neoforge;

import dev.minerux.guard.AuthorizationService;
import dev.minerux.guard.ServerGuardConfig;
import net.neoforged.bus.api.SubscribeEvent;
import net.neoforged.fml.common.Mod;
import net.neoforged.neoforge.common.NeoForge;
import net.neoforged.neoforge.event.server.ServerStartingEvent;
import net.neoforged.neoforge.event.server.ServerStoppedEvent;

/**
 * Server-only NeoForge lifecycle entrypoint.
 *
 * This is not a pre-login enforcement hook. Never advertise the server as
 * beta-protected until a verified pre-join denial mechanism is installed.
 */
@Mod("mineruxbetaguard")
public final class MineruxBetaGuardMod {
    private AuthorizationService authorization;

    public MineruxBetaGuardMod() {
        NeoForge.EVENT_BUS.register(this);
    }

    @SubscribeEvent
    public void onServerStarting(ServerStartingEvent event) {
        // A missing token or an invalid URL must prevent an apparently healthy
        // beta server from starting without its authorization backend.
        authorization = new AuthorizationService(ServerGuardConfig.fromEnvironment());
        System.err.println("[Minerux Beta Guard] Backend configuration validated. "
                + "WARNING: player login enforcement is NOT installed; do not open the beta.");
    }

    @SubscribeEvent
    public void onServerStopped(ServerStoppedEvent event) {
        if (authorization != null) {
            authorization.close();
            authorization = null;
        }
    }
}
