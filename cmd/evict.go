package cmd

import (
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"log"
	"time"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

//Version is set during build time
var Version = "unknown"

type EvictCmdOptions struct {
	configFlags *genericclioptions.ConfigFlags
	iostreams   genericclioptions.IOStreams

	args         []string
	namespace    string
	kubeclient   kubernetes.Interface
	printVersion bool
	graceSecs    int64
	cordonSecs   int64
}

func NewEvictCmdOptions(streams genericclioptions.IOStreams) *EvictCmdOptions {
	return &EvictCmdOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
		iostreams:   streams,
	}
}

func NewCmdEvict(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewEvictCmdOptions(streams)

	cmd := &cobra.Command{
		Use:          "evict pod ",
		Short:        "evict pod for low-version kubernetes, example 1.15.x",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if o.printVersion {
				fmt.Println(Version)
				os.Exit(0)
			}

			if err := o.Complete(c, args); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			if err := o.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().Int64VarP(&o.graceSecs, "grace-secords", "g", 120, "grace period seconds")
	cmd.Flags().Int64VarP(&o.cordonSecs, "cordon-secords", "c", 1, "cordon seconds")
	o.configFlags.AddFlags(cmd.Flags())

	return cmd
}

func (o *EvictCmdOptions) Run() error {

	for _, podName := range o.args {
		o.evictPod(o.namespace, podName)
	}

	return nil
}

func (o *EvictCmdOptions) evictPod(namespace string, podName string) {

	coreV1 := o.kubeclient.CoreV1()

	pods, err := coreV1.Pods(namespace).List(metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", podName),
	})
	ifError(err)

	nodeName := ""
	for _, item := range pods.Items {
		nodeName = item.Spec.NodeName
	}

	_, err = coreV1.Nodes().Patch(nodeName, types.StrategicMergePatchType, []byte("{\"spec\":{\"unschedulable\":true}}"))
	ifError(err)
	log.Printf("cordon %s", nodeName)

	err = coreV1.Pods(o.namespace).Delete(podName, &metav1.DeleteOptions{
		GracePeriodSeconds: &o.graceSecs,
	})
	ifError(err)
	log.Printf("deleting pod %s", podName)

	time.Sleep(time.Duration(o.cordonSecs) * time.Second)
	_, err = coreV1.Nodes().Patch(nodeName, types.StrategicMergePatchType, []byte("{\"spec\":{\"unschedulable\":null}}"))
	ifError(err)
	log.Printf("uncordon %s", nodeName)
}

func (o *EvictCmdOptions) Complete(cmd *cobra.Command, args []string) error {
	o.args = args

	config, err := o.configFlags.ToRESTConfig()
	if err != nil {
		return err
	}

	o.kubeclient, err = kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	o.namespace, _, err = o.configFlags.ToRawKubeConfigLoader().Namespace()
	if err != nil || len(o.namespace) == 0 {
		return err
	}

	return nil
}

func (o *EvictCmdOptions) Validate() error {
	if len(o.args) < 1 {
		return fmt.Errorf("arguments required")
	}

	return nil
}

func ifError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
