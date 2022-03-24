package provisioner

import (
	"encoding/json"
	"errors"
	"github.com/api7/amesh/pkg/apisix"
	"github.com/golang/protobuf/ptypes/any"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	clusterv3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	endpointv3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
)

func (p *xdsProvisioner) processRouteConfigurationV3(res *any.Any) ([]*apisix.Route, error) {
	var route routev3.RouteConfiguration
	err := anypb.UnmarshalTo(res, &route, proto.UnmarshalOptions{
		DiscardUnknown: true,
	})
	if err != nil {
		p.logger.Errorw("found invalid RouteConfiguration resource",
			zap.Error(err),
			zap.Any("resource", res),
		)
		return nil, err
	}

	routes, err := p.TranslateRouteConfiguration(&route, p.routeOwnership)
	if err != nil {
		p.logger.Errorw("failed to translate RouteConfiguration to APISIX routes",
			zap.Error(err),
			zap.Any("route", &route),
		)
		return nil, err
	}
	return routes, nil
}

func (p *xdsProvisioner) processStaticRouteConfigurations(rcs []*routev3.RouteConfiguration) ([]*apisix.Route, error) {
	var (
		routes []*apisix.Route
	)
	for _, rc := range rcs {
		route, err := p.TranslateRouteConfiguration(rc, p.routeOwnership)
		if err != nil {
			p.logger.Errorw("failed to translate RouteConfiguration to APISIX routes",
				zap.Error(err),
				zap.Any("route", &route),
			)
			return nil, err
		}
	}
	return routes, nil
}

func (p *xdsProvisioner) processClusterV3(res *any.Any) (*apisix.Upstream, error) {
	var cluster clusterv3.Cluster
	err := anypb.UnmarshalTo(res, &cluster, proto.UnmarshalOptions{
		DiscardUnknown: true,
	})
	if err != nil {
		p.logger.Errorw("found invalid Cluster resource",
			zap.Error(err),
			zap.Any("resource", res),
		)
		return nil, err
	}
	ups, err := p.TranslateCluster(&cluster)
	if err != nil {
		return nil, err
	}
	return ups, nil
}

func (p *xdsProvisioner) processClusterLoadAssignmentV3(res *any.Any) (*apisix.Upstream, error) {
	var cla endpointv3.ClusterLoadAssignment
	err := anypb.UnmarshalTo(res, &cla, proto.UnmarshalOptions{
		DiscardUnknown: true,
	})
	if err != nil {
		p.logger.Errorw("found invalid ClusterLoadAssignment resource",
			zap.Error(err),
			zap.Any("resource", res),
		)
		return nil, err
	}

	ups, ok := p.upstreams[cla.ClusterName]
	if !ok {
		p.logger.Warnw("found invalid ClusterLoadAssignment resource",
			zap.String("reason", "cluster unknown"),
			zap.Any("resource", res),
		)
		return nil, errors.New("UnknownClusterName")
	}

	nodes, err := p.TranslateClusterLoadAssignment(&cla)
	if err != nil {
		p.logger.Errorw("failed to translate ClusterLoadAssignment",
			zap.Error(err),
			zap.Any("resource", res),
		)
		return nil, err
	}

	// Do not set on the original ups to avoid race conditions.
	data, err := json.Marshal(ups)
	if err != nil {
		return nil, err
	}
	var newUps apisix.Upstream
	err = json.Unmarshal(data, &newUps)
	if err != nil {
		return nil, err
	}

	newUps.Nodes = nodes
	p.upstreams[cla.ClusterName] = &newUps
	return &newUps, nil
}