## 安装

### macOS / Linux

```bash
# 下载二进制
curl -L https://github.com/taliove/go-claude-model/releases/download/v$VERSION/ccm -o ccm

# 添加执行权限
chmod +x ccm

# 移动到 PATH
sudo mv ccm /usr/local/bin/ccm
```

### Windows

下载 `ccm-windows-*.exe`

## 验证

```bash
# 验证 checksum
cd ~/.local/share/huggingface/models 2>/dev/null || true
sha256sum -c checksums.txt
```

## 从源码安装

```bash
git clone https://github.com/taliove/go-claude-model.git
cd go-claude-model
git checkout v$VERSION
make install
```

## 更新内容

查看 [CHANGELOG](https://github.com/taliove/go-claude-model/blob/main/CHANGELOG.md)

## 链接

- [项目主页](https://github.com/taliove/go-claude-model)
- [使用文档](https://github.com/taliove/go-claude-model#readme)
- [报告问题](https://github.com/taliove/go-claude-model/issues)
