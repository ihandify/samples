import express from 'express';
import cors from 'cors';
import path from 'path';
import { fileURLToPath } from 'url';
import dotenv from 'dotenv';
import { generateScopedPublicKey } from './services.js';

dotenv.config();

const app = express();
const PORT = 3000;

// Setup __dirname equivalent for ES modules
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Enable CORS (Mirrors FastAPI's allow_origins=["*"])
app.use(cors({
    origin: '*',
    credentials: true,
    methods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
    allowedHeaders: ['*'] 
}));

// Built-in middleware to parse JSON bodies
app.use(express.json());

// Serve static files from the /static folder
app.use('/static', express.static(path.join(__dirname, 'static')));

// Serve the demo HTML file at the root route
app.get('/', (req, res) => {
    res.sendFile(path.join(__dirname, 'static', 'demo.html'));
});

// POST Endpoint
app.post('/api/generate-scoped-public-key', async (req, res) => {
    const { engines, expiresInSeconds } = req.body;

    const result = await generateScopedPublicKey(engines, expiresInSeconds);

    if (result === null) {
        return res.status(502).json({
            detail: "Upstream service error"
        });
    }

    return res.json(result);
});

// Fallback or backup root route metadata (Note: in original code this was duplicate, 
// but if you want it on a different path like /api, you can use this)
app.get('/status', (req, res) => {
    res.json({
        message: "Backend is running",
        endpoint: "/api/generate-scoped-public-key"
    });
});

app.listen(PORT, '0.0.0.0', () => {
    console.log(`Server is running on http://0.0.0.0:${PORT}`);
});
