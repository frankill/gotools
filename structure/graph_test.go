package structure_test

import (
	"fmt"
	"testing"

	"github.com/frankill/gotools/structure"
)

func TestGraphOperations(t *testing.T) {
	// 创建一个新的图，比较函数用于判断两个数据是否相等
	graph := structure.NewGraph[int](func(v, d int) bool {
		return v == d
	}, 100, func(v int) int { return v })

	// 添加顶点
	graph.AddVer(1)
	graph.AddVer(1)
	graph.AddVer(2)
	graph.AddVer(2)
	graph.AddVer(3)

	// 添加边
	if added := graph.Addedge(1, 2); !added {
		t.Errorf("Failed to add edge from 1 to 2")
	}
	if added := graph.Addedge(1, 3); !added {
		t.Errorf("Failed to add edge from 1 to 2")
	}
	if added := graph.Addedge(2, 3); !added {
		t.Errorf("Failed to add edge from 2 to 3")
	}
	if added := graph.Addedge(3, 2); !added {
		t.Errorf("Failed to add edge from 2 to 3")
	}
	if added := graph.Addedge(3, 2); !added {
		t.Errorf("Failed to add edge from 2 to 3")
	}

	fmt.Println(graph.DrawGraph())

	// 测试顶点索引获取
	index := graph.GetIndex(3)
	if index != 2 {
		t.Errorf("Expected index 2 for value 3, got %d", index)
	}

	// 验证顶点数量
	if loc := graph.Getloc(); loc != 3 {
		t.Errorf("Expected vertex count 3, got %d", loc)
	}

}
