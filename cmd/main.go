package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

func main() {
	var identifier = os.Getenv("POD_NAME")
	log.Println("My name is:", identifier)

	cfg, _ := rest.InClusterConfig()
	clientset, _ := kubernetes.NewForConfig(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var lock = &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      "my-lock",
			Namespace: "default",
		},
		Client: clientset.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: identifier,
		},
	}

	done := make(chan struct{})
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Shutting down")
		cancel()
		close(done)
	}()

	var ticker = time.NewTicker(time.Second)
	defer ticker.Stop()
	var leading int32
	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   2 * time.Second,
		RenewDeadline:   1 * time.Second,
		RetryPeriod:     500 * time.Millisecond,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				atomic.StoreInt32(&leading, 1)
				log.Println(identifier, "started leading")
			},
			OnStoppedLeading: func() {
				atomic.StoreInt32(&leading, 0)
				log.Println(identifier, "has stopped leading")
			},
			OnNewLeader: func(identity string) {
				log.Println(identity, "is the new leader")
			},
		},
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Listening...")
	srv.ListenAndServe()
	<-done
}
