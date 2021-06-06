package main

import (
	"context"
	"fmt"
	"github.com/healthcheck-exporter/cmd/api"
	"github.com/healthcheck-exporter/cmd/authentication"
	"github.com/healthcheck-exporter/cmd/configuration"
	"github.com/healthcheck-exporter/cmd/exporter"
	"github.com/healthcheck-exporter/cmd/healthcheck"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"net/http"
	"os"

	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
)

func main() {
	start()

	config := configuration.GetConfiguration()
	authClient := authentication.NewAuthClient(config)

	ex := exporter.NewExporter(config)

	hcClient := healthcheck.NewHealthCheck(config, authClient, ex)
	//
	//http.Handle("/metrics", promhttp.Handler())
	//
	//http.Handle("/probe", promhttp.Handler())

	// initialize api
	router := api.NewRouter(hcClient)

	// enable CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"GET"},
	})

	log.Info(fmt.Sprintf(http.ListenAndServe(":2112",
		corsHandler.Handler(router)).Error()))

	//log.Info(fmt.Sprintf(http.ListenAndServe(":2112", nil).Error()))

	//
	//
	//err = http.ListenAndServe(":2112", nil)
	//if err != nil {
	//	panic(err)
	//}
}

func start() error {
	//var kubeconfig *string
	//if home := homeDir(); home != "" {
	//	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	//} else {
	//	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	//}
	//flag.Parse()

	//// use the current context in kubeconfig
	//config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	//if err != nil {
	//	return err
	//}

	//buildV1Client, err := buildv1.NewForConfig(config)
	//if err != nil {
	//	return err
	//}

	namespace := "suep-omp-prod"

	//// get all builds
	//builds, err := buildV1Client.Builds(namespace).List(context.TODO(), metav1.ListOptions{})
	//if err != nil {
	//	return err
	//}
	//fmt.Printf("There are %d builds in project %s\n", len(builds.Items), namespace)
	//// List names of all builds
	//for i, build := range builds.Items {
	//	fmt.Printf("index %d: Name of the build: %s", i, build.Name)
	//}
	//
	//// get a specific build
	//build := "cakephp-ex-1"
	//myBuild, err := buildV1Client.Builds(namespace).Get(context.TODO(), build, metav1.GetOptions{})
	//if err != nil {
	//	return err
	//}
	//fmt.Printf("Found build %s in namespace %s\n", build, namespace)
	//fmt.Printf("Raw printout of the build %+v\n", myBuild)
	//// get details of the build
	//fmt.Printf("name %s, start time %s, duration (in sec) %.0f, and phase %s\n",
	//	myBuild.Name, myBuild.Status.StartTimestamp.String(),
	//	myBuild.Status.Duration.Seconds(), myBuild.Status.Phase)
	//
	//// trigger a build
	//buildConfig := "cakephp-ex"
	//myBuildConfig, err := buildV1Client.BuildConfigs(namespace).Get(context.TODO(), buildConfig, metav1.GetOptions{})
	//if err != nil {
	//	return err
	//}
	//fmt.Printf("Found BuildConfig %s in namespace %s\n", myBuildConfig.Name, namespace)
	//buildRequest := v1.BuildRequest{}
	//buildRequest.Kind = "BuildRequest"
	//buildRequest.APIVersion = "build.openshift.io/v1"
	//objectMeta := metav1.ObjectMeta{}
	//objectMeta.Name = "cakephp-ex"
	//buildRequest.ObjectMeta = objectMeta
	//buildTriggerCause := v1.BuildTriggerCause{}
	//buildTriggerCause.Message = "Manually triggered"
	//buildRequest.TriggeredBy = []v1.BuildTriggerCause{buildTriggerCause}
	//myBuild, err = buildV1Client.BuildConfigs(namespace).Instantiate(context.TODO(), objectMeta.Name, &buildRequest, metav1.CreateOptions{})
	//
	//if err != nil {
	//	return err
	//}
	//fmt.Printf("Name of the triggered build %s\n", myBuild.Name)

	// Instantiate loader for kubeconfig file.
	kubeconfig1 := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	// Get a rest.Config from the kubeconfig file.  This will be passed into all
	// the client objects we create.
	restconfig, err := kubeconfig1.ClientConfig()
	if err != nil {
		panic(err)
	}

	// Create a Kubernetes core/v1 client.
	coreclient, err := corev1client.NewForConfig(restconfig)
	if err != nil {
		panic(err)
	}
	// List all Pods in our current Namespace.
	pods, err := coreclient.Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Pods in namespace %s:\n", namespace)
	for _, pod := range pods.Items {
		fmt.Printf("  %s\n", pod.Name)
	}

	// List all Pods in our current Namespace.
	pods1, err := coreclient.Pods(namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: "app=json-server",
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Pods in namespace %s:\n", namespace)
	for _, pod := range pods1.Items {
		fmt.Printf("  %s\n", pod.Name)
	}

	err = coreclient.Pods(namespace).Delete(context.Background(), "json-server-57bbd69859-bcshr", metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	}

	return nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
