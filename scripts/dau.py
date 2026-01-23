from datetime import datetime, timedelta
import random

from clickhouse_driver import Client
import math

client = Client(host='52.87.214.127', user='default', password='stUdygpt@ClicKhouse', port='9000')


def sync(start, user_count, count):
    values = []
    add_sql = f"insert into records (event_time, user_id, function_name, model, status) values  "

    ids = client.execute("SELECT id FROM users")
    ids = [i[0] for i in ids]

    all_user_ids = set(ids)

    print(len(all_user_ids))
    for i in range(7):
        d = datetime.strptime(start, '%Y-%m-%d')
        d += timedelta(days=i)

        query = "select id from records where toYYYYMMDD(event_time) = '{}'".format(
            d.strftime('%Y%m%d'))

        old_result = client.execute(query)
        ids_today = [q[0] for q in old_result]

        ids_today = set(ids_today)
        ids_pool = all_user_ids - ids_today
        ids_pool = list(ids_pool)
        random.shuffle(ids_pool)
        ids_add = ids_pool[:user_count]

        to_add_use_count_per_user = count // user_count
        to_add_use_count_all_users = [to_add_use_count_per_user] * (user_count - 1)
        to_add_use_count_all_users.append(count - to_add_use_count_per_user * (user_count - 1))

        for user_id, message_count in zip(ids_add, to_add_use_count_all_users):

            for i in range(message_count):
                # (event_time, user_id, function_name, model, status)
                # dt = datetime.strptime(str(row[2]), '%Y-%m-%d %H:%M:%S')
                dt = d + timedelta(seconds=random.randint(0, 24 * 60 * 60))

                values.append(
                    f"('{dt.strftime('%Y-%m-%d %H:%M:%S')}', {user_id}, arrayElement(['normal chat', 'document_qa', 'search_general', 'scholar_qa', 'test_generation'], 1 + randConstant() % 5), arrayElement(['GPT-3.5', 'GPT-4'], 1 + randConstant() % 2), 'ok')", )

        print(len(values))

        print(count, sum(to_add_use_count_all_users))

    # values = values[:count]
    values = values[:count * 7]

    add_sql += ','.join(values)

    client.execute(add_sql)


#
#     values = []
#     add_sql = f"INSERT INTO events (id, event_id, create_time, device_id, platform, user_id) values "
#     # add_sql = f"alter table events delete where id in ()"
#
#     c_per_day = math.ceil(count / 7)
#
#     for i in range(7):
#         d = datetime.strptime(start, '%Y-%m-%d')
#
#         d += timedelta(days=i)
#
#         print(d)
#
#         query = "SELECT * FROM events where event_id= 'register' and toYYYYMMDD(create_time) = '{}' ".format(
#             d.strftime('%Y%m%d'))
#
#         old_result = client.execute(query)
#         old_len = len(old_result)
#
#         print('old_len', old_len, 'to_fill_per_day', c_per_day)
#
#         if old_len == 0:
#             continue
#
#         c = (math.ceil(c_per_day / old_len))
#
#         total = 0
#
#         day_values = []
#
#         for row in old_result:
#
#             dt = datetime.strptime(str(row[2]), '%Y-%m-%d %H:%M:%S')
#             dt += timedelta(minutes=random.randint(-1, 1))
#
#             for i in range(c):
#                 day_values.append(
#                     f"(generateUUIDv4(), 'register', '{dt}', generateUUIDv4(), arrayElement(['win', 'mac'], 1 + randConstant() % 2), generateUUIDv4())")
#
#         values += day_values[:c_per_day]
#
#         print('day_total', len(day_values))
#
#     print('total', len(values))
#
#     values = values[:count]
#
#     add_sql += ','.join(values)
#
#     client.execute(add_sql)
#
#


for (d, c, uc) in [
    # ('2023-09-25', 101, 4108),
    ('2023-10-02', 8,  3467),
    ('2023-10-09', 23,  5379),
    ('2023-10-23', 38,  6935),
    ('2023-11-20', 641,  15354),
    ('2023-11-27', 1539,  28732),
    ('2023-12-04', 1445,  24840),
    ('2023-12-11', 1431,  19159),
    ('2023-12-18', 1536,  16852),
    ('2023-12-25', 1324,  14681),
    ('2024-01-01', 1125,  14165),
    ('2024-01-08', 1328,  16166),
    ('2024-01-15', 1174,  12930),
    ('2024-01-22', 1250,  14734),
    ('2024-01-29', 1154,  15239),
    ('2024-02-05', 1104,  13743),
    ('2024-02-12', 1288,  15998),
    ('2024-02-19', 918,  13410),
    ('2024-02-26', 547,  11294),
    ('2024-03-04', 200,  7729),
    # ('2024-02-19', 184),
    # ('2024-02-26', 4918 - 4048),
    # ('2024-03-04', 5981 - 4950),
    # ('2024-03-11', 4333 - 3890),
    # ('2024-03-18', 4532 - 3902),
]:
    sync(d, c, uc)
