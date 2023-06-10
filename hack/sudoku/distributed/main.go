package main

import (
	"context"
	"fmt"
	"github.com/perdasilva/replee/pkg/deppy"
	"github.com/perdasilva/replee/pkg/deppy/resolution"
	"github.com/perdasilva/replee/pkg/deppy/resolver"
	"github.com/perdasilva/replee/pkg/deppy/variable_sources"
	"github.com/perdasilva/replee/pkg/deppy/variables"
	"log"
	"strconv"
	"strings"
)

const sudokuCellKind = "deppy.variable.cell"
const sudokuVariableKind = "deppy.variable.sudoku"

func main() {
	ctx := context.Background()
	sudokuVariableMaker := variable_sources.NewVariableSourceBuilder("variable-make").
		WithUpdateFn(func(ctx context.Context, problem deppy.MutableResolutionProblem, variable deppy.Variable) error {
			for row := 0; row < 9; row++ {
				for col := 0; col < 9; col++ {
					for num := 0; num < 9; num++ {
						varId := deppy.Identifierf("%d:%d:%d", row, col, num)
						v := variables.NewMutableVariable(varId, sudokuVariableKind, map[string]interface{}{
							"row": row,
							"col": col,
							"num": num,
						})
						if err := problem.ActivateVariable(v); err != nil {
							return err
						}
					}
				}
			}
			return nil
		}).Build(ctx)

	cellVariableMaker := variable_sources.NewVariableSourceBuilder("cell-make").
		WithVariableFilterFn(func(variable deppy.Variable) bool {
			return variable != nil && variable.Kind() == sudokuVariableKind
		}).
		WithUpdateFn(func(ctx context.Context, problem deppy.MutableResolutionProblem, variable deppy.Variable) error {
			row, _ := variable.GetProperty("row")
			col, _ := variable.GetProperty("col")
			num, _ := variable.GetProperty("num")
			v := variables.NewMutableVariable(deppy.Identifierf("%s:%s", row, col), sudokuCellKind, map[string]interface{}{
				"row": row,
				"col": col,
			})
			if err := v.AddMandatory("mandatory"); err != nil {
				return err
			}
			if err := v.AddDependency("pick-one", deppy.Identifierf("%d:%d:%d", row, col, num)); err != nil {
				return err
			}
			return problem.ActivateVariable(v)
		}).
		Build(ctx)

	cellRowColConstraintSource := variable_sources.NewVariableSourceBuilder("cell-row-col-constraint").
		WithVariableFilterFn(func(variable deppy.Variable) bool {
			return variable != nil && variable.Kind() == sudokuVariableKind
		}).
		WithUpdateFn(func(ctx context.Context, problem deppy.MutableResolutionProblem, variable deppy.Variable) error {
			row, _ := variable.GetProperty("row")
			col, _ := variable.GetProperty("col")
			num, _ := variable.GetProperty("num")

			v, err := problem.GetMutableVariable(variable.Identifier(), variable.Kind())
			if err != nil {
				return err
			}
			for i := 1; i < 9; i++ {
				rowCID := deppy.Identifierf("row-conflict:%d:%d:%d", row, col, i+1)
				rowVarID := deppy.Identifierf("%d:%d:%d", row, (col.(int)+i)%9, num)
				if err := v.AddConflict(rowCID, rowVarID); err != nil {
					return err
				}

				colCID := deppy.Identifierf("col-conflict:%d:%d:%d", row, col, i+1)
				colVarID := deppy.Identifierf("%d:%d:%d", (row.(int)+i)%9, col, num)
				if err := v.AddConflict(colCID, colVarID); err != nil {
					return err
				}
			}
			return nil
		}).
		Build(ctx)

	boxConstraintSource := variable_sources.NewVariableSourceBuilder("box-constraint").
		WithVariableFilterFn(func(variable deppy.Variable) bool {
			return variable != nil && variable.Kind() == sudokuVariableKind
		}).
		WithUpdateFn(func(ctx context.Context, problem deppy.MutableResolutionProblem, variable deppy.Variable) error {
			row, _ := variable.GetProperty("row")
			col, _ := variable.GetProperty("col")
			num, _ := variable.GetProperty("num")

			v, err := problem.GetMutableVariable(variable.Identifier(), variable.Kind())
			if err != nil {
				return err
			}

			// add box constraints
			if err := v.AddConflict("box-conflict-1", deppy.Identifierf("%d:%d:%d", (row.(int)+1)%3+3*(row.(int)/3), (col.(int)+1)%3+3*(col.(int)/3), num)); err != nil {
				return err
			}

			if err := v.AddConflict("box-conflict-2", deppy.Identifierf("%d:%d:%d", (row.(int)+1)%3+3*(row.(int)/3), (col.(int)+2)%3+3*(col.(int)/3), num)); err != nil {
				return err
			}

			if err := v.AddConflict("box-conflict-3", deppy.Identifierf("%d:%d:%d", (row.(int)+2)%3+3*(row.(int)/3), (col.(int)+1)%3+3*(col.(int)/3), num)); err != nil {
				return err
			}

			if err := v.AddConflict("box-conflict-4", deppy.Identifierf("%d:%d:%d", (row.(int)+2)%3+3*(row.(int)/3), (col.(int)+2)%3+3*(col.(int)/3), num)); err != nil {
				return err
			}

			return nil
		}).
		Build(ctx)

	problem, err := resolution.NewResolutionProblemBuilder("sudoku").
		WithVariableSources(
			sudokuVariableMaker,
			cellVariableMaker,
			boxConstraintSource,
			cellRowColConstraintSource,
		).
		Build(ctx)

	if err != nil {
		panic(err)
	}

	s := resolver.NewDeppyResolver()
	solution, err := s.Solve(ctx, problem, resolver.DisableOrderPreference())
	if err != nil {
		panic(err)
	}

	if len(solution.NotSatisfiable()) > 0 {
		log.Fatalln(solution.NotSatisfiable())
	}

	var board = [9][9]int{}
	for _, v := range solution.SelectedVariables() {
		if v.Kind() == sudokuVariableKind {
			cp := strings.Split(v.Identifier().String(), ":")
			row, _ := strconv.Atoi(cp[0])
			col, _ := strconv.Atoi(cp[1])
			num, _ := strconv.Atoi(cp[2])
			board[row][col] = num + 1
		}
	}
	for _, row := range board {
		for _, col := range row {
			fmt.Printf("%d ", col)
		}
		fmt.Println()
	}
}
