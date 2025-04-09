#!/bin/bash

# 登录获取令牌
echo "正在登录..."
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | jq -r '.token')

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
  echo "登录失败，请检查用户名和密码"
  exit 1
fi

echo "登录成功，获取商品列表..."

# 使用令牌获取商品列表
curl -s -X GET http://localhost:8080/api/products \
  -H "Authorization: Bearer $TOKEN" | jq

echo -e "\n操作完成"
