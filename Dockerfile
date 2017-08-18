# Use an official Python runtime as a parent image
#FROM golang:1.8-alpine
FROM hypriot/rpi-alpine

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
ADD . /app

# Run weatherai when the container launches
ENTRYPOINT ["/app/weatherai"]
