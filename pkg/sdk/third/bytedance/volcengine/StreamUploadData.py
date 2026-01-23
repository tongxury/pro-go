import datetime
import hashlib
import hmac
import os
import zlib
from urllib.parse import quote

import requests

# 以下参数视服务不同而不同，一个服务内通常是一致的
Service = "iccloud_muse"
Version = "2022-02-01"
Region = "cn-north-1"
Host = "icp.volcengineapi.com"
ContentType = "application/json"

# 请求的凭证，从IAM或者STS服务中获取
AK = ""
SK = "=="


def norm_query(params):
    query = ""
    for key in sorted(params.keys()):
        if type(params[key]) == list:
            for k in params[key]:
                if type(k) == int:  # 如果值是整数，直接转换为字符串
                    k = str(k)
                query = (
                        query + quote(key.encode('utf-8'), safe=b"-_.~") + "=" + quote(k.encode('utf-8'), safe=b"-_.~") + "&"
                )
        else:
            if type(params[key]) == int:  # 如果值是整数，直接转换为字符串
                params[key] = str(params[key])
            query = (query + quote(key.encode('utf-8'), safe=b"-_.~") + "=" + quote(params[key].encode('utf-8'), safe=b"-_.~") + "&")
    query = query[:-1]
    return query.replace("+", "%20")


# sha256 非对称加密
def hmac_sha256(key: bytes, content: str):
    return hmac.new(key, content.encode("utf-8"), hashlib.sha256).digest()


# sha256 hash算法
def hash_sha256(content: str):
    return hashlib.sha256(content.encode("utf-8")).hexdigest()


# 签名请求函数
def request(method, date, query, header, ak, sk, action, body):
    # 初始化身份证明结构体
    credential = {
        "access_key_id": ak,
        "secret_access_key": sk,
        "service": Service,
        "region": Region,
    }
    # 初始化签名结构体
    request_param = {
        "body": body,
        "host": Host,
        "path": "/",
        "method": method,
        "content_type": ContentType,
        "date": date,
        "query": {"Action": action, "Version": Version, **query},
    }
    if body is None:
        request_param["body"] = ""
    # 初始化签名结果的结构体
    x_date = request_param["date"].strftime("%Y%m%dT%H%M%SZ")
    short_x_date = x_date[:8]
    if isinstance(request_param["body"], bytes):  # 判断 body 是否为字节类型
        x_content_sha256 = hashlib.sha256(request_param["body"]).hexdigest()  # 处理字节类型的 body
    else:
        x_content_sha256 = hash_sha256(request_param["body"])
    sign_result = {
        "Host": request_param["host"],
        "X-Content-Sha256": x_content_sha256,
        "X-Date": x_date,
        "Content-Type": request_param["content_type"],
    }
    # 计算 Signature 签名
    signed_headers_str = ";".join(
        ["content-type", "host", "x-content-sha256", "x-date"]
    )
    # signed_headers_str = signed_headers_str + ";x-security-token"
    canonical_request_str = "\n".join(
        [request_param["method"].upper(),
         request_param["path"],
         norm_query(request_param["query"]),
         "\n".join(
             [
                 "content-type:" + request_param["content_type"],
                 "host:" + request_param["host"],
                 "x-content-sha256:" + x_content_sha256,
                 "x-date:" + x_date,
             ]
         ),
         "",
         signed_headers_str,
         x_content_sha256,
         ]
    )
    # 打印正规化的请求用于调试比对
    print(canonical_request_str)
    hashed_canonical_request = hash_sha256(canonical_request_str)
    # 打印hash值用于调试比对
    print(hashed_canonical_request)
    credential_scope = "/".join([short_x_date, credential["region"], credential["service"], "request"])
    string_to_sign = "\n".join(["HMAC-SHA256", x_date, credential_scope, hashed_canonical_request])
    # 打印最终计算的签名字符串用于调试比对
    print(string_to_sign)
    k_date = hmac_sha256(credential["secret_access_key"].encode("utf-8"), short_x_date)
    k_region = hmac_sha256(k_date, credential["region"])
    k_service = hmac_sha256(k_region, credential["service"])
    k_signing = hmac_sha256(k_service, "request")
    signature = hmac_sha256(k_signing, string_to_sign).hex()
    sign_result["Authorization"] = "HMAC-SHA256 Credential={}, SignedHeaders={}, Signature={}".format(
        credential["access_key_id"] + "/" + credential_scope,
        signed_headers_str,
        signature,
    )
    header = {**header, **sign_result}
    # header = {**header, **{"X-Security-Token": SessionToken}}
    # 将 Signature 签名写入 HTTP Header 中，并发送 HTTP 请求
    r = requests.request(method=method,
                         url="https://{}{}".format(request_param["host"], request_param["path"]),
                         headers=header,
                         params=request_param["query"],
                         data=request_param["body"],
                         )
    return r.json()


# 图片Size字节查询
def get_image_size_in_bytes(image_path):
    # 获取文件大小
    size_in_bytes = os.path.getsize(image_path)
    return size_in_bytes


#  获取图片CRC
def get_image_crc(image_path):
    with open(image_path, 'rb') as file:
        # 读取图片文件的二进制内容
        data = file.read()
        # 使用zlib计算CRC32值
        crc = zlib.crc32(data)
        # zlib.crc32返回的是一个无符号整数，可能是一个负数（由于Python的整数是补码形式）
        # 为了得到一个总是正数的CRC值（例如用于打印或比较），我们可以将其转换为无符号整数
        crc_unsigned = crc & 0xFFFFFFFF
        return crc_unsigned


# 获取图片MD5
def get_image_md5(image_path):
    # 创建一个md5 hash对象
    md5_hash = hashlib.md5()

    # 打开图片文件，以二进制模式读取
    with open(image_path, 'rb') as file:
        # 读取文件内容并更新到hash对象中
        chunk = file.read(4096)
        while chunk:
            md5_hash.update(chunk)
            chunk = file.read(4096)

    # 获取16进制哈希值
    md5_digest = md5_hash.hexdigest()
    return md5_digest


# 获取图片
image_path = "/Users/bytedance/Desktop/智能创作云/测试素材/图片/图片27.png"  # 替换为你的图片路径
crc = get_image_crc(image_path)
size = get_image_size_in_bytes(image_path)
md5 = get_image_md5(image_path)

# 上传数据，支持并发、幂等

if __name__ == "__main__":
    now = datetime.datetime.utcnow()

    # 以二进制模式打开文件
    file_handle = open(image_path, 'rb')
    # 读取文件的全部内容
    form = file_handle.read()
    # 关闭文件
    file_handle.close()

    query = {
        'Action': "StreamUploadData",
        'Version': "2022-02-01",
        "Md5": md5,
        "Size": size,
        "Offset": 0,
        "OwnerId": 7348381768925446179,
        "OwnerType": "PERSON"
    }
    print(query)

    response_body = request("POST", now, query, {}, AK, SK, "StreamUploadData", body=form)
    print(response_body)