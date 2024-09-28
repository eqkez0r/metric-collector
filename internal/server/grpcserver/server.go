package grpc_server

import (
	"context"
	"github.com/Eqke/metric-collector/internal/server/grpcserver/interceptors"
	store "github.com/Eqke/metric-collector/internal/storage"
	"github.com/Eqke/metric-collector/pkg/metric"
	pb "github.com/eqkez0r/metric-collector-grpc-api/grpc/metric_collector"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"strconv"
	"sync"
)

type StoreProvider interface {
	SetMetric(context.Context, metric.Metrics) error
	SetMetrics(context.Context, []metric.Metrics) error
	GetMetric(context.Context, metric.Metrics) (metric.Metrics, error)
	GetMetrics(context.Context) (map[string][]store.Metric, error)
}

type GRPCServer struct {
	logger     *zap.SugaredLogger
	store      StoreProvider
	grpcServer *grpc.Server
	host       string

	pb.UnimplementedMetricCollectorServer
}

func New(
	logger *zap.SugaredLogger,
	store StoreProvider,
	host string,
) *GRPCServer {
	grpcserver := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.LoggerInterceptor(logger),
		))
	server := &GRPCServer{
		logger:     logger.Named("grpc-server"),
		store:      store,
		host:       host,
		grpcServer: grpcserver,
	}
	pb.RegisterMetricCollectorServer(grpcserver, server)
	return server
}

func (g *GRPCServer) Run(ctx context.Context, wg *sync.WaitGroup) {
	g.logger.Info("Starting gRPC server")

	listen, err := net.Listen("tcp", g.host)
	if err != nil {
		g.logger.Errorw("Failed to listen", "host", g.host, "error", err)
		return
	}
	reflection.Register(g.grpcServer)
	go func() {
		err = g.grpcServer.Serve(listen)
		if err != nil {
			g.logger.Errorw("Failed to start gRPC server", "error", err)
		}
		wg.Done()
	}()
	defer g.grpcServer.GracefulStop()
	<-ctx.Done()
}

func (g *GRPCServer) Shutdown() {
	g.grpcServer.GracefulStop()
}

func (g *GRPCServer) ReceiveMetric(ctx context.Context, req *pb.ReceiveMetricRequest) (*pb.ReceiveMetricResponse, error) {
	const op = "grpcServer.ReceiveMetric"

	g.logger.Info("Receive metric request")
	m := metric.Metrics{
		ID:    req.Metric.MetricName,
		MType: req.Metric.MetricType,
		Value: req.Metric.Value,
		Delta: req.Metric.Delta,
	}
	g.logger.Infof("Receive metric: %v", m)
	err := g.store.SetMetric(ctx, m)
	if err != nil {
		g.logger.Error(op, err)
		return nil, err
	}
	g.logger.Info("Metric stored successfully")
	return &pb.ReceiveMetricResponse{}, nil
}

func (g *GRPCServer) ReceiveMetricBatch(ctx context.Context, req *pb.ReceiveMetricBatchRequest) (*pb.ReceiveMetricResponse, error) {
	const op = "grpcServer.ReceiveMetricBatch"
	g.logger.Info("Receive metric batch request")
	ms := make([]metric.Metrics, 0, len(req.Metrics))
	for _, m := range req.Metrics {
		ms = append(ms, metric.Metrics{
			ID:    m.MetricName,
			MType: m.MetricType,
			Value: m.Value,
			Delta: m.Delta,
		})
	}
	g.logger.Infof("Receive metric batch: %v", ms)
	err := g.store.SetMetrics(ctx, ms)
	if err != nil {
		g.logger.Error(op, err)
		return nil, err
	}
	g.logger.Info("Metric stored successfully")
	return &pb.ReceiveMetricResponse{}, nil
}

func (g *GRPCServer) ReadMetric(ctx context.Context, req *pb.ReadMetricRequest) (*pb.ReadMetricResponse, error) {
	const op = "grpcServer.ReadMetric"
	g.logger.Infof("Read metric request")
	m := metric.Metrics{
		ID:    req.MetricName,
		MType: req.MetricType,
	}
	newM, err := g.store.GetMetric(ctx, m)
	if err != nil {
		g.logger.Error(op, err)
		return nil, err
	}
	g.logger.Info("Metric read success")
	return &pb.ReadMetricResponse{
		Metric: &pb.Metric{
			MetricName: newM.ID,
			MetricType: newM.MType,
			Delta:      newM.Delta,
			Value:      newM.Value,
		},
	}, nil
}

func (g *GRPCServer) ReadAllMetric(ctx context.Context, _ *pb.ReadAllMetricRequest) (*pb.ReadAllMetricResponse, error) {
	const op = "grpcServer.ReadAllMetric"
	mp, err := g.store.GetMetrics(ctx)
	if err != nil {
		g.logger.Error(op, err)
		return nil, err
	}
	metrics := make([]*pb.Metric, 0, len(mp))
	for k, m := range mp {
		for _, mm := range m {
			metric := pb.Metric{
				MetricName: mm.Name,
				MetricType: k,
			}
			switch k {
			case "gauge":
				value, err := strconv.ParseFloat(mm.Value, 64)
				if err != nil {
					g.logger.Error(op, err)
					continue
				}
				metric.Value = &value
			case "counter":
				counter, err := strconv.ParseInt(mm.Value, 10, 64)
				if err != nil {
					g.logger.Error(op, err)
					continue
				}
				metric.Delta = &counter
			}
			metrics = append(metrics, &metric)
		}
	}
	g.logger.Infof("Read all metrics: %v", metrics)
	return &pb.ReadAllMetricResponse{Metrics: metrics}, nil
}
