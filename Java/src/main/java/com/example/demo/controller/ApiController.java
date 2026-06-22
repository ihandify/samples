package com.example.demo.controller;

import com.example.demo.dto.ScopedKeyRequest;
import com.example.demo.service.UpstreamService;
import jakarta.validation.Valid;
import org.springframework.core.io.ClassPathResource;
import org.springframework.core.io.Resource;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.Map;

@RestController
@CrossOrigin(origins = "*", allowedHeaders = "*", methods = {RequestMethod.GET, RequestMethod.POST})
public class ApiController {

    private final UpstreamService upstreamService;

    public ApiController(UpstreamService upstreamService) {
        this.upstreamService = upstreamService;
    }

    // Serve HTML page at the Root Route
    @GetMapping(value = "/", produces = MediaType.TEXT_HTML_VALUE)
    public Resource home() {
        return new ClassPathResource("static/demo.html");
    }

    // POST API Target
    @PostMapping("/api/generate-scoped-public-key")
    public ResponseEntity<Object> generateKey(@RequestBody ScopedKeyRequest request) {
        
        Map<String, Object> result = upstreamService.generateScopedPublicKey(
                request.engines(),
                request.expiresInSeconds()
        );

        if (result == null) {
            return ResponseEntity.status(HttpStatus.BAD_GATEWAY)
                    .body(Map.of("detail", "Upstream service error"));
        }

        return ResponseEntity.ok(result);
    }

    // Backup health validation route
    @GetMapping("/status")
    public Map<String, String> status() {
        return Map.of(
            "message", "Backend is running",
            "endpoint", "/api/generate-scoped-public-key"
        );
    }
}
