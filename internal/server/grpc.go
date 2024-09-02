package server

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

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
	s := grpc.NewServer()
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

	metric, err := convertFromProto(in.Metric)
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
		// TODO not found metric
	}
	var err error
	response.Metric, err = convertToProto(metric)
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
		m, _ := convertToProto(&metric)
		response.Metrics = append(response.Metrics, m)
	}

	return &response, nil
}

func convertFromProto(m *pb.Metric) (*metrics.Metric, error) {
	metric := &metrics.Metric{
		Name: m.Id,
	}

	switch m.MType {
	case pb.MType_COUNTER:
		metric.Type = metrics.Counter
		metric.Value = m.Delta
	case pb.MType_GAUGE:
		metric.Type = metrics.Gauge
		metric.Value = m.Value
	default:
		return nil, metrics.ErrIncorrectMetricTypeOrValue
	}
	return metric, nil
}

func convertToProto(m *metrics.Metric) (*pb.Metric, error) {
	metric := &pb.Metric{
		Id: m.Name,
	}

	switch m.Type {
	case metrics.Counter:
		metric.MType = pb.MType_COUNTER
		delta, ok := m.Value.(int64)
		if !ok {
			return nil, metrics.ErrIncorrectMetricTypeOrValue
		}
		metric.Delta = delta
	case metrics.Gauge:
		metric.MType = pb.MType_GAUGE
		value, ok := m.Value.(float64)
		if !ok {
			return nil, metrics.ErrIncorrectMetricTypeOrValue
		}
		metric.Value = value
	default:
		return nil, metrics.ErrIncorrectMetricTypeOrValue
	}
	return metric, nil
}
