# Storefront API 与询盘管理 Issue 草案

> 状态：本地草案，尚未发布到 GitHub。2026-10-12 使用当前授权尝试创建路线图 Task 时收到 `403 Resource not accessible by personal access token`。仓库当前无开放 issue，GitHub 支持 `Task`、`Bug`、`Feature` 类型，但未找到 issue 模板。
>
> 所有 issue 均遵守：`template_id: gas-grill`；目标页面/数据范围写明；扩展所选模板、不创建第二种 visitor-facing 视觉模板；每店独立部署、独立 PostgreSQL；按 `AGENTS.md` / `TESTING.md` 验证。GitHub labels 未探测到/不假设预先存在，因此先不依赖 labels。

## 本地 TODO 记录方式

本文件在仓库内作为 issue 草案与验收依据维护，不依赖 GitHub Issues。当前推进状态以 `docs/superpowers/plans/TODO.md` 为准；具体实现步骤和测试命令分别写在各阶段 plan 中。完成事项必须附实际验证结果后再勾选。

## 路线图父任务

**标题：** Roadmap: API-backed storefront and lightweight inquiry management  
**类型：** Task

### 目标

将当前静态 mock storefront 分阶段演进为每店独立部署、由 Go API + PostgreSQL 提供展示内容并持久化询盘的轻量外贸 MVP。详细基线见 `docs/superpowers/specs/2026-10-12-storefront-api-inquiry-design.md`，阶段路线见 `docs/superpowers/plans/2026-10-12-storefront-backend-roadmap.md`。

### 范围

- `template_id: gas-grill`
- 目标：home、collection、product detail、company content、inquiry form、internal inquiry management。
- 扩展当前模板；不引入第二种视觉模板。
- 每个店铺独立 Go API 与 PostgreSQL。

### 验收

- 分阶段 issue 有明确依赖和独立验收标准。
- 内容来源最终切换至店铺 API/数据库；询盘可安全持久化并由单人/少量用户跟进。
- README 仅在能力实现并验证后更新。
- 不包括 checkout、payment、共享多租户数据库或完整 CRM。

## Issue 1 — API / 数据契约与后端技术决策

**标题：** Define storefront catalog and inquiry API contracts  
**类型：** Task  
**依赖：** 无（本地设计基线已确认）

### 目标与范围

定义 Go 服务与 Next.js storefront 之间稳定、版本化的 catalog/company API 契约，以及 inquiry 创建和内部运营模型。`template_id: gas-grill`；目标数据覆盖首页、分类页、产品详情页、公司页和询盘；扩展现有模板，不引入新模板。

### 要点

- 映射当前 `src/lib/catalog/mock-data.ts` 与 `src/lib/catalog/company-content.ts` 到 API 内容 DTO，确认草稿/已发布、产品/分类关联、媒体资源、规格及可选外贸字段。
- 定义 Inquiry、活动/状态变更、隐私同意、来源字段；明确公开写接口禁止客户端设置运营字段。
- 产品询盘保存产品 ID 与提交时型号/名称快照。
- 给出请求/响应示例、校验、分页/排序（如需要）、未知记录、错误响应及兼容策略。
- 明确示例数据的 placeholder/verified 状态表达，禁止未经证实的认证、测试、工厂、MOQ/交期保证。
- 决定并记录 Go 版本、HTTP 路由方案、SQL 访问方案、迁移工具、PostgreSQL 主版本、本地连接现有 PostgreSQL 的配置、Next.js 取数策略、API 不可用体验和管理端认证边界。
- 附件默认不纳入首批；未来使用独立受限文件存储，不把二进制放进 PostgreSQL。

### 验收标准

- 契约文档描述 catalog/company read DTO 和 inquiry create DTO，含字段类型、必填/可选、长度/范围、枚举与信任边界。
- Inquiry 状态固定为 `new`、`in_progress`、`quoted`、`won`、`lost`、`on_hold`；状态/活动历史独立于访客提交载荷。
- API 版本、错误码类别、未知 slug/空数据行为、日期/时间格式和产品快照规则有明确约定。
- 技术选择有简短理由，符合独立部署、单店 PostgreSQL、轻量依赖和可测试原则。
- 对照当前首页/分类/详情/公司页面，内容契约不缺渲染必需数据。
- README 现状不被误写成后端已实现。

