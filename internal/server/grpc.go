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

	metric, err := convert(in.Metric)
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

func convert(m *pb.Metric) (*metrics.Metric, error) {
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
