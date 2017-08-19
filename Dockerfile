# Use an official Python runtime as a parent image
FROM golang:1.8-alpine
RUN apk add --no-cache ca-certificates apache2-utils

# Set the working directory to /opt
WORKDIR /opt

# Copy the current directory contents into the container at /opt
ADD ./weatherai /opt

# Run weatherai when the container launches
ENTRYPOINT ["weatherai"]
