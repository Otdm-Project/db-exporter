package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gocql/gocql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	aliveMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "db_alive",
		Help: "Apache Cassandra monitoring",
	})
)

func main() {
	cluster := gocql.NewCluster("localhost") // CassandraノードのIPアドレス
	cluster.Port = 9042                      // 使いたいKeyspace
	cluster.Consistency = gocql.Quorum       // Consistency Levelの設定
	cluster.Timeout = 5 * time.Second
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		http.ListenAndServe(":9300", nil)
	}()
	for range time.Tick(10 * time.Second) {
		err := connectToCassandra(cluster)
		if err != nil {
			fmt.Println("0", err)
			aliveMetric.Set(0)
		} else {
			fmt.Println("1")
			aliveMetric.Set(1)
		}
	}
}

func connectToCassandra(cluster *gocql.ClusterConfig) error {
	// Cassandraへのセッション接続を試みる
	session, err := cluster.CreateSession()
	if err != nil {
		return err
	}
	defer session.Close()
	return nil
}
