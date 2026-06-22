using System.Net.Http.Json;
using System.Text.Json.Serialization;

namespace CSharpSample.API;

// Payload matching the original JSON structure
public record ScopedKeyRequest(
    [property: JsonPropertyName("engines")]
    [property: JsonIgnore(Condition = JsonIgnoreCondition.WhenWritingNull)] 
    List<string>? Engines,

    [property: JsonPropertyName("expiresInSeconds")]
    [property: JsonIgnore(Condition = JsonIgnoreCondition.WhenWritingNull)] 
    int? ExpiresInSeconds
);

public class UpstreamService
{
    private readonly HttpClient _httpClient;
    private readonly string? _apiUrl;
    private readonly string? _apiKey;

    public UpstreamService(HttpClient httpClient)
    {
        _httpClient = httpClient;
        
        _apiUrl = Environment.GetEnvironmentVariable("API_URL") ?? "";
        _apiKey = Environment.GetEnvironmentVariable("API_KEY") ?? "";

        if (string.IsNullOrEmpty(_apiUrl))
            Console.WriteLine("⚠️ Warning: API_URL is not set in environment variables");
        if (string.IsNullOrEmpty(_apiKey))
            Console.WriteLine("⚠️ Warning: API_KEY is not set in environment variables");
    }

    public async Task<object?> GenerateScopedPublicKeyAsync(List<string>? engines = null, int? expiresInSeconds = null)
    {
        try
        {
            var requestUrl = $"{_apiUrl?.TrimEnd('/')}/plan/auth/generate-scoped-public-key";
            var payload = new ScopedKeyRequest(engines, expiresInSeconds);

            var request = new HttpRequestMessage(HttpMethod.Post, requestUrl)
            {
                Content = JsonContent.Create(payload)
            };

            // Set up headers
            request.Headers.Add("X-Api-Key", _apiKey);

            // Timeout is automatically managed on the HttpClient registration (30s)
            var response = await _httpClient.SendAsync(request);

            if (!response.IsSuccessStatusCode)
            {
                Console.WriteLine($"❌ Generate Scoped Public Key API Request Failed! Status: {response.StatusCode}");
                
                if (response.StatusCode == System.Net.HttpStatusCode.Unauthorized)
                    Console.WriteLine("💡 Hint: API Key is missing or invalid.");
                else if (response.StatusCode == System.Net.HttpStatusCode.Forbidden)
                    Console.WriteLine("💡 Hint: Your IP address is not included in the IP allowlist.");

                return null;
            }

            Console.WriteLine("✅ Generate Scoped Public Key API Request Successful!");
            
            // Read response content dynamically as an unstructured JSON object/dictionary
            return await response.Content.ReadFromJsonAsync<Dictionary<string, object>>();
        }
        catch (HttpRequestException ex)
        {
            Console.WriteLine("No response received from server. Please check your URL or ensure the backend server is running.");
            Console.WriteLine(ex.Message);
            return null;
        }
    }
}
