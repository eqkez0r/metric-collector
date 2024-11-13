package grpcagent

import (
	"context"
	"github.com/Eqke/metric-collector/internal/agent/config"
	"github.com/Eqke/metric-collector/internal/agent/poller"
	"github.com/Eqke/metric-collector/pkg/metric"
	pb "github.com/eqkez0r/metric-collector-grpc-api/grpc/metric_collector"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"strconv"
	"sync"
	"time"
)

type GRPCClient struct {
	host   string
	logger *zap.SugaredLogger

	reportInterval time.Duration

	poller poller.MetricPoller
}

func New(
	logger *zap.SugaredLogger,
	settings *config.AgentConfig,
	poller poller.MetricPoller,
) *GRPCClient {
	return &GRPCClient{
		host:           settings.GrpcServerHost,
		logger:         logger,
		reportInterval: time.Second * time.Duration(settings.ReportInterval),
		poller:         poller,
	}
}

func (gc *GRPCClient) Run(ctx context.Context, wg *sync.WaitGroup) {
	reportTicker := time.NewTicker(gc.reportInterval)
	defer reportTicker.Stop()
	for {
		select {
		case <-ctx.Done():
			{
				gc.logger.Infow("grpc agent was stopped")
				wg.Done()
				return
			}
		case <-reportTicker.C:
			{
				gc.logger.Infow("starting report")
				gc.Poll(ctx)
				gc.logger.Infow("finished report")
			}
		default:
			{

			}
		}
	}
}

func (gc *GRPCClient) Poll(ctx context.Context) {
	conn, err := grpc.NewClient(gc.host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		gc.logger.Errorw("failed to connect to grpc server", "host", gc.host, "error", err)
		return
	}
	defer conn.Close()
	grpcConn := pb.NewMetricCollectorClient(conn)
	metricMap := gc.poller.GetMetrics()
	metricList := make([]*pb.Metric, 0, len(metricMap["gauge"])+len(metricMap["counter"]))

	for mt, mm := range metricMap {

		for mn, v := range mm {
			pushMetric := &pb.Metric{
				MetricName: mn.String(),
				MetricType: mt.String(),
			}
			switch mt {
			case metric.TypeGauge:
				{
					value, err := strconv.ParseFloat(v, 64)
					if err != nil {
						gc.logger.Errorw("failed to parse value", "metric", mt, "value", v)
						continue
					}
					pushMetric.Value = &value

				}
			case metric.TypeCounter:
				{
					counter, err := strconv.ParseInt(v, 10, 64)
					if err != nil {
						gc.logger.Errorw("failed to parse value", "metric", mt, "value", v)
						continue
					}
					pushMetric.Delta = &counter
				}
			default:
				{
					gc.logger.Errorw("unknown metric type", "type", mt)
					continue
				}
			}
			_, err = grpcConn.ReceiveMetric(ctx, &pb.ReceiveMetricRequest{Metric: pushMetric})
			if err != nil {
				gc.logger.Errorw("failed to send metric", "metric", pushMetric.MetricName, "value", v)
				continue
			}
			gc.logger.Infow("send metric success")
			metricList = append(metricList, pushMetric)
		}

	}
	_, err = grpcConn.ReceiveMetricBatch(ctx, &pb.ReceiveMetricBatchRequest{
		Metrics: metricList,
	})
	if err != nil {
		gc.logger.Errorw("failed to send metric batch", "metricList", metricList)
		return
	}
	gc.logger.Infow("send metric batch success")
}
