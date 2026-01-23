--market-cache <MARKET_CACHE>
Jupiter Europa 的 URL、文件路径或远程文件路径，指定市场缓存的来源。未指定时将使用关联市场模式的默认值。需要注意某些 AMMs 的 params 字段是特定 AMM 类型所必需的。

--market-mode <MARKET_MODE>
切换市场模式。选择的模式将决定 API 如何获取市场信息。可用的值包括 europa（默认）、remote 和 file。

--rpc-url <RPC_URL>
用于轮询和获取用户帐户的 RPC URL。

--secondary-rpc-urls <SECONDARY_RPC_URLS>...
用于某些 RPC 调用的备用 RPC URL 列表。

-e, --yellowstone-grpc-endpoint <YELLOWSTONE_GRPC_ENDPOINT>
指定 Yellowstone 的 gRPC 终端节点，例如 https://jupiter.rpcpool.com。

-x, --yellowstone-grpc-x-token <YELLOWSTONE_GRPC_X_TOKEN>
Yellowstone gRPC 的 x 令牌，通常是主机名后的令牌。

--yellowstone-grpc-enable-ping
启用 ping 以检测 gRPC 服务器，适用于负载均衡的 Yellowstone gRPC 终端节点。

--snapshot-poll-interval-ms <SNAPSHOT_POLL_INTERVAL_MS>
定义 AMMs 相关账户的缓存间隔。黄石 gRPC 模式下，定期轮询快照确认的 AMM 账户状态。默认情况下，轮询模式为 200 毫秒，黄石 gRPC 模式下为 30000 毫秒。
-
-enable-external-amm-loading
启用从 keyedUiAccounts 加载外部 AMMs，以支持与交换相关的端点。

--disable-swap-cache-loading
禁用加载与交换相关功能所需的缓存，例如地址查找表，适合于仅报价的 API。

--allow-circular-arbitrage
允许输入和输出为相同铸币的套利报价和交换。

--sentry-dsn <SENTRY_DSN>
为 Sentry 指定数据源名称以发送错误。

--dex-program-ids <DEX_PROGRAM_IDS>...
指定要包含的 DEX 程序 ID 列表，未包含的程序 ID 将不会加载。

--exclude-dex-program-ids <EXCLUDE_DEX_PROGRAM_IDS>...
指定要排除的 DEX 程序 ID 列表，排除的程序 ID 不会加载。

--filter-markets-with-mints <FILTER_MARKETS_WITH_MINTS>...
用于过滤市场的铸币列表，只有那些至少包含两个来自该集合的铸币的市场会被包含。

-H, --host <HOST>
应用程序的主机地址，默认为 0.0.0.0。

-p, --port <PORT>
应用程序运行的端口号，默认为 8080。

--metrics-port <METRICS_PORT>
用于 Prometheus 监控的指标端口。

-s, --expose-quote-and-simulate
启用 /quote-and-simulate 端点，用于在单个请求中报价并模拟交换。

--enable-deprecated-indexed-route-maps
启用计算和提供 /indexed-route-map 端点，此选项已被弃用，不推荐启用，因其有较高的开销。

--enable-new-dexes
启用最近集成的新 DEX。

--enable-diagnostic
启用 /diagnostic 端点以报价。

--enable-add-market
启用 /add-market 端点以动态加载新市场。

--total-thread-count <TOTAL_THREAD_COUNT>
用于 jupiter-swap-api 进程的线程总数，默认为 16。

--webserver-thread-count <WEBSERVER_THREAD_COUNT>
用于 Web 服务器的线程数量，默认为 2。

--update-thread-count <UPDATE_THREAD_COUNT>
更新线程数，默认为 4。

--loki-url <LOKI_URL>
Loki 的 URL。

--loki-username <LOKI_USERNAME>
Loki 的用户名。

--loki-password <LOKI_PASSWORD>
Loki 的密码。

--loki-custom-labels <LOKI_CUSTOM_LABELS>...
添加自定义标签到 Loki 指标，例如 APP_NAME=jupiter-swap-api,ENVIRONMENT=production。

-h, --help
打印帮助信息。