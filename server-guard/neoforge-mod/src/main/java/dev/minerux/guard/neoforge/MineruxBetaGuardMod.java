package dev.minerux.guard.neoforge;

import dev.minerux.guard.AuthorizationService;
import dev.minerux.guard.ServerGuardConfig;
import net.neoforged.bus.api.SubscribeEvent;
import net.neoforged.fml.common.Mod;
import net.neoforged.neoforge.common.NeoForge;
import net.neoforged.neoforge.event.server.ServerStartingEvent;
import net.neoforged.neoforge.event.entity.player.PlayerNegotiationEvent;
import net.minecraft.network.chat.Component;
import net.neoforged.neoforge.event.server.ServerStoppedEvent;

/**
 * Server-only NeoForge lifecycle entrypoint.
 *
 * The PlayerNegotiationEvent is a pre-join handshake barrier. Verify the
 * disconnect behavior on the exact NeoForge server build before deployment.
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
        System.err.println("[Minerux Beta Guard] Backend configuration validated; negotiation access gate registered.");
    }

    /**
     * Negotiation is paused until the backend check finishes. Deny on any
     * missing profile, unavailable backend, invalid token, or denied grant.
     * Requires Minecraft online-mode=true for authenticated UUIDs.
     */
    @SubscribeEvent
    public void onPlayerNegotiation(PlayerNegotiationEvent event) {
        var profile = event.getProfile();
        var connection = event.getConnection();
        if (authorization == null || profile == null || profile.getId() == null) {
            connection.disconnect(Component.literal("Minerux: acceso beta no disponible."));
            return;
        }
        event.enqueueWork(authorization.check(profile.getId()).thenAccept(allowed -> {
            if (!allowed) {
                connection.disconnect(Component.literal("Minerux: acceso beta denegado."));
            }
        }));
    }

    @SubscribeEvent
    public void onServerStopped(ServerStoppedEvent event) {
        if (authorization != null) {
            authorization.close();
            authorization = null;
        }
    }
}
