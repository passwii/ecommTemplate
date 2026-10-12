# Storefront MVP 本地开发任务板

以此文件作为 GitHub Issues 的本地替代。每项完成后更新勾选状态，并在对应 PR/提交中保留验证证据。整体基线见：

- 设计：`docs/superpowers/specs/2026-10-12-storefront-api-inquiry-design.md`
- API 契约：`docs/superpowers/specs/2026-10-12-storefront-api-contract.md`
- 总路线：`docs/superpowers/plans/2026-10-12-storefront-backend-roadmap.md`
- 阶段事项与验收：`docs/superpowers/plans/2026-10-12-storefront-issues.md`
- 当前可执行计划：`docs/superpowers/plans/2026-10-12-go-postgres-foundation.md`

## 约定

- 所有 storefront 交付标注 `template_id: gas-grill`，明确页面/数据范围，并说明扩展现有模板、不新增模板。
- 每个店铺独立部署、使用独立 PostgreSQL；本地开发复用本机 PG，但只使用项目专属数据库 `ecommtemple_dev`。
- 不读取、覆盖、清空或迁移本机其他数据库；不提交 `.env`、凭据或数据库 dump。
- TODO 状态只在真实实现和验证后更新；README 只描述已实现并验证的能力。
- 应用实现前先确认在非 `main` 工作分支上；不要未经明确同意直接在 `main` 写应用代码。

## 当前状态

### 完成

- [x] 核对 `main` 与 `origin/main` 并 fast-forward pull；当时已同步。
- [x] 记录总体设计基线。
- [x] 记录 API / 数据模型契约草案。
- [x] 结合当前 catalog/company/admin 代码梳理分阶段路线。
- [x] 按“本地 TODO”方式保留原拟 GitHub issue 的范围和验收内容。
- [x] 确认本机 Go `1.27.1`、PostgreSQL client `18.6`，本机 `localhost:5432` 接受连接（未读取业务库数据）。

### 完成

- [x] 建立 `backend/` Go module 与环境配置校验。
- [x] 增加标准库 HTTP 服务骨架、liveness/readiness health endpoints。
- [x] 增加 pgx `database/sql` PostgreSQL 连接池与凭据脱敏错误处理。
- [x] 增加 Goose 显式 migration 命令及 `gas-grill` 初始 metadata migration。
- [x] 仅创建本机专属 `ecommtemple_dev` 空数据库并应用初始 migration；未读取或更改其他库。
- [x] 更新 root/backend README，明确后端基础已存在、catalog 和询盘能力仍未实现。
- [x] 在隔离分支 `feat/backend-foundation` 工作；原 `main` 仅有 worktree ignore 规则提交。

### 完成

- [x] Catalog/company schema migration（`00002`）：collections、products、company_content。
- [x] `catalog` 包：领域模型、安全资源 URL 规则、仓储接口与 PostgreSQL 只读实现。
- [x] 只读 API：`/api/v1/collections`、`/api/v1/collections/{slug}`、`/api/v1/products`、`/api/v1/products/{slug}`、`/api/v1/company-content`。
- [x] 幂等 seed 命令 `cmd/storefront-seed`，写入标记为 `sample` 的示例数据。
- [x] 共享本地数据库守卫 `database.ValidateLocalDatabaseURL`，迁移与 seed 均限制 loopback + `ecommtemple_dev`。
- [x] README（root/backend）更新 catalog 端点、seed 用法与能力边界。

### 验证记录（Catalog 阶段实测）

- backend：`gofmt`、`go build ./...`、`go test ./...`、`go vet ./...` 全部通过。
- catalog 集成测试在 `ecommtemple_dev` 事务内运行并回滚（草稿排除、未知 slug、过滤/limit、不安全 URL 丢弃、seed 幂等）。
- API 冒烟：collections/products/company-content 返回 seed 数据；未知 slug → 404；`limit=500` → 400。
- `go run ./cmd/storefront-seed` 连续执行两次无重复数据。

