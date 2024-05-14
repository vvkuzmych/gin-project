# Use a Golang base image
FROM golang:1.21.1

# Set the working directory inside the container
WORKDIR /app

# Copy the local code to the container
COPY . .
COPY ./db/migration/ ./db/migration/

COPY go.mod go.sum ./
RUN go mod download

# Build the Go application
RUN go build -o /gin-project cmd/main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["/gin-project"]




# # Final stage
# FROM gcr.io/distroless/base-debian10
# COPY --from=builder /goapp /
# ENTRYPOINT ["/goapp"]
