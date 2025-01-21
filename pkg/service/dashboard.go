package service

import (
	helper "ahava/pkg/helper"
)

type DashboardService interface {
}

type dashboardService struct {
	helper helper.Helper
}

func NewDashboardService(h helper.Helper) DashboardService {
	return &adminService{
		helper: h,
	}
}
