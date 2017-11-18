# BookShop

This repo try to demonstrates building a "real-world" web services in go. It is used only for learning purposes.

## Description

We focus on building the services in monolithic way first. We iterate over to split this monolithic to microservices that can be deployed separately.

Bookshop consists of following services
* Auth  - Login, Signup, ResetPassword, Impersonate and Forgot password
* Catalog - View, filter, search books
* Order  - Place, View and Cancel Orders
* Payment - Add/Edit payment method and Make payment.
* Notification - Email and SMS notifications.

### Roadmap
- [ ] Elegant monolitic exposing REST endpoints for all the services - v1.0
- [ ] Add async task queue support - v2.0
- [ ] Catalog search via elastic search - v3.0
- [ ] Split to microservices. grpc transport between microservices and single API gateway - v4.0
- [ ] Dockerize all the services - v5.0
- [ ] Make services deployable via kubernetes - v6.0
- [ ] Add service discovery support - v7.0
- [ ] Monitoring via Prometheus and Grafana - v8.0
- [ ] Log aggregation via "OK log" - v9.0

### Organization

We will stick to Domain Driven Design as much as possible.

```
bookshop
├── build                  # Compiled files
├── cmd                    # Main entry points
    ├── bookstore
		├── main.go
├── internal               # Domain related packages
    ├── auth
    ├── catalog
    └── order
    └── payment
    └── notification
└── resource               # High-level packages. Can be used by any domain 
	└── db
└── vendor
└── Gopkg.toml
└── Gopkg.yaml
```