## Issue 2 — Go/PostgreSQL 服务基础

**标题：** Add per-store Go API and PostgreSQL foundation  
**类型：** Feature  
**依赖：** Issue 1

### 目标与范围

建立本地可运行、可测试的 Go API 服务和 PostgreSQL 开发环境，为每店独立部署打基础。`template_id: gas-grill`；此 issue 仅服务端基础，不迁移 storefront 页面，也不实现 CRM UI。

### 验收标准

- Go module、服务入口、明确配置结构、结构化日志和优雅关闭（按选定技术方案）。
- PostgreSQL 可通过本地开发配置启动；无生产密钥硬编码；配置缺失/数据库不可用时有安全清晰的错误。
- 版本化 SQL migration 命令可重复执行且状态可检查；不得依赖启动时无提示 destructive auto-migrate。
- `/health/live` 与数据库就绪检查端点语义清晰；readiness 可报告 DB 不可用但不泄漏凭据。
- 单元测试和 PostgreSQL 集成测试的运行命令写入服务 README 或仓库 README 对应章节。
- Go 格式化、`go test ./...`、`go vet ./...`（或项目选定等价检查）通过并报告实测结果。
- 本地启动步骤可在干净环境按文档复现；店铺数据使用本地 PostgreSQL volume，不采用共享多租户 schema。

## Issue 3 — Catalog/company schema、seed 与只读 API

**标题：** Serve seeded catalog and company content from the Go API  
**类型：** Feature  
**依赖：** Issue 1、Issue 2

### 目标与范围

将当前 gas-grill 示例商品、分类与公司页面内容作为可审查的 seed 数据保存到数据库，并提供公开只读 API。`template_id: gas-grill`；目标数据为 home、collection、product、company content；继续扩展现有模板。

### 验收标准

- migration schema 支持当前 catalog/company 合同所需字段，公开 DTO 不暴露数据库内部细节。
- seed 数据可安全重复执行，不重复产生记录；seed 中事实与占位样例可辨识；不得暗示虚构认证/工厂/测试事实。
- API 至少支持已发布产品列表/详情、分类列表/详情、公司内容；筛选/排序保持有界且明确。
- 草稿不出现在公开 API；未知 slug 返回约定的 404；空分类/空结果可表示；媒体 URL/资源链接遵守安全规则。
- PostgreSQL 集成测试覆盖 seed、关联、published/draft、未知 slug、空数据和响应 DTO。
- OpenAPI 或等价 API 文档包含可执行/准确的请求响应样例。

## Issue 4 — Next.js 页面切换至 Catalog API

**标题：** Read storefront catalog and company pages from the API  
**类型：** Feature  
**依赖：** Issue 1、Issue 3

### 目标与范围

让 Next.js public storefront 的首页、分类、产品详情与公司页从店铺 Go API 读取数据。`template_id: gas-grill`；目标页面为 `/`、`/collections/[slug]`、`/products/[slug]`、`/company`；扩展当前模板。

### 验收标准

- route/page data access 使用独立 catalog API adapter；模板组件只消费共享内容契约，不访问数据库。
- 生产渲染数据由 API 提供；静态 mock 仅可作为 seed/test fixture，不能与 API 并列成为生产事实来源。
- API base URL 由服务端环境配置，不把服务端凭据放到 `NEXT_PUBLIC_*`。
- 定义渲染/缓存/重新验证策略；API 故障与空数据提供明确友好状态，不显示伪造成功数据。
- 未知商品/分类继续按 Next.js not-found 语义；SEO metadata、canonical、sitemap 只使用已发布 API 数据且行为已验证。
- `npm run lint`、`npm run typecheck`、`npm test`、`npm run build` 及相关桌面/移动/键盘检查通过并记录。
- README 更新本地前后端启动、seed、环境变量和能力边界；不得声称 admin localStorage 已迁移为 server publishing。

## Issue 5 — Public Inquiry API 与数据库持久化

**标题：** Persist storefront inquiries with validation and abuse controls  
**类型：** Feature  
**依赖：** Issue 1、Issue 2、Issue 3

### 目标与范围

