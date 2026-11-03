package dev.minerux.guard;

/** Thirty deterministic authorization response contract checks. */
public final class BetaAccessContractTest {
    public static void main(String[] args) {
        verify(1, "{\"allowed\":true,\"status\":\"beta\"}", true);
        verify(2, "{\"allowed\":true,\"status\":\"online\"}", true);
        verify(3, "{\"allowed\":false,\"status\":\"beta\"}", false);
        verify(4, "{\"allowed\":false,\"status\":\"online\"}", false);
        verify(5, "{\"allowed\":true,\"status\":\"maintenance\"}", false);
        verify(6, "{\"allowed\":false,\"status\":\"maintenance\"}", false);
        verify(7, "{\"allowed\":true}", false);
        verify(8, "{\"status\":\"beta\"}", false);
        verify(9, "{}", false);
        verify(10, "[]", false);
        verify(11, "null", false);
        verify(12, "true", false);
        verify(13, "", false);
        verify(14, "{\"allowed\":true,\"status\":\"unknown\"}", false);
        verify(15, "{\"allowed\":true,\"status\":\"BETA\"}", false);
        verify(16, "{\"allowed\":\"true\",\"status\":\"beta\"}", false);
        verify(17, "{\"allowed\":1,\"status\":\"beta\"}", false);
        verify(18, "{\"allowed\":null,\"status\":\"beta\"}", false);
        verify(19, "{\"status\":\"beta\",\"allowed\":true}", false);
        verify(20, "{\"allowed\":true,\"status\":\"beta\",\"allowed\":false}", false);
        verify(21, "{\"allowed\":true,\"status\":\"beta\",\"extra\":1}", false);
        verify(22, "{\"allowed\":true,\"status\":\"beta\"} trailing", false);
        verify(23, "prefix {\"allowed\":true,\"status\":\"beta\"}", false);
        verify(24, "{\"allowed\":true,\"status\":\"beta\"", false);
        verify(25, "{\"allowed\":true,\"status\":\"beta\"}}", false);
        verify(26, "{\"allowed\":true,\"status\":\"beta\"}{\"allowed\":true,\"status\":\"beta\"}", false);
        verify(27, "{\"allowed\":true,\"status\":\"beta\"}\u0000", false);
        verify(28, "{\"allowed\":true,\"status\":\"beta\"}\n", true);
        verify(29, " { \"allowed\" : true , \"status\" : \"beta\" } ", true);
        verify(30, "{\"allowed\":true,\"status\":\"beta\"}                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    ", false);
        System.out.println("30 authorization contract checks passed");
    }
    private static void verify(int id, String response, boolean expected) {
        boolean actual = BetaAccessClient.allowedFromJson(response);
        if (actual != expected) {
            throw new AssertionError("Contract case " + id + " expected " + expected + ", got " + actual);
        }
    }
}
