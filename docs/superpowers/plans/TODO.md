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

### 完成

- [x] 共享内容类型抽取到 `src/lib/catalog/types.ts`。
- [x] `catalog-view.mjs` 纯映射函数 + 单元测试（排序、缺省字段、空图片过滤、价格判定）。
- [x] 服务端 API 客户端 `src/lib/catalog/api.ts`：读取 `CATALOG_API_URL`，区分 `not_found` 与 `unavailable`，60 秒缓存，不回退到样例数据。
- [x] 首页、分类页、商品详情页、公司页改为运行时读取 Catalog API。
- [x] 不可用状态组件；未知商品/分类返回 404；公司页区分 `sample` 内容。
- [x] sitemap/robots 改为从 API 取已发布数据，并设 1 小时 revalidate。
- [x] 模板组件只消费视图模型（props），不再直接读取 mock 数据。

### 验证记录（前端接入阶段实测）

- `npm run lint`、`npm run typecheck`、`npm test`（9 pass）、`npm run build` 均通过且无警告。
- 配置 API 构建并启动后：首页/分类页/详情页渲染 seed 数据（Ridge 4-Burner、748 sq in、Brushed stainless steel 等）；`/company` 显示 sample 提示；未知商品 404；sitemap 列出 1 首页 + 3 分类 + 4 商品。
- 未配置 API 构建并启动后：构建仍成功；首页与公司页显示“暂时不可用”，页面中不含任何样例商品（计数 0）；sitemap 仅剩首页。

### 完成

- [x] Go API 可信代理支持：`TRUST_PROXY_HEADER`，经代理时限流按真实访客而非代理计数（含单元测试）。
- [x] honeypot 字段改为 `referenceCode`，避免浏览器自动填充误伤真实访客。
- [x] 同源代理路由 `src/app/api/inquiries/route.ts`：Content-Type 与体积校验，转发 `X-Forwarded-For`，不回传内部错误。
- [x] `InquiryForm` 客户端组件：必填校验、提交中/成功/失败状态、字段级错误摘要（role=alert + 聚焦）、隐私说明与同意、隐藏 honeypot、键盘与标签可访问性。
- [x] `/inquiries` 页面：`?product=<slug>` 预填产品并显示产品名。
- [x] 详情页 CTA 与站点头部/页脚 CTA 改为指向 `/inquiries`。

### 验证记录（询盘表单阶段实测）

- backend：`gofmt -l`（无输出）、`go build ./...`、`go test ./...`、`go vet ./...` 通过。
- frontend：`npm run lint`、`npm run typecheck`、`npm test`、`npm run build` 通过。
- 端到端（Next 代理 → Go API → PostgreSQL）：`/inquiries` 渲染表单；`?product=ridge-4-burner` 显示 "Enquiry about"；提交返回 201 并落库（`ridge-4-burner` / DE / 50 / `new`）；字段校验返回 400 + `fields`；同一 `X-Forwarded-For` 第 3 次提交返回 429。
- 冒烟产生的询盘行已清理，`inquiries` 计数为 0。

### 完成

- [x] 管理端 schema migration（`00004`）：`admin_users`、`admin_sessions`（令牌仅存哈希）。
- [x] `admin` 包：bcrypt 密码哈希与强度校验、不透明会话令牌、登录/校验/登出服务。
- [x] 管理端 API：会话端点 + 询盘列表/详情/更新/活动，全部经服务端鉴权中间件保护。
- [x] `cmd/storefront-admin create-admin`：首个管理员引导；默认仅限本地库，生产需 `ALLOW_REMOTE_DATABASE=true` 显式开启；拒绝弱密码与重复邮箱。
- [x] inquiry 包管理能力：List（状态/关键词筛选 + 总数）、Update（状态变更写入活动历史）、AddActivity、ListActivities。
- [x] Next.js 管理界面：`/admin/login`、`/admin/inquiries`（筛选/分页）、`/admin/inquiries/[id]`（详情、跟进、活动时间线、登出）。
- [x] 同源代理 `/api/admin/*`：仅允许 GET/POST/PATCH/DELETE，拒绝路径穿越，转发会话 cookie。
- [x] `SESSION_COOKIE_SECURE` 配置；cookie 为 HttpOnly + SameSite=Lax。

### 验证记录（管理后台阶段实测）

- backend：`gofmt -l`（无输出）、`go build ./...`、`go test ./...`、`go vet ./...` 通过。
- frontend：`npm run lint`、`npm run typecheck`、`npm test`、`npm run build` 通过。
- 端到端：未登录访问列表 307 → `/admin/login`；登录 200 且下发 HttpOnly cookie；错误密码 401；伪造 cookie 401；登出后 401；公开询盘 201 后出现在列表；状态更新 200 并生成 "Status changed from new to quoted" 历史；活动记录 201；详情页显示新状态与历史。
- 管理 CLI：创建成功；重复邮箱、弱密码、未开启远程开关的目标库均被拒绝。
- 冒烟产生的询盘与管理员账号已清理（`inquiries=0 admins=0`）。

### 后续（按真实需求再排）

- [ ] 附件上传与隔离存储。
- [ ] 可配置运营字段（出现真实需求后）。
- [ ] CSV/报表、密码轮转与备份恢复演练文档。

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