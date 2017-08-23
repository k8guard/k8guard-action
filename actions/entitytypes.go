package actions

import (
	libs "github.com/k8guard/k8guardlibs"
	"github.com/k8guard/k8guardlibs/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ActionableEntity interface {
	DoAction()
}

// See http://stackoverflow.com/questions/28800672/how-to-add-new-methods-to-an-existing-type-in-go
type ActionPod libs.Pod
type ActionDeployment libs.Deployment
type ActionIngress libs.Ingress
type ActionJob libs.Job
type ActionCronJob libs.CronJob

func (a ActionPod) DoAction() {
	clientset, err := k8s.LoadClientset()
	if err != nil {
		panic(err)
	}
	libs.Log.Debug("Deleting Pod ", a.Name, " in namesapce ", a.Namespace)
	err = clientset.CoreV1().Pods(a.Namespace).Delete(a.Name, &metav1.DeleteOptions{})
	if err != nil {
		libs.Log.Error(err)
	}
}

func (a ActionDeployment) DoAction() {
	clientset, err := k8s.LoadClientset()
	if err != nil {
		panic(err)
	}
	libs.Log.Debug("Scaling Deployment ", a.Name, " in namesapce ", a.Namespace)
	kd, err := clientset.AppsV1beta1().Deployments(a.Namespace).Get(a.Name, metav1.GetOptions{})
	if err != nil {
		libs.Log.Error(err)
		return
	}
	replicas := int32(0)
	kd.Spec.Replicas = &replicas
	_, err = clientset.AppsV1beta1().Deployments(a.Namespace).Update(kd)
	if err != nil {
		libs.Log.Error(err)
	}
}

func (a ActionIngress) DoAction() {
	clientset, err := k8s.LoadClientset()
	if err != nil {
		panic(err)
	}
	libs.Log.Debug("Deleting Ingress ", a.Name, " in namesapce ", a.Namespace)
	err = clientset.Ingresses(a.Namespace).Delete(a.Name, &metav1.DeleteOptions{})
	if err != nil {
		libs.Log.Error(err)
	}
}

func (a ActionJob) DoAction() {
	clientset, err := k8s.LoadClientset()
	if err != nil {
		panic(err)
	}
	libs.Log.Debug("Deleting Job ", a.Name, " in namesapce ", a.Namespace)
	err = clientset.BatchV1().Jobs(a.Namespace).Delete(a.Name, &metav1.DeleteOptions{})
	if err != nil {
		libs.Log.Error(err)
	}
}

func (a ActionCronJob) DoAction() {
	if libs.Cfg.IncludeAlpha == false {
		libs.Log.Debug("Ignoring CronJob action as alpha features are not enabled ")
		return
	}

	clientset, err := k8s.LoadClientset()
	if err != nil {
		panic(err)
	}
	libs.Log.Debug("Disabling CronJob ", a.Name, " in namesapce ", a.Namespace)

	kcj, err := clientset.BatchV2alpha1().CronJobs(a.Namespace).Get(a.Name, metav1.GetOptions{})
	if err != nil {
		libs.Log.Error(err)
		return
	}
	suspend := true
	kcj.Spec.Suspend = &suspend
	_, err = clientset.BatchV2alpha1().CronJobs(a.Namespace).Update(kcj)
	if err != nil {
		libs.Log.Error(err)
	}
}
