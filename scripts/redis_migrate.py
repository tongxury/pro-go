from redis import StrictRedis as Redis
from redis.exceptions import ConnectionError, ResponseError
from tqdm import tqdm
import logging
import time
from typing import Optional, Dict, Any
import os

# 配置日志
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s',
    filename='redis_migration.log'
)


class RedisMigration:
    def __init__(
            self,
            source_host: str,
            source_port: int,
            source_password: Optional[str],
            source_db: int,
            target_host: str,
            target_port: int,
            target_password: Optional[str],
            target_db: int,
            batch_size: int = 1000
    ):
        """初始化 Redis 迁移工具"""
        try:
            self.source = Redis(
                host=source_host,
                port=source_port,
                username='default',
                password=source_password,
                db=source_db,
                decode_responses=True,
                socket_timeout=5,
                socket_connect_timeout=5
            )

            self.target = Redis(
                host=target_host,
                port=target_port,
                username='default',
                password=target_password,
                db=target_db,
                decode_responses=True,
                socket_timeout=5,
                socket_connect_timeout=5
            )

            self.batch_size = batch_size

        except Exception as e:
            logging.error(f"初始化 Redis 连接失败: {str(e)}")
            raise

    def test_connection(self) -> bool:
        """测试源和目标 Redis 连接"""
        try:
            self.source.ping()
            logging.info("源 Redis 连接成功")
            self.target.ping()
            logging.info("目标 Redis 连接成功")
            return True
        except (ConnectionError, ResponseError) as e:
            logging.error(f"Redis 连接失败: {str(e)}")
            return False
        except Exception as e:
            logging.error(f"未知错误: {str(e)}")
            return False

    def get_key_type(self, key: str) -> str:
        """获取 key 的类型"""
        return self.source.type(key)

    def copy_string(self, key: str) -> bool:
        """复制字符串类型的键值"""
        try:
            value = self.source.get(key)
            ttl = self.source.ttl(key)

            pipeline = self.target.pipeline()
            pipeline.set(key, value)
            if ttl > 0:
                pipeline.expire(key, ttl)
            pipeline.execute()
            return True
        except Exception as e:
            logging.error(f"复制字符串键 {key} 失败: {str(e)}")
            return False

    def copy_hash(self, key: str) -> bool:
        """复制哈希类型的键值"""
        try:
            value = self.source.hgetall(key)
            ttl = self.source.ttl(key)

            pipeline = self.target.pipeline()
            if value:  # 只在有值时执行 hmset
                pipeline.hmset(key, value)
            if ttl > 0:
                pipeline.expire(key, ttl)
            pipeline.execute()
            return True
        except Exception as e:
            logging.error(f"复制哈希键 {key} 失败: {str(e)}")
            return False

    def copy_list(self, key: str) -> bool:
        """复制列表类型的键值"""
        try:
            value = self.source.lrange(key, 0, -1)
            ttl = self.source.ttl(key)

            pipeline = self.target.pipeline()
            if value:
                pipeline.rpush(key, *value)
            if ttl > 0:
                pipeline.expire(key, ttl)
            pipeline.execute()
            return True
        except Exception as e:
            logging.error(f"复制列表键 {key} 失败: {str(e)}")
            return False

    def copy_set(self, key: str) -> bool:
        """复制集合类型的键值"""
        try:
            value = self.source.smembers(key)
            ttl = self.source.ttl(key)

            pipeline = self.target.pipeline()
            if value:
                pipeline.sadd(key, *value)
            if ttl > 0:
                pipeline.expire(key, ttl)
            pipeline.execute()
            return True
        except Exception as e:
            logging.error(f"复制集合键 {key} 失败: {str(e)}")
            return False

    def copy_zset(self, key: str) -> bool:
        """复制有序集合类型的键值"""
        try:
            value = self.source.zrange(key, 0, -1, withscores=True)
            ttl = self.source.ttl(key)

            pipeline = self.target.pipeline()
            if value:
                pipeline.zadd(key, dict(value))
            if ttl > 0:
                pipeline.expire(key, ttl)
            pipeline.execute()
            return True
        except Exception as e:
            logging.error(f"复制有序集合键 {key} 失败: {str(e)}")
            return False

    def migrate_data(self) -> Dict[str, Any]:
        """迁移所有数据"""
        if not self.test_connection():
            return {"success": False, "message": "连接测试失败"}

        start_time = time.time()
        stats = {
            "total": 0,
            "success": 0,
            "failed": 0,
            "types": {}
        }

        try:
            # 获取所有键
            cursor = 0
            all_keys = set()

            logging.info("开始扫描源 Redis 的所有键...")
            while True:
                cursor, keys = self.source.scan(cursor, count=self.batch_size)
                all_keys.update(keys)
                if cursor == 0:
                    break

            total_keys = len(all_keys)
            logging.info(f"共发现 {total_keys} 个键需要迁移")

            # 使用进度条显示迁移进度
            with tqdm(total=total_keys, desc="迁移进度") as pbar:
                for key in all_keys:
                    try:
                        key_type = self.get_key_type(key)
                        stats["types"][key_type] = stats["types"].get(key_type, 0) + 1

                        success = False
                        if key_type == "string":
                            success = self.copy_string(key)
                        elif key_type == "hash":
                            success = self.copy_hash(key)
                        elif key_type == "list":
                            success = self.copy_list(key)
                        elif key_type == "set":
                            success = self.copy_set(key)
                        elif key_type == "zset":
                            success = self.copy_zset(key)
                        elif key_type == "none":
                            logging.warning(f"键 {key} 不存在或已过期")
                            success = True  # 跳过不存在的键

                        if success:
                            stats["success"] += 1
                        else:
                            stats["failed"] += 1

                        stats["total"] += 1
                        pbar.update(1)

                        if stats["total"] % 1000 == 0:
                            logging.info(f"已处理 {stats['total']}/{total_keys} 个键")

                    except Exception as e:
                        logging.error(f"处理键 {key} 时发生错误: {str(e)}")
                        stats["failed"] += 1
                        continue

        except Exception as e:
            logging.error(f"迁移过程发生错误: {str(e)}")
            return {
                "success": False,
                "message": str(e),
                "stats": stats
            }
        finally:
            try:
                self.source.close()
                self.target.close()
            except:
                pass

        end_time = time.time()
        duration = end_time - start_time

        result = {
            "success": True,
            "duration": f"{duration:.2f}秒",
            "stats": stats
        }

        logging.info(f"迁移完成，耗时: {duration:.2f}秒")
        logging.info(f"统计信息: {stats}")

        return result


