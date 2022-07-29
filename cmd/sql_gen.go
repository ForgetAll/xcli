package cmd

import (
	"fmt"

	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver" // sql parser driver
	"github.com/spf13/cobra"
)

var sql = ""

var sqlGenCmd = &cobra.Command{
	Use:   "sql-gen",
	Short: "generate code by sql",
	Long:  `generate code by sql, just generate default empty implement code`,
	Run: func(cmd *cobra.Command, args []string) {
		if sql == "" {
			return
		}
		result, err := parse(sql)
		if err != nil {
			cmd.Println("sql error!")
			return
		}

		cmd.Println(result)
	},
}

func init() {
	rootCmd.AddCommand(sqlGenCmd)
	sqlGenCmd.PersistentFlags().StringVarP(&sql, "sql", "s", "", "input sql param")
}

func parse(sql string) (string, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return "", err
	}

	for _, node := range stmtNodes {
		m := &Meta{}
		node.Accept(m)
		fmt.Printf("col name: %v\n", m.colNames)
		fmt.Printf("table name: %v\n", m.tableName)
	}

	return "", nil
}

type Meta struct {
	colNames  []string
	tableName string
}

func (c *Meta) Enter(in ast.Node) (ast.Node, bool) {
	if name, ok := in.(*ast.ColumnName); ok {
		c.colNames = append(c.colNames, name.Name.O)
	}

	if name, ok := in.(*ast.TableName); ok {
		c.tableName = name.Name.O
	}

	return in, false
}

func (c *Meta) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
