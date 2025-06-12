#!/bin/bash

# 该脚本用于从Git历史中移除大文件

echo "开始清理Git仓库中的大文件..."

# 使用git filter-branch删除大文件
git filter-branch --force --index-filter \
  "git rm --cached --ignore-unmatch data/blockchain/00003.mem data/blockchain/000005.vlog" \
  --prune-empty --tag-name-filter cat -- --all

# 清理和优化仓库
echo "清理和优化仓库..."
rm -rf .git/refs/original/
git reflog expire --expire=now --all
git gc --prune=now
git gc --aggressive --prune=now

echo "清理完成！"
echo "现在你可以尝试重新推送到GitHub:"
echo "git push origin --force"
echo ""
echo "注意：这是一个强制推送，会覆盖远程仓库的历史。如果有其他人在使用这个仓库，请先通知他们。" 