package main

import (
	"github.com/go-logr/logr"
	kcv1alpha1 "github.com/k14s/kapp-controller/pkg/apis/kappctrl/v1alpha1"
	ctlapp "github.com/k14s/kapp-controller/pkg/app"
	kcclient "github.com/k14s/kapp-controller/pkg/client/clientset/versioned"
	"github.com/k14s/kapp-controller/pkg/deploy"
	"github.com/k14s/kapp-controller/pkg/fetch"
	"github.com/k14s/kapp-controller/pkg/template"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type AppFactory struct {
	coreClient kubernetes.Interface
	appClient  kcclient.Interface
}

func (f *AppFactory) NewCRDAppFromName(request reconcile.Request, log logr.Logger) *ctlapp.CRDApp {
	return ctlapp.NewCRDAppFromName(request.NamespacedName, log, f.appClient)
}

func (f *AppFactory) NewCRDApp(app *kcv1alpha1.App, log logr.Logger) (*ctlapp.CRDApp, error) {
	fetchFactory := fetch.NewFactory(f.coreClient)
	templateFactory := template.NewFactory(f.coreClient, fetchFactory)
	deployFactory := deploy.NewFactory(f.coreClient)
	return ctlapp.NewCRDApp(app, log, f.appClient, fetchFactory, templateFactory, deployFactory)
}
