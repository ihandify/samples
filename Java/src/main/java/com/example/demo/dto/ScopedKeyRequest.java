package com.example.demo.dto;

import jakarta.validation.constraints.NotEmpty;
import jakarta.validation.constraints.NotNull;
import java.util.List;

public record ScopedKeyRequest(
    List<String> engines,
    Integer expiresInSeconds
) {}
