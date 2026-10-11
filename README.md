# Gas Grill Storefront MVP

Next.js 15 / React 19 单模板独立站 storefront MVP，以户外燃气烤炉为示例。当前交付的是用于展示商品与分类的前端，不包含购物车结算（checkout）或支付能力。

当前 storefront 使用的模板标识为 `template_id: gas-grill`。模板面向访客的页面由 Next.js 渲染；商品和分类内容目前来自仓库内的静态 mock 数据。

## 环境要求

- Node.js 18.18+、19.8+ 或 20+（与 lockfile 中 Next.js 的 engines 要求一致）；建议使用仍受支持的 Node.js LTS 版本。
- npm。仓库提交了 `package-lock.json`，安装依赖时使用 `npm ci`。
- 当前本地展示不需要环境变量，也不依赖 Go API 或数据库。

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

- 商品和分类来自静态 mock 数据；当前没有连接后台服务或数据库。
- Go API、每个店铺独立部署与独立 PostgreSQL 数据库属于后续规划，不是当前已交付能力。
- 当前不提供 checkout、订单处理或支付功能，也未集成第三方电商插件。
