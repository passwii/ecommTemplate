# Storefront Catalog 与 Inquiry API 契约草案

**状态：** 路线已确认；具体实现前可据 issue 做小幅修订  
**日期：** 2026-10-12  
**模板：** `template_id: gas-grill`，扩展现有模板，不创建新模板。  
**部署边界：** 每店一个 Go API、一个 PostgreSQL 数据库；无共享多租户数据库。

## 技术默认值

- Go：优先使用当前受支持的稳定 Go 版本；实施 issue 固定具体最低版本，并在 `go.mod`/CI/README 一致声明。
- HTTP：第一版优先 Go 标准库 `net/http` 与 `http.ServeMux`，只有出现明确路由/中间件需求时才引入框架。
- 数据库：PostgreSQL；访问优先 `database/sql` + PostgreSQL driver 与显式 SQL，避免首版引入重 ORM。
- Schema：采用版本化、显式、可审阅和可重复执行的 SQL migration；禁止把 ORM 自动改 schema 作为生产迁移策略。
- JSON API：统一 `/api/v1` 前缀，JSON 请求/响应，UTC RFC 3339 时间戳。
- 配置：服务端环境配置；数据库 URL、session/CSRF 密钥等不使用 `NEXT_PUBLIC_*`，不提交真实秘密。
- 本地开发：复用用户现有本机 PostgreSQL 服务，使用项目专属数据库 `ecommtemple_dev`；不额外启动第二套 PostgreSQL 容器。部署时每个店铺仍有独立 PostgreSQL。只有经确认本地服务无法满足需要时，才重新评估容器化方案。

## 内容实体

### Collection

- `id`: 稳定 ID
- `slug`: URL 稳定标识
- `name`
- `description`
- `status`: 至少区分 `draft` / `published`（如分类只通过已发布产品可见，可在服务层归一）
- `sort_order`

### Product

- `id`: 稳定 ID
- `slug`: URL 稳定标识
- `name`, `model`/`sku`（已有 sku 可兼容为型号，但实现时需确定语义）
- `collection_id`
- `short_description`, `description`
- `images[]`: `url`, `alt`, `sort_order`
- `specifications[]`: `name`, `value`；结构保留，但需要筛选的稳定规格未来可提升为类型字段
- `applications[]`
- `minimum_order_quantity?`, `lead_time?`
- `resources[]`: `title`, safe `url`, optional `file_type`, `description`
- 当前模板兼容展示字段（fuel、burners、cookingArea、finish、price 等）需在迁移映射表中逐一决定保留、归一或弃用；不得凭空补齐。
- `status`: `draft` / `published`
- `featured`, `sort_order`
- 内容可信状态：样例/占位内容与经核验内容须可区别；不可将占位的规格/认证作为事实输出。

价格字段当前存在于 storefront 类型中，但外贸询盘站是否公开价格尚未确定。初版契约应保留可配置/可为空的价格展示策略，不能假设所有商品都需公开零售价格。

### Company content

- About、manufacturing、quality、credentials、case studies、FAQ 可按当前 `CompanyContent` 类型映射。
- 展示内容使用 `status`/`verification_status` 或等效机制标识 draft/sample/verified。
- Credential 至少记录名称、签发方、详情；只有拥有可验证依据的资料允许公开呈现为认证事实。
- company 页面接口只返回允许公开且已发布的内容。

## Public Read API

建议端点：

- `GET /api/v1/collections`
- `GET /api/v1/collections/{slug}`
- `GET /api/v1/products?collection={slug}&featured=true`
- `GET /api/v1/products/{slug}`
- `GET /api/v1/company-content`
- `GET /health/live`
- `GET /health/ready`

约定：

- 仅公开已发布记录；产品下架或未知 slug 对匿名客户端统一返回 `404`，不泄露草稿存在性。
- 列表默认有界（例如 limit 24，最大 100），稳定排序；分页策略在接口实现前固定为 `limit/offset` 或 cursor，不支持无界拉取。
- 响应包装使用 `data` 与 `meta`；错误使用稳定机器码和可安全展示的人类信息，不包含 SQL、堆栈或内部文件路径。
- 空列表返回成功且 `data: []`；API 故障由 storefront 呈现服务不可用/重试体验，不能偷偷回退到过期静态 mock 并误称为数据库实时内容。
- 缓存策略由 Next.js server-side fetch/route 使用；公开已发布内容可短时缓存或按 revalidation 配置，草稿绝不进入缓存公开响应。

