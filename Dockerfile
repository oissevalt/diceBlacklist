FROM golang:latest

# Set the environment for cross-compilation
ENV CGO_ENABLED=1 \
    CC=x86_64-w64-mingw32-gcc \
    CXX=x86_64-w64-mingw32-g++

# Set up the working directory
WORKDIR /usr/src/program

# Install cross-compiler for 64-bit Windows
RUN apt-get update && \
    apt-get install -y --no-install-recommends mingw-w64 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Copy the local package files to the container's workspace.
COPY . .
RUN chmod +x ./build_script.sh

# Use an entrypoint script to handle building
ENTRYPOINT ["./build_script.sh"]
