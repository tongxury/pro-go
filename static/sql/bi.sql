CREATE TABLE juguang_events
(
    clickId      String,
    caidMd5      FixedString(32),
    paid         String,
    caid         String,
    ua           String,
    ua1          String,
    idfaMd5      FixedString(32),
    os           LowCardinality(String),
    ip           String,
    ts           DateTime,
    appId        String,
    advertiserId String,
    unitId       String,
    campaignId   String,
    creativityId String,
    placement    LowCardinality(String),
    androidId    String,
    oaid         String,
    oaidMd5      FixedString(32),
    imei         String
) ENGINE = MergeTree()
      PARTITION BY toYYYYMM(ts)
      ORDER BY (ts, advertiserId, campaignId)
      SETTINGS index_granularity = 8192;



DROP TABLE IF EXISTS events;
CREATE TABLE events
(
    event_name   LowCardinality(String),
    ip           IPv4,
    country_code LowCardinality(String),
    created_at   TIMESTAMP,
    user_id      String,
    device_id    String,
    channel      LowCardinality(String),
    platform     LowCardinality(String),
    INDEX events_idx1 (event_name, user_id, created_at) TYPE minmax GRANULARITY 1
) ENGINE = MergeTree
      PARTITION BY toYYYYMMDD(created_at)
      ORDER BY (event_name, created_at, user_id)
      SETTINGS index_granularity = 8192;


drop table if exists user_payments;
CREATE TABLE user_payments
(
    id         UInt64,
    created_at TIMESTAMP,
    platform   String,
    user_id    String,
    plan_id    String,
    amount     Float64,
    expire_at  TIMESTAMP,
    status     String
) ENGINE = MySQL('', 'pro', '', 'admin',
                 '');



drop table if exists users;
CREATE TABLE users
(
    id         UInt64,
    created_at TIMESTAMP,
    nickname   String,
    email      String,
    phone      String,
    channel    String
) ENGINE = MySQL('', 'pro', 'users', 'admin', '');


drop table if exists questions;
create table questions
(
    id         String,
    created_at TIMESTAMP,
    prompt_id  String,
    session_id String,
    scene      String,
    user_id    String,
    status     String,
    cost       INT,
    model      String,
    INDEX questions_idx1 (scene, prompt_id, created_at) TYPE minmax GRANULARITY 1
) ENGINE = MergeTree
      PARTITION BY toYYYYMMDD(created_at)
      ORDER BY (id)
      SETTINGS index_granularity = 8192;

