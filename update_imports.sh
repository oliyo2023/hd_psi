#!/bin/bash

# 在backend目录中查找所有Go文件
find backend -name "*.go" -type f | while read -r file; do
    # 将导入路径从 "hd_psi/xxx" 更新为 "hd_psi/backend/xxx"
    sed -i 's|"hd_psi/|"hd_psi/backend/|g' "$file"
    echo "Updated imports in $file"
done

echo "Import paths updated successfully!"
