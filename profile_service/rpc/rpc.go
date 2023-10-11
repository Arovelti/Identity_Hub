package rpc

import (
	"context"
	"fmt"

	pb "github.com/Arovelti/identityhub/api/proto"
	"github.com/Arovelti/identityhub/repository"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RPC struct {
	pb.UnimplementedProfilesServer
	profileRepo repository.Repository
}

func NewRPC(repo repository.Repository) *RPC {
	return &RPC{
		profileRepo: repo,
	}
}

func (r *RPC) SignIn(ctx context.Context, in *pb.LoginWithBasicAuthRequest) (*pb.LoginWithBasicAuthResponse, error) {
	// Considering that we have a BasicAuthInterseptor -
	// we don't need this endpoint
	return nil, status.Errorf(codes.Unimplemented, "method SingIn not implemented")
}

func (r *RPC) CreateProfile(ctx context.Context, in *pb.Profile) (*emptypb.Empty, error) {
	p := profileFromWithoutIDProtoBuf(in)
	p.ID = uuid.New()

	err := r.profileRepo.Create(p)
	return &emptypb.Empty{}, err
}

func (r *RPC) GetProfileByID(ctx context.Context, in *pb.GetProfileByIDRequest) (*pb.Profile, error) {
	id := uuid.MustParse(in.Id)
	p, err := r.profileRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get by id: %v", err)
	}

	return profileToProtobuf(p), nil
}

func (r *RPC) GetProfileByUsername(ctx context.Context, in *pb.GetProfileByUsernameRequest) (*pb.Profile, error) {
	p, err := r.profileRepo.GetByUsername(in.Username)
	if err != nil {
		return nil, fmt.Errorf("get by username: %v", err)
	}

	return profileToProtobuf(p), nil
}

func (r *RPC) ListProfiles(ctx context.Context, in *empty.Empty) (*pb.GetProfileListResponse, error) {
	list := r.profileRepo.List()
	profiles := make([]*pb.Profile, 0, len(list))

	for _, p := range list {
		profiles = append(profiles, profileToProtobuf(p))
	}

	return &pb.GetProfileListResponse{Profiles: profiles}, nil
}

func (r *RPC) UpdateProfile(ctx context.Context, in *pb.UpdateProfileRequest) (*pb.Profile, error) {
	// From Request
	id := uuid.MustParse(in.Id)
	name := in.Name

	// From Profile
	p := profileFromWithIDProtoBuf(in.Profile)

	if err := r.profileRepo.Update(id, name, p); err != nil {
		return nil, err
	}

	return in.Profile, nil
}

func (r *RPC) DeleteProfileByID(ctx context.Context, in *pb.DeleteProfileRequest) (*emptypb.Empty, error) {
	err := r.profileRepo.Delete(uuid.MustParse(in.Id))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
