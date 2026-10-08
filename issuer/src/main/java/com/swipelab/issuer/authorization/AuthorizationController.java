package com.swipelab.issuer.authorization;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class AuthorizationController {
    @PostMapping("authorizations")
    public ResponseEntity<String> authorize(@RequestBody AuthorizationRequest request) {
        System.out.println("Authorization received:" + request.authorizationId());

        return ResponseEntity.ok("Authorization received by issuer");
    }

}