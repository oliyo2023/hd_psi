# 服装进销存系统 (HD-PSI)

基于 Go + Gin + GORM 开发的服装行业进销存管理系统。

## 项目进度

### ✅ 已完成功能

1. 基础框架搭建
   - Gin Web 框架集成
   - GORM + MySQL 数据库集成
   - 基础项目结构

2. 数据模型设计
   - 商品管理 (`Product`)
   - 供应商管理 (`Supplier`)
   - 库存管理 (`Inventory`, `InventoryTransaction`, `InventoryAlert`, `InventoryThreshold`)
   - 库存盘点 (`InventoryCheck`, `InventoryCheckItem`, `InventoryCheckAdjustment`)
   - 会员管理 (`Member`)
   - 试衣记录 (`FittingRecord`)

3. 工具类开发
   - 商品二维码生成工具 (`utils/qrcode.go`)
   - 基于 SKU + 批次的加密二维码实现

### 🚧 进行中功能

1. 库存盘点模块
   - 盘点单创建
   - 盘点流程管理
   - 差异调整审批

### ⏳ 待开发功能

1. 采购管理
   - 供应商评级体系
   - 采购订单流程
   - 到货验收

2. 销售管理
   - POS 收银集成
   - 移动端扫码
   - 议价管理

3. 会员管理
   - 积分系统
   - 消费分析
   - 智能推荐

4. 前端开发
   - Vue3 + Element Plus 界面
   - 移动端 H5

5. 扩展功能
   - 微信小程序
   - 数据分析看板
   - 智能试衣镜集成

## 技术栈

- 后端：Go 1.21
- Web 框架：Gin 1.9.1
- ORM：GORM 1.25.5
- 数据库：MySQL 8.0

## 项目结构

```
hd_psi/
├── config/         # 配置文件
├── controllers/    # 控制器
├── models/         # 数据模型
├── routes/         # 路由配置
├── utils/          # 工具类
└── main.go         # 程序入口
```

## 开发团队

待补充

## 开发环境搭建

待补充

## 部署指南

待补充