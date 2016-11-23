# gnt
游戏数值导出工具（Game Numerical Tool）

### Excel文件规范

工作表格式：前三行为表头，第一行为描述，第二行为字段名，第三行为类型

默认支持三种数据类型：int, float, string

### 配置文件
YAML格式

```yaml
name: 配置名
template: 生成模版文件(lua.tpl)，路径同配置文件
sheet: 默认数据工作表索引（从1开始）
pkcols: 主键列数组
pksep: 主键分隔符（适用于复合主键情况）
list: 数据项列表
	input: excel文件名(example.xlsx)
	sheet: 当前数据项工作表索引（从1开始）
	pkcols: 当前数据项主键列数组
	cols: 数据列（负数代表不包含，正数代表包含，从1开始，不能同时包含正负列号)
	output: 输出文件路径名
```

### 模版文件
使用Go语言[text/template][1]模版语法描述，控制输出文件的生成内容和格式

### 命令行参数

**_-config_**: 配置文件路径

**_-input_**: excel文件夹路径

**_-output_**: 输出文件夹路径

**_-debug_**: 打开log输出

**_-force_**: 强制处理所有input

##### 例子
##### 例子
example目录


[1]: https://golang.org/pkg/text/template/
