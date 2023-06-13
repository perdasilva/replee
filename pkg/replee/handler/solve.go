package handler

import (
	"context"
	"encoding/json"
	"github.com/perdasilva/replee/pkg/deppy"
	"github.com/perdasilva/replee/pkg/deppy/resolver"
	"github.com/perdasilva/replee/pkg/replee/action"
	"github.com/perdasilva/replee/pkg/replee/store"
	"log"
)

var _ Handler = &SolveActionHandler{}

type SolveActionHandler struct {
}

func (h SolveActionHandler) HandleAction(a action.Action, s store.ResolutionProblemStore) error {
	if err := validate(a, action.ActionTypeSolve, []string{ParameterProblemID}); err != nil {
		return err
	}
	problemID := a.Parameters[ParameterProblemID].(string)
	if m, err := s.Get(deppy.Identifier(problemID)); err != nil {
		return err
	} else {
		r := resolver.NewDeppyResolver()
		if solution, err := r.Solve(context.Background(), m); err != nil {
			return err
		} else {
			s, err := json.MarshalIndent(solution, "", "  ")
			if err != nil {
				log.Fatal(err)
			}
			log.Println(string(s))
		}
	}
	return nil
}
