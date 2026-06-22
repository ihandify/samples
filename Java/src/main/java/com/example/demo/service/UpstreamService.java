package com.example.demo.service;

import com.example.demo.dto.ScopedKeyRequest;
import io.github.cdimascio.dotenv.Dotenv;
import org.springframework.http.HttpStatusCode;
import org.springframework.http.MediaType;
import org.springframework.http.client.SimpleClientHttpRequestFactory;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestClient;
import org.springframework.web.client.RestClientResponseException;

import java.util.List;
import java.util.Map;
import java.util.HashMap;

@Service
public class UpstreamService {

    private final RestClient restClient;
    private final String apiUrl;
    private final String apiKey;

    public UpstreamService() {
        // Load configurations from .env
        Dotenv dotenv = Dotenv.configure().ignoreIfMissing().load();
        this.apiUrl = dotenv.get("API_URL", "");
        this.apiKey = dotenv.get("API_KEY", "");

        if (this.apiUrl.isEmpty()) System.out.println("⚠️ Warning: API_URL is not set in .env");
        if (this.apiKey.isEmpty()) System.out.println("⚠️ Warning: API_KEY is not set in .env");

        // Set up connection timeouts (30 seconds)
        SimpleClientHttpRequestFactory requestFactory = new SimpleClientHttpRequestFactory();
        requestFactory.setConnectTimeout(30000);
        requestFactory.setReadTimeout(30000);

        this.restClient = RestClient.builder()
                .requestFactory(requestFactory)
                .build();
    }

    @SuppressWarnings("unchecked")
    public Map<String, Object> generateScopedPublicKey(List<String> engines, Integer expiresInSeconds) {
        String url = this.apiUrl + "/plan/auth/generate-scoped-public-key";
        
        Map<String, Object> payload = new HashMap<>();
        if (engines != null) {
            payload.put("engines", engines);
        }
        if (expiresInSeconds != null) {
            payload.put("expiresInSeconds", expiresInSeconds);
        }

        try {
            return this.restClient.post()
                    .uri(url)
                    .contentType(MediaType.APPLICATION_JSON)
                    .header("X-Api-Key", this.apiKey)
                    .body(payload)
                    .retrieve()
                    .body(Map.class); // Dynamically maps arbitrary upstream JSON responses

        } catch (RestClientResponseException e) {
            System.out.println("❌ Generate Scoped Public Key API Request Failed!");
            HttpStatusCode statusCode = e.getStatusCode();
            System.out.println("Error Status Code: " + statusCode.value());
            System.out.println("Error Details: " + e.getResponseBodyAsString());

            if (statusCode.value() == 401) {
                System.out.println("💡 Hint: API Key is missing or invalid.");
            } else if (statusCode.value() == 403) {
                System.out.println("💡 Hint: Your IP address is not included in the IP allowlist.");
            }
            return null;
        } catch (Exception e) {
            System.out.println("No response received from server. Please check your URL or ensure the backend server is running.");
            System.out.println(e.getMessage());
            return null;
        }
    }
}
