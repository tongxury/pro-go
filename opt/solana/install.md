### 安装 solana cli

https://docs.solanalabs.com/

### 生成 validator-keypair.json
```shell
solana-keygen new -o validator-keypair.json
# solana-keygen new -o ~/vote-account-keypair.json
```

### 启动命令
```shell
solana-validator \
    --identity /data/solana/validator-keypair.json \
    --known-validator 5D1fNXzvv5NjV1ysLjirC4WY92RNsVH18vjmcszZd8on \
    --known-validator dDzy5SR3AXdYWVqbDEkVFdvSPCtS9ihF5kJkHCtXoFs \
    --known-validator eoKpUABi59aT4rR9HGS3LcMecfut9x7zJyodWWP43YQ \
    --known-validator 7XSY3MrYnK8vq693Rju17bbPkCN3Z7KvvfvJx4kdrsSY \
    --known-validator Ft5fbkqNa76vnsjYNwjDZUXoTWpP7VYm3mtsaQckQADN \
    --known-validator 9QxCLckBiJc783jnMvXZubK4wH86Eqqvashtrwvcsgkv \
    --dynamic-port-range 8000-8020 \
    --entrypoint 15.204.241.61:8001 \
    --entrypoint 173.231.44.194:8001 \
    --entrypoint 178.32.184.117:8001 \
    --entrypoint 202.8.11.43:8001 \
    --entrypoint 216.144.245.106:8001 \
    --gossip-host 157.90.94.185 \
    --gossip-port 8001 \
    --expected-genesis-hash 5eykt4UsFv8P8NJdTREpY1vzqKqZKvdpKuc147dw2N9d \
    --rpc-port 8899 \
    --no-voting \
    --ledger /data/solana/ledger \
    --limit-ledger-size 200000000 \
    --disable-banking-trace \
    --enable-extended-tx-metadata-storage \
    --enable-rpc-transaction-history \
    --full-rpc-api \
    --only-known-rpc \
    --replay-slots-concurrently \
    --rpc-pubsub-enable-block-subscription \
    --rpc-scan-and-fix-roots \
    --account-index program-id \
    --account-index spl-token-owner \
    --account-index spl-token-mint \
    --bind-address 0.0.0.0 \
    --log /data/solana/solana-validator.log \
    --wal-recovery-mode absolute_consistency \
    --hard-fork 291900000
```

```shell

solana-validator \
    --identity /data/solana/validator-keypair.json \
    --known-validator 5D1fNXzvv5NjV1ysLjirC4WY92RNsVH18vjmcszZd8on \
    --known-validator 7XSY3MrYnK8vq693Rju17bbPkCN3Z7KvvfvJx4kdrsSY \
    --known-validator Ft5fbkqNa76vnsjYNwjDZUXoTWpP7VYm3mtsaQckQADN \
    --known-validator 9QxCLckBiJc783jnMvXZubK4wH86Eqqvashtrwvcsgkv \
    --only-known-rpc \
    --log /data/solana/solana-validator.log \
    --ledger /data/solana/ledger \
    --rpc-port 8899 \
    --dynamic-port-range 8000-8020 \
    --entrypoint entrypoint.testnet.solana.com:8001 \
    --entrypoint entrypoint2.testnet.solana.com:8001 \
    --entrypoint entrypoint3.testnet.solana.com:8001 \
    --expected-genesis-hash 4uhcVJyU9pJkvQyS88uRDiswHXSCkY3zQawwpjk2NsNY \
    --wal-recovery-mode skip_any_corrupted_record \
    --limit-ledger-size
```