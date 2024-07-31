package query_test

import (
	"testing"

	"github.com/frankill/gotools/query"
)

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
