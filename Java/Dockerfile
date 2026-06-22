# --- Stage 1: Build the Application jar file ---
FROM maven:3.9.6-eclipse-temurin-21-alpine AS builder
WORKDIR /app

# Copy the entire project context (including pom.xml and src)
COPY . .

# Compile and package everything into an executable production runner
RUN mvn clean package -DskipTests

# --- Stage 2: Runtime Environment ---
FROM eclipse-temurin:21-jre-alpine
WORKDIR /app

# Copy the compiled .jar asset out of the builder footprint
COPY --from=builder /app/target/*.jar app.jar

# Expose default application communication line
EXPOSE 3000

# Instruct Spring Boot to run on Port 3000 instead of default 8080
ENV SERVER_PORT=3000

ENTRYPOINT ["java", "-jar", "app.jar"]
