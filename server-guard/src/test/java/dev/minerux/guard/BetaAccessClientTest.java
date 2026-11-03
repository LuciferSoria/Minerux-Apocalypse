package dev.minerux.guard;

public final class BetaAccessClientTest {
    public static void main(String[] args) {
        check("{\"allowed\":true,\"status\":\"beta\"}", true);
        check("{\"allowed\":true,\"status\":\"online\"}", true);
        check("{\"allowed\":false,\"status\":\"beta\"}", false);
        check("{\"allowed\":false,\"status\":\"maintenance\"}", false);
        check("{\"allowed\":true,\"status\":\"maintenance\"}", false);
        check("{\"allowed\":true}", false);
        check("{\"allowed\":true,\"status\":\"beta\",\"allowed\":false}", false);
        check("error: \"allowed\":true", false);
        check(null, false);
        System.out.println("BetaAccessClient tests passed");
    }
    private static void check(String json, boolean expected) {
        boolean actual = BetaAccessClient.allowedFromJson(json);
        if (actual != expected) throw new AssertionError("expected " + expected + " for " + json);
    }
}
