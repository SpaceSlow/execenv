package client

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/SpaceSlow/execenv/internal/config"
	"github.com/SpaceSlow/execenv/internal/metrics"
	pb "github.com/SpaceSlow/execenv/internal/proto"
	"github.com/SpaceSlow/execenv/internal/utils"
)

var _ Sender = (*grpcSender)(nil)

type grpcSender struct {
	addr string
}

func newGrpcSender() (*grpcSender, error) {
	cfg, err := config.GetAgentConfig()
	if err != nil {
		return nil, err
	}
	s := &grpcSender{
		addr: cfg.ServerAddr.String(),
	}

	return s, nil
}

func (s *grpcSender) Send(metrics []metrics.Metric) error {
	cfg, err := config.GetAgentConfig()
	if err != nil {
		return err
	}
	md := metadata.New(map[string]string{"X-Real-IP": cfg.LocalIP})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	pbMetrics := make([]*pb.Metric, 0, len(metrics))
	for _, m := range metrics {
		metric, err := pb.ConvertToProto(&m)
		if err != nil {
			return err
		}
		pbMetrics = append(pbMetrics, metric)
	}

	conn, err := grpc.NewClient(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	c := pb.NewMetricServiceClient(conn)

	resCh := make(chan *pb.BatchAddMetricsResponse, 1)
	defer close(resCh)
	sendMetrics := func() error {
		var res *pb.BatchAddMetricsResponse
		res, err = c.BatchAddMetrics(ctx, &pb.BatchAddMetricsRequest{Metrics: pbMetrics})
		if err != nil {
			if len(resCh) > 0 {
				<-resCh
			}
			resCh <- res
			return err
		}
		if len(resCh) > 0 {
			<-resCh
		}
		resCh <- res
		return nil
	}

	return <-utils.RetryFunc(sendMetrics, cfg.Delays)
}
