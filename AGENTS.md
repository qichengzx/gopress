# gopress — 快速静态博客框架

## 构建与运行

```bash
go build ./...              # 编译
go run main.go -g           # 生成静态文件到 public/
go run main.go -s           # 启动 HTTP 服务（默认端口 8092）
go run main.go -p 8080 -s   # 指定端口
```

- 无 Makefile / Taskfile / CI / lint 配置
- 无测试文件（`*_test.go` 不存在）
- 全部标准 Go 工具链，无额外构建步骤

## 项目结构

| 目录 | 用途 |
|------|------|
| `main.go` | CLI 入口，`-g` 生成 / `-s` 提供服务 / `-p` 端口 |
| `config/provider.go` | 解析 `_config.yaml` 及主题 `_config.yaml` |
| `generator/` | 核心：文章解析、页面渲染、sitemap、atom、文件操作 |
| `plugins/sitemap/` | sitemap.xml 插件 |
| `themes/<theme>/` | Go `html/template` 布局文件（非 text/template） |
| `source/_posts/` | Markdown 文章（YAML front matter + `---` 分隔） |
| `public/` | 生成输出目录（`gitignore` 中） |

## 重要细节

- 模板文件路径：`themes/<theme>/layout/*.html`，通过 `template.ParseGlob` 加载
- 主布局入口：`layout.html` 中的 `{{ define "layout" }}`
- 文章格式：以 `---` 或 `;;;` 分隔 YAML front matter 与 Markdown 正文
- 文章日期格式必须是 `2006-01-02 15:04:05`
- `Build()` 会先清空 `public/` 再重新生成全部页面
- `Atom()` 方法已实现但未被 `Build()` 调用——是游离代码
- 资源复制：`source/images/` → `public/images/`；主题 `source/css|js|fancybox` → `public/css|js|fancybox`
- 主题配置由 `themes/<theme>/_config.yaml` 提供，结构体 `ThemeCfg`
- 默认 permalink：`:year/:month/:day/:title.html`
- 依赖：`gorilla/feeds`、`blackfriday/v2`、`gopkg.in/yaml.v2`