### 完成

- [x] 询盘 schema migration（`00003`）：`inquiries`、`inquiry_activities`，含状态约束与活动类型约束。
- [x] `inquiry` 包：字段校验、限流器、服务层（幂等 + 产品快照）、PostgreSQL 仓储（询盘与 `created` 活动同一事务）。
- [x] 公开接口 `POST /api/v1/inquiries`：严格 JSON 解码、64KB 体积限制、Content-Type 校验、honeypot、按客户端限流、可选 `requestId` 幂等、错误映射不泄漏内部细节。
- [x] 运营状态仅由服务端设置；访客伪造运营字段被拒绝。
- [x] 配置新增 `INQUIRY_RATE_LIMIT` / `INQUIRY_RATE_WINDOW`（默认 5 次 / 10 分钟）。
- [x] README（root/backend）与 API 契约记录询盘端点、校验规则与能力边界。

### 验证记录（询盘阶段实测）

- backend：`gofmt -l`（无输出）、`go build ./...`、`go test ./...`、`go vet ./...` 全部通过。
- 询盘集成测试在独立 schema 内运行，覆盖持久化、`created` 活动、约束失败回滚、幂等查询、未知 ID。
- API 冒烟：有效提交 201 且落库（含产品快照 `ridge-4-burner`、状态 `new`、国家 US、数量 120、`created` 活动）；幂等重放 200 且不重复建单；未知产品 404；伪造 `status` 400；honeypot 400；缺少同意 400（带字段）；非 JSON 415；超过限流 429。
- 冒烟产生的测试询盘行已清理，`inquiries` 计数为 0。

### 进行中

- [ ] Next.js 页面接入 Catalog API（前端仍读静态 mock）。
- [ ] 询盘表单（CTA、产品预填、提交反馈）。
- [ ] 询盘管理后台（认证、列表/详情、状态与跟进）。

### 验证记录（本阶段实测）

- backend：`gofmt`、`go test ./...`、`go vet ./...` 全部通过。
- 迁移：`storefront-migrate status` 显示 `pending`，`up` 应用 `00001`，再次 `status` 显示 `applied`；重复 `up` 报 `database schema already current`。
- API 冒烟：`/health/live` 与 `/health/ready` 对专属库均返回 200，响应体无凭据。
- frontend：`npm test`（3 pass）、`npm run lint`、`npm run typecheck`、`npm run build` 均通过。
- 安全检查：`git diff --check` 干净；无跟踪的 `.env`/dump 文件；`secret` 匹配均为测试占位字符串。

### 已知待办（不阻塞本阶段）

- migration runner 的失败回滚测试尚未覆盖（需一次性测试库）。
- Goose 目前对 SQL migration 无内置 checksum 校验；如需防篡改需后续方案。
- 干净克隆/新机器复现步骤尚未在另一环境实测。


### 队列

- [ ] 定稿 catalog/company/inquiry API 契约与端点 DTO。
- [ ] Go 服务基础：配置、`net/http`、PostgreSQL pool、health/readiness、显式 SQL migrations。
- [ ] Catalog/company schema、可审阅 seed 数据、公开只读 API。
- [ ] Next.js 首页、分类页、商品详情和公司页接入 Catalog API。
- [ ] 询盘公开 API：字段校验、来源/同意记录、反垃圾和数据库持久化。
- [ ] 店铺询盘表单：CTA、产品预填、隐私说明、成功/失败状态。
- [ ] 轻量管理端：认证、询盘列表/详情、状态、活动备注、下一步跟进。
- [ ] 按真实需求评估附件、可配置字段、CSV/报表；当前不预建通用 CRM 平台。

## 重要安全提醒

本机 PostgreSQL 当前只确认服务在运行，尚未确认项目数据库是否存在、当前角色权限或实际 server version。不要假设存在可用用户；不要打印 connection string；不要运行宽泛数据库枚举、drop/reset 或对现有库执行 migration。