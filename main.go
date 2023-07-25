package main

import (
	"context"
	"io/ioutil"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var s = gocron.NewScheduler(time.UTC).SingletonMode()

type Config struct {
	Version              string   `yaml:"version"`
	Namespace            string   `yaml:"namespace"`
	Name                 string   `yaml:"hpaName"`
	Schedule             Schedule `yaml:"schedule"`
	ScaleUpMaxReplicas   int32    `yaml:"scaleUpMaxReplicas"`
	ScaleUpMinReplicas   int32    `yaml:"scaleUpMinReplicas"`
	ScaleDownMaxReplicas int32    `yaml:"scaleDownMaxReplicas"`
	ScaleDownMinReplicas int32    `yaml:"scaleDownMinReplicas"`
}

type Schedule struct {
	Enabled       bool   `yaml:"enabled"`
	ScaleUpTime   string `yaml:"scaleUpTime"`
	ScaleDownTime string `yaml:"scaleDownTime"`
}

func main() {
	c := readConfig()
	log.Printf("starting custom HPA scheduler: %s\n", c.Version)

	saToken, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	client, err := kubernetes.NewForConfig(saToken)
	if err != nil {
		panic(err.Error())
	}

	ctx := context.Background()

	_, err = s.Every(1).Day().At(c.Schedule.ScaleUpTime).Do(scale, ctx, client, "up")
	_, err = s.Every(1).Day().At(c.Schedule.ScaleDownTime).Do(scale, ctx, client, "down")

	if err != nil {
		panic(err)
	}
	s.StartBlocking()

}

func readConfig() Config {
	config := Config{}
	confBytes, err := ioutil.ReadFile("./config.yaml")

	if err != nil {
		panic(err.Error())
	}

	err = yaml.Unmarshal(confBytes, &config)
	if err != nil {
		panic(err.Error())
	}

	return config
}

func scale(ctx context.Context, client *kubernetes.Clientset, direction string) error {
	c := readConfig()

	if !c.Schedule.Enabled {
		return nil
	}

	hpa, err := client.AutoscalingV2().HorizontalPodAutoscalers(c.Namespace).Get(ctx, c.Name, metav1.GetOptions{})

	if err != nil {
		return err
	}

	minReplicas := int32(1)
	maxReplicas := int32(1)

	switch direction {
	case "up":
		minReplicas = c.ScaleUpMinReplicas
		maxReplicas = c.ScaleUpMaxReplicas

	case "down":
		minReplicas = c.ScaleDownMinReplicas
		maxReplicas = c.ScaleDownMaxReplicas
	}

	hpa.Spec.MinReplicas = &minReplicas
	hpa.Spec.MaxReplicas = maxReplicas

	_, err = client.AutoscalingV2().HorizontalPodAutoscalers(c.Namespace).Update(ctx, hpa, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	log.Printf("horizontal pod autoscaler updated, name: %s, minReplicas: %d, maxReplicas: %d", c.Name, minReplicas, maxReplicas)
	return nil
}
