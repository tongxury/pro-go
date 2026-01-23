import time

from clickhouse_driver import Client
import math

if __name__ == '__main__':
    client = Client(host='52.87.214.127', user='default', password='stUdygpt@ClicKhouse', port='9000')

    while True:
        time.sleep(1)
        results = client.execute("""
        select device_id,
               ip,
               app_version,
               created_at,
               chrome_version,
               platform,
               user_id,
               channel,
               country
        from user_devices
        where empty(country) and notEmpty(ip)
        limit 1000;
        """)

        # if not results:
        #     break

        insert_sql = "insert into user_devices (device_id, ip, app_version, created_at, chrome_version, platform, user_id, channel, country) values "
        values = []

        if results is not None:
            for result in results:
                ip_detail = client.execute(f"""
                    select country_code
                    from ip_details
                    where IPv4StringToNum('{result[1]}') between ip_range_start_num and ip_range_end_num;
                """, )

                print(result[0], result[1], ip_detail[0][0])

                values.append(
                    f"('{result[0]}','{result[1]}','{result[2]}','{result[3]}','{result[4]}','{result[5]}','{result[6]}','{result[7]}','{ip_detail[0][0]}')")

        if len(values) > 0:
            insert_sql += ",".join(values)

        print(insert_sql)
        # break
        client.execute(insert_sql)
        client.execute("optimize table user_devices final")
