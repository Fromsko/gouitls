#!/bin/bash

# 获取最新的标签
latest_tag=$(git describe --tags --abbrev=0)

# 生成更新日志
changelog=$(git log ${latest_tag}..HEAD --oneline --pretty=format:"- %s")

# 保存更新日志到文件
echo "${changelog}" > CHANGELOG.md
