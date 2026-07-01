# iHandify Sample Application in Java

This sample application demonstrates how to integrate iHandify APIs into your application.

## Prerequisites

- Docker
- Docker Compose
- An active iHandify subscription
- A valid Secret API Key

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/ihandify/samples.git
cd samples/Java
```

### 2. Obtain Your Secret API Key

1. Sign in to the iHandify Dashboard.
2. Navigate to **API Key**.
3. Copy your **API Key**.

### 3. Configure Environment Variables

Open `.env` and replace the placeholder value with your Secret API Key:

```env
API_KEY="Paste_your_Secret_API_Key_here"
```

### 4. Start the Application

Run the application using Docker Compose:

```bash
docker compose up -d --build
```

Verify that the containers are running:

```bash
docker compose ps
```

### 5. Open the Demo Application

Open your browser and navigate to:

```text
http://localhost:8060
```

You can now try handwriting recognition using the sample application.

---

## Authentication Flow

This sample follows the recommended security architecture:

```text
Frontend
    ↓
Sample Backend
    ↓ (Secret API Key)
iHandify Authentication API
    ↓
Scoped Public API Key
    ↓
Frontend
    ↓
iHandify Recognition API
```

1. The frontend sends a request to the sample backend.
```javascript
    const response = await fetch(KEY_API_URL, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            expiresInSeconds: SCOPED_KEY_EXPIRES_IN
        })
    })
```
2. The backend uses the Secret API Key to request a Scoped Public API Key from iHandify.
```java
    String url = this.apiUrl + "/plan/auth/generate-scoped-public-key";
        
    Map<String, Object> payload = new HashMap<>();
    if (engines != null) {
        payload.put("engines", engines);
    }
    if (expiresInSeconds != null) {
        payload.put("expiresInSeconds", expiresInSeconds);
    }

    return this.restClient.post()
                .uri(url)
                .contentType(MediaType.APPLICATION_JSON)
                .header("X-Api-Key", this.apiKey)
                .body(payload)
                .retrieve()
                .body(Map.class); // Dynamically maps arbitrary upstream
```
3. The Scoped Public API Key is returned to the frontend.
```javascript
    const result = await response.json();
    scopedPublicKey = result.data.scopedPublicKey;
```
4. The frontend uses the Scoped Public API Key to call iHandify recognition APIs.
```javascript
    const fd = new FormData();
    fd.append("input", JSON.stringify(pattern));

    const response = await fetch(url, {
        method: 'POST',
        body: fd,
        signal: AbortSignal.timeout(REQUEST_TIMEOUT_MS),
        headers: {
            'x-scoped-key': scopedPublicKey || ''
        }
    });
```

> [!CAUTION]
**Never expose your Secret API Key in frontend applications.**

---

## Stopping the Application

```bash
docker compose down
```

## Troubleshooting

### Container failed to start

Check container logs:

```bash
docker compose logs
```

### Authentication failed

Verify that:

- The Secret API Key is valid.
- The API Key is correctly configured in `.env`.
- Your subscription is active.

### Port 8060 is already in use

Modify the port mapping in `docker-compose.yml` or stop the application currently using port 8060.

---

## Additional Resources

- Documentation: https://ihandify.com/docs
- Additional Sample Applications: https://github.com/ihandify

For questions or support, please contact the iHandify team.