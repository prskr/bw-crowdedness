package main

import (
	"net/http"

	"github.com/jasonlvhit/gocron"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/baez90/bw-crowdedness/internal/bw"
)

var (
	bws = []bw.BW{
		{
			Domain:    "www.boulderwelt-dortmund.de",
			ShortName: "dortmund",
		},
		{
			Domain:    "www.boulderwelt-muenchen-sued.de",
			ShortName: "munich_south",
		},
		{
			Domain:    "www.boulderwelt-muenchen-west.de",
			ShortName: "munich_west",
		},
		{
			Domain:    "www.boulderwelt-muenchen-ost.de",
			ShortName: "munich_east",
		},
		{
			Domain:    "www.boulderwelt-frankfurt.de",
			ShortName: "frankfurt",
		},
		{
			Domain:    "www.boulderwelt-regensburg.de",
			ShortName: "regensburg",
		},
	}

	metricLabels     = []string{"branch"}
	crowdednessGauge *prometheus.GaugeVec
	queueGauge       *prometheus.GaugeVec
	fetchStatTiming  *prometheus.HistogramVec
)

func main() {
	var (
		err    error
		logger *zap.Logger
	)
	logger, _ = zap.NewProduction()
	_ = gocron.Every(5).Second().From(gocron.NextTick()).Do(func() {
		fetchBWStats(logger)
	})
	initMetrics()

	go func() {
		<-gocron.Start()
	}()

	http.Handle("/metrics", promhttp.Handler())
	err = http.ListenAndServe(":9091", nil)
	logger.Info("Stopped HTTP server", zap.Error(err))
}

func fetchBWStats(logger *zap.Logger) {
	logger.Info("Start fetching current metrics")
	for _, instance := range bws {
		instance := instance
		go processBW(logger, instance)
	}
}

func processBW(logger *zap.Logger, instance bw.BW) {
	var (
		branchLogger = logger.With(zap.String("branch", instance.ShortName))
		timer        = prometheus.NewTimer(fetchStatTiming.WithLabelValues(instance.ShortName))
		stats        bw.Stats
		err          error
	)

	defer timer.ObserveDuration()
	if stats, err = bw.StatsForBW(instance.Domain); err != nil {
		branchLogger.Error(
			"failed to collect BW stats",
			zap.Error(err),
		)
		return
	}

	branchLogger.Info("Got current metrics")

	crowdednessGauge.WithLabelValues(instance.ShortName).Set(float64(stats.CrowdednessPercent))
	queueGauge.WithLabelValues(instance.ShortName).Set(float64(stats.Queue))
}

func initMetrics() {
	crowdednessGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "bw",
		Subsystem: "crowdedness",
		Name:      "percentage",
		Help:      "",
	}, metricLabels)

	queueGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "bw",
		Subsystem: "crowdedness",
		Name:      "queue",
		Help:      "",
	}, metricLabels)

	fetchStatTiming = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "bw",
		Subsystem: "crowdedness",
		Name:      "fetch_stat",
		Help:      "",
	}, metricLabels)

	prometheus.MustRegister(crowdednessGauge, queueGauge, fetchStatTiming)
}
