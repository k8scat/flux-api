# Flux-API

**免费渠道的 Flux 文生图 API，同时兼容 OpenAI 的 Image 和 Chat 接口，支持自动翻译 Prompt，可以直接接入 [one-api](https://github.com/songquanpeng/one-api) / [one-hub](https://github.com/MartialBE/one-hub) / [new-api](https://github.com/Calcium-Ion/new-api) 等中转平台**

演示地址：[openai-all.com](https://openai-all.com)

生成图片：

<img src="https://chat.ggemini.pro/a-cute-baby-sea-otter.png" width="240">

## 免费渠道

- [x] [SiliconFlow](https://docs.siliconflow.cn/reference/black-forest-labsflux1-schnell) 限时免费
- [x] [Getimg.AI](https://getimg.ai/pricing) 每个月100次免费

## 运行

```bash
docker run -d -p 8080:8080 --name flux-api k8scat/flux-api:latest
```

**Flux 对中文理解不好，支持开启自动翻译功能：**

```bash
docker run -d -p 8080:8080 --name flux-api \
    -e TRANSLATE_ENABLE=true \
    -e TRANSLATE_API_BASE=https://api.openai-all.com \
    -e TRANSLATE_API_KEY=sk-xxx \
    -e TRANSLATE_MODEL=gpt-4o \
    -e TRANSLATE_PROMPT_TEMPLATE="Translate into English: %s" \
    k8scat/flux-api:latest
```

## 使用说明

### 兼容 Image 接口

```bash
curl http://127.0.0.1:8080/v1/images/generations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer #AUTH" \
  -d '{
    "model": "#MODEL",
    "prompt": "A cute baby sea otter",
    "n": 1,
    "size": "1024x1024"
  }'
```

### 兼容 Chat 接口

```bash
curl http://127.0.0.1:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer #AUTH" \
  -d '{
    "model": "#MODEL",
    "messages": [
      {
        "role": "user",
        "content": "A cute baby sea otter"
      }
    ]
  }'
```

### SiliconFlow 配置

- #MODEL => 固定值 `FLUX.1-schnell`
- #AUTH => SiliconFlow 的 API Key，获取地址：https://cloud.siliconflow.cn/account/ak

### GetimgAI 配置

- #MODEL => 固定值 `flux-v1`
- #AUTH => [GetimgAI](https://getimg.ai/text-to-image) 浏览器登录后的完整 Cookie

## 交流群

<img src="https://chat.ggemini.pro/flux-api.jpg" width="240" />

## 开源许可

[MIT](./LICENSE)
