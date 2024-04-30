# Use a Golang base image
FROM golang:1.21.1

# Set the working directory inside the container
WORKDIR /app

# Copy the local code to the container
COPY . .

COPY go.mod go.sum ./
RUN go mod download

# Build the Go application
RUN go build -o /gin-project cmd/main.go

# Expose port 8080 to the outside world
EXPOSE 8080:80

# Command to run the executable
CMD ["go", "run", "cmd/main.go"]




# # Final stage
# FROM gcr.io/distroless/base-debian10
# COPY --from=builder /goapp /
# ENTRYPOINT ["/goapp"]
