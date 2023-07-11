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
  variableMaker := variable_sources.NewVariableSourceBuilder("variable-make").
    WithUpdateFn(func(ctx context.Context, problem deppy.MutableResolutionProblem, variable deppy.Variable) error {
      for row := 0; row < 9; row++ {
        for col := 0; col < 9; col++ {
          cell := variables.NewMutableVariable(deppy.Identifierf("%d:%d", row, col), sudokuCellKind, nil)
          if err := cell.AddMandatory("mandatory"); err != nil {
            return err
          }
          for num := 0; num < 9; num++ {
            varId := deppy.Identifier(fmt.Sprintf("%d:%d:%d", row, col, num))

            // add row/col/num variables
            v := variables.NewMutableVariable(varId, sudokuVariableKind, nil)

            // add row/col constraints
            for i := 1; i < 9; i++ {
              rowCID := deppy.Identifierf("row-conflict:%d:%d:%d", row, col, i+1)
              rowVarID := deppy.Identifierf("%d:%d:%d", row, (col+i)%9, num)
              if err := v.AddConflict(rowCID, rowVarID); err != nil {
                return err
              }

              colCID := deppy.Identifierf("col-conflict:%d:%d:%d", row, col, i+1)
              colVarID := deppy.Identifierf("%d:%d:%d", (row+i)%9, col, num)
              if err := v.AddConflict(colCID, colVarID); err != nil {
                return err
              }
            }

            // add box constraints
            if err := v.AddConflict("box-conflict-1", deppy.Identifierf("%d:%d:%d", (row+1)%3+3*(row/3), (col+1)%3+3*(col/3), num)); err != nil {
              return err
            }

            if err := v.AddConflict("box-conflict-2", deppy.Identifierf("%d:%d:%d", (row+1)%3+3*(row/3), (col+2)%3+3*(col/3), num)); err != nil {
              return err
            }

            if err := v.AddConflict("box-conflict-3", deppy.Identifierf("%d:%d:%d", (row+2)%3+3*(row/3), (col+1)%3+3*(col/3), num)); err != nil {
              return err
            }

            if err := v.AddConflict("box-conflict-4", deppy.Identifierf("%d:%d:%d", (row+2)%3+3*(row/3), (col+2)%3+3*(col/3), num)); err != nil {
              return err
            }

            if err := problem.ActivateVariable(v); err != nil {
              return err
            }
            if err := cell.AddDependency("pick-one", v.VariableID()); err != nil {
              return err
            }
          }
          if err := problem.ActivateVariable(cell); err != nil {
            return err
          }
        }
      }
      return nil
    }).Build(ctx)

  problem, err := resolution.NewResolutionProblemBuilder("sudoku").
    WithVariableSources(variableMaker).
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
      cp := strings.Split(v.VariableID().String(), ":")
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
