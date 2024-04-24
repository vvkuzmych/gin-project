# Use the official PostgreSQL image as a base
FROM postgres:latest

# (Optional) Add custom initialization scripts
# COPY ./init.sql /docker-entrypoint-initdb.d/

# (Optional) Set environment variables
ENV POSTGRES_DB=gin_pro1
ENV POSTGRES_USER=admin
ENV POSTGRES_PASSWORD=admin123

# The default port that PostgreSQL listens on
EXPOSE 5432

