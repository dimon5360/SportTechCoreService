package app

import (
	"app/main/internal/endpoint"
	authEndpoint "app/main/internal/endpoint/auth"
	profileEndpoint "app/main/internal/endpoint/profile"
	reportEndpoint "app/main/internal/endpoint/report"
	userEndpoint "app/main/internal/endpoint/user"
	middleware "app/main/internal/middleware/jwt"
	"app/main/internal/repository"
	profileRepository "app/main/internal/repository/profile"
	redisRepository "app/main/internal/repository/redis"
	reportRepository "app/main/internal/repository/report"
	userRepository "app/main/internal/repository/user"
	"app/main/internal/service"
	userService "app/main/internal/service/user"
	"app/main/pkg/env"
	"log"
	"os"
)

const (
	redisEnvKey = "REDIS_ENV"
	mongoEnvKey = "MONGO_ENV"
)

type ServiceProvider struct {
	service service.Interface

	eUser    endpoint.Interface
	eProfile endpoint.Interface
	eReport  endpoint.Interface
	eAuth    endpoint.Interface

	rUser    repository.Interface
	rProfile repository.Interface
	rReport  repository.Interface
	rToken   repository.Interface
}

func NewServiceProvider() *ServiceProvider {
	return &ServiceProvider{}
}

func (sp *ServiceProvider) Init() *ServiceProvider {
	env.Init()

	redisEnv := os.Getenv(redisEnvKey)
	if len(redisEnv) == 0 {
		log.Fatal("redis environment not found")
	}

	mongoEnv := os.Getenv(mongoEnvKey)
	if len(mongoEnv) == 0 {
		log.Fatal("mongo environment not found")
	}

	env.LoadFile(redisEnv, mongoEnv)
	sp.initUserService()
	return sp
}

func (sp *ServiceProvider) initUserService() {
	sp.service = userService.New(middleware.New(sp.redisRepository()),
		sp.getUserEndpoint(),
		sp.getProfileEndpoint(),
		sp.getReportEndpoint(),
		sp.getAuthEndpoint())

	if err := sp.service.Init(); err != nil {

		log.Fatal(err)
	}
}

func (sp *ServiceProvider) getUserEndpoint() endpoint.Interface {
	sp.eUser = userEndpoint.New(sp.userRepository(), sp.profileRepository())
	if sp.eUser == nil {
		log.Fatal("Failed endpoint creation")
	}
	return sp.eUser
}

func (sp *ServiceProvider) getProfileEndpoint() endpoint.Interface {
	if sp.eProfile == nil {
		sp.eProfile = profileEndpoint.New(sp.profileRepository())
	}
	return sp.eProfile
}

func (sp *ServiceProvider) getReportEndpoint() endpoint.Interface {
	if sp.eReport == nil {
		sp.eReport = reportEndpoint.New(sp.reportRepository())
	}
	return sp.eReport
}

func (sp *ServiceProvider) getAuthEndpoint() endpoint.Interface {
	if sp.eAuth == nil {
		sp.eAuth = authEndpoint.New(sp.userRepository(), sp.redisRepository())
	}
	return sp.eAuth
}

func (sp *ServiceProvider) userRepository() repository.Interface {
	if sp.rUser == nil {
		sp.rUser = userRepository.New()
	}
	return sp.rUser
}

func (sp *ServiceProvider) profileRepository() repository.Interface {
	if sp.rProfile == nil {
		sp.rProfile = profileRepository.New()
	}
	return sp.rProfile
}

func (sp *ServiceProvider) reportRepository() repository.Interface {
	if sp.rReport == nil {
		sp.rReport = reportRepository.New()
	}
	return sp.rReport
}

func (sp *ServiceProvider) redisRepository() repository.Interface {
	if sp.rToken == nil {
		sp.rToken = redisRepository.New()
	}
	return sp.rToken
}
