import * as Updates from "expo-updates";

import Constants from "expo-constants";

//用js的好处是 更改不需要重新打包
export const getConfig = () => {

  console.log('env', Constants.expoConfig?.extra?.["ENV"])
  const productionConfig = {
    ENV: "production",
    API_URL: "https://y.yoozyai.com",

    TOS_REGION: "cn-shanghai",
    TOS_ENDPOINT: "tos-cn-shanghai.volces.com",
    TOS_BUCKET: "yoozyres",

    STORAGE: "oss",

    OSS_BUCKET: "yoozy",
    OSS_REGION: "oss-cn-hangzhou",
    OSS_ENDPOINT: "https://oss-cn-hangzhou.aliyuncs.com",

    QINIU_ENDPOINT: "https://res.veogo.cn",
    QINIU_BUCKET: "veogo",

    REGION: "cn-north-1",
    BUCKET: "veogoresources",
    ENDPOINT: "https://veogoresources.s3.cn-north-1.amazonaws.com.cn",
  };

  const configs = {
    development: {
      ...productionConfig,
      // ENV: "development",
      // API_URL: "https://i-beta.yoozyai.com",
    },
    preview: productionConfig,
    production: productionConfig,
  } as any;

  if (Updates.channel) {
    return configs[Updates.channel] || configs["production"];
  } else {
    return configs[Constants.expoConfig?.extra?.["ENV"] || "production"];
  }
};