通过公开 API 接收买家询盘并持久化到当前店铺 PostgreSQL。`template_id: gas-grill`；目标数据包括产品详情询盘及全站询盘；不实现管理 UI，不接受附件（二期另议）。

### 验收标准

- 仅接收已约定访客业务字段：联系人姓名、邮箱、可选公司/电话、国家或地区、产品 ID、提交时产品型号/名称快照（服务端生成）、可选数量、定制留言、隐私同意、来源路径。
- 后端验证邮箱、长度、数量、国家/地区和已发布产品关联；拒绝未知/草稿产品；服务端生成 ID/时间/来源归一值。
- 客户端不能指定状态、负责人、跟进记录或其它运营字段；操作字段由后端信任边界控制。
- 同意文本版本/时间等隐私证据按契约记录；日志不记录不必要的联系方式或完整留言。
- 配置式速率限制、honeypot 或等效基础反垃圾措施；返回安全、稳定的字段错误和服务端错误，不暴露 SQL/堆栈。
- 对重复提交的处理策略明确（例如不承诺 exactly-once，客户端请求 ID/幂等 key 需按契约）；数据库失败不返回假成功。
- 自动测试覆盖有效提交、字段错误、未知产品、client-forged operational fields、反垃圾/限流和数据库错误。
- 公开 API 文档给出状态码、字段及隐私处理说明。

## Issue 6 — Storefront 询盘 CTA 与表单 UX

**标题：** Add product-prefilled inquiry form and submission feedback  
**类型：** Feature  
**依赖：** Issue 4、Issue 5

### 目标与范围

在详情页与全站提供显著询盘入口，提交到持久化 API。`template_id: gas-grill`；目标页面包括产品详情和询盘表单/弹层或专页（由实现计划确定）；扩展当前模板风格。

### 验收标准

- 产品详情 CTA 带入产品 ID 和型号；后端仍重新校验并生成可信快照。
- 表单字段与 API 一致，访客不接触运营状态字段；表单长度与必填策略避免不必要摩擦。
- 清晰隐私用途说明和同意控件；联系字段按最小必要原则收集。
- 有可访问的 loading、成功、字段失败、网络失败状态；保留可修正输入，不重复显示成功或丢失上下文。
- 键盘焦点、标签、错误关联和移动布局经过验证。
- 浏览器自动化/组件测试验证预填、校验、成功失败；`npm run lint`、`npm run typecheck`、`npm test`、`npm run build` 实测通过。

## Issue 7 — 轻量询盘后台

**标题：** Add authenticated lightweight inquiry follow-up management  
**类型：** Feature  
**依赖：** Issue 5（建议在 Issue 6 后交付或并行设计）

### 目标与范围

给店主/少量内部用户提供单店安全后台，不引入完整 CRM。`template_id: gas-grill`；目标页面为管理区询盘列表与详情，不改变访客模板视觉风格。

### 验收标准

- 服务端认证保护所有询盘管理读写端点；前端路由隐藏不是授权措施。
- 登录凭据/会话密钥从安全配置提供，密码安全存储，会话具备适当过期/撤销和 cookie 防护；禁止默认生产密码。
- 可列表、查看详情、按状态/时间/产品做基础搜索筛选，并支持必要分页。
- 可更新状态、记录备注/联系活动、设置最近联系与下一步跟进动作/日期；所有写操作带操作者和时间审计。
- 能查看原始访客提交及产品快照，访客数据与运营数据权限边界明确。
- 测试覆盖未认证拒绝、有效授权、非法状态、审计记录和基本列表筛选。
- 文档说明初始管理员引导、密码轮转、备份/恢复与本地开发；小团队暂不做复杂角色、审批或分配规则。

## 延后 Issue 队列（不要现在创建实现票）

- 附件上传与隔离存储：等待明确存储服务、体积/格式限制、恶意文件扫描和访问策略。
- 可配置运营字段：等待至少一个真实字段需求；保留核心 typed columns，不预建通用对象引擎。
- CSV、报表、评分、自动分配、营销自动化：由真实线索量和流程证据触发。
- 浏览器本地商品 admin 的服务端发布迁移：目录读 API 稳定后单独评估权限/编辑冲突/发布审核。
