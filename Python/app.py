from typing import List, Optional
from dotenv import load_dotenv
from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import FileResponse
from fastapi.staticfiles import StaticFiles
from pydantic import BaseModel, Field

from services import generate_scoped_public_key

load_dotenv()

app = FastAPI()

# Enable CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # change in production
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


class GenerateKeyRequest(BaseModel):
    engines: Optional[List[str]] = Field(default=None)
    expiresInSeconds: Optional[int] = Field(default=None)

# Serve static files
app.mount("/static", StaticFiles(directory="static"), name="static")

@app.get("/")
async def home():
    return FileResponse("static/demo.html")

@app.post("/api/generate-scoped-public-key")
async def generate_key(data: GenerateKeyRequest):
    print(data)
    result = generate_scoped_public_key(
        data.engines,
        data.expiresInSeconds
    )

    if result is None:
        raise HTTPException(
            status_code=502,
            detail="Upstream service error"
        )

    return result


@app.get("/")
async def root():
    return {
        "message": "Backend is running",
        "endpoint": "/api/generate-scoped-public-key"
    }


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(
        "app:app",
        host="0.0.0.0",
        port=3000,
        reload=True
    )