示例：

```json
{
  "data": {
    "id": "grill-01",
    "slug": "ridge-4-burner",
    "name": "Ridge 4-Burner",
    "model": "Ridge 4-Burner",
    "collection": { "slug": "freestanding", "name": "Freestanding grills" },
    "images": [{ "url": "/images/grill-ridge.jpg", "alt": "...", "sortOrder": 0 }],
    "specifications": [{ "name": "Burners", "value": "4" }],
    "applications": [],
    "resources": [],
    "contentStatus": "sample"
  }
}
```

## Inquiry Create API

建议端点：`POST /api/v1/inquiries`。访客提交 DTO：

```json
{
  "contact": {
    "name": "Alex Morgan",
    "email": "buyer@example.com",
    "company": "Example Importers",
    "phone": "+1..."
  },
  "destinationCountry": "US",
  "productId": "grill-01",
  "quantity": 120,
  "customizationRequest": "Private-label packaging options",
  "privacyConsent": {
    "accepted": true,
    "noticeVersion": "inquiry-privacy-v1"
  },
  "sourcePath": "/products/ridge-4-burner"
}
```

服务端行为：

- 从可信数据库记录生成 `productSnapshot: {id, slug, model, name}`；不信任客户端传来的产品名称/型号。
- 生成 Inquiry ID、received time、status=`new`、source canonicalization；运营字段不能从公开 DTO 接收。
- 访客可选字段省略或 `null` 的规则一致；对 Unicode、超长文本、无效邮箱、非正数数量、无效国家代码进行服务端校验。
- 第一阶段不接受附件字段；文件上传需单独端点、存储策略和安全审查。
- 记录隐私声明版本、同意结果、同意时间；不将完整 PII 写入普通应用日志。

成功响应 `201`：

```json
{
  "data": {
    "id": "inq_01J...",
    "status": "received",
    "receivedAt": "2026-10-12T12:00:00Z"
  }
}
```

公开错误类型：`validation_error`（字段错误）、`not_found`（产品不可询盘/未知）、`rate_limited`、`service_unavailable`、`internal_error`。对提交失败不得返回成功；错误内容不得暴露数据库或堆栈细节。

## 运营模型

### Inquiry

- ID、创建/更新时间、联系人姓名/邮箱、公司、电话、目的国家/地区。
- 关联 product ID + 不可变提交快照。
- 数量、定制需求、privacy consent 证据、source path / referrer（按隐私最小化原则）。
- 状态：`new`、`in_progress`、`quoted`、`won`、`lost`、`on_hold`。
- 可选 owner；`last_contacted_at`；`next_action` 与 `next_action_due_at`；可选 `loss_reason`。

### Inquiry activity

- ID、inquiry ID、activity type、内容、发生时间、操作者（系统/内部用户）。
- 状态更改写入活动历史，不能仅覆盖当前 status。
- 备注和联系记录不覆写访客原始提交。

## 安全与运维基线

- 对公开 POST 设置 body 上限、Content-Type 检查、输入校验、IP/网络级或等效速率限制和 honeypot/时间陷阱等基础反垃圾措施；代理头仅信任明确配置的可信代理。
- 生产部署须 HTTPS；管理端认证、防 CSRF（若使用 cookie session）、安全 cookie、密码哈希和管理员初始化流程由管理端 issue 细化。
- SQL 一律使用参数化查询；查询使用 request context/timeouts；数据库 readiness 不泄漏连接字符串。
- 后续生产化 issue 需定义数据库备份/恢复演练、迁移回滚/前滚策略、日志 PII 脱敏和保留期限。

## 明确暂不做

- 附件存储/上传、自动评分、线索分配、邮件序列、公共跨店 CRM、通用自定义对象。
- 通过 JSON 字段承载全部 core schema。
- 未经业务验证就强制要求公司名、电话、数量、预算或 BANT/MEDDIC 信息。
