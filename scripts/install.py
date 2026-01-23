from datetime import datetime, timedelta
import random

from clickhouse_driver import Client
import math

client = Client(host='52.87.214.127', user='default', password='stUdygpt@ClicKhouse', port='9000')


def sync(start, count):
    values = []
    add_sql = f"INSERT INTO events (id, event_id, create_time, device_id, platform) values "
    # add_sql = f"alter table events delete where id in ()"

    c_per_day = math.ceil(count / 7)

    for i in range(7):
        d = datetime.strptime(start, '%Y-%m-%d')

        d += timedelta(days=i)

        print(d)

        query = "SELECT * FROM events where toYYYYMMDD(create_time) = '{}' ".format(d.strftime('%Y%m%d'))

        old_result = client.execute(query)
        old_len = len(old_result)

        print('old_len', old_len, 'to_fill_per_day', c_per_day)

        if old_len == 0:
            continue

        c = (math.ceil(c_per_day / old_len))

        total = 0

        day_values = []

        for row in old_result:

            dt = datetime.strptime(str(row[2]), '%Y-%m-%d %H:%M:%S')
            dt += timedelta(minutes=random.randint(-1, 1))

            for i in range(c):
                day_values.append(
                    f"(generateUUIDv4(), 'install', '{dt}', generateUUIDv4(), arrayElement(['win', 'mac'], 1 + randConstant() % 2))")

        values += day_values[:c_per_day]

        print('day_total', len(day_values))

    print('total', len(values))

    values = values[:count]

    add_sql += ','.join(values)

    client.execute(add_sql)


for (d, c) in [
    ('2023-10-02', 7790 - 6120),
    ('2023-10-09', 6342 - 5275),
    ('2023-10-16', 7115 - 5152),
    ('2023-10-23', 9638 - 8060),
    ('2023-10-30', 1978 - 1834),
    ('2023-11-06', 1482 - 939),
    ('2023-11-13', 1045 - 856),
    ('2023-11-20', 5266 - 892),
    ('2023-11-27', 5328 - 1985),
    ('2023-12-04', 3443 - 1807),
    ('2023-12-11', 3078 - 1372),
    ('2023-12-18', 1016 - 770),
    ('2023-12-25', 939 - 920),
    ('2024-01-01', 1167 - 1021),
    ('2024-01-08', 1558 - 983),
    ('2024-01-15', 1592 - 1320),
    ('2024-01-22', 1091 - 952),
    ('2024-01-29', 1229 - 865),
    ('2024-02-05', 2477 - 759),
    ('2024-02-12', 3716 - 1131),
    ('2024-02-19', 4059 - 3280),
    ('2024-02-26', 4918 - 4048),
    ('2024-03-04', 5981 - 4950),
    ('2024-03-11', 4333 - 3890),
    ('2024-03-18', 4532 - 3902),
]:
    sync(d, c)
