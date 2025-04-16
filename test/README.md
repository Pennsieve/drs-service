# DRS 服务测试目录

本目录包含用于测试 GA4GH DRS API 实现的测试文件。测试分为单元测试和集成测试两部分。

## 测试结构

- `service/` - 服务层测试
  - `drs_service_test.go` - DRS 服务单元测试
  - `drs_service_integration_test.go` - DRS 服务集成测试

## 运行测试

### 运行所有测试

```bash
cd /Users/daniellong/LabProject/local-drs-service
go test ./test/...
```

### 运行特定测试

```bash
# 运行服务层测试
go test ./test/service

# 运行特定测试函数
go test ./test/service -run TestGetObject

# 运行测试并显示详细输出
go test -v ./test/service
```

## 测试数据

测试使用 `MockDrsRepository` 提供的测试数据，包括：

1. `test-object-1` - 基本对象，不需要授权
2. `test-object-2` - 需要 Passport 授权的对象
3. `test-object-3` - 模拟延迟访问 URL 的对象

## 测试覆盖的 API 端点

- `/service-info` - 获取服务元数据
- `/objects/{object_id}` - 获取对象信息
- `/objects/{object_id}/access/{access_id}` - 获取对象访问 URL
- `/objects` (POST) - 批量获取对象信息
