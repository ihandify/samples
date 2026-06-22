using CSharpSample.API;

var builder = WebApplication.CreateBuilder(args);

// Load .env variables manually using DotNetEnv
DotNetEnv.Env.Load();

// Configure CORS (Mirrors FastAPI's allow_origins=["*"])
builder.Services.AddCors(options =>
{
    options.AddDefaultPolicy(policy =>
    {
        policy.AllowAnyOrigin()
              .AllowAnyHeader()
              .AllowAnyMethod();
    });
});

// Register UpstreamService with a configured HttpClient + 30-second timeout
builder.Services.AddHttpClient<UpstreamService>(client =>
{
    client.Timeout = TimeSpan.FromSeconds(30);
});

var app = builder.Build();

app.UseCors();

// Serve static files from the wwwroot or custom folder. 
// To match your project directory, we map /static URL path to the ./static directory.
app.UseStaticFiles(new StaticFileOptions
{
    FileProvider = new Microsoft.Extensions.FileProviders.PhysicalFileProvider(
        Path.Combine(builder.Environment.ContentRootPath, "static")),
    RequestPath = "/static"
});

// Serve the demo HTML at the root URL
app.MapGet("/", async context =>
{
    var filePath = Path.Combine(builder.Environment.ContentRootPath, "static", "demo.html");
    await context.Response.SendFileAsync(filePath);
});

// POST Endpoint
app.MapPost("/api/generate-scoped-public-key", async (ScopedKeyRequest request, UpstreamService upstreamService) =>
{
    var result = await upstreamService.GenerateScopedPublicKeyAsync(request.Engines, request.ExpiresInSeconds);

    if (result == null)
    {
        return Results.Json(new { detail = "Upstream service error" }, statusCode: 502);
    }

    return Results.Ok(result);
});

// Optional Health Status Route
app.MapGet("/status", () => Results.Ok(new 
{
    message = "Backend is running",
    endpoint = "/api/generate-scoped-public-key"
}));

// Run on port 3000 (0.0.0.0 handles container binding mapping)
app.Run("http://0.0.0.0:3000");
