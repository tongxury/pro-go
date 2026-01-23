FROM --platform=linux/amd64 debian:bookworm-slim

ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /app

# Update package lists and install wget and tar
RUN apt-get update && apt-get install -y \
    wget \
    unzip \
    openssl \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

RUN wget https://github.com/jup-ag/jupiter-swap-api/releases/download/v6.0.33/jupiter-swap-api-x86_64-unknown-linux-gnu.zip \
    && unzip jupiter-swap-api-x86_64-unknown-linux-gnu.zip \
    && chmod +x jupiter-swap-api

ENV RUST_LOG=info

#ENTRYPOINT ["/app/jupiter-swap-api", "--enable-new-dexes", "--rpc-url https://dark-radial-seed.solana-mainnet.quiknode.pro/95494d59ce7464c3b374dcae1d25c0a3cba837f2"]


