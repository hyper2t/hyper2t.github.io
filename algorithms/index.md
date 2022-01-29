# 刷 Leetcode 的算法总结


这篇文章总结了我刷 Leetcode 时候用到的算法模板，包括回溯算法，二分法用到的双指针，动态规划等等。

<!--more-->

## 1. 回溯算法

### 1.1 解题过程

关于回溯算法，一般是通过在集合中递归查找子集，集合的大小构成树的宽度，递归的深度构成树的深度。

以下是回溯"三部曲"：
1. 确定回溯函数的返回值和参数
2. 确定回溯函数的终止条件
3. 确定回溯搜索的遍历过程

### 1.2 回溯模板

```c++
void backtracking(参数){
  if (终止条件){
    存放结果;
    return;
  }
  
  for (选择：本层集合中的元素，即树中节点孩子的数量就是集合的大小){
    处理节点;
    backtracking(路径,选择列表);
    回溯, 撤销处理结果;
  }
}
```

### 1.3 Leetcode 实战

#### 全排列 (Leetcode 46)

给定一个不含重复数字的数组`nums` ，返回其 所有可能的全排列 。你可以 按任意顺序 返回答案。

示例 1：
```markdown
输入：nums = [1,2,3]
输出：[[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]]
```

示例 2：
```markdown
输入：nums = [0,1]
输出：[[0,1],[1,0]]
```

示例 3：
```markdown
输入：nums = [1]
输出：[[1]]
```

提示：

- 1 $\leq$ nums.length $\leq$ 6
- -10 $\leq$ nums[i] $\leq$ 10
- `nums` 中的所有整数互不相同

解法如下：
```go
func permute(nums []int) [][]int {
	var res [][]int
	if len(nums) == 0 {
		return res
	}
	var tmp []int
	var visited = make([]bool, len(nums))
	backtracking(nums, &res, tmp, visited)
	return res
}

func backtracking(arr []int, res *[][]int, tmp []int, visited []bool) {
	// 成功找到一组
	if len(tmp) == len(arr) {
		var c = make([]int, len(tmp))
		copy(c, tmp)
		*res = append(*res, c)
		return
	}
	// 回溯
	for i := 0; i < len(arr); i++ {
		if !visited[i] {
			visited[i] = true
			tmp = append(tmp, arr[i])
			backtracking(arr, res, tmp, visited)
			tmp = tmp[:len(tmp)-1]
			visited[i] = false
		}
	}
}
```