def main():
    try:
        # Redis 配置
        source_config = {
            "host": "118.196.63.209",  # 源 Redis 主机地址
            "port": 16379,  # 源 Redis 端口
            "password": os.getenv("REDIS_PASSWORD"),  # 源 Redis 密码
            "db": 0  # 源 Redis 数据库编号
        }

        target_config = {
            "host": "101.132.192.41",  # 目标 Redis 主机地址
            "port": 6379,  # 目标 Redis 端口
            "password": os.getenv("REDIS_PASSWORD"),  # 目标 Redis 密码
            "db": 0  # 目标 Redis 数据库编号
        }

        print("开始 Redis 数据迁移...")
        print(f"源 Redis: {source_config['host']}:{source_config['port']}")
        print(f"目标 Redis: {target_config['host']}:{target_config['port']}")

        # 创建迁移工具实例
        migrator = RedisMigration(
            source_host=source_config["host"],
            source_port=source_config["port"],
            source_password=source_config["password"],
            source_db=source_config["db"],
            target_host=target_config["host"],
            target_port=target_config["port"],
            target_password=target_config["password"],
            target_db=target_config["db"],
            batch_size=1000
        )

        # 执行迁移
        result = migrator.migrate_data()

        # 打印结果
        if result["success"]:
            print("\n迁移成功！")
            print(f"总耗时: {result['duration']}")
            print(f"总键数: {result['stats']['total']}")
            print(f"成功: {result['stats']['success']}")
            print(f"失败: {result['stats']['failed']}")
            print("\n各类型键统计:")
            for key_type, count in result['stats']['types'].items():
                print(f"  {key_type}: {count}")
        else:
            print(f"\n迁移失败: {result['message']}")

    except KeyboardInterrupt:
        print("\n用户中断操作，正在退出...")
    except Exception as e:
        print(f"\n发生错误: {str(e)}")
        logging.error(f"程序执行错误: {str(e)}")
    finally:
        print("\n程序执行完成，详细日志请查看 redis_migration.log")


if __name__ == "__main__":
    main()