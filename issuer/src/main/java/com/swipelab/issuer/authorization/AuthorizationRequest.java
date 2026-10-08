package com.swipelab.issuer.authorization;

import com.fasterxml.jackson.annotation.JsonProperty;

public record AuthorizationRequest(
                @JsonProperty("authorization_id") String authorizationId,
                @JsonProperty("merchant_id") String merchantId,
                @JsonProperty("card_number") String cardNumber,
                long amount,
                String currency) {
}
