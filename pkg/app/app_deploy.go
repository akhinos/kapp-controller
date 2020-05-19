package app

import (
	"fmt"

	ctldep "github.com/k14s/kapp-controller/pkg/deploy"
	"github.com/k14s/kapp-controller/pkg/exec"
)

func (a *App) deploy(tplOutput string, changedFunc func(exec.CmdRunResult)) exec.CmdRunResult {
	err := a.blockDeletion()
	if err != nil {
		result := exec.CmdRunResult{}
		result.AttachErrorf("Blocking for deploy: %s", err)
		return result
	}

	if len(a.app.Spec.Deploy) != 1 {
		return exec.NewCmdRunResultWithErr(fmt.Errorf("Expected exactly one deploy option"))
	}

	var result exec.CmdRunResult

	for _, dep := range a.app.Spec.Deploy {
		switch {
		case dep.Kapp != nil:
			kapp := a.deployFactory.NewKapp(*dep.Kapp, a.app.Spec.Cluster, a.deployGenericOpts(), a.newCancelCh())
			result = kapp.Deploy(tplOutput, changedFunc)
		default:
			result.AttachErrorf("%s", fmt.Errorf("Unsupported way to deploy"))
		}
		if result.Error != nil {
			break
		}
	}

	return result
}

func (a *App) delete(changedFunc func(exec.CmdRunResult)) exec.CmdRunResult {
	var result exec.CmdRunResult

	for _, dep := range a.app.Spec.Deploy {
		switch {
		case dep.Kapp != nil:
			kapp := a.deployFactory.NewKapp(*dep.Kapp, a.app.Spec.Cluster, a.deployGenericOpts(), a.newCancelCh())
			result = kapp.Delete(changedFunc)
		default:
			result.AttachErrorf("%s", fmt.Errorf("Unsupported way to delete"))
		}
		if result.Error != nil {
			break
		}
	}

	if result.Error == nil {
		err := a.unblockDeletion()
		if err != nil {
			result := exec.CmdRunResult{}
			result.AttachErrorf("Unblocking for deploy: %s", err)
			return result
		}
	}

	return result
}

func (a *App) inspect() exec.CmdRunResult {
	var result exec.CmdRunResult

	for _, dep := range a.app.Spec.Deploy {
		switch {
		case dep.Kapp != nil:
			kapp := a.deployFactory.NewKapp(*dep.Kapp, a.app.Spec.Cluster, a.deployGenericOpts(), a.newCancelCh())
			result = kapp.Inspect()
		default:
			result.AttachErrorf("%s", fmt.Errorf("Unsupported way to inspect"))
		}
		if result.Error != nil {
			break
		}
	}

	return result
}

func (a *App) deployGenericOpts() ctldep.GenericOpts {
	return ctldep.GenericOpts{Name: a.app.Name, Namespace: a.app.Namespace}
}
