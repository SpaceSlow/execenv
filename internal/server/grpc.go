package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/SpaceSlow/execenv/internal/interceptors"
	"github.com/SpaceSlow/execenv/internal/metrics"
	pb "github.com/SpaceSlow/execenv/internal/proto"
	"github.com/SpaceSlow/execenv/internal/storages"
)

var _ ShutdownRunner = (*grpcStrategy)(nil)

type grpcStrategy struct {
	srv      *grpc.Server
	listener *net.Listener
}

func newGrpcStrategy(address string, storage storages.MetricStorage) *grpcStrategy {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err) // ?
	}
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.LogUnaryInterceptor,
			interceptors.WithCheckingTrustedSubnetUnaryInterceptor,
		),
	)
	pb.RegisterMetricServiceServer(s, &MetricServiceServer{storage: storage})

	runner := &grpcStrategy{
		srv:      s,
		listener: &listen,
	}

	return runner
}

func (s grpcStrategy) Run() error {
	return s.srv.Serve(*s.listener)
}

func (s grpcStrategy) Shutdown(_ context.Context) error {
	s.srv.GracefulStop()
	return nil
}

// MetricServiceServer поддерживает все необходимые методы сервера.
type MetricServiceServer struct {
	pb.UnimplementedMetricServiceServer

	storage storages.MetricStorage
}

// AddMetric реализует интерфейс добавления метрики.
func (s *MetricServiceServer) AddMetric(ctx context.Context, in *pb.AddMetricRequest) (*pb.AddMetricResponse, error) {
	var response pb.AddMetricResponse

	metric, err := pb.ConvertFromProto(in.Metric)
	if err != nil {
		response.Error = err.Error()
		return &response, nil
	}
	_, err = s.storage.Add(metric)
	if err != nil {
		response.Error = err.Error()
	}

	return &response, nil
}

// BatchAddMetrics реализует интерфейс добавления нескольких метрик.
func (s *MetricServiceServer) BatchAddMetrics(ctx context.Context, in *pb.BatchAddMetricsRequest) (*pb.BatchAddMetricsResponse, error) {
	var response pb.BatchAddMetricsResponse

	metricSlice := make([]metrics.Metric, 0, len(in.Metrics))

	for _, metric := range in.Metrics {
		m, err := pb.ConvertFromProto(metric)
		if err != nil {
			response.Error = err.Error()
			return &response, nil
		}
		metricSlice = append(metricSlice, *m)
	}

	err := s.storage.Batch(metricSlice)
	if err != nil {
		response.Error = err.Error()
	}

	return &response, nil
}

// GetMetric реализует интерфейс получения метрики.
func (s *MetricServiceServer) GetMetric(ctx context.Context, in *pb.GetMetricRequest) (*pb.GetMetricResponse, error) {
	var (
		response pb.GetMetricResponse
		metric   *metrics.Metric
		ok       bool
	)

	switch in.MType {
	case pb.MType_COUNTER:
		metric, ok = s.storage.Get(metrics.Counter, in.Id)
	case pb.MType_GAUGE:
		metric, ok = s.storage.Get(metrics.Gauge, in.Id)
	default:
		metric, ok = nil, false
	}

	if !ok {
		response.Error = status.Error(codes.NotFound, "").Error()
		return &response, nil
	}
	var err error
	response.Metric, err = pb.ConvertToProto(metric)
	if err != nil {
		response.Error = err.Error()
	}

	return &response, nil
}

// ListMetrics реализует интерфейс получения метрики.
func (s *MetricServiceServer) ListMetrics(ctx context.Context, in *pb.ListMetricsRequest) (*pb.ListMetricsResponse, error) {
	var response pb.ListMetricsResponse

	metricSlice := s.storage.List()
	response.Metrics = make([]*pb.Metric, 0, len(metricSlice))
	for _, metric := range metricSlice {
		m, _ := pb.ConvertToProto(&metric)
		response.Metrics = append(response.Metrics, m)
	}

	return &response, nil
}
