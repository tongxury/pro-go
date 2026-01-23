from pymongo import MongoClient
import os
from clickhouse_driver import Client
import datetime
import urllib.parse


class DatabaseMigrator:
    def __init__(self):
        # MongoDB配置
        self.mongo_config = {
            'username': 'root',
            'password': os.getenv('MONGO_PASSWORD'),
            'host': '54.222.192.154',
            'port': 27017,
            'database': 'pro',
            'collection': 'questions'
        }

        # ClickHouse配置
        self.clickhouse_config = {
            'host': '52.80.182.238',
            'port': 9000,
            'user': 'default',
            'password': os.getenv('CLICKHOUSE_PASSWORD'),
            'database': 'default'
        }

        self.batch_size = 10
        self.mongo_client = None
        self.ch_client = None

    def get_mongo_uri(self):
        """构建MongoDB URI"""
        username = urllib.parse.quote_plus(self.mongo_config['username'])
        password = urllib.parse.quote_plus(self.mongo_config['password'])
        return f"mongodb://{username}:{password}@{self.mongo_config['host']}:{self.mongo_config['port']}/{self.mongo_config['database']}?authSource=admin"

    def connect_mongo(self):
        """连接MongoDB"""
        try:
            self.mongo_client = MongoClient(self.get_mongo_uri())
            # 验证连接
            self.mongo_client.server_info()
            print("MongoDB连接成功!")
            return True
        except Exception as e:
            print(f"MongoDB连接失败: {e}")
            return False

    def connect_clickhouse(self):
        """连接ClickHouse"""
        try:
            self.ch_client = Client(
                host=self.clickhouse_config['host'],
                port=self.clickhouse_config['port'],
                user=self.clickhouse_config['user'],
                password=self.clickhouse_config['password'],
                database=self.clickhouse_config['database']
            )
            print("ClickHouse连接成功!")
            return True
        except Exception as e:
            print(f"ClickHouse连接失败: {e}")
            return False

    def migrate_data(self):
        """执行数据迁移"""
        try:
            if not self.connect_mongo() or not self.connect_clickhouse():
                return False

            db = self.mongo_client[self.mongo_config['database']]
            collection = db[self.mongo_config['collection']]

            batch_data = []
            processed_count = 0
            error_count = 0
            total_count = collection.count_documents({})
            print(f"开始迁移，总数据量: {total_count}")

            # 先打印一条数据看看结构
            sample_doc = collection.find_one()
            print("示例数据结构:", sample_doc)

            print(str(sample_doc['_id']))

            for doc in collection.find():
                if not doc.get('prompt'):
                    continue

                if not doc.get('prompt').get('id'):
                    continue

                try:
                    # 转换为列表格式
                    row = [
                        str(doc['_id']),  # id
                        datetime.datetime.fromtimestamp(doc['createdAt']),  # created_at
                        doc['prompt']['id'],  # prompt_id
                        str(doc['session']['_id']),  # session_id
                        doc['session'].get('scene', ''),  # scene
                        doc['userId'],  # user_id
                        doc['status']  # status
                    ]

                    batch_data.append(row)
                    processed_count += 1

                    if len(batch_data) >= self.batch_size:
                        try:
                            print('_insert_batch')

                            self._insert_batch(batch_data)
                            print(
                                f"进度: {processed_count}/{total_count} ({(processed_count / total_count * 100):.2f}%)")
                        except Exception as e:
                            print(f"插入批次数据失败: {e}")
                            print("失败的数据示例:", batch_data[0])
                            # 如果批量插入失败，尝试逐条插入
                            for single_row in batch_data:
                                try:
                                    self._insert_single(single_row)
                                except Exception as e:
                                    print(f"单条插入失败: {e}")
                                    print("失败的数据:", single_row)
                                    error_count += 1
                        finally:
                            batch_data = []

                except Exception as e:
                    print(f"处理文档时出错 {doc.get('_id')}: {e}")
                    error_count += 1
                    continue

            # 处理剩余数据
            if batch_data:
                try:
                    self._insert_batch(batch_data)
                    print(f"进度: {processed_count}/{total_count} (100%)")
                except Exception as e:
                    print(f"插入最后批次数据失败: {e}")
                    # 尝试逐条插入剩余数据
                    for single_row in batch_data:
                        try:
                            self._insert_single(single_row)
                        except Exception as e:
                            print(f"单条插入失败: {e}")
                            print("失败的数据:", single_row)
                            error_count += 1

            print(f"数据迁移完成! 总处理: {processed_count}, 错误: {error_count}")
            return True

        except Exception as e:
            print(f"迁移过程中出错: {e}")
            return False

        finally:
            self.close_connections()

    def _insert_batch(self, batch_data):
        """批量插入数据到ClickHouse"""
        self.ch_client.execute(
            '''
            INSERT INTO questions
            (id, created_at, prompt_id, session_id, scene, user_id, status)
            VALUES
            ''',
            batch_data
        )

    def _insert_single(self, row_data):
        """单条插入数据到ClickHouse"""
        self.ch_client.execute(
            '''
            INSERT INTO questions
            (id, created_at, prompt_id, session_id, scene, user_id, status)
            VALUES
            ''',
            [row_data]
        )

    def verify_migration(self):
        """验证迁移结果"""
        try:
            if not self.connect_mongo() or not self.connect_clickhouse():
                return False

            mongo_count = self.mongo_client[self.mongo_config['database']][
                self.mongo_config['collection']].count_documents({})

            # 检查ClickHouse数据
            ch_count = self.ch_client.execute('SELECT count() FROM questions')[0][0]
            print("\n数据验证结果:")
            print(f"MongoDB数据总量: {mongo_count}")
            print(f"ClickHouse数据总量: {ch_count}")
            print(f"数据是否一致: {'✓' if mongo_count == ch_count else '✗'}")

            # 抽样检查一些数据
            sample_data = self.ch_client.execute('SELECT * FROM questions LIMIT 5')
            print("\n抽样数据:")
            for row in sample_data:
                print(row)

            # 检查数据分布
            date_distribution = self.ch_client.execute('''
                SELECT
                    toDate(created_at) as date,
                    count() as count
                FROM questions
                GROUP BY date
                ORDER BY date
                LIMIT 5
            ''')
            print("\n数据时间分布(前5天):")
            for date, count in date_distribution:
                print(f"{date}: {count}条")

            return mongo_count == ch_count

        except Exception as e:
            print(f"验证过程中出错: {e}")
            return False

    def close_connections(self):
        """关闭数据库连接"""
        if self.mongo_client:
            self.mongo_client.close()
        print("数据库连接已关闭")


def main():
    migrator = DatabaseMigrator()
    try:
        if migrator.migrate_data():
            migrator.verify_migration()
    except Exception as e:
        print(f"执行失败: {e}")
    finally:
        migrator.close_connections()


if __name__ == "__main__":
    main()
