# Gas Grill Storefront MVP

Next.js 15 / React 19 单模板独立站 storefront MVP，以户外燃气烤炉为示例。当前交付的是用于展示商品与分类的前端，不包含购物车结算（checkout）或支付能力。

当前 storefront 使用的模板标识为 `template_id: gas-grill`。模板面向访客的页面由 Next.js 渲染；商品和分类内容目前来自仓库内的静态 mock 数据。

## 环境要求

- Node.js 18.18+、19.8+ 或 20+（与 lockfile 中 Next.js 的 engines 要求一致）；建议使用仍受支持的 Node.js LTS 版本。
- npm。仓库提交了 `package-lock.json`，安装依赖时使用 `npm ci`。
- 当前本地展示不需要环境变量，也不依赖 Go API 或数据库。

## 站点域名与 SEO

设置 `NEXT_PUBLIC_SITE_URL` 为正式 storefront 的站点 origin（例如 `https://shop.example.com`），不要填写含路由的基础路径。正式域名确认前可以留空：此时页面保留 title/description 和 OG/Twitter 文案，但不输出 canonical 或绝对社交图片 URL，sitemap 为空，避免把演示域名或 localhost 发布成索引地址。配置 origin 后，首页、已发布商品/分类页生成 canonical 与社交预览；`/sitemap.xml` 只列首页以及至少有一个已发布商品的分类和已发布商品，`/robots.txt` 提供 sitemap 地址并禁止抓取管理页。草稿与未知 slug 不进入公开索引。

## 商品发布数据约定

公开商品页面、分类页和 SEO 清单只读取 `src/lib/catalog/mock-data.ts` 中构建时可用的静态目录。新增或修改公开商品后，需更新该目录并重新构建/部署；已部署站点不支持在后台即时发布新 URL。`/admin/products` 仅供当前浏览器内编辑和预览：其 localStorage 数据不会被服务端、其他访客、静态路由、metadata 或 sitemap 读取，页面中的“published”只表示本地预览状态，不是线上发布信号。请勿将此 MVP 的本地保存描述为生产级后台发布能力。

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

- `/`：首页，展示精选商品与分类入口。
- `/collections/[slug]`：商品分类页；当前 mock 分类包括 `built-in`、`freestanding` 和 `compact`。
- `/products/[slug]`：商品详情页；商品详情由 mock catalog 中的 slug 匹配。

## 商品、分类与模板入口

- `src/lib/catalog/mock-data.ts`：`Product` / `Collection` 内容类型、mock 商品与分类列表，以及按 slug 查询的辅助函数。新增商品或分类时从这里维护数据，并确保商品的 `collectionSlug` 指向有效分类。
- `src/lib/templates/registry.ts`：模板注册表、模板标识和默认模板。首版 `DEFAULT_TEMPLATE` 为 `gas-grill`。
- `src/app/`：首页、分类页、商品详情页等路由入口。
- `src/components/templates/gas-grill/`：当前 gas-grill 模板专属的页面组件。

新增商品或分类前，请同时检查对应路由和模板组件对 catalog 字段的使用，保持 slug、分类关联及展示内容一致。首版只实现并使用 `gas-grill` 模板。

## 当前能力边界

- 后端基础已新增于 `backend/`：Go API、本机 PostgreSQL 项目专属数据库和显式 migration，提供 `/health/live`、`/health/ready`，以及只读的 `/api/v1/collections`、`/api/v1/products`、`/api/v1/company-content` 端点。本地启动细节见 [`backend/README.md`](./backend/README.md)。
- 样例商品/分类/公司内容已可通过 `go run ./cmd/storefront-seed` 写入店铺数据库并由上述 API 读取；这些记录标记为 `sample`，不是已核验的真实产品或公司事实。
- **前端页面仍读取仓库内静态 mock 数据**，尚未接入 Catalog API；询盘 API、询盘持久化和 CRM 管理页也尚未实现。
- 管理员商品/公司页面写入浏览器 localStorage 的内容仅用于当前浏览器预览，不会发布到共享 storefront、创建公开商品 URL 或进入 sitemap。
- Go API、每个店铺独立部署与独立 PostgreSQL 数据库属于后续规划，不是当前已交付能力。
- 当前不提供 checkout、订单处理或支付功能，也未集成第三方电商插件。
