package internal

import (
	"context"
	"strconv"

	porter "github.com/tuihub/protos/pkg/librarian/porter/v1"
	librarian "github.com/tuihub/protos/pkg/librarian/v1"
	"github.com/tuihub/tuihub-go"
	"github.com/tuihub/tuihub-steam/internal/biz"
)

type Handler struct {
	porter.UnimplementedLibrarianPorterServiceServer
	steam *biz.SteamUseCase
}

func NewHandler(apiKey string) *Handler {
	return &Handler{
		UnimplementedLibrarianPorterServiceServer: porter.UnimplementedLibrarianPorterServiceServer{},
		steam: biz.NewSteamUseCase(apiKey),
	}
}

func (h Handler) PullAccount(ctx context.Context, req *porter.PullAccountRequest) (
	*porter.PullAccountResponse, error) {
	u, err := h.steam.GetUser(ctx, req.GetAccountId().GetPlatformAccountId())
	if err != nil {
		return nil, err
	}
	return &porter.PullAccountResponse{Account: &librarian.Account{
		Id:                nil,
		Platform:          req.GetAccountId().GetPlatform(),
		PlatformAccountId: req.GetAccountId().GetPlatformAccountId(),
		Name:              u.Name,
		ProfileUrl:        u.ProfileURL,
		AvatarUrl:         u.AvatarURL,
		LatestUpdateTime:  nil,
	}}, nil
}

func (h Handler) PullAppInfo(ctx context.Context, req *porter.PullAppInfoRequest) (
	*porter.PullAppInfoResponse, error) {
	appID, err := strconv.Atoi(req.GetAppInfoId().GetSourceAppId())
	if err != nil {
		return nil, err
	}
	a, err := h.steam.GetAppDetails(ctx, appID)
	if err != nil {
		return nil, err
	}
	return &porter.PullAppInfoResponse{AppInfo: &librarian.AppInfo{
		Id:          nil,
		Internal:    false,
		Source:      req.GetAppInfoId().GetSource(),
		SourceAppId: req.GetAppInfoId().GetSourceAppId(),
		SourceUrl:   &a.StoreURL,
		Details: &librarian.AppInfoDetails{ // TODO
			Description: a.Description,
			ReleaseDate: a.ReleaseDate,
			Developer:   a.Developer,
			Publisher:   a.Publisher,
			Version:     "",
			ImageUrls:   nil,
		},
		Name:               a.Name,
		Type:               ToPBAppType(a.Type),
		ShortDescription:   a.ShortDescription,
		IconImageUrl:       "",
		BackgroundImageUrl: a.BackgroundImageURL,
		CoverImageUrl:      a.CoverImageURL,
		Tags:               nil,
		AltNames:           nil,
	}}, nil
}

func (h Handler) PullAccountAppInfoRelation(ctx context.Context, req *porter.PullAccountAppInfoRelationRequest) (
	*porter.PullAccountAppInfoRelationResponse, error) {
	al, err := h.steam.GetOwnedGames(ctx, req.GetAccountId().GetPlatformAccountId())
	if err != nil {
		return nil, err
	}
	appList := make([]*librarian.AppInfo, len(al))
	for i, a := range al {
		appList[i] = &librarian.AppInfo{ // TODO
			Id:       nil,
			Internal: false,
			Source: tuihub.WellKnownToString(
				librarian.WellKnownAppInfoSource_WELL_KNOWN_APP_INFO_SOURCE_STEAM,
			),
			SourceAppId:        strconv.Itoa(int(a.AppID)),
			SourceUrl:          nil,
			Details:            nil,
			Name:               a.Name,
			Type:               0,
			ShortDescription:   "",
			IconImageUrl:       a.IconImageURL,
			BackgroundImageUrl: a.BackgroundImageURL,
			CoverImageUrl:      a.CoverImageURL,
			Tags:               nil,
			AltNames:           nil,
		}
	}
	return &porter.PullAccountAppInfoRelationResponse{AppInfos: appList}, nil
}

func ToPBAppType(t biz.AppType) librarian.AppType {
	switch t {
	case biz.AppTypeGame:
		return librarian.AppType_APP_TYPE_GAME
	default:
		return librarian.AppType_APP_TYPE_UNSPECIFIED
	}
}

func (h Handler) SearchAppInfo(ctx context.Context, req *porter.SearchAppInfoRequest) (
	*porter.SearchAppInfoResponse, error) {
	al, err := h.steam.SearchAppByName(ctx, req.GetName())
	if err != nil {
		return nil, err
	}
	appList := make([]*librarian.AppInfo, len(al))
	for i, a := range al {
		appList[i] = &librarian.AppInfo{ // TODO
			Id:       nil,
			Internal: false,
			Source: tuihub.WellKnownToString(
				librarian.WellKnownAppInfoSource_WELL_KNOWN_APP_INFO_SOURCE_STEAM,
			),
			SourceAppId:        strconv.Itoa(int(a.AppID)),
			SourceUrl:          nil,
			Details:            nil,
			Name:               a.Name,
			Type:               0,
			ShortDescription:   "",
			IconImageUrl:       a.IconImageURL,
			BackgroundImageUrl: a.BackgroundImageURL,
			CoverImageUrl:      a.CoverImageURL,
			Tags:               nil,
			AltNames:           nil,
		}
	}
	return &porter.SearchAppInfoResponse{AppInfos: appList}, nil
}
