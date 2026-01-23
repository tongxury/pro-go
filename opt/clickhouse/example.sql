
drop table if exists dex_trades_sharded on cluster default;
create table dex_trades_sharded on cluster default
(
    id           String,
    wallet       String,
    side         String,
    token        String,
    amount       Float64,
    price        Float64,
    quote_token  String,
    quote_amount Float64,
    quote_price  Float64,
    created_at   DATETIME,
    block_index  Int64,
    tx_result    BOOLEAN default True,
    INDEX dex_trades_idx1 (token, created_at) TYPE minmax GRANULARITY 1,
    INDEX dex_trades_idx2 (wallet, created_at) TYPE minmax GRANULARITY 1,
    INDEX dex_trades_idx3 (created_at) TYPE minmax GRANULARITY 1
    ) ENGINE = MergeTree
    PARTITION BY toYYYYMMDD(created_at)
    ORDER BY (token, created_at, block_index)
    SETTINGS index_granularity = 8192;



create materialized view trending_tokens_1h on cluster default engine ReplacingMergeTree() order by (token, time) as
select token, count(1) as buyCount, min(created_at) as time
from dex_trades
where created_at > now() - interval 1 hour
  and token not in
    ('EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v', 'Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB',
    'So11111111111111111111111111111111111111112')
  and side = 'buy'
group by token;
-- order by count(1) desc
-- limit ?

select *
from trending_tokens_1h
         limit 1;

drop table if exists dex_trades_sharded on cluster default;
create table dex_trades_sharded on cluster default
(
    id           String,
    wallet       String,
    side         String,
    token        String,
    amount       Float64,
    price        Float64,
    quote_token  String,
    quote_amount Float64,
    quote_price  Float64,
    created_at   DATETIME,
    block_index  Int64,
    tx_result    BOOLEAN default True,
    INDEX dex_trades_idx1 (token, created_at) TYPE minmax GRANULARITY 1,
    INDEX dex_trades_idx2 (wallet, created_at) TYPE minmax GRANULARITY 1,
    INDEX dex_trades_idx3 (created_at) TYPE minmax GRANULARITY 1
    ) ENGINE = MergeTree
    PARTITION BY toYYYYMMDD(created_at)
    ORDER BY (token, created_at, block_index)
    SETTINGS index_granularity = 8192;


select count()
from dex_trades;

drop table if exists dex_trades on cluster default;
CREATE TABLE dex_trades ON CLUSTER default
AS dex_trades_sharded ENGINE = Distributed('default', 'default', 'dex_trades_sharded', cityHash64(token));

select cityHash64('3Cyrb9NABtEMJWv5yJ5pnccJPt3UYb1TmyiACu98NQts')

-- CREATE MATERIALIZED VIEW [IF NOT EXISTS] [db.]table_name [ON CLUSTER cluster_name] [TO[db.]name] [ENGINE = engine] [POPULATE] AS SELECT ...

drop view ohlc_6m_mv;
CREATE MATERIALIZED VIEW ohlc_6m_mv engine ReplacingMergeTree() order by token primary key token as
select token,
       max(price)                                                                 as high,
       min(price)                                                                 as low,
       arrayFirst(x -> 1, arraySort(x ->
                                    x.2, groupArray((price, block_index)))).1 as open,
       arrayLast(x -> 1, arraySort(x ->
                                       x.2, groupArray((price, block_index)))).1  as close,
       min(created_at)
from dex_trades
where created_at > now() - (INTERVAL 6 MINUTE)
group by token;

drop view ohlc_24h_mv;

CREATE MATERIALIZED VIEW ohlc_24h_mv engine ReplacingMergeTree() order by token primary key token as
select token,
       max(price)                                                                 as high,
       min(price)                                                                 as low,
       arrayFirst(x -> 1, arraySort(x ->
                                    x.2, groupArray((price, block_index)))).1 as open,
       arrayLast(x -> 1, arraySort(x ->
                                       x.2, groupArray((price, block_index)))).1  as close,
       min(created_at)
from dex_trades
where created_at > now() - (INTERVAL 24 HOUR)
group by token;


drop view trades_all;
create materialized view trades_all
    engine Log
as
select id, wallet, side, created_at, token, block_index
from trades;

create table dex_trades
(
    id           String,
    wallet       String,
    side         String,
    token        String,
    amount       Float64,
    price        Float64,
    quote_token  String,
    quote_amount Float64,
    quote_price  Float64,
    created_at   DATETIME,
    block_index  Int64,
    tx_result    BOOLEAN default True,
    INDEX dex_trades_idx1 (token, created_at) TYPE minmax GRANULARITY 1,
    INDEX dex_trades_idx2 (wallet, created_at) TYPE minmax GRANULARITY 1,
    INDEX dex_trades_idx3 (created_at) TYPE minmax GRANULARITY 1
) ENGINE = MergeTree
      PARTITION BY toYYYYMMDD(created_at)
      ORDER BY (token, created_at, block_index)
      SETTINGS index_granularity = 8192;



select now() - (INTERVAL 24 HOUR);



select *
from dex_trades
         limit 1;

select *
from system.clusters;

drop table if exists dex_trades_stream on cluster default;
create table dex_trades_stream
(
    id           String,
    wallet       String,
    side         String,
    token        String,
    amount       Float64,
    price        Float64,
    quote_token  String,
    quote_amount Float64,
    quote_price  Float64,
    created_at   Int64,
    block_index  Int64,
    tx_result    BOOLEAN
) ENGINE = Kafka SETTINGS kafka_broker_list = 'kafka-headless.prod:9092',
    kafka_topic_list = 'dexTradesStream',
    kafka_group_name = 'dexTradesStreamConsumer',
    kafka_format = 'JSONEachRow',
    kafka_num_consumers = 10;


drop table if exists dex_trades;
create table dex_trades
(
    id           String,
    wallet       String,
    side         String,
    token        String,
    amount       Float64,
    price        Float64,
    quote_token  String,
    quote_amount Float64,
    quote_price  Float64,
    created_at   Int64,
    block_index  Int64,
    tx_result    BOOLEAN default True,
    INDEX dex_trades_idx1 (token, created_at) TYPE minmax GRANULARITY 1,
    INDEX dex_trades_idx2 (wallet, created_at) TYPE minmax GRANULARITY 1,
    INDEX dex_trades_idx3 (created_at) TYPE minmax GRANULARITY 1
) ENGINE = MergeTree
      PARTITION BY toYYYYMMDD(toDateTime(created_at))
      ORDER BY (token, created_at, block_index)
      SETTINGS index_granularity = 8192;

drop view dex_trades_stream_migration;
CREATE MATERIALIZED VIEW dex_trades_stream_migration to dex_trades AS
SELECT *
FROM dex_trades_stream;


select toDateTime(1736180163);

select *
from trades;

select *
from system.kafka_consumers
where table = 'dex_trades_stream'
    format Vertical;

set stream_like_engine_allow_direct_select = 1;
select count()
from dex_trades_stream
;

select *
from trending_tokens_1h
order by buyCount desc
    limit 10;


select id, toDateTime(created_at)
from dex_trades
order by created_at desc
    limit 1;

helm install kafka oci://registry-1.docker.io/bitnamicharts/kafka -n prod
helm install clickhouse oci://registry-1.docker.io/bitnamicharts/clickhouse -n prod

WUr8Iz3pt2