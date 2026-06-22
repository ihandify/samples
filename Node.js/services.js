import axios from 'axios';
import dotenv from 'dotenv';

dotenv.config();

const API_URL = process.env.API_URL || "";
const API_KEY = process.env.API_KEY || "";

if (!API_URL) {
    console.warn("⚠️ Warning: API_URL is not set in .env");
}
if (!API_KEY) {
    console.warn("⚠️ Warning: API_KEY is not set in .env");
}

/**
 * Calls the upstream API to generate a scoped public key.
 * @param {string[]} engines 
 * @param {number} expiresInSeconds 
 * @returns {Promise<object|null>}
 */
export async function generateScopedPublicKey(engines, expiresInSeconds) {
    try {
        const payload = {};
        if (engines != null) {
            payload.engines = engines;
        }
        if (expiresInSeconds != null) {
            payload.expiresInSeconds = expiresInSeconds;
        }

        const headers = {
            "Content-Type": "application/json",
            "X-Api-Key": API_KEY
        };

        // 30 second timeout (30000ms)
        const response = await axios.post(
            `${API_URL}/plan/auth/generate-scoped-public-key`,
            payload,
            { headers, timeout: 30000 }
        );

        console.log("✅ Generate Scoped Public Key API Request Successful!");
        return response.data;

    } catch (error) {
        if (error.response) {
            // The server responded with a status code outside the 2xx range
            const statusCode = error.response.status;
            console.error("❌ Generate Scoped Public Key API Request Failed!");
            console.error(`Error Status Code: ${statusCode}`);
            console.error(`Error Details:`, error.response.data);

            if (statusCode === 401) {
                console.log("💡 Hint: API Key is missing or invalid.");
            } else if (statusCode === 403) {
                console.log("💡 Hint: Your IP address is not included in the IP allowlist.");
            }
        } else if (error.request) {
            // The request was made but no response was received
            console.error(
                "No response received from server. " +
                "Please check your URL or ensure the backend server is running."
            );
            console.error(error.message);
        } else {
            // Something else triggered the error
            console.error("Error setting up request:", error.message);
        }
        
        return null;
    }
}
