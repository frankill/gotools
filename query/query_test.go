package query_test

import (
	"encoding/json"
	"testing"

	"github.com/frankill/gotools/query"
)

func TestEsQuery_MultipleQueries(t *testing.T) {
	// 创建一个新的EsQuery实例
	query := query.NewESMustQuery()

	// 使用Eq方法添加一个查询条件
	query.Eq("fieldName1", "value1")

	// 使用In方法添加一个查询条件
	query.In("fieldName2", "value2", "value3")

	// 使用Gt方法添加一个查询条件
	query.Gt("fieldName3", 10)

	// 构建查询
	buildQuery := query.Build()

	// 断言构建的查询不为空
	if buildQuery == nil {
		t.Errorf("Expected non-nil query, got nil")
	}

	// 将查询序列化为JSON字符串，以便检查
	source, err := buildQuery.Source()
	if err != nil {
		t.Fatalf("Failed to get query source: %v", err)
	}
	b, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("Failed to marshal query source: %v", err)
	}

	// 期望的查询条件JSON字符串
	expected := `{"bool":{"must":[{"term":{"fieldName1":"value1"}},{"terms":{"fieldName2":["value2","value3"]}},{"range":{"fieldName3":{"from":10,"include_lower":false,"include_upper":true,"to":null}}}]}}`
	if string(b) != expected {
		t.Errorf("Expected query %s, got %s", expected, string(b))
	}
}

func TestEsQuery_Eq(t *testing.T) {
	// 创建一个新的EsQuery实例
	query := query.NewESMustQuery()

	// 使用Eq方法添加一个查询条件
	query.Eq("fieldName", "value")

	// 构建查询
	buildQuery := query.Build()

	// 断言构建的查询不为空
	if buildQuery == nil {
		t.Errorf("Expected non-nil query, got nil")
	}

	// 将查询序列化为JSON字符串，以便检查
	source, err := buildQuery.Source()
	if err != nil {
		t.Fatalf("Failed to get query source: %v", err)
	}
	b, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("Failed to marshal query source: %v", err)
	}

	// 期望的查询条件JSON字符串
	expected := `{"bool":{"must":{"term":{"fieldName":"value"}}}}`
	if string(b) != expected {
		t.Errorf("Expected query %s, got %s", expected, string(b))
	}
}

// 示例测试函数 TestSQLBuilder_Build 测试 SQLBuilder 的功能
func TestSQLBuilder_Build(t *testing.T) {
	tests := []struct {
		name     string
		builder  *query.SQLBuilder
		expected string
	}{
		{
			name: "Basic SELECT query",
			builder: query.NewSQLBuilder().
				Select("id", "name", "age").
				From("users").
				Where("age > 18", "name != 'Admin'").
				OrderBy("name"),
			expected: "SELECT id, name, age FROM users WHERE age > 18 AND name != 'Admin' ORDER BY name",
		},
		{
			name: "SELECT query with IN subquery",
			builder: query.NewSQLBuilder().
				Select("id", "name", "age").
				From("users").
				In("id", query.NewSQLBuilder().
					Select("user_id").
					From("permissions").
					Where("role = 'admin'")).
				OrderBy("name"),
			expected: "SELECT id, name, age FROM users WHERE id IN (SELECT user_id FROM permissions WHERE role = 'admin') ORDER BY name",
		},
		{
			name: "SELECT query with BETWEEN clause",
			builder: query.NewSQLBuilder().
				Select("id", "name", "age").
				From("users").
				Between("age", 18, 30).
				OrderBy("name"),
			expected: "SELECT id, name, age FROM users WHERE age BETWEEN 18 AND 30 ORDER BY name",
		},
		{
			name: "SELECT query with AND and OR conditions",
			builder: query.NewSQLBuilder().
				Select("id", "name", "age").
				From("users").
				Where("test is not null").
				And("age > 18", "name != 'Admin'").
				Or("age < 10", "name = 'Guest'").
				OrderBy("name"),
			expected: "SELECT id, name, age FROM users WHERE test is not null AND (age > 18 AND name != 'Admin') AND (age < 10 OR name = 'Guest') ORDER BY name",
		},
		{
			name: "SELECT query with multiple GROUP BY fields",
			builder: query.NewSQLBuilder().
				Select("category", "COUNT(*) as count").
				From("products").
				GroupBy("category", "sub_category"),
			expected: "SELECT category, COUNT(*) as count FROM products GROUP BY category, sub_category",
		},
		{
			name: "SELECT query with multiple ORDER BY fields",
			builder: query.NewSQLBuilder().
				Select("id", "name").
				From("users").
				OrderBy("name ASC", "id DESC"),
			expected: "SELECT id, name FROM users ORDER BY name ASC, id DESC",
		},
		// Add more test cases as needed
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.builder.Build()
			if actual != tc.expected {
				t.Errorf("Expected SQL: %s, but got: %s", tc.expected, actual)
			}
		})
	}
}
