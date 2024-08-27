# Flux-API

**免费渠道的 Flux 文生图 API，适配 OpenAI 兼容性的接口**

演示地址：[openai-all.com](https://openai-all.com)

## 免费渠道

- [x] [SiliconFlow](https://docs.siliconflow.cn/reference/black-forest-labsflux1-schnell) 限时免费
- [x] [Getimg.AI](https://getimg.ai/pricing) 每个月100次免费

## 运行

```bash
docker run -d -p 8080:8080 --name flux-api k8scat/flux-api:latest
```

## 使用

```bash
curl http://127.0.0.1:8080/v1/images/generations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -d '{
    "model": "$MODEL",
    "prompt": "A cute baby sea otter",
    "n": 1,
    "size": "1024x1024"
  }'
```

### SiliconFlow 配置

$MODEL => 固定值 `siliconflow/FLUX.1-schnell`
$OPENAI_API_KEY => SiliconFlow 的 API Key，获取地址：https://cloud.siliconflow.cn/account/ak

### GetimgAI 配置

$MODEL => 固定值 `getimgai/flux-v1`
$OPENAI_API_KEY => [GetimgAI](https://getimg.ai/text-to-image) 浏览器登录后的完整 Cookie

## 交流群

<img src="https://chat.ggemini.pro/flux-api.jpg" width="240" />

## 开源许可

[MIT](./LICENSE)
