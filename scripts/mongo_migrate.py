from urllib.parse import quote_plus

from pymongo import MongoClient
from tqdm import tqdm
import logging
import time
import os

# 配置日志
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s',
    filename='mongodb_migration.log'
)


def connect_mongodb(uri, db_name):
    """连接MongoDB数据库"""
    try:
        client = MongoClient(uri)
        db = client[db_name]
        return client, db
    except Exception as e:
        logging.error(f"数据库连接失败: {str(e)}")
        raise


def copy_database(source_uri, source_db_name, target_uri, target_db_name, batch_size=1000):
    source_client = None
    target_client = None
    try:
        logging.info(f"开始连接源数据库: {source_db_name}")
        source_client = MongoClient(source_uri)
        source_db = source_client[source_db_name]
        logging.info("源数据库连接成功")

        logging.info(f"开始连接目标数据库: {target_db_name}")
        target_client = MongoClient(target_uri)
        target_db = target_client[target_db_name]
        logging.info("目标数据库连接成功")

        collections = source_db.list_collection_names()
        logging.info(f"发现 {len(collections)} 个集合需要复制: {', '.join(collections)}")

        # 创建会话
        with source_client.start_session() as session:
            for collection_name in collections:
                try:
                    logging.info(f"\n{'=' * 50}")
                    logging.info(f"开始处理集合: {collection_name}")

                    source_collection = source_db[collection_name]
                    target_collection = target_db[collection_name]

                    total_docs = source_collection.count_documents({})
                    logging.info(f"集合 {collection_name} 包含 {total_docs} 个文档")

                    # 索引复制
                    logging.info(f"开始复制集合 {collection_name} 的索引")
                    indexes = list(source_collection.list_indexes())
                    logging.info(f"发现 {len(indexes)} 个索引需要复制")

                    for index in indexes:
                        if index['name'] != '_id_':
                            try:
                                index_keys = [(k, v) for k, v in index['key'].items()]
                                index_options = {'name': index['name']}
                                if 'unique' in index:
                                    index_options['unique'] = index['unique']
                                if 'sparse' in index:
                                    index_options['sparse'] = index['sparse']
                                if 'expireAfterSeconds' in index:
                                    index_options['expireAfterSeconds'] = index['expireAfterSeconds']

                                target_collection.create_index(index_keys, **index_options)
                                logging.info(f"索引 {index['name']} 创建成功")
                            except Exception as e:
                                logging.warning(f"创建索引失败 {index['name']}: {str(e)}")

                    # 使用会话进行数据复制
                    logging.info(f"开始复制集合 {collection_name} 的数据")
                    copied_docs = 0
                    error_count = 0

                    with tqdm(total=total_docs, desc=f"复制 {collection_name}") as pbar:
                        # 使用会话查询数据
                        cursor = source_collection.find(
                            {},
                            session=session,
                            no_cursor_timeout=True,
                            batch_size=batch_size
                        )

                        batch = []
                        try:
                            for doc in cursor:
                                batch.append(doc)
                                if len(batch) >= batch_size:
                                    try:
                                        target_collection.insert_many(batch, ordered=False)
                                        copied_docs += len(batch)
                                        pbar.update(len(batch))
                                        if copied_docs % (batch_size * 10) == 0:
                                            logging.info(f"已复制 {copied_docs}/{total_docs} 个文档")
                                        batch = []
                                    except Exception as e:
                                        error_count += 1
                                        logging.error(f"插入批次失败 (第 {error_count} 次错误): {str(e)}")
                                        if error_count >= 5:
                                            logging.error(f"批次大小: {len(batch)}, 已复制文档数: {copied_docs}")
                                        batch = []
                                        continue

                            # 处理最后一批
                            if batch:
                                try:
                                    target_collection.insert_many(batch, ordered=False)
                                    copied_docs += len(batch)
                                    pbar.update(len(batch))
                                    logging.info(f"最后一批 {len(batch)} 个文档复制完成")
                                except Exception as e:
                                    logging.error(f"插入最后一批文档失败: {str(e)}")

                        finally:
                            cursor.close()

                    # 验证文档数量
                    source_count = source_collection.count_documents({})
                    target_count = target_collection.count_documents({})

                    if source_count == target_count:
                        logging.info(f"集合 {collection_name} 复制完成，文档数量匹配: {target_count}")
                    else:
                        logging.warning(
                            f"集合 {collection_name} 文档数量不匹配! "
                            f"源: {source_count}, 目标: {target_count}, "
                            f"差异: {abs(source_count - target_count)}"
                        )

                except Exception as e:
                    logging.error(f"处理集合 {collection_name} 时发生错误: {str(e)}")
                    continue

    except Exception as e:
        logging.error(f"数据库复制过程发生错误: {str(e)}")
        raise

    finally:
        if source_client:
            source_client.close()
            logging.info("源数据库连接已关闭")
        if target_client:
            target_client.close()
            logging.info("目标数据库连接已关闭")

def main():
    # 配置参数
    SOURCE_URI = os.getenv("MONGO_SOURCE_URI")
    TARGET_URI = os.getenv("MONGO_TARGET_URI")
    SOURCE_DB = "yoozy_pro"
    TARGET_DB = "admin"
    BATCH_SIZE = 1000

    start_time = time.time()
    logging.info(f"开始数据迁移任务")
    logging.info(f"源数据库: {SOURCE_DB}")
    logging.info(f"目标数据库: {TARGET_DB}")
    logging.info(f"批处理大小: {BATCH_SIZE}")

    try:
        copy_database(
            source_uri=SOURCE_URI,
            source_db_name=SOURCE_DB,
            target_uri=TARGET_URI,
            target_db_name=TARGET_DB,
            batch_size=BATCH_SIZE
        )

        end_time = time.time()
        duration = end_time - start_time
        logging.info(f"数据迁移完成，总耗时: {duration:.2f} 秒")
        logging.info(f"平均处理速度: {duration / 60:.2f} 分钟/GB")

    except Exception as e:
        logging.error(f"迁移过程失败: {str(e)}")

if __name__ == "__main__":
    main()