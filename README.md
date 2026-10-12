# Gas Grill Storefront MVP

Next.js 15 / React 19 单模板独立站 storefront MVP，以户外燃气烤炉为示例。当前交付的是用于展示商品与分类的前端，不包含购物车结算（checkout）或支付能力。

当前 storefront 使用的模板标识为 `template_id: gas-grill`。模板面向访客的页面由 Next.js 渲染；商品、分类与公司内容在运行时从本店铺自己的 Go API 读取（见 `CATALOG_API_URL`）。

## 环境要求

- Node.js 18.18+、19.8+ 或 20+（与 lockfile 中 Next.js 的 engines 要求一致）；建议使用仍受支持的 Node.js LTS 版本。
- npm。仓库提交了 `package-lock.json`，安装依赖时使用 `npm ci`。
- 公开页面需要能访问本店铺的 Go API（见 `backend/README.md`）。`npm run lint`、`npm run typecheck`、`npm test` 与 `npm run build` 都不需要 API 在运行。

## 环境变量

- `CATALOG_API_URL`：本店铺 Go API 的 origin（例如 `http://127.0.0.1:8080`）。仅服务端读取，不要使用 `NEXT_PUBLIC_*` 前缀，也不要放入服务端密钥。
- `NEXT_PUBLIC_SITE_URL`：正式 storefront 的站点 origin，用于 canonical、OG 与 sitemap（见下节）。

未配置 `CATALOG_API_URL` 时，公开页面会显示明确的“暂时不可用”状态，**不会**回退到静态样例数据。

## 站点域名与 SEO

设置 `NEXT_PUBLIC_SITE_URL` 为正式 storefront 的站点 origin（例如 `https://shop.example.com`），不要填写含路由的基础路径。正式域名确认前可以留空：此时页面保留 title/description 和 OG/Twitter 文案，但不输出 canonical 或绝对社交图片 URL，sitemap 为空，避免把演示域名或 localhost 发布成索引地址。配置 origin 后，首页、已发布商品/分类页生成 canonical 与社交预览；`/sitemap.xml` 只列首页以及至少有一个已发布商品的分类和已发布商品，`/robots.txt` 提供 sitemap 地址并禁止抓取管理页。草稿与未知 slug 不进入公开索引。

## 商品数据来源与发布约定

公开商品页、分类页、公司页与 SEO 清单在运行时读取本店铺 Go API 返回的**已发布**记录，并使用短时缓存（60 秒）刷新；API 不可用时显示明确的不可用状态，不会用静态样例数据冒充实时内容。

`src/lib/catalog/mock-data.ts` 现仅作为 Go 后端的样例种子来源参考与本地管理预览使用，不再是公开页面的运行时数据源。新增或修改公开商品应写入店铺数据库（参见 `backend/README.md` 的 seed 与 migration 说明）。`/admin/products` 与 `/admin/company` 仅供当前浏览器内编辑和预览：其 localStorage 数据不会被服务端、其他访客、路由、metadata 或 sitemap 读取，页面中的“published”只表示本地预览状态，不是线上发布信号。请勿将此 MVP 的本地保存描述为生产级后台发布能力。

## 本地开发

在项目根目录执行：

```bash
npm ci
npm run dev
```

然后打开终端输出中显示的本地地址（通常为 <http://localhost:3000>）。开发服务器支持热更新。

## 校验与生产运行

以下脚本均来自根目录 `package.json`：

```bash
npm run typecheck  # TypeScript 类型检查：tsc --noEmit
npm run lint       # ESLint：eslint .
npm run build      # Next.js 生产构建：next build
npm run start      # 启动生产服务器：next start（先运行 build）
```

## 页面路由

- `/`：首页，展示精选商品与分类入口（来自 Catalog API）。
- `/collections/[slug]`：商品分类页；分类由店铺数据库提供，样例种子包含 `built-in`、`freestanding` 和 `compact`。
- `/products/[slug]`：商品详情页；按已发布商品的 slug 匹配，未知 slug 返回 404。
- `/company`：公司/信任信息页，内容来自 Catalog API 的 `company-content`。
- `/inquiries`：询盘/报价请求页；`?product=<slug>` 预填产品。
- `/api/inquiries`：同源询盘提交代理，转发到 Go API（浏览器不直接访问后端）。

## 商品、分类与模板入口

- `src/lib/catalog/api.ts`：服务端 Catalog API 客户端，读取 `CATALOG_API_URL`，返回 `not_found` / `unavailable` 结果而不是抛错或回退。
- `src/lib/catalog/catalog-view.mjs`：把 API 载荷映射为模板视图模型的纯函数（有单元测试）。
- `src/lib/catalog/types.ts`：模板共享的内容类型；`src/lib/catalog/mock-data.ts` 现仅作为样例种子参考与本地管理预览。
- `src/lib/templates/registry.ts`：模板注册表、模板标识和默认模板。首版 `DEFAULT_TEMPLATE` 为 `gas-grill`。
- `src/app/`：首页、分类页、商品详情页、公司页等路由入口。
- `src/components/templates/gas-grill/`：当前 gas-grill 模板专属的页面组件，只消费视图模型，不直接访问 API 或数据库。

新增商品或分类应写入店铺数据库；模板组件与共享契约保持一致。首版只实现并使用 `gas-grill` 模板。

## 当前能力边界

- 后端基础已新增于 `backend/`：Go API、本机 PostgreSQL 项目专属数据库和显式 migration，提供 `/health/live`、`/health/ready`，以及只读的 `/api/v1/collections`、`/api/v1/products`、`/api/v1/company-content` 端点。本地启动细节见 [`backend/README.md`](./backend/README.md)。
- 样例商品/分类/公司内容已可通过 `go run ./cmd/storefront-seed` 写入店铺数据库并由上述 API 读取；这些记录标记为 `sample`，不是已核验的真实产品或公司事实。
- 公开询盘接口 `POST /api/v1/inquiries` 已实现：字段校验、产品快照、隐私同意记录、honeypot、按客户端限流、可选幂等键，并持久化到店铺数据库（询盘与 `created` 活动在同一事务内写入）。
- 首页、分类页、商品详情页与公司页已在运行时从 Catalog API 读取数据（60 秒缓存）；API 不可用时显示明确的不可用状态，不回退到样例数据。未知商品/分类返回 404。
- 询盘表单已上线于 `/inquiries`：产品详情页与全站 CTA 均指向该页，可按 `?product=<slug>` 预填产品；表单包含校验、提交中/成功/失败状态、隐私说明与隐藏 honeypot。浏览器提交到同源 `/api/inquiries`，由服务端代理到 Go API。
- 尚未实现：询盘管理后台（认证、列表、状态更新与跟进）。
- 管理员商品/公司页面写入浏览器 localStorage 的内容仅用于当前浏览器预览，不会发布到共享 storefront、创建公开商品 URL 或进入 sitemap。
- 每个店铺独立部署与独立 PostgreSQL 数据库已由 `backend/` 的部署边界体现；本地开发复用本机 PostgreSQL 的项目专属库。
- 当前不提供 checkout、订单处理或支付功能，也未集成第三方电商插件。
