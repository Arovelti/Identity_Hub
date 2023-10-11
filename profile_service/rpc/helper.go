package rpc

import (
	pb "github.com/Arovelti/identityhub/api/proto"
	"github.com/Arovelti/identityhub/profile_service/models"
	"github.com/google/uuid"
)

func profileToProtobuf(p *models.Profile) *pb.Profile {
	return &pb.Profile{
		Id:       p.ID.String(),
		Name:     p.Name,
		Email:    p.Email,
		Password: p.Password,
		Admin:    p.Admin,
	}
}

func profileFromWithoutIDProtoBuf(in *pb.Profile) *models.Profile {
	return &models.Profile{
		Name:     in.Name,
		Email:    in.Email,
		Password: in.Password,
		Admin:    in.Admin,
	}
}

func profileFromWithIDProtoBuf(in *pb.Profile) *models.Profile {
	return &models.Profile{
		ID:       uuid.MustParse(in.Id),
		Name:     in.Name,
		Email:    in.Email,
		Password: in.Password,
		Admin:    in.Admin,
	}
}
