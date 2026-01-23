以下是对 solana-genesis 命令及其选项的详细解释：

用法
solana-genesis 是 Solana 区块链的一部分，用于生成区块链的创世区块，重建可扩展性。该命令需要启动验证者的身份和公钥以及账本目录。

标志 (FLAGS)
--enable-warmup-epochs：启用后，初始的纪元将很短，然后逐渐增长。此选项适用于快速让质押在开发过程中热身。

-h, --help：打印帮助信息，显示命令的所有可用选项和用法。

-V, --version：打印版本信息，显示当前 solana-genesis 的版本。

选项 (OPTIONS)
--bootstrap-stake-authorized-pubkey <BOOTSTRAP STAKE AUTHORIZED PUBKEY>：指定一个文件路径，该文件包含被授权管理启动验证者质押的公钥。默认等于 --bootstrap-validator IDENTITY_PUBKEY。

-b, --bootstrap-validator <IDENTITY_PUBKEY VOTE_PUBKEY STAKE_PUBKEY>...：启动验证者的身份公钥、投票公钥和质押公钥。

--bootstrap-validator-lamports <LAMPORTS>：指定分配给启动验证者的 lamport 数量，默认为 500,000,000,000 lamports。

--bootstrap-validator-stake-lamports <LAMPORTS>：指定分配给启动验证者质押账户的 lamport 数量，默认为 500,000,000 lamports。

--bpf-program <ADDRESS LOADER SBF_PROGRAM.SO>...：安装一个 SBF 程序到指定的地址。

--cluster-type <cluster_type>：选择将要启用的集群特性。默认为 mainnet-beta，可能值有 development, devnet, testnet, mainnet-beta。

--creation-time <RFC3339 DATE TIME>：启动验证者开始集群的时间；默认为当前系统时间。

-t, --faucet-lamports <LAMPORTS>：指定分配给水龙头（faucet）的 lamport 数量。

-m, --faucet-pubkey <PUBKEY>：指定一个文件路径，该文件包含水龙头公钥，默认路径为 /root/.config/solana/id.json。

--fee-burn-percentage <NUMBER>：收取的费用中用于销毁的百分比；默认为 50%。

--hashes-per-tick <NUM_HASHES|"auto"|"sleep">：每个滴答时间内计算的 PoH 哈希数量。如果是 "auto"，则根据 --target-tick-duration 和计算机的哈希率确定。如果是 "sleep"，则在开发中按照 --target-tick-duration 睡眠而不是计算；默认值是 "auto"。

--inflation <inflation>：选择通货膨胀选项；可能值有 pico, full, none。

--lamports-per-byte-year <LAMPORTS>：每年每字节数据收取的 lamport 费用；默认值为 3480 lamports。

-l, --ledger <DIR>：将指定目录用作持久账本位置。

--max-genesis-archive-unpacked-size <NUMBER>：创建的创世归档的最大总未压缩文件大小；默认值为 10,485,760 字节。

--primordial-accounts-file <FILENAME>...：包含初始账户和余额公钥的文件路径。

--rent-burn-percentage <NUMBER>：收取的租金中用于销毁的百分比；默认为 50%。

--rent-exemption-threshold <NUMBER>：账户要被视为免租的资金余额保持时间（单位：年）；默认为 2 年。

--slots-per-epoch <SLOTS>：每个纪元的插槽数量。

--target-lamports-per-signature <LAMPORTS>：集群按目标签名每个插槽收费的 lamport 数量；默认值为 10,000 lamports。

--target-signatures-per-slot <NUMBER>：用于估算集群所需处理能力的值。当最新插槽处理的签名数量少于/多于该值时，下一插槽的每个签名费用将下降/上升。值为 0 时禁用基于签名的费用调整；默认值为 20,000。

--target-tick-duration <MILLIS>：集群的目标滴答率，以毫秒为单位。

--ticks-per-slot <TICKS>：每个插槽的滴答数量；默认为 64。

--upgradeable-program <ADDRESS UPGRADEABLE_LOADER SBF_PROGRAM.SO UPGRADE_AUTHORITY>...：以给定地址安装一个可升级的 SBF 程序，以及相应的升级授权者（或“none”）。

--vote-commission-percentage <NUMBER>：投票佣金的百分比；默认为 100%。

这些选项共同配置了 Solana 区块链的基本构建，提供了灵活性和可定制性，以适应开发和生产环境的不同需求。