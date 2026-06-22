<?php

function generateScopedPublicKey(?array $engines = null, ?int $expiresInSeconds = null)
{
    $apiUrl = getenv('API_URL');
    $apiKey = getenv('API_KEY');

    if (!$apiUrl || !$apiKey) {
        error_log('API_URL or API_KEY is not configured.');
        return null;
    }

    $data = [];
    if ($engines !== null) {
        $data['engines'] = $engines;
    }
        if ($expiresInSeconds !== null) {
        $data['expiresInSeconds'] = $expiresInSeconds;
    }

    $payload = json_encode($data, JSON_FORCE_OBJECT);

    $ch = curl_init();

    curl_setopt_array($ch, [
        CURLOPT_URL => rtrim($apiUrl, '/') . '/plan/auth/generate-scoped-public-key',
        CURLOPT_RETURNTRANSFER => true,
        CURLOPT_POST => true,
        CURLOPT_POSTFIELDS => $payload,
        CURLOPT_HTTPHEADER => [
            'Content-Type: application/json',
            'X-Api-Key: ' . $apiKey
        ],
        CURLOPT_TIMEOUT => 30
    ]);

    $response = curl_exec($ch);
    $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);

    if (curl_errno($ch)) {
        error_log('Request failed: ' . curl_error($ch));
        curl_close($ch);
        return null;
    }

    curl_close($ch);

    if ($httpCode < 200 || $httpCode >= 300) {
        error_log("Generate Scoped Public Key API failed. HTTP {$httpCode}");
        error_log($response);

        if ($httpCode === 401) {
            error_log('Hint: API Key is missing or invalid.');
        } elseif ($httpCode === 403) {
            error_log('Hint: Your IP address is not included in the IP allowlist.');
        }

        return null;
    }

    return json_decode($response, true);
}
