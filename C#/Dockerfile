# --- Stage 1: Build the app using .NET 10.0 SDK ---
FROM mcr.microsoft.com/dotnet/sdk:10.0-alpine AS build-env
WORKDIR /app

# Copy csproj and restore dependencies
COPY *.csproj ./
RUN dotnet restore

# Copy remaining source code and publish
COPY . ./
RUN dotnet publish -c Release -o out

# --- Stage 2: Runtime environment using .NET 10.0 ASP.NET ---
FROM mcr.microsoft.com/dotnet/aspnet:10.0-alpine
WORKDIR /app

# Copy published outputs from Stage 1
COPY --from=build-env /app/out .
COPY ./static ./static

EXPOSE 3000

ENTRYPOINT ["dotnet", "csharp-sample.dll"]
