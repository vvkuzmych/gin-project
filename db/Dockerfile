# Use the official PostgreSQL image as a base
FROM postgres:latest

# (Optional) Add custom initialization scripts
# COPY ./init.sql /docker-entrypoint-initdb.d/

# (Optional) Set environment variables
ENV POSTGRES_DB=gin_pro
ENV POSTGRES_USER=vladimirkuzmich
ENV POSTGRES_PASSWORD=yourpassword

# The default port that PostgreSQL listens on
EXPOSE 5432

